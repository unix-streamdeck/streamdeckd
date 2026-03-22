package examples

import (
	"image"
	"log"
	"strconv"

	"github.com/unix-streamdeck/api/v2"
)

type CounterHandler struct {
	Count    int
	Running  bool
	Callback func(image image.Image)
	Update   chan int
}

func (c *CounterHandler) Start(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	if c.Callback == nil {
		c.Callback = callback
	}
	if c.Update == nil {
		log.Println("Test")
		c.Update = make(chan int)
	}
	if c.Running {
		for {
			select {
			case delta := <-c.Update:
				c.Count = c.Count + delta
				width, height := info.GetDimensions(handlerType)
				img := image.NewRGBA(image.Rect(0, 0, width, height))
				Count := strconv.Itoa(c.Count)
				imgParsed, err := api.DrawText(img, Count, api.DrawTextOptions{})
				if err != nil {
					log.Println(err)
				} else {
					callback(imgParsed)
				}
			default:
				continue
			}
		}

	}
}

func (c *CounterHandler) IsRunning() bool {
	return c.Running
}

func (c *CounterHandler) SetRunning(running bool) {
	c.Running = running
}

func (c *CounterHandler) Stop() {
	c.Running = false
}

func (c *CounterHandler) Input(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1, event api.InputEvent) {
	var delta int
	if event.EventType == api.KEY_PRESS || event.EventType == api.KNOB_PRESS || event.EventType == api.SCREEN_SHORT_TAP {
		delta = 1
	}

	if event.EventType == api.KNOB_CW {
		delta = int(event.RotateNotches)
	}

	if event.EventType == api.KNOB_CCW {
		delta = int(event.RotateNotches) * -1
	}

	if c.Update != nil {
		c.Update <- delta
	}
}

func RegisterCounter() api.Module {
	return api.Module{
		NewForeground: func() api.ForegroundHandler {
			return &CounterHandler{Running: true, Count: 0}
		},
		NewInput: func() api.InputHandler {
			return &CounterHandler{Running: true, Count: 0}
		}, Name: "Counter"}
}
