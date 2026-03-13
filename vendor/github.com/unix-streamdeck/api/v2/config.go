package api

import (
	"image"
	"time"
)

type LcdBackgrounder interface {
	GetTouchPanelBackground() string
	GetTouchPanelBackgroundBuff() image.Image
	SetTouchPanelBackgroundBuff(img image.Image)
}

type KeyGridBackgrounder interface {
	GetKeyGridBackground() string
	GetKeyGridBackgroundBuff() image.Image
	SetKeyGridBackgroundBuff(img image.Image)
}

type ConfigV3 struct {
	Modules           []string            `json:"modules,omitempty"`
	Decks             []DeckV3            `json:"decks"`
	ObsConnectionInfo ObsConnectionInfoV2 `json:"obs_connection_info,omitempty"`
}

type DeckV3 struct {
	Serial                   string      `json:"serial"`
	Pages                    []PageV3    `json:"pages"`
	TouchPanelBackground     string      `json:"touch_panel_background"`
	TouchPanelBackgroundBuff image.Image `json:"-"`
	KeyGridBackground        string      `json:"key_grid_background"`
	KeyGridBackgroundBuff    image.Image `json:"-"`
}

func (d *DeckV3) GetTouchPanelBackground() string {
	return d.TouchPanelBackground
}

func (d *DeckV3) GetTouchPanelBackgroundBuff() image.Image {
	return d.TouchPanelBackgroundBuff
}

func (d *DeckV3) SetTouchPanelBackgroundBuff(img image.Image) {
	d.TouchPanelBackgroundBuff = img
}

func (d *DeckV3) GetKeyGridBackground() string {
	return d.KeyGridBackground
}

func (d *DeckV3) GetKeyGridBackgroundBuff() image.Image {
	return d.KeyGridBackgroundBuff
}

func (d *DeckV3) SetKeyGridBackgroundBuff(img image.Image) {
	d.KeyGridBackgroundBuff = img
}

type PageV3 struct {
	Keys                     []KeyV3     `json:"keys"`
	Knobs                    []KnobV3    `json:"knobs"`
	TouchPanelBackground     string      `json:"touch_panel_background"`
	TouchPanelBackgroundBuff image.Image `json:"-"`
	KeyGridBackground        string      `json:"key_grid_background"`
	KeyGridBackgroundBuff    image.Image `json:"-"`
}

func (p *PageV3) GetTouchPanelBackground() string {
	return p.TouchPanelBackground
}

func (p *PageV3) GetTouchPanelBackgroundBuff() image.Image {
	return p.TouchPanelBackgroundBuff
}

func (p *PageV3) SetTouchPanelBackgroundBuff(img image.Image) {
	p.TouchPanelBackgroundBuff = img
}

func (p *PageV3) GetKeyGridBackground() string {
	return p.KeyGridBackground
}

func (p *PageV3) GetKeyGridBackgroundBuff() image.Image {
	return p.KeyGridBackgroundBuff
}

func (p *PageV3) SetKeyGridBackgroundBuff(img image.Image) {
	p.KeyGridBackgroundBuff = img
}

type KeyV3 struct {
	Application             map[string]*KeyConfigV3 `json:"application,omitempty"`
	ActiveBuff              image.Image             `json:"-"`
	ActiveIconHandlerStruct *IconHandler            `json:"-"`
	ActiveKeyHandlerStruct  *KeyHandler             `json:"-"`
	ActiveApplication       string                  `json:"-"`
}

type KnobV3 struct {
	Application              map[string]*KnobConfigV3 `json:"application,omitempty"`
	ActiveBuff               image.Image              `json:"-"`
	ActiveApplication        string                   `json:"-"`
	TouchPanelBackground     string                   `json:"touch_panel_background"`
	TouchPanelBackgroundBuff image.Image              `json:"-"`
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

type KeyConfigV3 struct {
	Icon                string            `json:"icon,omitempty"`
	SwitchPage          int               `json:"switch_page,omitempty"`
	Text                string            `json:"text,omitempty"`
	TextSize            int               `json:"text_size,omitempty"`
	TextAlignment       string            `json:"text_alignment,omitempty"`
	FontFace            string            `json:"font_face,omitempty"`
	TextColour          string            `json:"text_colour,omitempty"`
	Keybind             string            `json:"keybind,omitempty"`
	Command             string            `json:"command,omitempty"`
	Brightness          int               `json:"brightness,omitempty"`
	Url                 string            `json:"url,omitempty"`
	KeyHold             int               `json:"key_hold,omitempty"`
	ObsCommand          string            `json:"obs_command,omitempty"`
	ObsCommandParams    map[string]string `json:"obs_command_params,omitempty"`
	IconHandler         string            `json:"icon_handler,omitempty"`
	KeyHandler          string            `json:"key_handler,omitempty"`
	IconHandlerFields   map[string]any    `json:"icon_handler_fields,omitempty"`
	KeyHandlerFields    map[string]any    `json:"key_handler_fields,omitempty"`
	SharedHandlerFields map[string]any    `json:"shared_handler_fields,omitempty"`
	Buff                image.Image       `json:"-"`
	IconHandlerStruct   IconHandler       `json:"-"`
	KeyHandlerStruct    KeyHandler        `json:"-"`
	SharedState         map[string]any    `json:"-"`
}

type KnobConfigV3 struct {
	Icon                     string             `json:"icon,omitempty"`
	Text                     string             `json:"text,omitempty"`
	TextSize                 int                `json:"text_size,omitempty"`
	TextAlignment            string             `json:"text_alignment,omitempty"`
	FontFace                 string             `json:"font_face,omitempty"`
	TextColour               string             `json:"text_colour,omitempty"`
	LcdHandler               string             `json:"lcd_handler,omitempty"`
	KnobOrTouchHandler       string             `json:"knob_or_touch_handler,omitempty"`
	Buff                     image.Image        `json:"-"`
	LcdHandlerStruct         LcdHandler         `json:"-"`
	KnobOrTouchHandlerStruct KnobOrTouchHandler `json:"-"`
	LcdHandlerFields         map[string]any     `json:"lcd_handler_fields,omitempty"`
	KnobOrTouchHandlerFields map[string]any     `json:"knob_or_touch_handler_fields,omitempty"`
	SharedHandlerFields      map[string]any     `json:"shared_handler_fields,omitempty"`
	KnobPressAction          KnobActionV3       `json:"knob_press_action,omitempty"`
	KnobTurnUpAction         KnobActionV3       `json:"knob_turn_up_action,omitempty"`
	KnobTurnDownAction       KnobActionV3       `json:"knob_turn_down_action,omitempty"`
	SharedState              map[string]any     `json:"-"`
	TouchPanelBackground     string             `json:"touch_panel_background"`
	TouchPanelBackgroundBuff image.Image        `json:"-"`
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

type KnobActionV3 struct {
	SwitchPage       int               `json:"switch_page,omitempty"`
	Keybind          string            `json:"keybind,omitempty"`
	Command          string            `json:"command,omitempty"`
	Brightness       int               `json:"brightness,omitempty"`
	Url              string            `json:"url,omitempty"`
	ObsCommand       string            `json:"obs_command,omitempty"`
	ObsCommandParams map[string]string `json:"obs_command_params,omitempty"`
}

type StreamDeckInfoV1 struct {
	Cols             int       `json:"cols,omitempty"`
	Rows             int       `json:"rows,omitempty"`
	IconSize         int       `json:"icon_size,omitempty"`
	Page             int       `json:"page"`
	Serial           string    `json:"serial,omitempty"`
	Name             string    `json:"name,omitempty"`
	Connected        bool      `json:"connected"`
	LastConnected    time.Time `json:"last_connected,omitempty"`
	LastDisconnected time.Time `json:"last_disconnected,omitempty"`
	LcdWidth         int       `json:"lcd_width,omitempty"`
	LcdHeight        int       `json:"lcd_height,omitempty"`
	LcdCols          int       `json:"lcd_cols,omitempty"`
	KnobCols         int       `json:"knob_cols,omitempty"`
}

type ObsConnectionInfoV2 struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}
