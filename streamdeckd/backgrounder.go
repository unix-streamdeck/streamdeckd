package streamdeckd

import (
	"image"
	"log"

	"github.com/unix-streamdeck/api/v3"
)

type IBackgrounder interface {
	SetLcdBackground(backgrounder api.LcdBackgrounder)
	SetKeyBackground(backgrounder api.KeyGridBackgrounder)
	SetIndividualLcdBackground(backgrounder api.LcdSegmentBackgrounder, index int)
	SetIndividualKeyBackground(backgrounder api.KeyBackgrounder, index int)
	AttachPageChangeListener()
}

type Backgrounder struct {
	vdev IVirtualDev
}

func (bg *Backgrounder) SetLcdBackground(backgrounder api.LcdBackgrounder) {
	if backgrounder.GetTouchPanelBackground() == "" {
		return
	}

	if backgrounder.GetTouchPanelBackgroundHandler() == nil {
		var handler api.BackgroundHandler

		module, err := GetModule(backgrounder.GetTouchPanelBackground(), BACKGROUND)
		if err != nil {
			log.Println(err)
			return
		}

		handler = module.NewBackground()

		backgrounder.SetTouchPanelBackgroundHandler(handler)
	}

	if backgrounder.GetTouchPanelBackgroundHandlerFields() != nil {
		go backgrounder.GetTouchPanelBackgroundHandler().StartGrid(backgrounder.GetTouchPanelBackgroundHandlerFields(), api.LCD, *bg.vdev.SdInfo(), func(imgs []image.Image) {
			if len(imgs) == bg.vdev.SdInfo().KnobCols {
				backgrounder.SetTouchPanelBackgroundBuff(imgs)

				for u := range imgs {
					bg.vdev.SetPanelBackground(u, bg.vdev.PageManager().GetPage())
				}
			}
		})
		return
	}
}

func (bg *Backgrounder) SetKeyBackground(backgrounder api.KeyGridBackgrounder) {
	if backgrounder.GetKeyGridBackground() == "" {
		return
	}

	if backgrounder.GetKeyGridBackgroundHandler() == nil {
		var handler api.BackgroundHandler

		module, err := GetModule(backgrounder.GetKeyGridBackground(), BACKGROUND)
		if err != nil {
			log.Println(err)
			return
		}

		handler = module.NewBackground()

		backgrounder.SetKeyGridBackgroundHandler(handler)
	}

	if backgrounder.GetKeyGridBackgroundHandler() != nil {
		go backgrounder.GetKeyGridBackgroundHandler().StartGrid(backgrounder.GetKeyGridBackgroundHandlerFields(), api.KEY, *bg.vdev.SdInfo(), func(imgs []image.Image) {
			if len(imgs) == bg.vdev.SdInfo().Cols*bg.vdev.SdInfo().Rows {
				backgrounder.SetKeyGridBackgroundBuff(imgs)

				for u := range imgs {
					bg.vdev.SetKeyBackground(u, bg.vdev.PageManager().GetPage())
				}
			}
		})
		return
	}
}

func (bg *Backgrounder) SetIndividualLcdBackground(backgrounder api.LcdSegmentBackgrounder, index int) {
	if backgrounder.GetTouchPanelBackground() == "" {
		return
	}

	if backgrounder.GetTouchPanelBackgroundHandler() == nil {
		var handler api.BackgroundHandler

		module, err := GetModule(backgrounder.GetTouchPanelBackground(), BACKGROUND)
		if err != nil {
			log.Println(err)
			return
		}

		handler = module.NewBackground()

		backgrounder.SetTouchPanelBackgroundHandler(handler)
	}

	if backgrounder.GetTouchPanelBackgroundHandlerFields() != nil {
		go backgrounder.GetTouchPanelBackgroundHandler().Start(backgrounder.GetTouchPanelBackgroundHandlerFields(), api.LCD, *bg.vdev.SdInfo(), func(img image.Image) {
			backgrounder.SetTouchPanelBackgroundBuff(img)

			bg.vdev.SetPanelBackground(index, bg.vdev.PageManager().GetPage())
		})
		return
	}
}

func (bg *Backgrounder) SetIndividualKeyBackground(backgrounder api.KeyBackgrounder, index int) {
	if backgrounder.GetKeyBackground() == "" {
		return
	}

	if backgrounder.GetKeyBackgroundHandler() == nil {
		var handler api.BackgroundHandler

		module, err := GetModule(backgrounder.GetKeyBackground(), BACKGROUND)
		if err != nil {
			log.Println(err)
			return
		}

		handler = module.NewBackground()

		backgrounder.SetKeyBackgroundHandler(handler)
	}
}

func (bg *Backgrounder) AttachPageChangeListener() {

	bg.vdev.PageManager().AttachListener(func(newPage, _ int) {

		currentPage := bg.vdev.Config().Pages[newPage]

		go bg.SetKeyBackground(&currentPage)

		go bg.SetLcdBackground(&currentPage)

		for i := range currentPage.Keys {
			key := &currentPage.Keys[i]

			go bg.SetIndividualKeyBackground(key, i)
		}

		for i := range currentPage.Knobs {
			knob := &currentPage.Knobs[i]

			go bg.SetIndividualLcdBackground(knob, i)
		}
	})
}
