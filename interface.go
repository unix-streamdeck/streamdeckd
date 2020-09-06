package main

import (
	"context"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/streamdeckd/handlers"
	"golang.org/x/sync/semaphore"
	"image"
	"image/draw"
	"log"
	"os"
	"strings"
)

var p int
var sem = semaphore.NewWeighted(int64(1))

func LoadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return api.ResizeImage(img, sDInfo.IconSize), nil
}

func SetImage(img image.Image, i int, page int) {
	ctx := context.Background()
	err := sem.Acquire(ctx, 1)
	if err != nil {
		log.Println(err)
		return
	}
	defer sem.Release(1)
	if p == page && isOpen {
		err := dev.SetImage(uint8(i), img)
		if err != nil {
			if strings.Contains(err.Error(), "hidapi") {
				disconnect()
			} else {
				log.Println(err)
			}
		}
	}
}

func SetKeyImage(currentKey *api.Key, i int) {
	if currentKey.Buff == nil {
		if currentKey.Icon == "" {
			img := image.NewRGBA(image.Rect(0, 0, int(dev.Pixels), int(dev.Pixels)))
			draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)
			currentKey.Buff = img
		} else {
			img, err := LoadImage(currentKey.Icon)
			if err != nil {
				log.Println(err)
				return
			}
			currentKey.Buff = img
		}
		if currentKey.Text != "" {
			img, err := api.DrawText(currentKey.Buff, currentKey.Text)
			if err != nil {
				log.Println(err)
			} else {
				currentKey.Buff = img
			}
		}
	}
	if currentKey.Buff != nil {
		SetImage(currentKey.Buff, i, p)
	}
}

func SetPage(config *api.Config, page int) {
	p = page
	currentPage := config.Pages[page]
	for i := 0; i < len(currentPage); i++ {
		currentKey := &currentPage[i]
		go SetKey(currentKey, i, page)
	}
	EmitPage(p)
}

func SetKey(currentKey *api.Key, i int, page int) {
	if currentKey.Buff == nil {
		if currentKey.IconHandler == "" {
			SetKeyImage(currentKey, i)

		} else if currentKey.IconHandlerStruct == nil {
			var handler api.IconHandler
			if currentKey.IconHandler == "Gif" {
				handler = &handlers.GifIconHandler{Running:true}
			} else if currentKey.IconHandler == "Counter" {
				handler = &handlers.CounterIconHandler{Count:0, Running: true}
			} else if currentKey.IconHandler == "Time" {
				handler = &handlers.TimeIconHandler{Running:true}
			}
			if handler == nil {
				return
			}
			handler.Start(*currentKey, sDInfo, func(image image.Image) {
				if image.Bounds().Max.X != 72 || image.Bounds().Max.Y != 72 {
					image = api.ResizeImage(image, sDInfo.IconSize)
				}
				SetImage(image, i, page)
				currentKey.Buff = image
			})
			currentKey.IconHandlerStruct = handler
		}
	} else {
		SetImage(currentKey.Buff, i, p)
	}
	if currentKey.IconHandlerStruct != nil && !currentKey.IconHandlerStruct.IsRunning() {
		currentKey.IconHandlerStruct.SetRunning(true)
		currentKey.IconHandlerStruct.Start(*currentKey, sDInfo, func(image image.Image) {
			if image.Bounds().Max.X != 72 || image.Bounds().Max.Y != 72 {
				image = api.ResizeImage(image, sDInfo.IconSize)
			}
			SetImage(image, i, page)
			currentKey.Buff = image
		})
	}
}

func HandleInput(key *api.Key, page int) {
	if key.Command != "" {
		runCommand(key.Command)
	}
	if key.Keybind != "" {
		runCommand("xdotool key " + key.Keybind)
	}
	if key.SwitchPage != 0 {
		page = key.SwitchPage - 1
		SetPage(config, page)
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
		if key.KeyHandlerStruct == nil {
			var handler api.KeyHandler
			if key.KeyHandler == "Counter" {
				handler = handlers.CounterKeyHandler{}
			}
			if handler == nil {
				return
			}
			key.KeyHandlerStruct = handler
		}
		key.KeyHandlerStruct.Key(*key, sDInfo)
	}
}

func Listen() {
	kch, err := dev.ReadKeys()
	if err != nil {
		log.Println(err)
	}
	for isOpen {
		select {
		case k, ok := <-kch:
			if !ok {
				disconnect()
				return
			}
			if k.Pressed == true {
				if len(config.Pages)-1 >= p && len(config.Pages[p])-1 >= int(k.Index) {
					HandleInput(&config.Pages[p][k.Index], p)
				}
			}
		}
	}
}
