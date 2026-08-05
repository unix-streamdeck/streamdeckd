package api

import "image"

type KeyBackgrounder interface {
	GetKeyBackground() string
	GetKeyBackgroundBuff() image.Image
	SetKeyBackgroundBuff(img image.Image)
	GetKeyBackgroundHandler() BackgroundHandler
	SetKeyBackgroundHandler(handler BackgroundHandler)
	GetKeyBackgroundHandlerFields() map[string]any
}

type Key struct {
	Application                map[string]*KeyConfig `json:"application,omitempty"`
	ActiveBuff                 image.Image           `json:"-"`
	ActiveApplication          string                `json:"-"`
	KeyBackground              string                `json:"background"`
	KeyBackgroundBuff          image.Image           `json:"-"`
	KeyBackgroundHandler       BackgroundHandler     `json:"-"`
	KeyBackgroundHandlerFields map[string]any        `json:"key_background_handler_fields"`
	Layers                     []*Layer              `json:"layers"`
}

func (k *Key) GetLayers() []*Layer {
	return k.Layers
}

func (k *Key) GetKeyBackgroundHandlerFields() map[string]any {
	return k.KeyBackgroundHandlerFields
}

func (k *Key) GetKeyBackground() string {
	return k.KeyBackground
}

func (k *Key) GetKeyBackgroundBuff() image.Image {
	return k.KeyBackgroundBuff
}

func (k *Key) SetKeyBackgroundBuff(img image.Image) {
	k.KeyBackgroundBuff = img
}

func (k *Key) SetKeyBackgroundHandler(handler BackgroundHandler) {
	k.KeyBackgroundHandler = handler
}

func (k *Key) GetKeyBackgroundHandler() BackgroundHandler {
	return k.KeyBackgroundHandler
}

type KeyConfig struct {
	Layers                     []*Layer          `json:"layers"`
	KeyBackground              string            `json:"background"`
	KeyBackgroundBuff          image.Image       `json:"-"`
	KeyBackgroundHandler       BackgroundHandler `json:"-"`
	KeyBackgroundHandlerFields map[string]any    `json:"key_background_handler_fields"`
}

func (kc *KeyConfig) GetLayers() []*Layer {
	return kc.Layers
}

func (kc *KeyConfig) GetKeyBackgroundHandlerFields() map[string]any {
	return kc.KeyBackgroundHandlerFields
}

func (kc *KeyConfig) GetKeyBackground() string {
	return kc.KeyBackground
}

func (kc *KeyConfig) GetKeyBackgroundBuff() image.Image {
	return kc.KeyBackgroundBuff
}

func (kc *KeyConfig) SetKeyBackgroundBuff(img image.Image) {
	kc.KeyBackgroundBuff = img
}

func (kc *KeyConfig) SetKeyBackgroundHandler(handler BackgroundHandler) {
	kc.KeyBackgroundHandler = handler
}

func (kc *KeyConfig) GetKeyBackgroundHandler() BackgroundHandler {
	return kc.KeyBackgroundHandler
}
