package api

import "image"

type Handler interface {
}

type VisualHandler interface {
	Handler
	IsRunning() bool
	SetRunning(running bool)
	Stop()
}

type InputHandler interface {
	Handler
	Input(fields map[string]any, handlerType HandlerType, info StreamDeckInfoV1, event InputEvent)
}

type ForegroundHandler interface {
	VisualHandler
	Start(fields map[string]any, handlerType HandlerType, info StreamDeckInfoV1, callback func(image image.Image))
}

type CombinedHandler interface {
	VisualHandler
	InputHandler
}

type BackgroundHandler interface {
	ForegroundHandler
	StartGrid(fields map[string]any, handlerType HandlerType, info StreamDeckInfoV1, callback func(images []image.Image))
}

type HandlerType uint8

const (
	LCD HandlerType = iota
	KEY
)

type InputEventType uint8

const (
	KNOB_CCW InputEventType = iota
	KNOB_CW
	KNOB_PRESS
	SCREEN_SHORT_TAP
	SCREEN_LONG_TAP
	SCREEN_SWIPE
	KEY_PRESS
	KEY_RELEASE
)

type InputEvent struct {
	EventType     InputEventType
	RotateNotches uint8
}
