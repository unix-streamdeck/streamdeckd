package streamdeckd

import (
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"

	"github.com/unix-streamdeck/api/v2"
)

var applicationManager IApplicationManager = &ApplicationManager{}

var locked = false

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

func mergeSharedConfig(sharedConfig map[string]any, individualConfig map[string]any) map[string]any {
	merged := make(map[string]any)

	for k, v := range sharedConfig {
		merged[k] = v
	}

	for k, v := range individualConfig {
		merged[k] = v
	}

	return merged
}
