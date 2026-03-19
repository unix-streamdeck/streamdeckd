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
}

type IconHandler interface {
	VisualHandler
	Start(key KeyConfigV3, info StreamDeckInfoV1, callback func(image image.Image))
}

type KeyHandler interface {
	InputHandler
	Key(key KeyConfigV3, info StreamDeckInfoV1)
}

type LcdHandler interface {
	VisualHandler
	Start(key KnobConfigV3, info StreamDeckInfoV1, callback func(image image.Image))
}

type KnobOrTouchHandler interface {
	InputHandler
	Input(key KnobConfigV3, info StreamDeckInfoV1, event InputEvent)
}

type BackgroundHandler interface {
	VisualHandler
	Start(fields map[string]any, info StreamDeckInfoV1, callback func(images []image.Image))
	StartIndividual(fields map[string]any, info StreamDeckInfoV1, callback func(img image.Image))
}

type KeyGridBackgroundHandler interface {
	BackgroundHandler
}

type TouchPanelBackgroundHandler interface {
	BackgroundHandler
}

type InputEventType uint8

const (
	KNOB_CCW InputEventType = iota
	KNOB_CW
	KNOB_PRESS
	SCREEN_SHORT_TAP
	SCREEN_LONG_TAP
)

type InputEvent struct {
	EventType     InputEventType
	RotateNotches uint8
}
