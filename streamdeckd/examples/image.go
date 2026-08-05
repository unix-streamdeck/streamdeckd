package examples

import (
	"image"
	"log"
	"os"
	"os/exec"

	"github.com/unix-streamdeck/api/v3"
)

type ImageHandler struct {
	Running bool
}

func (i *ImageHandler) IsRunning() bool {
	return i.Running
}

func (i *ImageHandler) SetRunning(running bool) {
	i.Running = running
}

func (i *ImageHandler) Stop() {
	i.Running = false
}

func (i *ImageHandler) Start(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfo, callback func(image image.Image)) {
	i.Running = true

	icon, ok := fieldString(fields, "path")
	if !ok {
		log.Println("[WARN] Image handler: no icon specified")
		return
	}
	img, err := LoadImage(icon)
	if err != nil {
		log.Println("[ERROR] Image handler: error loading image", err)
		return
	}
	callback(img)
}

func RegisterImage() api.Module {
	return api.Module{
		Name: "Image",
		NewCombined: func() api.Handler {
			return &ImageHandler{}
		},
		CombinedFields: []api.Field{
			{Title: "Path", Name: "path", Type: api.File, FileTypes: []string{".jpg", ".jpeg", ".png", ".bmp"}},
		},
	}
}

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
