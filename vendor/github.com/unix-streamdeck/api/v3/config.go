package api

import (
	"image"
	"math"
	"time"
)

type LcdBackgrounder interface {
	GetTouchPanelBackground() string
	GetTouchPanelBackgroundBuff() []image.Image
	SetTouchPanelBackgroundBuff(img []image.Image)
	GetTouchPanelBackgroundHandler() BackgroundHandler
	SetTouchPanelBackgroundHandler(handler BackgroundHandler)
	GetTouchPanelBackgroundHandlerFields() map[string]any
}

type KeyGridBackgrounder interface {
	GetKeyGridBackground() string
	GetKeyGridBackgroundBuff() []image.Image
	SetKeyGridBackgroundBuff(img []image.Image)
	GetKeyGridBackgroundHandler() BackgroundHandler
	SetKeyGridBackgroundHandler(handler BackgroundHandler)
	GetKeyGridBackgroundHandlerFields() map[string]any
}

type Config struct {
	Modules           []string          `json:"modules,omitempty"`
	Decks             []Deck            `json:"decks"`
	ObsConnectionInfo ObsConnectionInfo `json:"obs_connection_info"`
}

type Deck struct {
	Serial                            string            `json:"serial"`
	Pages                             []Page            `json:"pages"`
	TouchPanelBackground              string            `json:"touch_panel_background"`
	TouchPanelBackgroundBuff          []image.Image     `json:"-"`
	TouchPanelBackgroundHandler       BackgroundHandler `json:"-"`
	TouchPanelBackgroundHandlerFields map[string]any    `json:"touch_panel_background_handler_fields"`
	KeyGridBackground                 string            `json:"key_grid_background"`
	KeyGridBackgroundBuff             []image.Image     `json:"-"`
	KeyGridBackgroundHandler          BackgroundHandler `json:"-"`
	KeyGridBackgroundHandlerFields    map[string]any    `json:"key_grid_background_handler_fields"`
}

func (d *Deck) SetTouchPanelBackgroundHandler(handler BackgroundHandler) {
	d.TouchPanelBackgroundHandler = handler
}

func (d *Deck) GetTouchPanelBackgroundHandler() BackgroundHandler {
	return d.TouchPanelBackgroundHandler
}

func (d *Deck) GetTouchPanelBackground() string {
	return d.TouchPanelBackground
}

func (d *Deck) GetTouchPanelBackgroundBuff() []image.Image {
	return d.TouchPanelBackgroundBuff
}

func (d *Deck) SetTouchPanelBackgroundBuff(img []image.Image) {
	d.TouchPanelBackgroundBuff = img
}

func (d *Deck) GetTouchPanelBackgroundHandlerFields() map[string]any {
	return d.TouchPanelBackgroundHandlerFields
}

func (d *Deck) GetKeyGridBackground() string {
	return d.KeyGridBackground
}

func (d *Deck) GetKeyGridBackgroundBuff() []image.Image {
	return d.KeyGridBackgroundBuff
}

func (d *Deck) SetKeyGridBackgroundBuff(img []image.Image) {
	d.KeyGridBackgroundBuff = img
}

func (d *Deck) SetKeyGridBackgroundHandler(handler BackgroundHandler) {
	d.KeyGridBackgroundHandler = handler
}

func (d *Deck) GetKeyGridBackgroundHandler() BackgroundHandler {
	return d.KeyGridBackgroundHandler
}

func (d *Deck) GetKeyGridBackgroundHandlerFields() map[string]any {
	return d.KeyGridBackgroundHandlerFields
}

type Page struct {
	Keys                              []Key             `json:"keys"`
	Knobs                             []Knob            `json:"knobs"`
	TouchPanelBackground              string            `json:"touch_panel_background"`
	TouchPanelBackgroundBuff          []image.Image     `json:"-"`
	TouchPanelBackgroundHandler       BackgroundHandler `json:"-"`
	TouchPanelBackgroundHandlerFields map[string]any    `json:"touch_panel_background_handler_fields"`
	KeyGridBackground                 string            `json:"key_grid_background"`
	KeyGridBackgroundBuff             []image.Image     `json:"-"`
	KeyGridBackgroundHandler          BackgroundHandler `json:"-"`
	KeyGridBackgroundHandlerFields    map[string]any    `json:"key_grid_background_handler_fields"`
}

func (p *Page) SetTouchPanelBackgroundHandler(handler BackgroundHandler) {
	p.TouchPanelBackgroundHandler = handler
}

func (p *Page) GetTouchPanelBackgroundHandler() BackgroundHandler {
	return p.TouchPanelBackgroundHandler
}

func (p *Page) GetTouchPanelBackground() string {
	return p.TouchPanelBackground
}

func (p *Page) GetTouchPanelBackgroundBuff() []image.Image {
	return p.TouchPanelBackgroundBuff
}

func (p *Page) SetTouchPanelBackgroundBuff(img []image.Image) {
	p.TouchPanelBackgroundBuff = img
}

func (p *Page) GetTouchPanelBackgroundHandlerFields() map[string]any {
	return p.TouchPanelBackgroundHandlerFields
}

func (p *Page) GetKeyGridBackground() string {
	return p.KeyGridBackground
}

func (p *Page) GetKeyGridBackgroundBuff() []image.Image {
	return p.KeyGridBackgroundBuff
}

func (p *Page) SetKeyGridBackgroundBuff(img []image.Image) {
	p.KeyGridBackgroundBuff = img
}

func (p *Page) SetKeyGridBackgroundHandler(handler BackgroundHandler) {
	p.KeyGridBackgroundHandler = handler
}

func (p *Page) GetKeyGridBackgroundHandler() BackgroundHandler {
	return p.KeyGridBackgroundHandler
}

func (p *Page) GetKeyGridBackgroundHandlerFields() map[string]any {
	return p.KeyGridBackgroundHandlerFields
}

type Layer struct {
	Handler     string           `json:"handler,omitempty"`
	Fields      map[string]any   `json:"fields,omitempty"`
	Visual      bool             `json:"visual,omitempty"`
	Input       bool             `json:"input,omitempty"`
	InputEvents []InputEventType `json:"input_events"`
	Global      bool             `json:"global,omitempty"`
	Label       string           `json:"label,omitempty"`

	HandlerStruct Handler     `json:"-"`
	Buff          image.Image `json:"-"`
}

type Layerable interface {
	GetLayers() []*Layer
}

type StreamDeckInfo struct {
	Cols                    int       `json:"cols,omitempty"`
	Rows                    int       `json:"rows,omitempty"`
	IconSize                int       `json:"icon_size,omitempty"`
	Page                    int       `json:"page"`
	Serial                  string    `json:"serial,omitempty"`
	Name                    string    `json:"name,omitempty"`
	Connected               bool      `json:"connected"`
	LastConnected           time.Time `json:"last_connected"`
	LastDisconnected        time.Time `json:"last_disconnected"`
	LcdWidth                int       `json:"lcd_width,omitempty"`
	LcdHeight               int       `json:"lcd_height,omitempty"`
	LcdCols                 int       `json:"lcd_cols,omitempty"`
	KnobCols                int       `json:"knob_cols,omitempty"`
	PaddingX                int       `json:"padding_x"`
	PaddingY                int       `json:"padding_y"`
	KeyGridBackgroundWidth  int       `json:"key_grid_background_width"`
	KeyGridBackgroundHeight int       `json:"key_grid_background_height"`
	LcdBackgroundWidth      int       `json:"lcd_background_width"`
	LcdBackgroundHeight     int       `json:"lcd_background_height"`
}

func (info StreamDeckInfo) GetDimensions(handlerType HandlerType) (int, int) {
	if handlerType == LCD {
		return info.LcdWidth, info.LcdHeight
	}
	return info.IconSize, info.IconSize
}

func (info StreamDeckInfo) GetGridDimensions(handlerType HandlerType) (int, int) {
	if handlerType == LCD {
		return info.LcdBackgroundWidth, info.LcdBackgroundHeight
	}
	return info.KeyGridBackgroundWidth, info.KeyGridBackgroundHeight
}

func (info StreamDeckInfo) SplitBackgroundImage(background image.Image, handlerType HandlerType) []image.Image {
	var frameArr []image.Image

	if handlerType == KEY {
		for keyIndex := range info.Cols * info.Rows {
			keyX := keyIndex % info.Cols
			keyY := int(math.Floor(float64(keyIndex) / float64(info.Cols)))

			x0, y0 := keyX*(info.IconSize+info.PaddingX), keyY*(info.IconSize+info.PaddingY)
			x1, y1 := keyX*(info.IconSize+info.PaddingX)+info.IconSize, keyY*(info.IconSize+info.PaddingY)+info.IconSize

			frameArr = append(frameArr, SubImage(background, x0, y0, x1, y1))
		}
	} else {
		for lcdIndex := range info.LcdCols {
			x0, y0 := info.LcdWidth*lcdIndex, 0
			x1, y1 := info.LcdWidth*(lcdIndex+1), info.LcdHeight

			subImage := SubImage(background, x0, y0, x1, y1)

			frameArr = append(frameArr, subImage)
		}
	}

	return frameArr
}

type ObsConnectionInfo struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}
