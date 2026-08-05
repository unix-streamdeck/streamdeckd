package api

import "image"

type LcdSegmentBackgrounder interface {
	GetTouchPanelBackground() string
	GetTouchPanelBackgroundBuff() image.Image
	SetTouchPanelBackgroundBuff(img image.Image)
	GetTouchPanelBackgroundHandler() ForegroundHandler
	SetTouchPanelBackgroundHandler(handler ForegroundHandler)
	GetTouchPanelBackgroundHandlerFields() map[string]any
}

type Knob struct {
	Application                       map[string]*KnobConfig `json:"application,omitempty"`
	ActiveBuff                        image.Image            `json:"-"`
	ActiveApplication                 string                 `json:"-"`
	Layers                            []*Layer               `json:"layers"`
	TouchPanelBackground              string                 `json:"touch_panel_background"`
	TouchPanelBackgroundBuff          image.Image            `json:"-"`
	TouchPanelBackgroundHandler       ForegroundHandler      `json:"-"`
	TouchPanelBackgroundHandlerFields map[string]any         `json:"touch_panel_background_handler_fields"`
}

func (k *Knob) GetLayers() []*Layer {
	return k.Layers
}

func (k *Knob) GetTouchPanelBackgroundHandlerFields() map[string]any {
	return k.TouchPanelBackgroundHandlerFields
}

func (k *Knob) SetTouchPanelBackgroundHandler(handler ForegroundHandler) {
	k.TouchPanelBackgroundHandler = handler
}

func (k *Knob) GetTouchPanelBackgroundHandler() ForegroundHandler {
	return k.TouchPanelBackgroundHandler
}

func (k *Knob) GetTouchPanelBackground() string {
	return k.TouchPanelBackground
}

func (k *Knob) GetTouchPanelBackgroundBuff() image.Image {
	return k.TouchPanelBackgroundBuff
}

func (k *Knob) SetTouchPanelBackgroundBuff(img image.Image) {
	k.TouchPanelBackgroundBuff = img
}

type KnobConfig struct {
	Layers                            []*Layer          `json:"layers"`
	TouchPanelBackground              string            `json:"touch_panel_background"`
	TouchPanelBackgroundBuff          image.Image       `json:"-"`
	TouchPanelBackgroundHandler       ForegroundHandler `json:"-"`
	TouchPanelBackgroundHandlerFields map[string]any    `json:"touch_panel_background_handler_fields"`
}

func (kc *KnobConfig) GetLayers() []*Layer {
	return kc.Layers
}

func (kc *KnobConfig) GetTouchPanelBackgroundHandlerFields() map[string]any {
	return kc.TouchPanelBackgroundHandlerFields
}

func (kc *KnobConfig) SetTouchPanelBackgroundHandler(handler ForegroundHandler) {
	kc.TouchPanelBackgroundHandler = handler
}

func (kc *KnobConfig) GetTouchPanelBackgroundHandler() ForegroundHandler {
	return kc.TouchPanelBackgroundHandler
}

func (kc *KnobConfig) GetTouchPanelBackground() string {
	return kc.TouchPanelBackground
}

func (kc *KnobConfig) GetTouchPanelBackgroundBuff() image.Image {
	return kc.TouchPanelBackgroundBuff
}

func (kc *KnobConfig) SetTouchPanelBackgroundBuff(img image.Image) {
	kc.TouchPanelBackgroundBuff = img
}
