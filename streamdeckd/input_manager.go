package streamdeckd

import (
	"github.com/unix-streamdeck/api/v2"
	streamdeck "github.com/unix-streamdeck/driver"
)

type IInputManager interface {
	HandleKeyInput(key *api.KeyV3, event streamdeck.InputEvent)
	HandleKnobInput(knob *api.KnobV3, event streamdeck.InputEvent)
}

type InputManager struct {
	vdev IVirtualDev
}

func (im *InputManager) HandleKeyInput(key *api.KeyV3, event streamdeck.InputEvent) {
	keyConfig, ok := key.Application[key.ActiveApplication]
	if !ok {
		im.vdev.Logger().Println("Err getting correct config for key")
		return
	}
	if event.EventType == streamdeck.KEY_PRESS {
		im.handleStandardActions(keyConfig)

		if keyConfig.KeyHandler != "" {
			var deckInfo api.StreamDeckInfoV1
			deckInfo = *im.vdev.SdInfo()
			if keyConfig.KeyHandlerStruct == nil {
				var handler api.InputHandler
				modules := AvailableModules()
				for _, module := range modules {
					if module.Name == keyConfig.KeyHandler {
						handler = module.NewInput()
					}
				}
				if handler == nil {
					im.vdev.Logger().Println("Could not find handler:", keyConfig.KeyHandler)
					return
				}
				keyConfig.KeyHandlerStruct = handler
			}
			trimmedKeyConfig := api.KeyConfigV3{KeyHandlerFields: keyConfig.KeyHandlerFields}
			if keyConfig.IconHandler == keyConfig.KeyHandler {
				trimmedKeyConfig.KeyHandlerFields = mergeSharedConfig(keyConfig.SharedHandlerFields, keyConfig.KeyHandlerFields)
				iconHandler := keyConfig.IconHandlerStruct
				comboHandler, ok := iconHandler.(api.CombinedHandler)
				if ok {
					keyConfig.KeyHandlerStruct = comboHandler
				}
			}
			keyConfig.KeyHandlerStruct.Input(trimmedKeyConfig.KeyHandlerFields, api.KEY, deckInfo, api.InputEvent{
				EventType:     api.InputEventType(event.EventType),
				RotateNotches: event.RotateNotches,
			})
		}
	}
	if keyConfig.KeyHold != 0 {
		if event.EventType == streamdeck.KEY_PRESS {
			err := kb.KeyDown(keyConfig.KeyHold)
			if err != nil {
				im.vdev.Logger().Println(err)
			}
		} else {
			err := kb.KeyUp(keyConfig.KeyHold)
			if err != nil {
				im.vdev.Logger().Println(err)
			}
		}
	}
}

func (im *InputManager) HandleKnobInput(knob *api.KnobV3, event streamdeck.InputEvent) {
	knobConfig, ok := knob.Application[knob.ActiveApplication]
	if !ok {
		im.vdev.Logger().Println("Err getting correct config for knob")
		return
	}
	if knobConfig.KnobOrTouchHandler != "" {
		var deckInfo api.StreamDeckInfoV1
		deckInfo = *im.vdev.SdInfo()
		if knobConfig.KnobOrTouchHandlerStruct == nil {
			var handler api.InputHandler
			modules := AvailableModules()
			for _, module := range modules {
				if module.Name == knobConfig.KnobOrTouchHandler {
					handler = module.NewInput()
				}
			}
			if handler == nil {
				im.vdev.Logger().Println("Could not find handler:", knobConfig.KnobOrTouchHandler)
				return
			}
			knobConfig.KnobOrTouchHandlerStruct = handler
		}
		trimmedKnobConfig := api.KnobConfigV3{KnobOrTouchHandlerFields: knobConfig.KnobOrTouchHandlerFields}
		if knobConfig.LcdHandler == knobConfig.KnobOrTouchHandler {
			trimmedKnobConfig.KnobOrTouchHandlerFields = mergeSharedConfig(knobConfig.SharedHandlerFields, knobConfig.KnobOrTouchHandlerFields)
			iconHandler := knobConfig.LcdHandlerStruct
			comboHandler, ok := iconHandler.(api.CombinedHandler)
			if ok {
				knobConfig.KnobOrTouchHandlerStruct = comboHandler
			}
		}
		knobConfig.KnobOrTouchHandlerStruct.Input(trimmedKnobConfig.KnobOrTouchHandlerFields, api.LCD, deckInfo, api.InputEvent{
			EventType:     api.InputEventType(event.EventType),
			RotateNotches: event.RotateNotches,
		})
	}
	var actions api.KnobActionV3
	if event.EventType == streamdeck.KNOB_PRESS {
		actions = knobConfig.KnobPressAction
	} else if event.EventType == streamdeck.KNOB_CCW {
		actions = knobConfig.KnobTurnDownAction
	} else if event.EventType == streamdeck.KNOB_CW {
		actions = knobConfig.KnobTurnUpAction
	}
	im.handleStandardActions(&actions)
}

func (im *InputManager) handleStandardActions(ia api.InputActions) {
	if ia.GetCommand() != "" {
		RunCommand(ia.GetCommand())
	}
	if ia.GetKeyBind() != "" {
		err := ExecuteKeybind(ia.GetKeyBind())
		if err != nil {
			im.vdev.Logger().Println("[ERROR] Failed to execute keybind:", err)
		}
	}
	if ia.GetSwitchPage() != 0 {
		page := ia.GetSwitchPage() - 1
		im.vdev.PageManager().SetPage(page)
	}
	if ia.GetBrightness() != 0 {
		err := im.vdev.SetBrightness(uint8(ia.GetBrightness()))
		if err != nil {
			im.vdev.Logger().Println(err)
		}
	}
	if ia.GetUrl() != "" {
		RunCommand("xdg-open " + ia.GetUrl())
	}
	if ia.GetObsCommand() != "" {
		runObsCommand(ia.GetObsCommand(), ia.GetObsCommandParams())
	}
}
