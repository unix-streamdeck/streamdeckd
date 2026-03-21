package api

import "image"

type KeyBackgrounder interface {
	GetKeyBackground() string
	GetKeyBackgroundBuff() image.Image
	SetKeyBackgroundBuff(img image.Image)
	GetKeyBackgroundHandler() KeyGridBackgroundHandler
	SetKeyBackgroundHandler(handler KeyGridBackgroundHandler)
	GetKeyBackgroundHandlerFields() map[string]any
}

type KeyV3 struct {
	Application                map[string]*KeyConfigV3  `json:"application,omitempty"`
	ActiveBuff                 image.Image              `json:"-"`
	ActiveIconHandlerStruct    *IconHandler             `json:"-"`
	ActiveKeyHandlerStruct     *KeyHandler              `json:"-"`
	ActiveApplication          string                   `json:"-"`
	KeyBackground              string                   `json:"background"`
	KeyBackgroundBuff          image.Image              `json:"-"`
	KeyBackgroundHandler       KeyGridBackgroundHandler `json:"-"`
	KeyBackgroundHandlerFields map[string]any           `json:"key_background_handler_fields"`
}

func (k *KeyV3) GetKeyBackgroundHandlerFields() map[string]any {
	return k.KeyBackgroundHandlerFields
}

func (k *KeyV3) GetKeyBackground() string {
	return k.KeyBackground
}

func (k *KeyV3) GetKeyBackgroundBuff() image.Image {
	return k.KeyBackgroundBuff
}

func (k *KeyV3) SetKeyBackgroundBuff(img image.Image) {
	k.KeyBackgroundBuff = img
}

func (k *KeyV3) SetKeyBackgroundHandler(handler KeyGridBackgroundHandler) {
	k.KeyBackgroundHandler = handler
}

func (k *KeyV3) GetKeyBackgroundHandler() KeyGridBackgroundHandler {
	return k.KeyBackgroundHandler
}

type KeyConfigV3 struct {
	Icon                       string                   `json:"icon,omitempty"`
	Text                       string                   `json:"text,omitempty"`
	TextSize                   int                      `json:"text_size,omitempty"`
	TextAlignment              VerticalAlignment        `json:"text_alignment,omitempty"`
	FontFace                   string                   `json:"font_face,omitempty"`
	TextColour                 string                   `json:"text_colour,omitempty"`
	SwitchPage                 int                      `json:"switch_page,omitempty"`
	Keybind                    string                   `json:"keybind,omitempty"`
	Command                    string                   `json:"command,omitempty"`
	Brightness                 int                      `json:"brightness,omitempty"`
	Url                        string                   `json:"url,omitempty"`
	KeyHold                    int                      `json:"key_hold,omitempty"`
	ObsCommand                 string                   `json:"obs_command,omitempty"`
	ObsCommandParams           map[string]string        `json:"obs_command_params,omitempty"`
	IconHandler                string                   `json:"icon_handler,omitempty"`
	KeyHandler                 string                   `json:"key_handler,omitempty"`
	IconHandlerFields          map[string]any           `json:"icon_handler_fields,omitempty"`
	KeyHandlerFields           map[string]any           `json:"key_handler_fields,omitempty"`
	SharedHandlerFields        map[string]any           `json:"shared_handler_fields,omitempty"`
	IconHandlerStruct          IconHandler              `json:"-"`
	KeyHandlerStruct           KeyHandler               `json:"-"`
	SharedState                map[string]any           `json:"-"`
	KeyBackground              string                   `json:"background"`
	KeyBackgroundBuff          image.Image              `json:"-"`
	KeyBackgroundHandler       KeyGridBackgroundHandler `json:"-"`
	KeyBackgroundHandlerFields map[string]any           `json:"key_background_handler_fields"`
}

func (kc *KeyConfigV3) GetIcon() string {
	return kc.Icon
}

func (kc *KeyConfigV3) GetText() string {
	return kc.Text
}

func (kc *KeyConfigV3) GetTextSize() int {
	return kc.TextSize
}

func (kc *KeyConfigV3) GetTextAlignment() VerticalAlignment {
	return kc.TextAlignment
}

func (kc *KeyConfigV3) GetFontFace() string {
	return kc.FontFace
}

func (kc *KeyConfigV3) GetTextColour() string {
	return kc.TextColour
}

func (kc *KeyConfigV3) GetSwitchPage() int {
	return kc.SwitchPage
}

func (kc *KeyConfigV3) GetKeyBind() string {
	return kc.Keybind
}

func (kc *KeyConfigV3) GetCommand() string {
	return kc.Command
}

func (kc *KeyConfigV3) GetBrightness() int {
	return kc.Brightness
}

func (kc *KeyConfigV3) GetUrl() string {
	return kc.Url
}

func (kc *KeyConfigV3) GetObsCommand() string {
	return kc.ObsCommand
}

func (kc *KeyConfigV3) GetObsCommandParams() map[string]string {
	return kc.ObsCommandParams
}

func (kc *KeyConfigV3) GetKeyBackgroundHandlerFields() map[string]any {
	return kc.KeyBackgroundHandlerFields
}

func (kc *KeyConfigV3) GetKeyBackground() string {
	return kc.KeyBackground
}

func (kc *KeyConfigV3) GetKeyBackgroundBuff() image.Image {
	return kc.KeyBackgroundBuff
}

func (kc *KeyConfigV3) SetKeyBackgroundBuff(img image.Image) {
	kc.KeyBackgroundBuff = img
}

func (kc *KeyConfigV3) SetKeyBackgroundHandler(handler KeyGridBackgroundHandler) {
	kc.KeyBackgroundHandler = handler
}

func (kc *KeyConfigV3) GetKeyBackgroundHandler() KeyGridBackgroundHandler {
	return kc.KeyBackgroundHandler
}
