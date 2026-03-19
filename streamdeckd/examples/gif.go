package examples

import (
	"github.com/unix-streamdeck/api/v2"
	"golang.org/x/sync/semaphore"

	"github.com/unix-streamdeck/streamdeckd/streamdeckd/examples/gif"
)

func RegisterGif() api.Module {
	return api.Module{
		Name: "Gif",
		NewIcon: func() api.IconHandler {
			return &gif.GifIconHandler{Running: true, Lock: semaphore.NewWeighted(1)}
		},
		IconFields: []api.Field{
			{Title: "Icon", Name: "icon", Type: api.File, FileTypes: []string{".gif"}},
			{Title: "Text", Name: "text", Type: api.Text},
			{Title: "Text Size", Name: "text_size", Type: api.Number},
			{Title: "Text Alignment", Name: "text_alignment", Type: api.TextAlignment},
			{Title: "Font Face", Name: "font_face", Type: api.FontFace},
			{Title: "Text Colour", Name: "text_colour", Type: api.Colour},
		},
		NewLcd: func() api.LcdHandler {
			return &gif.GifLcdHandler{Running: true, Lock: semaphore.NewWeighted(1)}
		},
		LcdFields: []api.Field{
			{Title: "Icon", Name: "icon", Type: api.File, FileTypes: []string{".gif"}},
			{Title: "Text", Name: "text", Type: api.Text},
			{Title: "Text Size", Name: "text_size", Type: api.Number},
			{Title: "Text Alignment", Name: "text_alignment", Type: api.TextAlignment},
			{Title: "Font Face", Name: "font_face", Type: api.FontFace},
			{Title: "Text Colour", Name: "text_colour", Type: api.Colour},
		},
		NewKeyGridBackground: func() api.KeyGridBackgroundHandler {
			return &gif.GifKeyGridBackgroundHandler{Running: true, Lock: semaphore.NewWeighted(1)}
		},
		KeyGridBackgroundFields: []api.Field{
			{Title: "Icon", Name: "icon", Type: api.File, FileTypes: []string{".gif"}},
			{Title: "Text", Name: "text", Type: api.Text},
			{Title: "Text Size", Name: "text_size", Type: api.Number},
			{Title: "Text Alignment", Name: "text_alignment", Type: api.TextAlignment},
			{Title: "Font Face", Name: "font_face", Type: api.FontFace},
			{Title: "Text Colour", Name: "text_colour", Type: api.Colour},
		},
		NewTouchPanelBackgroundHandler: func() api.TouchPanelBackgroundHandler {
			return &gif.GifTouchPanelBackgroundHandler{Running: true, Lock: semaphore.NewWeighted(1)}
		},
		TouchPanelBackgroundFields: []api.Field{
			{Title: "Icon", Name: "icon", Type: api.File, FileTypes: []string{".gif"}},
			{Title: "Text", Name: "text", Type: api.Text},
			{Title: "Text Size", Name: "text_size", Type: api.Number},
			{Title: "Text Alignment", Name: "text_alignment", Type: api.TextAlignment},
			{Title: "Font Face", Name: "font_face", Type: api.FontFace},
			{Title: "Text Colour", Name: "text_colour", Type: api.Colour},
		},
	}
}
