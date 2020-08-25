package main

import (
	"context"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/driver"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/sync/semaphore"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
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

	return ResizeImage(img), nil
}

func ResizeImage(img image.Image) image.Image {
	return resize.Resize(dev.Pixels, dev.Pixels, img, resize.Lanczos3)
}

func SetImage(img image.Image, i int, page int, dev streamdeck.Device) {
	ctx := context.Background()
	err := sem.Acquire(ctx, 1)
	if err != nil {
		log.Println(err)
		return
	}
	defer sem.Release(1)
	if p == page {
		dev.SetImage(uint8(i), img)
	}
}

func SetKeyImage(currentKey *api.Key, i int) {
	if currentKey.Buff == nil {
		if currentKey.Icon == "" {
			img := image.NewRGBA(image.Rect(0, 0, int(dev.Pixels), int(dev.Pixels)))
			draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{0, 0, 0, 255}), image.ZP, draw.Src)
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
			img := gg.NewContextForImage(currentKey.Buff)
			img.SetRGB(1, 1, 1)
			img.SetFontFace(inconsolata.Regular8x16)
			img.DrawStringAnchored(currentKey.Text, 72/2, 72/2, 0.5, 0.5)
			img.Clip()
			currentKey.Buff = img.Image()
		}
	}
	if currentKey.Buff != nil {
		SetImage(currentKey.Buff, i, p, dev)
	}
}

func SetPage(config *api.Config, page int, dev streamdeck.Device) {
	p = page
	currentPage := config.Pages[page]
	for i := 0; i < len(currentPage); i++ {
		currentKey := &currentPage[i]
		go SetKey(currentKey, i, page, dev)
	}
	EmitPage(p)
}

func SetKey(currentKey *api.Key, i int, page int, dev streamdeck.Device) {
	if currentKey.Buff == nil {
		if currentKey.IconHandler == "" {
			SetKeyImage(currentKey, i)

		} else if currentKey.IconHandlerStruct == nil {
			var handler api.IconHandler
			if currentKey.IconHandler == "Gif" {
				handler = &GifIconHandler{true}
			} else if currentKey.IconHandler == "Counter" {
				handler = &CounterIconHandler{0, true}
			} else if currentKey.IconHandler == "Time" {
				handler = &TimeIconHandler{true}
			}
			if handler == nil {
				return
			}
			handler.Icon(page, i, currentKey, dev)
			currentKey.IconHandlerStruct = handler
		}
	} else {
		SetImage(currentKey.Buff, i, p, dev)
	}
}

func HandleInput(key *api.Key, page int, index int, dev streamdeck.Device) {
	if key.Command != "" {
		runCommand(key.Command)
	}
	if key.Keybind != "" {
		runCommand("xdotool key " + key.Keybind)
	}
	if key.SwitchPage != 0 {
		page = key.SwitchPage - 1
		SetPage(config, page, dev)
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
				handler = CounterKeyHandler{}
			}
			if handler == nil {
				return
			}
			key.KeyHandlerStruct = handler
		}
		key.KeyHandlerStruct.Key(page, index, key, dev)
	}
}

func Listen() {
	kch, err := dev.ReadKeys()
	if err != nil {
		log.Println(err)
	}
	for {
		select {
		case k, ok := <-kch:
			if !ok {
				err = dev.Open()
				if err != nil {
					log.Println(err)
				}
				continue
			}
			if k.Pressed == true {
				if len(config.Pages)-1 >= p && len(config.Pages[p])-1 >= int(k.Index) {
					HandleInput(&config.Pages[p][k.Index], p, int(k.Index), dev)
				}
			}
		}
	}
}
