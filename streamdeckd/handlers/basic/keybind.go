package basic

import (
	"log"

	"github.com/unix-streamdeck/api/v3"
)

type KeybindHandler struct{}

func (k *KeybindHandler) Input(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfo, event api.InputEvent, host api.IHost) {
	if !shouldAct(event) {
		return
	}

	keybind, ok := fieldString(fields, "keybind")
	if !ok {
		log.Println("[WARN] Keybind handler: no keybind specified")
		return
	}

	if err := host.ExecuteKeybind(keybind); err != nil {
		log.Println("[ERROR] Keybind handler: failed to execute keybind:", err)
	}
}

func RegisterKeybind() api.Module {
	return api.Module{
		Name:        "Keybind",
		Description: "Simulates a keyboard shortcut on press",
		NewCombined: func() api.Handler {
			return &KeybindHandler{}
		},
		CombinedFields: []api.Field{
			{Title: "Keybind", Name: "keybind", Type: api.Text},
		},
	}
}
