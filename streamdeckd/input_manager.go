package streamdeckd

import (
	"github.com/unix-streamdeck/api/v2"
	streamdeck "github.com/unix-streamdeck/driver"
)

type IInputManager interface {
	HandleKeyInput(key *api.KeyV3, event streamdeck.InputEvent)
	HandleKnobInput(knob *api.KnobV3, event streamdeck.InputEvent)
	GetKeyState(index int) bool
}

type InputManager struct {
	vdev      IVirtualDev
	KeyStates []bool
}

func (im *InputManager) HandleKeyInput(key *api.KeyV3, event streamdeck.InputEvent) {
	keyConfig, ok := key.Application[key.ActiveApplication]
	if !ok {
		im.vdev.Logger().Println("Err getting correct config for key")
		return
	}
	if event.EventType == streamdeck.KEY_PRESS {
		im.KeyStates[event.Index] = true
		im.vdev.RedrawKey(int(event.Index))

		im.handleStandardActions(keyConfig)

		im.handleHandlerAction(keyConfig, api.KEY, event)

	} else {
		im.KeyStates[event.Index] = false
		im.vdev.RedrawKey(int(event.Index))
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
	im.handleHandlerAction(knobConfig, api.LCD, event)
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

func (im *InputManager) handleHandlerAction(foregroundActions api.ForegroundAndInputHandlerConfig, handlerType api.HandlerType, event streamdeck.InputEvent) {
	if foregroundActions.GetInputHandler() != "" {
		var deckInfo api.StreamDeckInfoV1
		deckInfo = *im.vdev.SdInfo()
		if foregroundActions.GetInputHandlerInstance() == nil {
			var handler api.InputHandler
			modules := AvailableModules()
			for _, module := range modules {
				if module.Name == foregroundActions.GetInputHandler() {
					handler = module.NewInput()
				}
			}
			if handler == nil {
				im.vdev.Logger().Println("Could not find handler:", foregroundActions.GetInputHandler())
				return
			}
			foregroundActions.SetInputHandlerInstance(handler)
		}
		fields := foregroundActions.GetInputHandlerFields()

		if foregroundActions.GetForegroundHandler() == foregroundActions.GetInputHandler() {
			fields = mergeSharedConfig(foregroundActions.GetSharedHandlerFields(), foregroundActions.GetInputHandlerFields())
			foregroundHandler := foregroundActions.GetForegroundHandlerInstance()
			comboHandler, ok := foregroundHandler.(api.CombinedHandler)
			if ok {
				foregroundActions.SetInputHandlerInstance(comboHandler)
			}
		}
		foregroundActions.GetInputHandlerInstance().Input(fields, handlerType, deckInfo, api.InputEvent{
			EventType:     api.InputEventType(event.EventType),
			RotateNotches: event.RotateNotches,
		})
	}
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

func (im *InputManager) GetKeyState(index int) bool {
	return im.KeyStates[index]
}
