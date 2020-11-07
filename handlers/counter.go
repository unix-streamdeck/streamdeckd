package handlers

import (
	"github.com/unix-streamdeck/api"
	"image"
	"image/draw"
	"log"
	"strconv"
)

func (c *CounterIconHandler) Start(k api.Key, info api.StreamDeckInfo, callback func(image image.Image)) {
	if c.Callback == nil {
		c.Callback = callback
	}
	if c.Running {
		img := image.NewRGBA(image.Rect(0, 0, info.IconSize, info.IconSize))
		draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)
		Count := strconv.Itoa(c.Count)
		imgParsed, err := api.DrawText(img, Count, k.TextSize, k.TextAlignment)
		if err != nil {
			log.Println(err)
		} else {
			callback(imgParsed)
		}
	}
}

func (c *CounterIconHandler) IsRunning() bool {
	return c.Running
}

func (c *CounterIconHandler) SetRunning(running bool)  {
	c.Running = running
}

func (c CounterIconHandler) Stop()  {
	c.Running = false
}

type CounterKeyHandler struct{}

func (CounterKeyHandler) Key(key api.Key, info api.StreamDeckInfo) {
	if key.IconHandler != "Counter" {
		return
	}
	handler := key.IconHandlerStruct.(*CounterIconHandler)
	handler.Count += 1
	if handler.Callback != nil {
		handler.Start(key, info, handler.Callback)
	}
}
