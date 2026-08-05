package examples

import (
	"log"

	"github.com/unix-streamdeck/api/v3"
)

type UrlHandler struct{}

func (u *UrlHandler) Input(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfo, event api.InputEvent, host api.IHost) {
	if !shouldAct(event) {
		return
	}

	url, ok := fieldString(fields, "url")
	if !ok {
		log.Println("[WARN] Url handler: no url specified")
		return
	}

	host.RunCommand("xdg-open " + url)
}

func RegisterUrl() api.Module {
	return api.Module{
		Name: "Url",
		NewCombined: func() api.Handler {
			return &UrlHandler{}
		},
		CombinedFields: []api.Field{
			{Title: "Url", Name: "url", Type: api.Text},
		},
	}
}
