package streamdeckd

import (
	"image"
	"log"

	"github.com/unix-streamdeck/api/v2"
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

		for _, module := range modules {
			if module.Name == backgrounder.GetTouchPanelBackground() {
				handler = module.NewBackground()
			}
		}

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

	if backgrounder.GetTouchPanelBackgroundBuff() != nil {
		return
	}

	img, err := LoadImage(backgrounder.GetTouchPanelBackground())
	if err != nil {
		bg.vdev.Logger().Println(err)
		return
	}

	img = api.ResizeImageWH(img, bg.vdev.SdInfo().LcdBackgroundWidth, bg.vdev.SdInfo().LcdBackgroundHeight)

	imgs := bg.vdev.SdInfo().SplitBackgroundImage(img, api.LCD)

	backgrounder.SetTouchPanelBackgroundBuff(imgs)

	for index, _ := range imgs {
		bg.vdev.SetPanelBackground(index, bg.vdev.PageManager().GetPage())
	}
}

func (bg *Backgrounder) SetKeyBackground(backgrounder api.KeyGridBackgrounder) {
	if backgrounder.GetKeyGridBackground() == "" {
		return
	}

	if backgrounder.GetKeyGridBackgroundHandler() == nil {
		var handler api.BackgroundHandler

		for _, module := range modules {
			if module.Name == backgrounder.GetKeyGridBackground() {
				handler = module.NewBackground()
			}
		}

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

	if backgrounder.GetKeyGridBackgroundBuff() != nil {
		return
	}

	img, err := LoadImage(backgrounder.GetKeyGridBackground())
	if err != nil {
		bg.vdev.Logger().Println(err)
		return
	}

	img = api.ResizeImageWH(img, bg.vdev.SdInfo().KeyGridBackgroundWidth, bg.vdev.SdInfo().KeyGridBackgroundHeight)

	imgs := bg.vdev.SdInfo().SplitBackgroundImage(img, api.KEY)

	backgrounder.SetKeyGridBackgroundBuff(imgs)

	for index, _ := range imgs {
		bg.vdev.SetKeyBackground(index, bg.vdev.PageManager().GetPage())
	}
}

func (bg *Backgrounder) SetIndividualLcdBackground(backgrounder api.LcdSegmentBackgrounder, index int) {
	if backgrounder.GetTouchPanelBackground() == "" {
		return
	}

	if backgrounder.GetTouchPanelBackgroundHandler() == nil {
		var handler api.BackgroundHandler

		for _, module := range modules {
			if module.Name == backgrounder.GetTouchPanelBackground() {
				handler = module.NewBackground()
			}
		}

		backgrounder.SetTouchPanelBackgroundHandler(handler)
	}

	if backgrounder.GetTouchPanelBackgroundHandlerFields() != nil {
		go backgrounder.GetTouchPanelBackgroundHandler().Start(backgrounder.GetTouchPanelBackgroundHandlerFields(), api.LCD, *bg.vdev.SdInfo(), func(img image.Image) {
			backgrounder.SetTouchPanelBackgroundBuff(img)

			bg.vdev.SetPanelBackground(index, bg.vdev.PageManager().GetPage())
		})
		return
	}

	if backgrounder.GetTouchPanelBackgroundBuff() != nil {
		return
	}

	img, err := LoadImage(backgrounder.GetTouchPanelBackground())
	if err != nil {
		bg.vdev.Logger().Println(err)
		return
	}

	img = api.ResizeImageWH(img, bg.vdev.SdInfo().LcdWidth, bg.vdev.SdInfo().LcdHeight)

	backgrounder.SetTouchPanelBackgroundBuff(img)

	bg.vdev.SetPanelBackground(index, bg.vdev.PageManager().GetPage())
}

func (bg *Backgrounder) SetIndividualKeyBackground(backgrounder api.KeyBackgrounder, index int) {
	if backgrounder.GetKeyBackground() == "" {
		return
	}

	if backgrounder.GetKeyBackgroundHandler() == nil {
		var handler api.BackgroundHandler

		for _, module := range modules {
			if module.Name == backgrounder.GetKeyBackground() {
				handler = module.NewBackground()
			}
		}

		backgrounder.SetKeyBackgroundHandler(handler)
	}

	if backgrounder.GetKeyBackgroundHandler() != nil {
		go backgrounder.GetKeyBackgroundHandler().Start(backgrounder.GetKeyBackgroundHandlerFields(), api.KEY, *bg.vdev.SdInfo(), func(img image.Image) {
			backgrounder.SetKeyBackgroundBuff(img)

			bg.vdev.SetKeyBackground(index, bg.vdev.PageManager().GetPage())
		})
		return
	}

	if backgrounder.GetKeyBackgroundBuff() != nil {
		return
	}

	img, err := LoadImage(backgrounder.GetKeyBackground())
	if err != nil {
		log.Println(err)
		return
	}

	img = api.ResizeImage(img, bg.vdev.SdInfo().IconSize)

	bg.vdev.SetKeyBackground(index, bg.vdev.PageManager().GetPage())
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
