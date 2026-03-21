package streamdeckd

import (
	"fmt"
	"image"

	"github.com/unix-streamdeck/api/v2"
)

type Foregrounder struct {
	vdev *VirtualDev
}

func (f *Foregrounder) SetKnob(currentKnobConfig *api.KnobConfigV3, knobIndex int, page int, activeApp string) {
	defer HandlePanic(func() {
		f.vdev.logger.Println("Restarting SetKnob")
		go f.SetKnob(currentKnobConfig, knobIndex, page, activeApp)
	})
	if currentKnobConfig.LcdHandler != "" {
		f.SetKnobHandler(currentKnobConfig, knobIndex, page, activeApp)
	}
	if currentKnobConfig.LcdHandlerStruct == nil {
		img := f.loadStaticImage(currentKnobConfig, f.vdev.sdInfo.LcdWidth, f.vdev.sdInfo.LcdHeight)
		if img != nil {
			f.vdev.SetKeyForeground(img, knobIndex, page)
		}
	}
}

func (f *Foregrounder) SetKnobHandler(currentKnobConfig *api.KnobConfigV3, knobIndex int, page int, activeApp string) {
	if currentKnobConfig.LcdHandlerStruct == nil {
		var handler api.LcdHandler
		modules := AvailableModules()
		for _, module := range modules {
			if module.Name == currentKnobConfig.LcdHandler {
				handler = module.NewLcd()
			}
		}
		if handler == nil {
			f.vdev.logger.Println("Could not find handler:", currentKnobConfig.LcdHandler)
			return
		}
		f.vdev.logger.Printf("Created %s\n", currentKnobConfig.LcdHandler)
		currentKnobConfig.LcdHandlerStruct = handler
	}
	f.vdev.logger.Printf("Started %s on knob %d with app profile %s\n", currentKnobConfig.LcdHandler, knobIndex, activeApp)
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

	go currentKnobConfig.LcdHandlerStruct.Start(trimmedKnobConfig, f.vdev.sdInfo, func(image image.Image) {
		if image.Bounds().Max.X != f.vdev.sdInfo.LcdWidth || image.Bounds().Max.Y != f.vdev.sdInfo.LcdHeight {
			image = api.ResizeImageWH(image, f.vdev.sdInfo.LcdWidth, f.vdev.sdInfo.LcdHeight)
		}
		f.vdev.SetPanelForeground(image, knobIndex, page)
	})
}

func (f *Foregrounder) SetKey(currentKeyConfig *api.KeyConfigV3, keyIndex int, page int, activeApp string) {
	defer HandlePanic(func() {
		f.vdev.logger.Println("Restarting SetKey")
		go f.SetKey(currentKeyConfig, keyIndex, page, activeApp)
	})
	if currentKeyConfig.IconHandler != "" {
		f.SetKeyHandler(currentKeyConfig, keyIndex, page, activeApp)
	}
	if currentKeyConfig.IconHandlerStruct == nil {
		img := f.loadStaticImage(currentKeyConfig, f.vdev.sdInfo.IconSize, f.vdev.sdInfo.IconSize)
		if img != nil {
			f.vdev.SetKeyForeground(img, keyIndex, page)
		}
	}
}

func (f *Foregrounder) SetKeyHandler(currentKeyConfig *api.KeyConfigV3, keyIndex int, page int, activeApp string) {
	if currentKeyConfig.IconHandlerStruct == nil {
		var handler api.IconHandler
		modules := AvailableModules()
		for _, module := range modules {
			if module.Name == currentKeyConfig.IconHandler {
				handler = module.NewIcon()
			}
		}
		if handler == nil {
			f.vdev.logger.Println("Could not find handler:", currentKeyConfig.IconHandler)
			return
		}
		f.vdev.logger.Printf("Created %s\n", currentKeyConfig.IconHandler)
		currentKeyConfig.IconHandlerStruct = handler
	}
	f.vdev.logger.Printf("Started %s on key %d with app profile `%s`\n", currentKeyConfig.IconHandler, keyIndex, activeApp)
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
	currentKeyConfig.IconHandlerStruct.Start(trimmedKeyConfig, f.vdev.sdInfo, func(image image.Image) {
		if image.Bounds().Max.X != f.vdev.sdInfo.IconSize || image.Bounds().Max.Y != f.vdev.sdInfo.IconSize {
			image = api.ResizeImage(image, f.vdev.sdInfo.IconSize)
		}
		f.vdev.SetKeyForeground(image, keyIndex, page)
	})
}

func (f *Foregrounder) loadStaticImage(fa api.ForegroundActions, w, h int) image.Image {
	var img image.Image
	if fa.GetIcon() == "" {
		img = image.NewRGBA(image.Rect(0, 0, w, h))
	} else {
		var err error
		img, err = LoadImage(fa.GetIcon())
		if err != nil {
			f.vdev.logger.Println(err)
			return nil
		}
	}
	if fa.GetText() != "" {
		var err error
		img, err = api.DrawText(img, fa.GetText(), api.DrawTextOptions{
			FontSize:          int64(fa.GetTextSize()),
			VerticalAlignment: fa.GetTextAlignment(),
			FontFace:          fa.GetFontFace(),
			Colour:            fa.GetTextColour(),
		})
		if err != nil {
			f.vdev.logger.Println(err)
			return nil
		}
	}
	return img
}

func (f *Foregrounder) AttachPageChangeListener() {
	f.vdev.pageManager.AttachListener(func(newPage, _ int) {
		currentPage := f.vdev.Config.Pages[newPage]

		for i, _ := range currentPage.Keys {
			key := &currentPage.Keys[i]

			if key.Application == nil {
				key.Application = map[string]*api.KeyConfigV3{}
				key.Application[""] = &api.KeyConfigV3{}
				currentPage.Keys[i] = *key
				f.vdev.logger.Println(fmt.Sprintf("Setting empty application on key: %d on page: %d", i, newPage))
				SaveConfig()
			}
			_, keyHasApp := key.Application[applicationManager.activeApplication]
			if key.ActiveApplication != "" && !keyHasApp {
				key.ActiveApplication = ""
			}
			if keyHasApp {
				key.ActiveApplication = applicationManager.activeApplication
			}
			go f.SetKey(key.Application[key.ActiveApplication], i, newPage, key.ActiveApplication)
		}
		for i, _ := range currentPage.Knobs {
			knob := &currentPage.Knobs[i]

			if knob.Application == nil {
				knob.Application = map[string]*api.KnobConfigV3{}
				knob.Application[""] = &api.KnobConfigV3{}
				currentPage.Knobs[i] = *knob
				f.vdev.logger.Println(fmt.Sprintf("Setting empty application on knob: %d on page: %d", i, newPage))
				SaveConfig()
			}
			_, knobHasApp := knob.Application[applicationManager.activeApplication]
			if knob.ActiveApplication != "" && !knobHasApp {
				knob.ActiveApplication = ""
			}
			if knobHasApp {
				knob.ActiveApplication = applicationManager.activeApplication
			}
			go f.SetKnob(knob.Application[knob.ActiveApplication], i, newPage, knob.ActiveApplication)
		}
	})
}

func (f *Foregrounder) AttachAppChangeListener() {
	applicationManager.AttachListener(func(application string) {
		page := f.vdev.Config.Pages[f.vdev.pageManager.page]
		for i := range page.Keys {
			key := &page.Keys[i]
			_, keyHasApp := key.Application[application]
			activeApp := key.ActiveApplication
			if key.Application[key.ActiveApplication].KeyHold != 0 && (keyHasApp || key.ActiveApplication != "") {
				kb.KeyUp(key.Application[key.ActiveApplication].KeyHold)
			}
			if key.ActiveApplication != "" && !keyHasApp {
				key.ActiveApplication = ""
			}
			if keyHasApp {
				key.ActiveApplication = application
			}
			if key.ActiveApplication != activeApp {
				go f.SetKey(key.Application[key.ActiveApplication], i, f.vdev.pageManager.page, key.ActiveApplication)
			}
		}
		for i := range page.Knobs {
			knob := &page.Knobs[i]
			_, keyHasApp := knob.Application[application]
			activeApp := knob.ActiveApplication
			if knob.ActiveApplication != "" && !keyHasApp {
				knob.ActiveApplication = ""
			}
			if keyHasApp {
				knob.ActiveApplication = application
			}
			if knob.ActiveApplication != activeApp {
				go f.SetKnob(knob.Application[knob.ActiveApplication], i, f.vdev.pageManager.page, knob.ActiveApplication)
			}
		}
	})
}
