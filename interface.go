package main

import (
	"context"
	"github.com/unix-streamdeck/api"
	_ "github.com/unix-streamdeck/driver"
	"github.com/unix-streamdeck/streamdeckd/handlers"
	"golang.org/x/sync/semaphore"
	"image"
	"image/draw"
	"log"
	"os"
	"strings"
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

func SetImage(dev *VirtualDev, img image.Image, i int, page int) {
	ctx := context.Background()
	err := sem.Acquire(ctx, 1)
	if err != nil {
		log.Println(err)
		return
	}
	defer sem.Release(1)
	if dev.Page == page && dev.IsOpen {
		err := dev.Deck.SetImage(uint8(i), img)
		if err != nil {
			if strings.Contains(err.Error(), "hidapi") {
				disconnect(dev)
			} else if strings.Contains(err.Error(), "dimensions") {
				log.Println(err)
			}else {
				log.Println(err)
			}
		}
	}
}

func SetKeyImage(dev *VirtualDev, currentKey *api.Key, i int, page int) {
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
		SetImage(dev, currentKey.Buff, i, page)
	}
}

func SetPage(dev *VirtualDev, page int) {
	if page != dev.Page {
		unmountPageHandlers(dev.Config[dev.Page])
	}
	dev.Page = page
	currentPage := dev.Config[page]
	for i := 0; i < len(currentPage); i++ {
		currentKey := &currentPage[i]
		go SetKey(dev, currentKey, i, page)
	}
	EmitPage(dev, page)
}

func SetKey(dev *VirtualDev, currentKey *api.Key, i int, page int) {
	var deckInfo api.StreamDeckInfo
	for i := range sDInfo {
		if sDInfo[i].Serial == dev.Deck.Serial {
			deckInfo = sDInfo[i]
		}
	}
	if currentKey.Buff == nil {
		if currentKey.IconHandler == "" {
			SetKeyImage(dev, currentKey, i, page)

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
				SetImage(dev, image, i, page)
				currentKey.Buff = image
			})
			currentKey.IconHandlerStruct = handler
		}
	} else {
		SetImage(dev, currentKey.Buff, i, page)
	}
	if currentKey.IconHandlerStruct != nil && !currentKey.IconHandlerStruct.IsRunning() {
		log.Printf("Started %s\n", currentKey.IconHandler)
		currentKey.IconHandlerStruct.Start(*currentKey, deckInfo, func(image image.Image) {
			if image.Bounds().Max.X != int(dev.Deck.Pixels) || image.Bounds().Max.Y != int(dev.Deck.Pixels) {
				image = api.ResizeImage(image, int(dev.Deck.Pixels))
			}
			SetImage(dev, image, i, page)
			currentKey.Buff = image
		})
	}
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
		SetPage(dev, page)
	}
	if key.Brightness != 0 {
		err := dev.Deck.SetBrightness(uint8(key.Brightness))
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

func Listen(dev *VirtualDev) {
	kch, err := dev.Deck.ReadKeys()
	if err != nil {
		log.Println(err)
	}
	for dev.IsOpen {
		select {
		case k, ok := <-kch:
			if !ok {
				disconnect(dev)
				return
			}
			if k.Pressed == true {
				if len(dev.Config)-1 >= dev.Page && len(dev.Config[dev.Page])-1 >= int(k.Index) {
					HandleInput(dev, &dev.Config[dev.Page][k.Index], dev.Page)
				}
			}
		}
	}
}
