package api

import "image"

type LcdSegmentBackgrounder interface {
	GetTouchPanelBackground() string
	GetTouchPanelBackgroundBuff() image.Image
	SetTouchPanelBackgroundBuff(img image.Image)
	GetTouchPanelBackgroundHandler() TouchPanelBackgroundHandler
	SetTouchPanelBackgroundHandler(handler TouchPanelBackgroundHandler)
	GetTouchPanelBackgroundHandlerFields() map[string]any
}

type KnobV3 struct {
	Application                       map[string]*KnobConfigV3    `json:"application,omitempty"`
	ActiveBuff                        image.Image                 `json:"-"`
	ActiveApplication                 string                      `json:"-"`
	TouchPanelBackground              string                      `json:"touch_panel_background"`
	TouchPanelBackgroundBuff          image.Image                 `json:"-"`
	TouchPanelBackgroundHandler       TouchPanelBackgroundHandler `json:"-"`
	TouchPanelBackgroundHandlerFields map[string]any              `json:"touch_panel_background_handler_fields"`
}

func (k *KnobV3) GetTouchPanelBackgroundHandlerFields() map[string]any {
	return k.TouchPanelBackgroundHandlerFields
}

func (k *KnobV3) SetTouchPanelBackgroundHandler(handler TouchPanelBackgroundHandler) {
	k.TouchPanelBackgroundHandler = handler
}

func (k *KnobV3) GetTouchPanelBackgroundHandler() TouchPanelBackgroundHandler {
	return k.TouchPanelBackgroundHandler
}

func (k *KnobV3) GetTouchPanelBackground() string {
	return k.TouchPanelBackground
}

func (k *KnobV3) GetTouchPanelBackgroundBuff() image.Image {
	return k.TouchPanelBackgroundBuff
}

func (k *KnobV3) SetTouchPanelBackgroundBuff(img image.Image) {
	k.TouchPanelBackgroundBuff = img
}

type KnobActionV3 struct {
	SwitchPage       int               `json:"switch_page,omitempty"`
	Keybind          string            `json:"keybind,omitempty"`
	Command          string            `json:"command,omitempty"`
	Brightness       int               `json:"brightness,omitempty"`
	Url              string            `json:"url,omitempty"`
	ObsCommand       string            `json:"obs_command,omitempty"`
	ObsCommandParams map[string]string `json:"obs_command_params,omitempty"`
}

type KnobConfigV3 struct {
	Icon                              string                      `json:"icon,omitempty"`
	Text                              string                      `json:"text,omitempty"`
	TextSize                          int                         `json:"text_size,omitempty"`
	TextAlignment                     string                      `json:"text_alignment,omitempty"`
	FontFace                          string                      `json:"font_face,omitempty"`
	TextColour                        string                      `json:"text_colour,omitempty"`
	LcdHandler                        string                      `json:"lcd_handler,omitempty"`
	KnobOrTouchHandler                string                      `json:"knob_or_touch_handler,omitempty"`
	LcdHandlerStruct                  LcdHandler                  `json:"-"`
	KnobOrTouchHandlerStruct          KnobOrTouchHandler          `json:"-"`
	LcdHandlerFields                  map[string]any              `json:"lcd_handler_fields,omitempty"`
	KnobOrTouchHandlerFields          map[string]any              `json:"knob_or_touch_handler_fields,omitempty"`
	SharedHandlerFields               map[string]any              `json:"shared_handler_fields,omitempty"`
	KnobPressAction                   KnobActionV3                `json:"knob_press_action,omitempty"`
	KnobTurnUpAction                  KnobActionV3                `json:"knob_turn_up_action,omitempty"`
	KnobTurnDownAction                KnobActionV3                `json:"knob_turn_down_action,omitempty"`
	SharedState                       map[string]any              `json:"-"`
	TouchPanelBackground              string                      `json:"touch_panel_background"`
	TouchPanelBackgroundBuff          image.Image                 `json:"-"`
	TouchPanelBackgroundHandler       TouchPanelBackgroundHandler `json:"-"`
	TouchPanelBackgroundHandlerFields map[string]any              `json:"touch_panel_background_handler_fields"`
}

func (kc *KnobConfigV3) GetTouchPanelBackgroundHandlerFields() map[string]any {
	return kc.TouchPanelBackgroundHandlerFields
}

func (kc *KnobConfigV3) SetTouchPanelBackgroundHandler(handler TouchPanelBackgroundHandler) {
	kc.TouchPanelBackgroundHandler = handler
}

func (kc *KnobConfigV3) GetTouchPanelBackgroundHandler() TouchPanelBackgroundHandler {
	return kc.TouchPanelBackgroundHandler
}

func (kc *KnobConfigV3) GetTouchPanelBackground() string {
	return kc.TouchPanelBackground
}

func (kc *KnobConfigV3) GetTouchPanelBackgroundBuff() image.Image {
	return kc.TouchPanelBackgroundBuff
}

func (kc *KnobConfigV3) SetTouchPanelBackgroundBuff(img image.Image) {
	kc.TouchPanelBackgroundBuff = img
}
