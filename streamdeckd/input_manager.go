package streamdeckd

import (
	"log"
	"slices"

	"github.com/unix-streamdeck/api/v3"
	streamdeck "github.com/unix-streamdeck/driver"
)

type IInputManager interface {
	HandleKeyInput(key *api.Key, event streamdeck.InputEvent)
	HandleKnobInput(knob *api.Knob, event streamdeck.InputEvent)
	GetKeyState(index int) bool
	AttachPageChangeListener()
	AttachAppChangeListener()
}

type InputManager struct {
	vdev      IVirtualDev
	KeyStates []bool
}

func (im *InputManager) HandleKeyInput(key *api.Key, event streamdeck.InputEvent) {
	keyConfig, ok := key.Application[key.ActiveApplication]
	if !ok {
		im.vdev.Logger().Println("Err getting correct config for key")
		return
	}

	if event.EventType == streamdeck.KEY_PRESS {
		im.KeyStates[event.Index] = true
		im.vdev.RedrawKey(int(event.Index))

	} else {
		im.KeyStates[event.Index] = false
		im.vdev.RedrawKey(int(event.Index))
	}

	for _, layer := range MergeLayers(keyConfig, key) {
		im.handleHandlerAction(layer, api.KEY, event)
	}
}

func (im *InputManager) HandleKnobInput(knob *api.Knob, event streamdeck.InputEvent) {
	knobConfig, ok := knob.Application[knob.ActiveApplication]
	if !ok {
		im.vdev.Logger().Println("Err getting correct config for knob")
		return
	}
	for _, layer := range MergeLayers(knobConfig, knob) {
		im.handleHandlerAction(layer, api.LCD, event)
	}
}

func (im *InputManager) handleHandlerAction(layer *api.Layer, handlerType api.HandlerType, event streamdeck.InputEvent) {
	if !layer.Input {
		return
	}
	if len(layer.InputEvents) > 0 && !slices.Contains(layer.InputEvents, api.InputEventType(event.EventType)) {
		return
	}
	if layer.Handler != "" {
		var deckInfo api.StreamDeckInfo
		deckInfo = *im.vdev.SdInfo()
		if layer.HandlerStruct == nil {
			var handler api.Handler

			module, err := GetModule(layer.Handler, INPUT)
			if err != nil {
				im.vdev.Logger().Println(err)
				return
			}

			handler = module.NewCombined()
			if handler == nil {
				im.vdev.Logger().Println("Could not find handler:", layer.Handler)
				return
			}
			layer.HandlerStruct = handler
		}
		fields := layer.Fields

		inputEvent := api.InputEvent{
			EventType:     api.InputEventType(event.EventType),
			RotateNotches: event.RotateNotches,
		}
		if event.EventType == streamdeck.SCREEN_SHORT_TAP || event.EventType == streamdeck.SCREEN_LONG_TAP {
			inputEvent.ScreenTapY = event.ScreenY
			inputEvent.ScreenTapX = event.ScreenX - uint16(int(event.Index)*deckInfo.LcdWidth)
		}
		inHandler, ok := layer.HandlerStruct.(api.InputHandler)
		if !ok {
			log.Println("[Warn] Handler", layer.Handler, "is not valid InputHandler")
			return
		}
		inHandler.Input(fields, handlerType, deckInfo, inputEvent, im.vdev.Host())
	}
}

func (im *InputManager) GetKeyState(index int) bool {
	return im.KeyStates[index]
}

func (im *InputManager) AttachPageChangeListener() {
	im.vdev.PageManager().AttachListener(func(_, oldPage int) {
		currentPage := im.vdev.Config().Pages[oldPage]

		for i, _ := range currentPage.Keys {
			key := &currentPage.Keys[i]
			im.cleanupKeyHolds(MergeLayers(key.Application[key.ActiveApplication], key), i)
		}
	})
}

func (im *InputManager) AttachAppChangeListener() {
	applicationManager.AttachListener(func(application string) {
		page := im.vdev.Config().Pages[im.vdev.PageManager().GetPage()]
		for i := range page.Keys {
			key := &page.Keys[i]
			_, keyHasApp := key.Application[application]
			activeApp := key.ActiveApplication
			appChange := false
			if key.ActiveApplication != "" && !keyHasApp {
				appChange = true
			}
			if keyHasApp {
				appChange = true
			}
			if appChange {
				im.cleanupKeyHolds(MergeLayers(key.Application[activeApp], key), i)
			}
		}
	})
}

func (im *InputManager) cleanupKeyHolds(layers []*api.Layer, index int) {
	if im.KeyStates[index] {
		for _, layer := range layers {
			if layer.Input && layer.HandlerStruct != nil {
				inHandler, ok := layer.HandlerStruct.(api.InputHandler)
				if !ok {
					continue
				}
				inHandler.Input(layer.Fields, api.KEY, *im.vdev.SdInfo(), api.InputEvent{
					EventType: api.KEY_RELEASE,
				}, im.vdev.Host())
			}
		}
	}
}
