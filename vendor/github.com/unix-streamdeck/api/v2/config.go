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

type InputActions interface {
	GetSwitchPage() int
	GetKeyBind() string
	GetCommand() string
	GetBrightness() int
	GetUrl() string
	GetObsCommand() string
	GetObsCommandParams() map[string]string
}

type ForegroundActions interface {
	GetIcon() string
	GetText() string
	GetTextSize() int
	GetTextAlignment() VerticalAlignment
	GetFontFace() string
	GetTextColour() string
}

type ForegroundAndInputHandlerConfig interface {
	GetForegroundHandler() string
	GetForegroundHandlerInstance() ForegroundHandler
	SetForegroundHandlerInstance(handler ForegroundHandler)
	GetForegroundHandlerFields() map[string]any
	GetInputHandler() string
	GetInputHandlerInstance() InputHandler
	SetInputHandlerInstance(handler InputHandler)
	GetInputHandlerFields() map[string]any
	GetSharedHandlerFields() map[string]any
}

type ConfigV3 struct {
	Modules           []string            `json:"modules,omitempty"`
	Decks             []DeckV3            `json:"decks"`
	ObsConnectionInfo ObsConnectionInfoV2 `json:"obs_connection_info,omitempty"`
}

type DeckV3 struct {
	Serial                            string            `json:"serial"`
	Pages                             []PageV3          `json:"pages"`
	TouchPanelBackground              string            `json:"touch_panel_background"`
	TouchPanelBackgroundBuff          []image.Image     `json:"-"`
	TouchPanelBackgroundHandler       BackgroundHandler `json:"-"`
	TouchPanelBackgroundHandlerFields map[string]any    `json:"touch_panel_background_handler_fields"`
	KeyGridBackground                 string            `json:"key_grid_background"`
	KeyGridBackgroundBuff             []image.Image     `json:"-"`
	KeyGridBackgroundHandler          BackgroundHandler `json:"-"`
	KeyGridBackgroundHandlerFields    map[string]any    `json:"key_grid_background_handler_fields"`
}

func (d *DeckV3) SetTouchPanelBackgroundHandler(handler BackgroundHandler) {
	d.TouchPanelBackgroundHandler = handler
}

func (d *DeckV3) GetTouchPanelBackgroundHandler() BackgroundHandler {
	return d.TouchPanelBackgroundHandler
}

func (d *DeckV3) GetTouchPanelBackground() string {
	return d.TouchPanelBackground
}

func (d *DeckV3) GetTouchPanelBackgroundBuff() []image.Image {
	return d.TouchPanelBackgroundBuff
}

func (d *DeckV3) SetTouchPanelBackgroundBuff(img []image.Image) {
	d.TouchPanelBackgroundBuff = img
}

func (d *DeckV3) GetTouchPanelBackgroundHandlerFields() map[string]any {
	return d.TouchPanelBackgroundHandlerFields
}

func (d *DeckV3) GetKeyGridBackground() string {
	return d.KeyGridBackground
}

func (d *DeckV3) GetKeyGridBackgroundBuff() []image.Image {
	return d.KeyGridBackgroundBuff
}

func (d *DeckV3) SetKeyGridBackgroundBuff(img []image.Image) {
	d.KeyGridBackgroundBuff = img
}

func (d *DeckV3) SetKeyGridBackgroundHandler(handler BackgroundHandler) {
	d.KeyGridBackgroundHandler = handler
}

func (d *DeckV3) GetKeyGridBackgroundHandler() BackgroundHandler {
	return d.KeyGridBackgroundHandler
}

func (d *DeckV3) GetKeyGridBackgroundHandlerFields() map[string]any {
	return d.KeyGridBackgroundHandlerFields
}

type PageV3 struct {
	Keys                              []KeyV3           `json:"keys"`
	Knobs                             []KnobV3          `json:"knobs"`
	TouchPanelBackground              string            `json:"touch_panel_background"`
	TouchPanelBackgroundBuff          []image.Image     `json:"-"`
	TouchPanelBackgroundHandler       BackgroundHandler `json:"-"`
	TouchPanelBackgroundHandlerFields map[string]any    `json:"touch_panel_background_handler_fields"`
	KeyGridBackground                 string            `json:"key_grid_background"`
	KeyGridBackgroundBuff             []image.Image     `json:"-"`
	KeyGridBackgroundHandler          BackgroundHandler `json:"-"`
	KeyGridBackgroundHandlerFields    map[string]any    `json:"key_grid_background_handler_fields"`
}

func (p *PageV3) SetTouchPanelBackgroundHandler(handler BackgroundHandler) {
	p.TouchPanelBackgroundHandler = handler
}

func (p *PageV3) GetTouchPanelBackgroundHandler() BackgroundHandler {
	return p.TouchPanelBackgroundHandler
}

func (p *PageV3) GetTouchPanelBackground() string {
	return p.TouchPanelBackground
}

func (p *PageV3) GetTouchPanelBackgroundBuff() []image.Image {
	return p.TouchPanelBackgroundBuff
}

func (p *PageV3) SetTouchPanelBackgroundBuff(img []image.Image) {
	p.TouchPanelBackgroundBuff = img
}

func (p *PageV3) GetTouchPanelBackgroundHandlerFields() map[string]any {
	return p.TouchPanelBackgroundHandlerFields
}

func (p *PageV3) GetKeyGridBackground() string {
	return p.KeyGridBackground
}

func (p *PageV3) GetKeyGridBackgroundBuff() []image.Image {
	return p.KeyGridBackgroundBuff
}

func (p *PageV3) SetKeyGridBackgroundBuff(img []image.Image) {
	p.KeyGridBackgroundBuff = img
}

func (p *PageV3) SetKeyGridBackgroundHandler(handler BackgroundHandler) {
	p.KeyGridBackgroundHandler = handler
}

func (p *PageV3) GetKeyGridBackgroundHandler() BackgroundHandler {
	return p.KeyGridBackgroundHandler
}

func (p *PageV3) GetKeyGridBackgroundHandlerFields() map[string]any {
	return p.KeyGridBackgroundHandlerFields
}

type StreamDeckInfoV1 struct {
	Cols                    int       `json:"cols,omitempty"`
	Rows                    int       `json:"rows,omitempty"`
	IconSize                int       `json:"icon_size,omitempty"`
	Page                    int       `json:"page"`
	Serial                  string    `json:"serial,omitempty"`
	Name                    string    `json:"name,omitempty"`
	Connected               bool      `json:"connected"`
	LastConnected           time.Time `json:"last_connected,omitempty"`
	LastDisconnected        time.Time `json:"last_disconnected,omitempty"`
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

func (info StreamDeckInfoV1) GetDimensions(handlerType HandlerType) (int, int) {
	if handlerType == LCD {
		return info.LcdWidth, info.LcdHeight
	}
	return info.IconSize, info.IconSize
}

func (info StreamDeckInfoV1) GetGridDimensions(handlerType HandlerType) (int, int) {
	if handlerType == LCD {
		return info.LcdBackgroundWidth, info.LcdBackgroundHeight
	}
	return info.KeyGridBackgroundWidth, info.KeyGridBackgroundHeight
}

func (info StreamDeckInfoV1) SplitBackgroundImage(background image.Image, handlerType HandlerType) []image.Image {
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

type ObsConnectionInfoV2 struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}
