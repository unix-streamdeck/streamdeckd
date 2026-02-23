package api

import "image"

type Handler interface {
}

type IconHandler interface {
	Handler
	Start(key KeyConfigV3, info StreamDeckInfoV1, callback func(image image.Image))
	IsRunning() bool
	SetRunning(running bool)
	Stop()
}

type KeyHandler interface {
	Handler
	Key(key KeyConfigV3, info StreamDeckInfoV1)
}

type LcdHandler interface {
	Handler
	Start(key KnobConfigV3, info StreamDeckInfoV1, callback func(image image.Image))
	IsRunning() bool
	SetRunning(running bool)
	Stop()
}

type KnobOrTouchHandler interface {
	Handler
	Input(key KnobConfigV3, info StreamDeckInfoV1, event InputEvent)
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
