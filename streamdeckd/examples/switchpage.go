package examples

import (
	"log"

	"github.com/unix-streamdeck/api/v3"
)

type SwitchPageHandler struct{}

func (s *SwitchPageHandler) Input(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfo, event api.InputEvent, host api.IHost) {
	if !shouldAct(event) {
		return
	}

	page, ok := fieldInt(fields, "page")
	page -= 1
	if !ok {
		log.Println("[WARN] SwitchPage handler: no page specified")
		return
	}

	if page < 0 {
		log.Printf("[WARN] SwitchPage handler: page %d is out of range\n", page)
		return
	}

	host.SetPage(page)
}

func RegisterSwitchPage() api.Module {
	return api.Module{
		Name: "SwitchPage",
		NewCombined: func() api.Handler {
			return &SwitchPageHandler{}
		},
		CombinedFields: []api.Field{
			{Title: "Page", Name: "page", Type: api.Number},
		},
	}
}
