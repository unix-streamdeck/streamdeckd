package main

import (
	"encoding/json"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/unix-streamdeck/streamdeck-lib"
	"golang.org/x/image/font/inconsolata"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

var page = 0
var dev streamdeck.Device
var config Config

func main() {
	d, err := streamdeck.Devices()
	if err != nil {
		log.Fatal(err)
	}
	if len(d) == 0 {
		log.Fatal("No Stream Deck devices found.")
	}
	dev = d[0]
	err = dev.Open()
	if err != nil {
		log.Fatal(err)
	}
	config, err = readConfig()
	if err != nil {
		log.Fatal(err)
	}
	cleanupHook()
	setPage()
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
				handleInput(config.Pages[page][k.Index])
			}
		}
	}
}

func readConfig() (Config, error) {
	data, err := ioutil.ReadFile(os.Getenv("HOME") + "/.streamdeck-config.json")
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func setImage(img image.Image, i int, p int) {
	if p == page {
		dev.SetImage(uint8(i), img)
	}
}

func setPage() {
	currentPage := config.Pages[page]
	for i := range currentPage {
		currentKey := currentPage[i]
		if currentKey.buff == nil {
			if currentKey.Icon == "" {
				img := image.NewRGBA(image.Rect(0, 0, int(dev.Pixels), int(dev.Pixels)))
				draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{0, 0, 0, 255}), image.ZP, draw.Src)
				currentKey.buff = img
			} else {
				img, err := loadImage(currentKey.Icon)
				if err != nil {
					log.Fatal(err)
				}
				currentKey.buff = img
			}
			if currentKey.Text != "" {
				img := gg.NewContextForImage(currentKey.buff)
				img.SetRGB(0, 0, 0)
				img.Clear()
				img.SetRGB(1, 1, 1)
				img.SetFontFace(inconsolata.Regular8x16)
				img.DrawStringAnchored(currentKey.Text, 72/2, 72/2, 0.5, 0.5)
				img.Clip()
				currentKey.buff = img.Image()
			}
		}
		setImage(currentKey.buff, i, page)
	}
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return resize.Resize(72, 72, img, resize.Lanczos3), nil
}

func handleInput(key Key) {
	if key.Command != "" {
		runCommand(key.Command)
	}
	if key.Keybind != "" {
		runCommand("xdotool key " + key.Keybind)
	}
	if key.SwitchPage != nil {
		page = (*key.SwitchPage) - 1
		setPage()
	}
	if key.Brightness != nil {
		_ = dev.SetBrightness(uint8(*key.Brightness))
	}
	if key.Url != "" {
		runCommand("xdg-open " + key.Url)
	}
}

func runCommand(command string) {
	args := strings.Split(command, " ")
	c := exec.Command(args[0], args[1:]...)
	if err := c.Start(); err != nil {
		panic(err)
	}
	err := c.Wait()
	if err != nil {
		log.Printf("command failed: %s", err)
	}
}

func cleanupHook() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		<-sigs
		_ = dev.Reset()
		os.Exit(0)
	}()
}
