package examples

import (
	"log"
	"sync"

	"github.com/unix-streamdeck/api/v3"
)

type KeyHoldHandler struct {
	mu   sync.Mutex
	held []int
}

func (k *KeyHoldHandler) Input(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfo, event api.InputEvent, host api.IHost) {
	switch event.EventType {
	case api.KEY_PRESS:
		keycodes, ok := keyHoldCodes(fields)
		if !ok {
			return
		}
		k.press(keycodes, host)
	case api.KEY_RELEASE:
		k.release(host)
	}
}

func (k *KeyHoldHandler) press(keycodes []int, host api.IHost) {
	k.mu.Lock()
	defer k.mu.Unlock()

	if len(k.held) > 0 {
		return
	}

	for index, keycode := range keycodes {
		if err := host.KeyDown(keycode); err != nil {
			log.Printf("[ERROR] KeyHold handler: failed to press key %d: %s\n", keycode, err)

			// Let go of whatever did go down, so a partial combination is not left held.
			for i := index - 1; i >= 0; i-- {
				if err := host.KeyUp(keycodes[i]); err != nil {
					log.Printf("[WARN] KeyHold handler: failed to release key %d: %s\n", keycodes[i], err)
				}
			}

			k.held = nil
			return
		}

		k.held = append(k.held, keycode)
	}
}

func (k *KeyHoldHandler) release(host api.IHost) {
	k.mu.Lock()
	defer k.mu.Unlock()

	for i := len(k.held) - 1; i >= 0; i-- {
		if err := host.KeyUp(k.held[i]); err != nil {
			log.Printf("[WARN] KeyHold handler: failed to release key %d: %s\n", k.held[i], err)
		}
	}

	k.held = nil
}

func keyHoldCodes(fields map[string]any) ([]int, bool) {
	if keybind, ok := fieldString(fields, "keybind"); ok {
		keycodes, err := api.ParseXDoToolKeybindString(keybind)
		if err != nil {
			log.Println("[ERROR] KeyHold handler: failed to parse keybind:", err)
			return nil, false
		}
		return keycodes, true
	}

	if keycode, ok := fieldInt(fields, "keycode"); ok {
		return []int{keycode}, true
	}

	log.Println("[WARN] KeyHold handler: no keybind or keycode specified")
	return nil, false
}

func RegisterKeyHold() api.Module {
	return api.Module{
		Name: "KeyHold",
		NewCombined: func() api.Handler {
			return &KeyHoldHandler{}
		},
		CombinedFields: []api.Field{
			{Title: "Keybind", Name: "keybind", Type: api.Text},
			{Title: "Keycode", Name: "keycode", Type: api.Number},
		},
	}
}
