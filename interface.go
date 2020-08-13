package main

import (
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/unix-streamdeck/streamdeck"
	"golang.org/x/image/font/inconsolata"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
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
	return resize.Resize(72, 72, img, resize.Lanczos3)
}

func SetImage(img image.Image, i int, page int, dev streamdeck.Device) {
	if p == page {
		dev.SetImage(uint8(i), img)
	}
}

func SetKey(currentKey *Key, i int) {
	if currentKey.Buff == nil {
		if currentKey.Icon == "" {
			img := image.NewRGBA(image.Rect(0, 0, int(dev.Pixels), int(dev.Pixels)))
			draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{0, 0, 0, 255}), image.ZP, draw.Src)
			currentKey.Buff = img
		} else {
			img, err := LoadImage(currentKey.Icon)
			if err != nil {
				log.Panic(err)
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
	SetImage(currentKey.Buff, i, p, dev)
}

func SetPage(config *Config, page int, dev streamdeck.Device) {
	p = page
	currentPage := config.Pages[page]
	for i := 0; i < 15; i++ {
		currentKey := &currentPage[i]
		if currentKey.Buff == nil {
			if currentKey.IconHandler == "" {
				SetKey(currentKey, i)

			} else if currentKey.IconHandlerStruct == nil {
				var handler IconHandler
				if currentKey.IconHandler == "Gif" {
					handler = &GifIconHandler{true}
				} else if currentKey.IconHandler == "Counter" {
					handler = &CounterIconHandler{0, true}
				} else if currentKey.IconHandler == "Time" {
					handler = &TimeIconHandler{true}
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

func HandleInput(key *Key, page int, index int, dev streamdeck.Device) {
	if key.Command != "" {
		runCommand(key.Command)
	}
	if key.Keybind != "" {
		runCommand("xdotool key " + key.Keybind)
	}
	if key.SwitchPage != nil {
		page = (*key.SwitchPage) - 1
		SetPage(config, page, dev)
	}
	if key.Brightness != nil {
		_ = dev.SetBrightness(uint8(*key.Brightness))
	}
	if key.Url != "" {
		runCommand("xdg-open " + key.Url)
	}
	if key.KeyHandler != "" {
		if key.KeyHandlerStruct == nil {
			var handler KeyHandler
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
		log.Fatal(err)
	}
	for {
		select {
		case k, ok := <-kch:
			if !ok {
				err = dev.Open()
				if err != nil {
					log.Fatal(err)
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
