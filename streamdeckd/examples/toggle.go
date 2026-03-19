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

type ToggleIconHandler struct {
	Status       bool
	Running      bool
	Lock         *semaphore.Weighted
	Callback     func(image image.Image)
	Quit         chan bool
	UpIconBuff   image.Image
	DownIconBuff image.Image
	FirstLoop    bool
}

func (c *ToggleIconHandler) Start(k api.KeyConfigV3, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	if c.Lock == nil {
		c.Lock = semaphore.NewWeighted(1)
	}
	if c.Quit == nil {
		c.Quit = make(chan bool)
	}
	if c.UpIconBuff == nil {
		c.UpIconBuff = c.GetImage("up_icon", k, info)
	}
	if c.DownIconBuff == nil {
		c.DownIconBuff = c.GetImage("down_icon", k, info)
	}
	c.FirstLoop = true
	go c.loop(k, callback)
}

func (c *ToggleIconHandler) GetImage(index string, k api.KeyConfigV3, info api.StreamDeckInfoV1) image.Image {
	path, ok := k.IconHandlerFields[index]
	if !ok {
		log.Println("image missing: " + index)
		return image.NewNRGBA(image.Rect(0, 0, info.IconSize, info.IconSize))
	}
	f, err := os.Open(path.(string))
	defer f.Close()
	if err != nil {
		log.Println(err)
		return image.NewNRGBA(image.Rect(0, 0, info.IconSize, info.IconSize))
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Println(err)
		return image.NewNRGBA(image.Rect(0, 0, info.IconSize, info.IconSize))
	}
	return api.ResizeImage(img, info.IconSize)
}

func (c *ToggleIconHandler) loop(k api.KeyConfigV3, callback func(image image.Image)) {
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
			command, ok := k.IconHandlerFields["check_command"]

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
			sharedStatus, ok := k.SharedState["status"].(bool)
			if !ok {
				sharedStatus = false
			}
			if status == sharedStatus && !c.FirstLoop {
				time.Sleep(250 * time.Millisecond)
				continue
			}
			sharedStatus = status
			k.SharedState["status"] = sharedStatus
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

func (c *ToggleIconHandler) IsRunning() bool {
	return c.Running
}

func (c *ToggleIconHandler) SetRunning(running bool) {
	c.Running = running
}

func (c *ToggleIconHandler) Stop() {
	c.Running = false
	c.Quit <- true
}

type ToggleKeyHandler struct{}

func (ToggleKeyHandler) Key(key api.KeyConfigV3, info api.StreamDeckInfoV1) {
	sharedStatus := key.SharedState["status"].(bool)
	index := "down_command"
	if !sharedStatus {
		index = "up_command"
	}
	command, ok := key.KeyHandlerFields[index]
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
		NewIcon: func() api.IconHandler {
			return &ToggleIconHandler{Running: true, Lock: semaphore.NewWeighted(1), FirstLoop: true}
		},
		NewKey: func() api.KeyHandler { return &ToggleKeyHandler{} },
		IconFields: []api.Field{
			{Title: "Up Icon", Name: "up_icon", Type: api.File, FileTypes: []string{".png", ".jpg", ".jpeg"}},
			{Title: "Down Icon", Name: "down_icon", Type: api.File, FileTypes: []string{".png", ".jpg", ".jpeg"}},
			{Title: "Check Command", Name: "check_command", Type: api.Text},
		},
		KeyFields: []api.Field{
			{Title: "Up Command", Name: "up_command", Type: api.Text},
			{Title: "Down Command", Name: "down_command", Type: api.Text},
		},
	}
}
