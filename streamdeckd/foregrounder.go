package streamdeckd

import (
	"fmt"
	"image"

	"github.com/unix-streamdeck/api/v2"
)

type IForegrounder interface {
	SetKnob(currentKnobConfig *api.KnobConfigV3, knobIndex int, page int, activeApp string)
	SetKey(currentKeyConfig *api.KeyConfigV3, keyIndex int, page int, activeApp string)
	AttachPageChangeListener()
	AttachAppChangeListener()
}

type Foregrounder struct {
	vdev IVirtualDev
}

func (f *Foregrounder) SetKnob(currentKnobConfig *api.KnobConfigV3, knobIndex int, page int, activeApp string) {
	defer HandlePanic(func() {
		f.vdev.Logger().Println("Restarting SetKnob")
		go f.SetKnob(currentKnobConfig, knobIndex, page, activeApp)
	})
	if currentKnobConfig.LcdHandler != "" {
		f.setHandler(currentKnobConfig, api.LCD, knobIndex, activeApp, func(img image.Image) {
			if img.Bounds().Dx() != f.vdev.SdInfo().LcdWidth || img.Bounds().Dy() != f.vdev.SdInfo().LcdHeight {
				img = api.ResizeImageWH(img, f.vdev.SdInfo().LcdWidth, f.vdev.SdInfo().LcdHeight)
			}
			f.vdev.SetPanelForeground(img, knobIndex, page)
		})
	}
	if currentKnobConfig.LcdHandlerStruct == nil {
		img := f.loadStaticImage(currentKnobConfig, f.vdev.SdInfo().LcdWidth, f.vdev.SdInfo().LcdHeight)
		if img != nil {
			f.vdev.SetKeyForeground(img, knobIndex, page)
		}
	}
}

func (f *Foregrounder) SetKey(currentKeyConfig *api.KeyConfigV3, keyIndex int, page int, activeApp string) {
	defer HandlePanic(func() {
		f.vdev.Logger().Println("Restarting SetKey")
		go f.SetKey(currentKeyConfig, keyIndex, page, activeApp)
	})
	if currentKeyConfig.IconHandler != "" {
		f.setHandler(currentKeyConfig, api.KEY, keyIndex, activeApp, func(img image.Image) {
			if img.Bounds().Dx() != f.vdev.SdInfo().IconSize || img.Bounds().Dy() != f.vdev.SdInfo().IconSize {
				img = api.ResizeImage(img, f.vdev.SdInfo().IconSize)
			}
			f.vdev.SetKeyForeground(img, keyIndex, page)
		})
	}
	if currentKeyConfig.IconHandlerStruct == nil {
		img := f.loadStaticImage(currentKeyConfig, f.vdev.SdInfo().IconSize, f.vdev.SdInfo().IconSize)
		if img != nil {
			f.vdev.SetKeyForeground(img, keyIndex, page)
		}
	}
}

func (f *Foregrounder) setHandler(foregroundActions api.ForegroundAndInputHandlerConfig, handlerType api.HandlerType, index int, activeApp string, callback func(img image.Image)) {
	if foregroundActions.GetForegroundHandlerInstance() == nil {
		var handler api.ForegroundHandler
		modules := AvailableModules()
		for _, module := range modules {
			if module.Name == foregroundActions.GetForegroundHandler() {
				handler = module.NewForeground()
			}
		}
		if handler == nil {
			f.vdev.Logger().Println("Could not find handler:", foregroundActions.GetForegroundHandler())
			return
		}
		f.vdev.Logger().Printf("Created %s\n", foregroundActions.GetForegroundHandler())
		foregroundActions.SetForegroundHandlerInstance(handler)
	}
	f.vdev.Logger().Printf("Started %s on key %d with app profile `%s`\n", foregroundActions.GetForegroundHandler(), index, activeApp)

	fields := foregroundActions.GetForegroundHandlerFields()

	if foregroundActions.GetForegroundHandler() == foregroundActions.GetInputHandler() {
		fields = mergeSharedConfig(foregroundActions.GetSharedHandlerFields(), foregroundActions.GetForegroundHandlerFields())
	}

	foregroundActions.GetForegroundHandlerInstance().Start(fields,
		handlerType, *f.vdev.SdInfo(), callback)
}

func (f *Foregrounder) loadStaticImage(fa api.ForegroundActions, w, h int) image.Image {
	var img image.Image
	if fa.GetIcon() == "" {
		img = image.NewRGBA(image.Rect(0, 0, w, h))
	} else {
		var err error
		img, err = LoadImage(fa.GetIcon())
		if err != nil {
			f.vdev.Logger().Println(err)
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
			f.vdev.Logger().Println(err)
			return nil
		}
	}
	return img
}

func (f *Foregrounder) AttachPageChangeListener() {
	f.vdev.PageManager().AttachListener(func(newPage, _ int) {
		currentPage := f.vdev.Config().Pages[newPage]

		for i, _ := range currentPage.Keys {
			key := &currentPage.Keys[i]

			if key.Application == nil {
				key.Application = map[string]*api.KeyConfigV3{}
				key.Application[""] = &api.KeyConfigV3{}
				currentPage.Keys[i] = *key
				f.vdev.Logger().Println(fmt.Sprintf("Setting empty application on key: %d on page: %d", i, newPage))
				SaveConfig()
			}
			_, keyHasApp := key.Application[applicationManager.GetApplication()]
			if key.ActiveApplication != "" && !keyHasApp {
				key.ActiveApplication = ""
			}
			if keyHasApp {
				key.ActiveApplication = applicationManager.GetApplication()
			}
			go f.SetKey(key.Application[key.ActiveApplication], i, newPage, key.ActiveApplication)
		}
		for i, _ := range currentPage.Knobs {
			knob := &currentPage.Knobs[i]

			if knob.Application == nil {
				knob.Application = map[string]*api.KnobConfigV3{}
				knob.Application[""] = &api.KnobConfigV3{}
				currentPage.Knobs[i] = *knob
				f.vdev.Logger().Println(fmt.Sprintf("Setting empty application on knob: %d on page: %d", i, newPage))
				SaveConfig()
			}
			_, knobHasApp := knob.Application[applicationManager.GetApplication()]
			if knob.ActiveApplication != "" && !knobHasApp {
				knob.ActiveApplication = ""
			}
			if knobHasApp {
				knob.ActiveApplication = applicationManager.GetApplication()
			}
			go f.SetKnob(knob.Application[knob.ActiveApplication], i, newPage, knob.ActiveApplication)
		}
	})
}

func (f *Foregrounder) AttachAppChangeListener() {
	applicationManager.AttachListener(func(application string) {
		page := f.vdev.Config().Pages[f.vdev.PageManager().GetPage()]
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
				go f.SetKey(key.Application[key.ActiveApplication], i, f.vdev.PageManager().GetPage(), key.ActiveApplication)
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
				go f.SetKnob(knob.Application[knob.ActiveApplication], i, f.vdev.PageManager().GetPage(), knob.ActiveApplication)
			}
		}
	})
}
