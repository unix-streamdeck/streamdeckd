package streamdeckd

import (
	"github.com/unix-streamdeck/api/v2"
	streamdeck "github.com/unix-streamdeck/driver"
	"image"
	"image/draw"
	"log"
	"os"
	"os/exec"
	"syscall"
)

var currentApplication = ""
var locked = false

func LoadImage(dev *VirtualDev, path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return api.ResizeImage(img, int(dev.Deck.Pixels)), nil
}

func SetKey(dev *VirtualDev, currentKeyConfig *api.KeyConfigV3, keyIndex int, page int, activeApp string) {
	defer HandlePanic(func() {
		log.Println("Restarting SetKey")
		go SetKey(dev, currentKeyConfig, keyIndex, page, activeApp)
	})
	if currentKeyConfig.IconHandler == "" {
		SetKeyImageHandlerless(dev, currentKeyConfig, keyIndex, page)
	} else {
		SetKeyImageHandler(dev, currentKeyConfig, keyIndex, page, activeApp)
	}
}
func SetKeyImageHandler(dev *VirtualDev, currentKeyConfig *api.KeyConfigV3, keyIndex int, page int, activeApp string) {
	if currentKeyConfig.IconHandlerStruct == nil {
		var handler api.IconHandler
		modules := AvailableModules()
		for _, module := range modules {
			if module.Name == currentKeyConfig.IconHandler {
				handler = module.NewIcon()
			}
		}
		if handler == nil {
			return
		}
		log.Printf("Created %s\n", currentKeyConfig.IconHandler)
		currentKeyConfig.IconHandlerStruct = handler
	}
	log.Printf("Started %s on key %d with app profile %s\n", currentKeyConfig.IconHandler, keyIndex, activeApp)
	trimmedKeyConfig := api.KeyConfigV3{IconHandlerFields: currentKeyConfig.IconHandlerFields}
	if currentKeyConfig.IconHandler == currentKeyConfig.KeyHandler {
		if currentKeyConfig.SharedState == nil {
			currentKeyConfig.SharedState = make(map[string]any)
		}
		trimmedKeyConfig.SharedState = currentKeyConfig.SharedState
		trimmedKeyConfig.IconHandlerFields = mergeSharedConfig(currentKeyConfig.SharedHandlerFields, currentKeyConfig.IconHandlerFields)
	} else {
		trimmedKeyConfig.SharedState = make(map[string]any)
	}
	currentKeyConfig.IconHandlerStruct.Start(trimmedKeyConfig, dev.sdInfo, func(image image.Image) {
		if image.Bounds().Max.X != int(dev.Deck.Pixels) || image.Bounds().Max.Y != int(dev.Deck.Pixels) {
			image = api.ResizeImage(image, int(dev.Deck.Pixels))
		}
		dev.SetImage(image, keyIndex, page)
		currentKeyConfig.Buff = image
	})
}

func SetKeyImageHandlerless(dev *VirtualDev, currentKeyConfig *api.KeyConfigV3, keyIndex int, page int) {
	if currentKeyConfig.Buff == nil {
		if currentKeyConfig.Icon == "" {
			img := image.NewRGBA(image.Rect(0, 0, int(dev.Deck.Pixels), int(dev.Deck.Pixels)))
			draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)
			currentKeyConfig.Buff = img
		} else {
			img, err := LoadImage(dev, currentKeyConfig.Icon)
			if err != nil {
				log.Println(err)
				return
			}
			currentKeyConfig.Buff = img
		}
		if currentKeyConfig.Text != "" {
			img, err := api.DrawText(currentKeyConfig.Buff, currentKeyConfig.Text, currentKeyConfig.TextSize, currentKeyConfig.TextAlignment)
			if err != nil {
				log.Println(err)
			} else {
				currentKeyConfig.Buff = img
			}
		}
	}
	dev.SetImage(currentKeyConfig.Buff, keyIndex, page)
}

func SetKnob(dev *VirtualDev, currentKnobConfig *api.KnobConfigV3, knobIndex int, page int, activeApp string) {
	defer HandlePanic(func() {
		log.Println("Restarting SetKnob")
		go SetKnob(dev, currentKnobConfig, knobIndex, page, activeApp)
	})
	if currentKnobConfig.LcdHandler == "" {
		SetKnobHandlerless(dev, currentKnobConfig, knobIndex, page)
	} else {
		SetKnobHandler(dev, currentKnobConfig, knobIndex, page, activeApp)
	}
}

func SetKnobHandlerless(dev *VirtualDev, currentKnobConfig *api.KnobConfigV3, knobIndex int, page int) {
	if currentKnobConfig.Buff == nil {
		if currentKnobConfig.Icon == "" {
			img := image.NewRGBA(image.Rect(0, 0, 200, 100))
			draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)
			currentKnobConfig.Buff = img
		} else {
			img, err := LoadImage(dev, currentKnobConfig.Icon)
			if err != nil {
				log.Println(err)
				return
			}
			currentKnobConfig.Buff = img
		}
		if currentKnobConfig.Text != "" {
			img, err := api.DrawText(currentKnobConfig.Buff, currentKnobConfig.Text, currentKnobConfig.TextSize, currentKnobConfig.TextAlignment)
			if err != nil {
				log.Println(err)
			} else {
				currentKnobConfig.Buff = img
			}
		}
		currentKnobConfig.Buff = api.ResizeImageWH(currentKnobConfig.Buff, 200, 100)
	}
	dev.SetPanelImage(currentKnobConfig.Buff, knobIndex, page)
}

func SetKnobHandler(dev *VirtualDev, currentKnobConfig *api.KnobConfigV3, knobIndex int, page int, activeApp string) {
	if currentKnobConfig.LcdHandlerStruct == nil {
		var handler api.LcdHandler
		modules := AvailableModules()
		for _, module := range modules {
			if module.Name == currentKnobConfig.LcdHandler {
				handler = module.NewLcd()
			}
		}
		if handler == nil {
			return
		}
		log.Printf("Created %s\n", currentKnobConfig.LcdHandler)
		currentKnobConfig.LcdHandlerStruct = handler
	}
	log.Printf("Started %s on knob %d with app profile %s\n", currentKnobConfig.LcdHandler, knobIndex, activeApp)
	trimmedKnobConfig := api.KnobConfigV3{LcdHandlerFields: currentKnobConfig.LcdHandlerFields}
	if currentKnobConfig.LcdHandler == currentKnobConfig.KnobOrTouchHandler {
		if currentKnobConfig.SharedState == nil {
			currentKnobConfig.SharedState = make(map[string]any)
		}
		trimmedKnobConfig.SharedState = currentKnobConfig.SharedState
		trimmedKnobConfig.LcdHandlerFields = mergeSharedConfig(currentKnobConfig.SharedHandlerFields, currentKnobConfig.LcdHandlerFields)
	} else {
                trimmedKnobConfig.SharedState = make(map[string]any)
        }
	currentKnobConfig.LcdHandlerStruct.Start(trimmedKnobConfig, dev.sdInfo, func(image image.Image) {
		if image.Bounds().Max.X != int(dev.Deck.LcdWidth) || image.Bounds().Max.Y != int(dev.Deck.LcdHeight) {
			image = api.ResizeImageWH(image, int(dev.Deck.LcdWidth), int(dev.Deck.LcdHeight))
		}
		dev.SetPanelImage(image, knobIndex, page)
		currentKnobConfig.Buff = image
	})
}

func RunCommand(command string) {
	go func() {
		cmd := exec.Command("/bin/sh", "-c", "/usr/bin/nohup "+command)

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid:   true,
			Pgid:      0,
			Pdeathsig: syscall.SIGHUP,
		}
		if err := cmd.Start(); err != nil {
			log.Println("There was a problem running ", command, ":", err)
		} else {
			pid := cmd.Process.Pid
			err := cmd.Process.Release()
			if err != nil {
				log.Println(err)
			}
			log.Println(command, " has been started with pid", pid)
		}
	}()
}

func HandleKeyInput(dev *VirtualDev, key *api.KeyV3, keyDown bool) {
	log.Println(key.ActiveApplication)
	keyConfig, ok := key.Application[key.ActiveApplication]
	if !ok {
		log.Println("Err getting correct config for key")
		return
	}
	if keyDown {
		if keyConfig.Command != "" {
			RunCommand(keyConfig.Command)
		}
		if keyConfig.Keybind != "" {
			RunCommand("xdotool key " + keyConfig.Keybind)
		}
		if keyConfig.SwitchPage != 0 {
			page := keyConfig.SwitchPage - 1
			dev.SetPage(page)
		}
		if keyConfig.Brightness != 0 {
			err := dev.SetBrightness(uint8(keyConfig.Brightness))
			if err != nil {
				log.Println(err)
			}
		}
		if keyConfig.Url != "" {
			RunCommand("xdg-open " + keyConfig.Url)
		}
		if keyConfig.ObsCommand != "" {
			runObsCommand(keyConfig.ObsCommand, keyConfig.ObsCommandParams)
		}
		if keyConfig.KeyHandler != "" {
			var deckInfo api.StreamDeckInfoV1
			deckInfo = dev.sdInfo
			if keyConfig.KeyHandlerStruct == nil {
				var handler api.KeyHandler
				modules := AvailableModules()
				for _, module := range modules {
					if module.Name == keyConfig.KeyHandler {
						handler = module.NewKey()
					}
				}
				if handler == nil {
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
				log.Println(err)
			}
		} else {
			err := kb.KeyUp(keyConfig.KeyHold)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func HandleKnobInput(dev *VirtualDev, knob *api.KnobV3, event streamdeck.InputEvent) {
	knobConfig, ok := knob.Application[knob.ActiveApplication]
	if !ok {
		log.Println("Err getting correct config for knob")
		return
	}
	if knobConfig.KnobOrTouchHandler != "" {
		var deckInfo api.StreamDeckInfoV1
		deckInfo = dev.sdInfo
		if knobConfig.KnobOrTouchHandlerStruct == nil {
			var handler api.KnobOrTouchHandler
			modules := AvailableModules()
			for _, module := range modules {
				if module.Name == knobConfig.KnobOrTouchHandler {
					handler = module.NewKnobOrTouch()
				}
			}
			if handler == nil {
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
	} else {
		var actions api.KnobActionV3
		if event.EventType == streamdeck.KNOB_PRESS {
			actions = knobConfig.KnobPressAction
		} else if event.EventType == streamdeck.KNOB_CCW {
			actions = knobConfig.KnobTurnDownAction
		} else if event.EventType == streamdeck.KNOB_CW {
			actions = knobConfig.KnobTurnUpAction
		}
		if actions.Command != "" {
			RunCommand(actions.Command)
		}
		if actions.Keybind != "" {
			RunCommand("xdotool key " + actions.Keybind)
		}
		if actions.SwitchPage != 0 {
			page := actions.SwitchPage - 1
			dev.SetPage(page)
		}
		if actions.Brightness != 0 {
			err := dev.SetBrightness(uint8(actions.Brightness))
			if err != nil {
				log.Println(err)
			}
		}
		if actions.Url != "" {
			RunCommand("xdg-open " + actions.Url)
		}
		if actions.ObsCommand != "" {
			runObsCommand(actions.ObsCommand, actions.ObsCommandParams)
		}
	}
}

func HandlePanic(cback func()) {
	if err := recover(); err != nil {
		log.Println("panic occurred:", err)
		cback()
	}
}

func mergeSharedConfig(sharedConfig map[string]any, individualConfig map[string]any) map[string]any {
	merged := make(map[string]any)

	// copy map1 into merged
	for k, v := range sharedConfig {
		merged[k] = v
	}

	// copy map2 into merged (overwrites if key already exists)
	for k, v := range individualConfig {
		merged[k] = v
	}

	return merged
}
