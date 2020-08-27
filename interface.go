package main

import (
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/driver"
	"github.com/unix-streamdeck/streamdeckd/handlers"
	"golang.org/x/image/font/inconsolata"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"strings"
)

var p int

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

func SetImage(img image.Image, i int, page int, _ streamdeck.Device) {
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

func SetKey(currentKey *api.Key, i int) {
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

func SetPage(config *api.Config, page int) {
	p = page
	currentPage := config.Pages[page]
	for i := 0; i < len(currentPage); i++ {
		currentKey := &currentPage[i]
		if currentKey.Buff == nil {
			if currentKey.IconHandler == "" {
				SetKey(currentKey, i)

			} else if currentKey.IconHandlerStruct == nil {
				var handler api.IconHandler
				if currentKey.IconHandler == "Gif" {
					handler = &handlers.GifIconHandler{Running:true, OnSetImage: SetImage}
				} else if currentKey.IconHandler == "Counter" {
					handler = &handlers.CounterIconHandler{Count:0, Running: true, OnSetImage: SetImage}
				} else if currentKey.IconHandler == "Time" {
					handler = &handlers.TimeIconHandler{Running:true, OnSetImage: SetImage}
				}
				if handler == nil {
					continue
				}
				handler.Icon(page, i, currentKey, dev)
				currentKey.IconHandlerStruct = handler
			}
		} else {
			SetImage(currentKey.Buff, i, p, dev)
		}
	}
	EmitPage(p)
}

func HandleInput(key *api.Key, page int, index int) {
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
		_ = dev.SetBrightness(uint8(key.Brightness))
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
		key.KeyHandlerStruct.Key(page, index, key, streamdeck.Device{})
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
					HandleInput(&config.Pages[p][k.Index], p, int(k.Index))
				}
			}
		}
	}
}
