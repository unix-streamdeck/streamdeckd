package examples

import (
	"context"
	"image"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/unix-streamdeck/api/v2"
	"golang.org/x/sync/semaphore"
)

type ToggleHandler struct {
	Status       bool
	Running      bool
	Lock         *semaphore.Weighted
	Callback     func(image image.Image)
	Quit         chan bool
	UpIconBuff   image.Image
	DownIconBuff image.Image
	FirstLoop    bool
}

func (c *ToggleHandler) Start(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	if c.Lock == nil {
		c.Lock = semaphore.NewWeighted(1)
	}
	if c.Quit == nil {
		c.Quit = make(chan bool)
	}
	if c.UpIconBuff == nil {
		c.UpIconBuff = c.GetImage("up_icon", fields, handlerType, info)
	}
	if c.DownIconBuff == nil {
		c.DownIconBuff = c.GetImage("down_icon", fields, handlerType, info)
	}
	c.FirstLoop = true
	go c.loop(fields, callback)
}

func (c *ToggleHandler) GetImage(index string, fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1) image.Image {
	path, ok := fields[index]
	w, h := info.GetDimensions(handlerType)
	if !ok {
		log.Println("image missing: " + index)
		return image.NewNRGBA(image.Rect(0, 0, w, h))
	}
	f, err := os.Open(path.(string))
	defer f.Close()
	if err != nil {
		log.Println(err)
		return image.NewNRGBA(image.Rect(0, 0, w, h))
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Println(err)
		return image.NewNRGBA(image.Rect(0, 0, w, h))
	}
	return api.ResizeImageWH(img, w, h)
}

func (c *ToggleHandler) loop(fields map[string]any, callback func(image image.Image)) {
	ctx := context.Background()
	err := c.Lock.Acquire(ctx, 1)
	if err != nil {
		return
	}
	defer c.Lock.Release(1)
	for {
		select {
		case <-c.Quit:
			return
		default:
			command, ok := fields["check_command"]

			if !ok {
				break
			}
			cmd := exec.Command("/bin/sh", "-c", command.(string))
			status := true
			if err := cmd.Start(); err != nil {
				status = false
			}
			err := cmd.Wait()
			if err != nil {
				status = false
			}
			sharedStatus, ok := fields["status"].(bool)
			if !ok {
				sharedStatus = false
			}
			if status == sharedStatus && !c.FirstLoop {
				time.Sleep(250 * time.Millisecond)
				continue
			}
			sharedStatus = status
			fields["status"] = sharedStatus
			c.FirstLoop = false
			img := c.UpIconBuff
			if sharedStatus == false {
				img = c.DownIconBuff
			}
			callback(img)
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func (c *ToggleHandler) IsRunning() bool {
	return c.Running
}

func (c *ToggleHandler) SetRunning(running bool) {
	c.Running = running
}

func (c *ToggleHandler) Stop() {
	c.Running = false
	c.Quit <- true
}

func (t *ToggleHandler) Input(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1, event api.InputEvent) {
	sharedStatus := fields["status"].(bool)
	index := "down_command"
	if !sharedStatus {
		index = "up_command"
	}
	command, ok := fields[index]
	commandString := command.(string)
	if !ok {
		return
	}
	go func() {
		cmd := exec.Command("/bin/sh", commandString)

		if err := cmd.Start(); err != nil {
			log.Println("There was a problem running ", commandString, ":", err)
		} else {
			pid := cmd.Process.Pid
			err := cmd.Process.Release()
			if err != nil {
				log.Println(err)
			}
			log.Println(commandString, " has been started with pid", pid)
		}
	}()
}

func RegisterToggle() api.Module {

	return api.Module{
		Name: "Toggle",
		NewForeground: func() api.ForegroundHandler {
			return &ToggleHandler{Running: true, Lock: semaphore.NewWeighted(1), FirstLoop: true}
		},
		NewInput: func() api.InputHandler { return &ToggleHandler{} },
		ForegroundFields: []api.Field{
			{Title: "Up Icon", Name: "up_icon", Type: api.File, FileTypes: []string{".png", ".jpg", ".jpeg"}},
			{Title: "Down Icon", Name: "down_icon", Type: api.File, FileTypes: []string{".png", ".jpg", ".jpeg"}},
			{Title: "Check Command", Name: "check_command", Type: api.Text},
		},
		InputFields: []api.Field{
			{Title: "Up Command", Name: "up_command", Type: api.Text},
			{Title: "Down Command", Name: "down_command", Type: api.Text},
		},
	}
}
