package examples

import (
	"image"
	"log"
	"time"

	"github.com/unix-streamdeck/api/v2"
)

type TimeHandler struct {
	Running bool
	Quit    chan bool
}

func (t *TimeHandler) Start(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	t.Running = true
	if t.Quit == nil {
		t.Quit = make(chan bool)
	}
	go t.timeLoop(fields, handlerType, info, callback)
}

func (t *TimeHandler) IsRunning() bool {
	return t.Running
}

func (t *TimeHandler) SetRunning(running bool) {
	t.Running = running
}

func (t *TimeHandler) Stop() {
	t.Running = false
	t.Quit <- true
}

func (t *TimeHandler) timeLoop(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	for {
		select {
		case <-t.Quit:
			return
		default:
			w, h := info.GetDimensions(handlerType)
			img := image.NewRGBA(image.Rect(0, 0, w, h))
			t := time.Now()
			tString := t.Format("15:04:05")
			imgParsed, err := api.DrawText(img, tString, api.DrawTextOptions{})
			if err != nil {
				log.Println(err)
			} else {
				callback(imgParsed)
			}
			time.Sleep(time.Second)
		}
	}
}

func RegisterTime() api.Module {
	return api.Module{NewForeground: func() api.ForegroundHandler {
		return &TimeHandler{Running: true}
	}, Name: "Time"}
}
