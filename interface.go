package main

import (
	"fmt"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/streamdeckd/handlers"
	"golang.org/x/sync/semaphore"
	"image"
	"image/draw"
	"log"
	"os"
	"os/exec"
	"syscall"
)


var sem = semaphore.NewWeighted(int64(1))

func LoadImage(dev *VirtualDev, path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return api.ResizeImage(img, int(dev.Deck.Pixels)), nil
}

func SetKeyImage(dev *VirtualDev, currentKey *api.Key, keyIndex int, page int) {
	if currentKey.Buff == nil {
		if currentKey.Icon == "" {
			img := image.NewRGBA(image.Rect(0, 0, int(dev.Deck.Pixels), int(dev.Deck.Pixels)))
			draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)
			currentKey.Buff = img
		} else {
			img, err := LoadImage(dev, currentKey.Icon)
			if err != nil {
				log.Println(err)
				return
			}
			currentKey.Buff = img
		}
		if currentKey.Text != "" {
			img, err := api.DrawText(currentKey.Buff, currentKey.Text, currentKey.TextSize, currentKey.TextAlignment)
			if err != nil {
				log.Println(err)
			} else {
				currentKey.Buff = img
			}
		}
	}
	if currentKey.Buff != nil {
		dev.SetImage(currentKey.Buff, keyIndex, page)
	}
}

func SetKey(dev *VirtualDev, currentKey *api.Key, keyIndex int, page int) {
	var deckInfo api.StreamDeckInfo
	for i := range sDInfo {
		if sDInfo[i].Serial == dev.Deck.Serial {
			deckInfo = sDInfo[i]
		}
	}
	if currentKey.Buff == nil {
		if currentKey.IconHandler == "" {
			SetKeyImage(dev, currentKey, keyIndex, page)

		} else if currentKey.IconHandlerStruct == nil {
			var handler api.IconHandler
			modules := handlers.AvailableModules()
			for _, module := range modules {
				if module.Name == currentKey.IconHandler {
					handler = module.NewIcon()
				}
			}
			if handler == nil {
				return
			}
			log.Printf("Created & Started %s\n", currentKey.IconHandler)
			handler.Start(*currentKey, deckInfo, func(image image.Image) {
				if image.Bounds().Max.X != int(dev.Deck.Pixels) || image.Bounds().Max.Y != int(dev.Deck.Pixels) {
					image = api.ResizeImage(image, int(dev.Deck.Pixels))
				}
				dev.SetImage(image, keyIndex, page)
				currentKey.Buff = image
			})
			currentKey.IconHandlerStruct = handler
		}
	} else {
		dev.SetImage(currentKey.Buff, keyIndex, page)
	}
	if currentKey.IconHandlerStruct != nil && !currentKey.IconHandlerStruct.IsRunning() {
		log.Printf("Started %s\n", currentKey.IconHandler)
		currentKey.IconHandlerStruct.Start(*currentKey, deckInfo, func(image image.Image) {
			if image.Bounds().Max.X != int(dev.Deck.Pixels) || image.Bounds().Max.Y != int(dev.Deck.Pixels) {
				image = api.ResizeImage(image, int(dev.Deck.Pixels))
			}
			dev.SetImage(image, keyIndex, page)
			currentKey.Buff = image
		})
	}
}

func runCommand(command string) {
	go func() {
		cmd := exec.Command("/bin/sh", "-c", "/usr/bin/nohup "+command)

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid:   true,
			Pgid:      0,
			Pdeathsig: syscall.SIGHUP,
		}
		if err := cmd.Start(); err != nil {
			fmt.Println("There was a problem running ", command, ":", err)
		} else {
			pid := cmd.Process.Pid
			cmd.Process.Release()
			fmt.Println(command, " has been started with pid", pid)
		}
	}()
}

func HandleInput(dev *VirtualDev, key *api.Key, page int) {
	if key.Command != "" {
		runCommand(key.Command)
	}
	if key.Keybind != "" {
		runCommand("xdotool key " + key.Keybind)
	}
	if key.SwitchPage != 0 {
		page = key.SwitchPage - 1
		dev.SetPage(page)
	}
	if key.Brightness != 0 {
		err := dev.SetBrightness(uint8(key.Brightness))
		if err != nil {
			log.Println(err)
		}
	}
	if key.Url != "" {
		runCommand("xdg-open " + key.Url)
	}
	if key.KeyHandler != "" {
		var deckInfo api.StreamDeckInfo
		found := false
		for i := range sDInfo {
			if sDInfo[i].Serial == dev.Deck.Serial {
				deckInfo = sDInfo[i]
				found = true
			}
		}
		if !found {
			return
		}
		if key.KeyHandlerStruct == nil {
			var handler api.KeyHandler
			modules := handlers.AvailableModules()
			for _, module := range modules {
				if module.Name == key.KeyHandler {
					handler = module.NewKey()
				}
			}
			if handler == nil {
				return
			}
			key.KeyHandlerStruct = handler
		}
		key.KeyHandlerStruct.Key(*key, deckInfo)
	}
}
