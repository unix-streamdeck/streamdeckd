package examples

import (
	"image"
	"log"

	"github.com/unix-streamdeck/api/v3"
)

type TextHandler struct {
	Running bool
}

func (t *TextHandler) IsRunning() bool {
	return t.Running
}

func (t *TextHandler) SetRunning(running bool) {
	t.Running = running
}

func (t *TextHandler) Stop() {
	t.Running = false
}

func (t *TextHandler) Start(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfo, callback func(image image.Image)) {
	t.Running = true

	text, ok := fieldString(fields, "text")
	if !ok {
		log.Println("[WARN] Text handler: no text specified")
		return
	}

	options := api.DrawTextOptions{}

	if size, ok := fieldInt(fields, "size"); ok {
		options.FontSize = int64(size)
	}
	if alignment, ok := fieldString(fields, "alignment"); ok {
		options.VerticalAlignment = api.VerticalAlignment(alignment)
	}
	if fontFace, ok := fieldString(fields, "font_face"); ok {
		options.FontFace = fontFace
	}
	if colour, ok := fieldString(fields, "colour"); ok {
		options.Colour = colour
	}

	width, height := info.GetDimensions(handlerType)

	// Drawn onto a transparent base so the layers underneath show through.
	img, err := api.DrawText(image.NewRGBA(image.Rect(0, 0, width, height)), text, options)
	if err != nil {
		log.Println("[ERROR] Text handler:", err)
		return
	}

	callback(img)
}

func RegisterText() api.Module {
	return api.Module{
		Name: "Text",
		NewCombined: func() api.Handler {
			return &TextHandler{}
		},
		CombinedFields: []api.Field{
			{Title: "Text", Name: "text", Type: api.Text},
			{Title: "Size", Name: "size", Type: api.Number},
			{Title: "Alignment", Name: "alignment", Type: api.TextAlignment},
			{Title: "Font Face", Name: "font_face", Type: api.FontFace},
			{Title: "Colour", Name: "colour", Type: api.Colour},
		},
	}
}
