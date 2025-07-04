package key

import (
    "github.com/unix-streamdeck/api"
    "github.com/unix-streamdeck/streamdeckd/streamdeckd"
    "image"
    "image/draw"
    "log"
    "strconv"
)

type CounterIconHandler struct {
	Count    int
	Running  bool
	Callback func(image image.Image)
}

func (c *CounterIconHandler) Start(k api.KeyConfigV3, info api.StreamDeckInfoV1, callback func(image image.Image)) {
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

func (CounterKeyHandler) Key(key api.KeyConfigV3, info api.StreamDeckInfoV1) {
	if key.IconHandler != "Counter" {
		return
	}
	handler := key.IconHandlerStruct.(*CounterIconHandler)
	handler.Count += 1
	if handler.Callback != nil {
		handler.Start(key, info, handler.Callback)
	}
}

func RegisterCounter() streamdeckd.Module {
	return streamdeckd.Module{NewIcon: func() api.IconHandler {
		return &CounterIconHandler{Running: true, Count: 0}
	}, NewKey: func() api.KeyHandler {
		return &CounterKeyHandler{}
	}, Name: "Counter"}
}