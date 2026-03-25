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

type KnobV3 struct {
	Application                       map[string]*KnobConfigV3 `json:"application,omitempty"`
	ActiveBuff                        image.Image              `json:"-"`
	ActiveApplication                 string                   `json:"-"`
	TouchPanelBackground              string                   `json:"touch_panel_background"`
	TouchPanelBackgroundBuff          image.Image              `json:"-"`
	TouchPanelBackgroundHandler       ForegroundHandler        `json:"-"`
	TouchPanelBackgroundHandlerFields map[string]any           `json:"touch_panel_background_handler_fields"`
}

func (k *KnobV3) GetTouchPanelBackgroundHandlerFields() map[string]any {
	return k.TouchPanelBackgroundHandlerFields
}

func (k *KnobV3) SetTouchPanelBackgroundHandler(handler ForegroundHandler) {
	k.TouchPanelBackgroundHandler = handler
}

func (k *KnobV3) GetTouchPanelBackgroundHandler() ForegroundHandler {
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

func (k *KnobActionV3) GetSwitchPage() int {
	return k.SwitchPage
}

func (k *KnobActionV3) GetKeyBind() string {
	return k.Keybind
}

func (k *KnobActionV3) GetCommand() string {
	return k.Command
}

func (k *KnobActionV3) GetBrightness() int {
	return k.Brightness
}

func (k *KnobActionV3) GetUrl() string {
	return k.Url
}

func (k *KnobActionV3) GetObsCommand() string {
	return k.ObsCommand
}

func (k *KnobActionV3) GetObsCommandParams() map[string]string {
	return k.ObsCommandParams
}

type KnobConfigV3 struct {
	Icon                              string            `json:"icon,omitempty"`
	Text                              string            `json:"text,omitempty"`
	TextSize                          int               `json:"text_size,omitempty"`
	TextAlignment                     VerticalAlignment `json:"text_alignment,omitempty"`
	FontFace                          string            `json:"font_face,omitempty"`
	TextColour                        string            `json:"text_colour,omitempty"`
	LcdHandler                        string            `json:"lcd_handler,omitempty"`
	KnobOrTouchHandler                string            `json:"knob_or_touch_handler,omitempty"`
	LcdHandlerStruct                  ForegroundHandler `json:"-"`
	KnobOrTouchHandlerStruct          InputHandler      `json:"-"`
	LcdHandlerFields                  map[string]any    `json:"lcd_handler_fields,omitempty"`
	KnobOrTouchHandlerFields          map[string]any    `json:"knob_or_touch_handler_fields,omitempty"`
	SharedHandlerFields               map[string]any    `json:"shared_handler_fields,omitempty"`
	KnobPressAction                   KnobActionV3      `json:"knob_press_action,omitempty"`
	KnobTurnUpAction                  KnobActionV3      `json:"knob_turn_up_action,omitempty"`
	KnobTurnDownAction                KnobActionV3      `json:"knob_turn_down_action,omitempty"`
	TouchPanelBackground              string            `json:"touch_panel_background"`
	TouchPanelBackgroundBuff          image.Image       `json:"-"`
	TouchPanelBackgroundHandler       ForegroundHandler `json:"-"`
	TouchPanelBackgroundHandlerFields map[string]any    `json:"touch_panel_background_handler_fields"`
}

func (kc *KnobConfigV3) SetForegroundHandlerInstance(handler ForegroundHandler) {
	kc.LcdHandlerStruct = handler
}

func (kc *KnobConfigV3) GetForegroundHandlerFields() map[string]any {
	return kc.LcdHandlerFields
}

func (kc *KnobConfigV3) GetForegroundHandlerInstance() ForegroundHandler {
	return kc.LcdHandlerStruct
}

func (kc *KnobConfigV3) GetForegroundHandler() string {
	return kc.LcdHandler
}

func (kc *KnobConfigV3) GetInputHandler() string {
	return kc.KnobOrTouchHandler
}

func (kc *KnobConfigV3) GetInputHandlerInstance() InputHandler {
	return kc.KnobOrTouchHandlerStruct
}

func (kc *KnobConfigV3) SetInputHandlerInstance(handler InputHandler) {
	kc.KnobOrTouchHandlerStruct = handler
}

func (kc *KnobConfigV3) GetInputHandlerFields() map[string]any {
	return kc.KnobOrTouchHandlerFields
}

func (kc *KnobConfigV3) GetSharedHandlerFields() map[string]any {
	return kc.SharedHandlerFields
}

func (kc *KnobConfigV3) GetIcon() string {
	return kc.Icon
}

func (kc *KnobConfigV3) GetText() string {
	return kc.Text
}

func (kc *KnobConfigV3) GetTextSize() int {
	return kc.TextSize
}

func (kc *KnobConfigV3) GetTextAlignment() VerticalAlignment {
	return kc.TextAlignment
}

func (kc *KnobConfigV3) GetFontFace() string {
	return kc.FontFace
}

func (kc *KnobConfigV3) GetTextColour() string {
	return kc.TextColour
}

func (kc *KnobConfigV3) GetTouchPanelBackgroundHandlerFields() map[string]any {
	return kc.TouchPanelBackgroundHandlerFields
}

func (kc *KnobConfigV3) SetTouchPanelBackgroundHandler(handler ForegroundHandler) {
	kc.TouchPanelBackgroundHandler = handler
}

func (kc *KnobConfigV3) GetTouchPanelBackgroundHandler() ForegroundHandler {
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
