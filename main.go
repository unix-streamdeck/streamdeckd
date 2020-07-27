package main

import (
	"encoding/json"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/unix-streamdeck/streamdeck-lib"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
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
				fmt.Print(config.Pages[page][k.Index])
				handleInput(config.Pages[page][k.Index])
			}
		}
	}
}

func readConfig() (Config, error) {
	data, err := ioutil.ReadFile("/home/jonsey/.streamdeck-config.json")
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
		if currentKey.Icon == "" {
			currentKey.Icon = "blank.png"
		}
		img, err := loadImage(currentKey.Icon)
		if err != nil {
			log.Fatal(err)
		}
		setImage(img, i, page)
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
		fmt.Print("Switching to: " + string(*key.SwitchPage))
		page = (*key.SwitchPage) -1
		setPage()
	}
	if key.Brightness != nil {
		_ = dev.SetBrightness(uint8(*key.Brightness))
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
