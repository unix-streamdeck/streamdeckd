package streamdeckd

import (
	"github.com/unix-streamdeck/api/v2"
	streamdeck "github.com/unix-streamdeck/driver"
)

type InputManager struct {
	vdev *VirtualDev
}

func (im *InputManager) HandleKeyInput(key *api.KeyV3, keyDown bool) {
	keyConfig, ok := key.Application[key.ActiveApplication]
	if !ok {
		im.vdev.logger.Println("Err getting correct config for key")
		return
	}
	if keyDown {
		im.HandleStandardActions(keyConfig)

		if keyConfig.KeyHandler != "" {
			var deckInfo api.StreamDeckInfoV1
			deckInfo = im.vdev.sdInfo
			if keyConfig.KeyHandlerStruct == nil {
				var handler api.KeyHandler
				modules := AvailableModules()
				for _, module := range modules {
					if module.Name == keyConfig.KeyHandler {
						handler = module.NewKey()
					}
				}
				if handler == nil {
					im.vdev.logger.Println("Could not find handler:", keyConfig.KeyHandler)
					return
				}
				keyConfig.KeyHandlerStruct = handler
			}
			trimmedKeyConfig := api.KeyConfigV3{KeyHandlerFields: keyConfig.KeyHandlerFields}
			if keyConfig.IconHandler == keyConfig.KeyHandler {
				trimmedKeyConfig.SharedState = keyConfig.SharedState
				trimmedKeyConfig.KeyHandlerFields = mergeSharedConfig(keyConfig.SharedHandlerFields, keyConfig.KeyHandlerFields)
			} else {
				trimmedKeyConfig.SharedState = make(map[string]any)
			}
			keyConfig.KeyHandlerStruct.Key(trimmedKeyConfig, deckInfo)
		}
	}
	if keyConfig.KeyHold != 0 {
		if keyDown {
			err := kb.KeyDown(keyConfig.KeyHold)
			if err != nil {
				im.vdev.logger.Println(err)
			}
		} else {
			err := kb.KeyUp(keyConfig.KeyHold)
			if err != nil {
				im.vdev.logger.Println(err)
			}
		}
	}
}

func (im *InputManager) HandleKnobInput(knob *api.KnobV3, event streamdeck.InputEvent) {
	knobConfig, ok := knob.Application[knob.ActiveApplication]
	if !ok {
		im.vdev.logger.Println("Err getting correct config for knob")
		return
	}
	if knobConfig.KnobOrTouchHandler != "" {
		var deckInfo api.StreamDeckInfoV1
		deckInfo = im.vdev.sdInfo
		if knobConfig.KnobOrTouchHandlerStruct == nil {
			var handler api.KnobOrTouchHandler
			modules := AvailableModules()
			for _, module := range modules {
				if module.Name == knobConfig.KnobOrTouchHandler {
					handler = module.NewKnobOrTouch()
				}
			}
			if handler == nil {
				im.vdev.logger.Println("Could not find handler:", knobConfig.KnobOrTouchHandler)
				return
			}
			knobConfig.KnobOrTouchHandlerStruct = handler
		}
		trimmedKnobConfig := api.KnobConfigV3{KnobOrTouchHandlerFields: knobConfig.KnobOrTouchHandlerFields}
		if knobConfig.LcdHandler == knobConfig.KnobOrTouchHandler {
			trimmedKnobConfig.SharedState = knobConfig.SharedState
			trimmedKnobConfig.KnobOrTouchHandlerFields = mergeSharedConfig(knobConfig.SharedHandlerFields, knobConfig.KnobOrTouchHandlerFields)
		} else {
			trimmedKnobConfig.SharedState = make(map[string]any)
		}
		knobConfig.KnobOrTouchHandlerStruct.Input(trimmedKnobConfig, deckInfo, api.InputEvent{
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
	im.HandleStandardActions(&actions)
}

func (im *InputManager) HandleStandardActions(ia api.InputActions) {
	if ia.GetCommand() != "" {
		RunCommand(ia.GetCommand())
	}
	if ia.GetKeyBind() != "" {
		err := ExecuteKeybind(ia.GetKeyBind())
		if err != nil {
			im.vdev.logger.Println("[ERROR] Failed to execute keybind:", err)
		}
	}
	if ia.GetSwitchPage() != 0 {
		page := ia.GetSwitchPage() - 1
		im.vdev.pageManager.SetPage(page)
	}
	if ia.GetBrightness() != 0 {
		err := im.vdev.SetBrightness(uint8(ia.GetBrightness()))
		if err != nil {
			im.vdev.logger.Println(err)
		}
	}
	if ia.GetUrl() != "" {
		RunCommand("xdg-open " + ia.GetUrl())
	}
	if ia.GetObsCommand() != "" {
		runObsCommand(ia.GetObsCommand(), ia.GetObsCommandParams())
	}
}
