package examples

import (
	"log"

	"github.com/unix-streamdeck/api/v3"
)

const (
	minBrightness = 0
	maxBrightness = 100
)

type BrightnessHandler struct{}

func (b *BrightnessHandler) Input(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfo, event api.InputEvent, host api.IHost) {
	if !shouldAct(event) {
		return
	}

	brightness, ok := fieldInt(fields, "brightness")
	if !ok {
		log.Println("[WARN] Brightness handler: no brightness specified")
		return
	}

	if brightness < minBrightness || brightness > maxBrightness {
		log.Printf("[WARN] Brightness handler: %d is outside 0-100, clamping\n", brightness)
		brightness = min(max(brightness, minBrightness), maxBrightness)
	}

	if err := host.SetBrightness(uint8(brightness)); err != nil {
		log.Println("[ERROR] Brightness handler:", err)
	}
}

func RegisterBrightness() api.Module {
	return api.Module{
		Name: "Brightness",
		NewCombined: func() api.Handler {
			return &BrightnessHandler{}
		},
		CombinedFields: []api.Field{
			{Title: "Brightness", Name: "brightness", Type: api.Number},
		},
	}
}
