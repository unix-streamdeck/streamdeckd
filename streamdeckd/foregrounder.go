package streamdeckd

import (
	"fmt"
	"image"
	"log"

	"github.com/unix-streamdeck/api/v3"
)

type IForegrounder interface {
	SetKnob(knob *api.Knob, knobIndex int, page int, activeApp string)
	SetKey(key *api.Key, keyIndex int, page int, activeApp string)
	AttachPageChangeListener()
	AttachAppChangeListener()
}

type Foregrounder struct {
	vdev IVirtualDev
}

func (f *Foregrounder) SetKnob(knob *api.Knob, knobIndex int, page int, activeApp string) {
	defer HandlePanic(func() {
		f.vdev.Logger().Println("Restarting SetKnob")
		go f.SetKnob(knob, knobIndex, page, activeApp)
	})
	for _, layer := range MergeLayers(knob.Application[knob.ActiveApplication], knob) {
		f.setHandler(layer, api.LCD, knobIndex, activeApp, func(img image.Image) {
			if img.Bounds().Dx() != f.vdev.SdInfo().LcdWidth || img.Bounds().Dy() != f.vdev.SdInfo().LcdHeight {
				img = api.ResizeImageWH(img, f.vdev.SdInfo().LcdWidth, f.vdev.SdInfo().LcdHeight)
			}
			layer.Buff = img
			f.vdev.SetPanelForeground(knobIndex, page)
		})
	}
}

func (f *Foregrounder) SetKey(key *api.Key, keyIndex int, page int, activeApp string) {
	defer HandlePanic(func() {
		f.vdev.Logger().Println("Restarting SetKey")
		go f.SetKey(key, keyIndex, page, activeApp)
	})
	for _, layer := range MergeLayers(key.Application[key.ActiveApplication], key) {
		f.setHandler(layer, api.KEY, keyIndex, activeApp, func(img image.Image) {
			if img.Bounds().Dx() != f.vdev.SdInfo().IconSize || img.Bounds().Dy() != f.vdev.SdInfo().IconSize {
				img = api.ResizeImage(img, f.vdev.SdInfo().IconSize)
			}
			layer.Buff = img
			f.vdev.SetKeyForeground(keyIndex, page)
		})
	}
}

func (f *Foregrounder) setHandler(layer *api.Layer, handlerType api.HandlerType, index int, activeApp string, callback func(img image.Image)) {
	if !layer.Visual {
		return
	}
	if layer.HandlerStruct == nil {
		var handler api.Handler
		modules := AvailableModules()
		for _, module := range modules {
			if module.Name == layer.Handler {
				handler = module.NewCombined()
			}
		}
		if handler == nil {
			f.vdev.Logger().Println("Could not find handler:", layer.Handler)
			return
		}
		f.vdev.Logger().Printf("Created %s\n", layer.Handler)
		layer.HandlerStruct = handler
	}
	f.vdev.Logger().Printf("Started %s on key %d with app profile `%s`\n", layer.Handler, index, activeApp)
	fgHandler, ok := layer.HandlerStruct.(api.ForegroundHandler)
	if !ok {
		log.Println("[Warn] Handler", layer.Handler, "is not valid ForegroundHandler")
		return
	}
	if !fgHandler.IsRunning() {
		fgHandler.Start(layer.Fields, handlerType, *f.vdev.SdInfo(), callback)
	}
}

func (f *Foregrounder) AttachPageChangeListener() {
	f.vdev.PageManager().AttachListener(func(newPage, _ int) {
		currentPage := f.vdev.Config().Pages[newPage]

		for i, _ := range currentPage.Keys {
			key := &currentPage.Keys[i]

			if key.Application == nil {
				key.Application = map[string]*api.KeyConfig{}
				key.Application[""] = &api.KeyConfig{}
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
			go f.SetKey(key, i, newPage, key.ActiveApplication)
		}
		for i, _ := range currentPage.Knobs {
			knob := &currentPage.Knobs[i]

			if knob.Application == nil {
				knob.Application = map[string]*api.KnobConfig{}
				knob.Application[""] = &api.KnobConfig{}
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
			go f.SetKnob(knob, i, newPage, knob.ActiveApplication)
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
			if key.ActiveApplication != "" && !keyHasApp {
				key.ActiveApplication = ""
			}
			if keyHasApp {
				key.ActiveApplication = application
			}
			if key.ActiveApplication != activeApp {
				go f.SetKey(key, i, f.vdev.PageManager().GetPage(), key.ActiveApplication)
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
				go f.SetKnob(knob, i, f.vdev.PageManager().GetPage(), knob.ActiveApplication)
			}
		}
	})
}
