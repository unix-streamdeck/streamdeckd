package streamdeckd

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"os/exec"

	"github.com/unix-streamdeck/api/v3"
)

var applicationManager IApplicationManager = &ApplicationManager{}

var locked = false

//go:embed rounded-corners.png
var roundedCorners []byte

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

	return img, nil
}

func RunCommand(command string) {
	go func() {
		cmd := exec.Command("/bin/sh", "-c", command)

		if err := cmd.Start(); err != nil {
			log.Println("There was a problem running ", command, ":", err)
		} else {
			pid := cmd.Process.Pid
			err := cmd.Process.Release()
			if err != nil {
				log.Println(err)
			}
			log.Println(command, " has been started with pid", pid)
		}
	}()
}

func ExecuteKeybind(keybind string) error {
	keys, err := api.ParseXDoToolKeybindString(keybind)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to parse keybind: %s", err))
	}

	for _, key := range keys {
		if err := kb.KeyDown(key); err != nil {
			for i := len(keys) - 1; i >= 0; i-- {
				keyUpErr := kb.KeyUp(keys[i])
				log.Printf("[WARN] Failed to release key %d: %v", keys[i], keyUpErr)
			}
			return errors.New(fmt.Sprintf("failed to press key %d: %s", key, err))
		}
	}

	for i := len(keys) - 1; i >= 0; i-- {
		if err := kb.KeyUp(keys[i]); err != nil {
			log.Printf("[WARN] Failed to release key %d: %v", keys[i], err)
		}
	}

	return nil
}

func HandlePanic(cback func()) {
	if err := recover(); err != nil {
		log.Println("panic occurred:", err)
		cback()
	}
}

func getRoundedCornersImage(iconSize int) (image.Image, error) {
	img, err := png.Decode(bytes.NewBuffer(roundedCorners))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return api.ResizeImage(img, iconSize), nil
}

func MergeLayers(applicationLayers api.Layerable, globalLayers api.Layerable) []*api.Layer {
	var layers []*api.Layer
	globalsAdded := false
	for _, layer := range applicationLayers.GetLayers() {
		if layer == nil {
			continue
		}
		if layer.Global {
			if globalsAdded {
				log.Println("[Warn] Attempt to add global layers more than once ignored")
				continue
			}
			layers = append(layers, globalLayers.GetLayers()...)
			globalsAdded = true
			continue
		}
		layers = append(layers, layer)
	}
	if !globalsAdded {
		globals := globalLayers.GetLayers()
		merged := make([]*api.Layer, 0, len(globals)+len(layers))
		merged = append(merged, globals...)
		return append(merged, layers...)
	}
	return layers
}
