package streamdeckd

import (
	"image"
	"log"

	"github.com/unix-streamdeck/api/v2"
)

type Backgrounder struct {
	vdev   *VirtualDev
	sdInfo api.StreamDeckInfoV1
}

func (bg *Backgrounder) setLcdBackground(backgrounder api.LcdBackgrounder) {
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
		go backgrounder.GetTouchPanelBackgroundHandler().StartGrid(backgrounder.GetTouchPanelBackgroundHandlerFields(), api.LCD, bg.sdInfo, func(imgs []image.Image) {
			if len(imgs) == bg.sdInfo.KnobCols {
				backgrounder.SetTouchPanelBackgroundBuff(imgs)

				for u := range imgs {
					bg.vdev.SetPanelBackground(u, bg.vdev.pageManager.page)
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
		bg.vdev.logger.Println(err)
		return
	}

	img = api.ResizeImageWH(img, bg.sdInfo.LcdBackgroundWidth, bg.sdInfo.LcdBackgroundHeight)

	imgs := bg.sdInfo.SplitBackgroundImage(img, api.LCD)

	backgrounder.SetTouchPanelBackgroundBuff(imgs)

	for index, _ := range imgs {
		bg.vdev.SetPanelBackground(index, bg.vdev.pageManager.page)
	}
}

func (bg *Backgrounder) setKeyBackground(backgrounder api.KeyGridBackgrounder) {
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
		go backgrounder.GetKeyGridBackgroundHandler().StartGrid(backgrounder.GetKeyGridBackgroundHandlerFields(), api.KEY, bg.sdInfo, func(imgs []image.Image) {
			if len(imgs) == bg.sdInfo.Cols*bg.sdInfo.Rows {
				backgrounder.SetKeyGridBackgroundBuff(imgs)

				for u := range imgs {
					bg.vdev.SetKeyBackground(u, bg.vdev.pageManager.page)
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
		bg.vdev.logger.Println(err)
		return
	}

	img = api.ResizeImageWH(img, bg.sdInfo.KeyGridBackgroundWidth, bg.sdInfo.KeyGridBackgroundHeight)

	imgs := bg.sdInfo.SplitBackgroundImage(img, api.KEY)

	backgrounder.SetKeyGridBackgroundBuff(imgs)

	for index, _ := range imgs {
		bg.vdev.SetKeyBackground(index, bg.vdev.pageManager.page)
	}
}

func (bg *Backgrounder) setIndividualLcdBackground(backgrounder api.LcdSegmentBackgrounder, index int) {
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
		go backgrounder.GetTouchPanelBackgroundHandler().Start(backgrounder.GetTouchPanelBackgroundHandlerFields(), api.LCD, bg.sdInfo, func(img image.Image) {
			backgrounder.SetTouchPanelBackgroundBuff(img)

			bg.vdev.SetPanelBackground(index, bg.vdev.pageManager.page)
		})
		return
	}

	if backgrounder.GetTouchPanelBackgroundBuff() != nil {
		return
	}

	img, err := LoadImage(backgrounder.GetTouchPanelBackground())
	if err != nil {
		bg.vdev.logger.Println(err)
		return
	}

	img = api.ResizeImageWH(img, bg.sdInfo.LcdWidth, bg.sdInfo.LcdHeight)

	backgrounder.SetTouchPanelBackgroundBuff(img)

	bg.vdev.SetPanelBackground(index, bg.vdev.pageManager.page)
}

func (bg *Backgrounder) setIndividualKeyBackground(backgrounder api.KeyBackgrounder, index int) {
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
		go backgrounder.GetKeyBackgroundHandler().Start(backgrounder.GetKeyBackgroundHandlerFields(), api.KEY, bg.sdInfo, func(img image.Image) {
			backgrounder.SetKeyBackgroundBuff(img)

			bg.vdev.SetKeyBackground(index, bg.vdev.pageManager.page)
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

	img = api.ResizeImage(img, bg.sdInfo.IconSize)

	bg.vdev.SetKeyBackground(index, bg.vdev.pageManager.page)
}

func (bg *Backgrounder) AttachPageChangeListener() {

	bg.vdev.pageManager.AttachListener(func(newPage, _ int) {

		currentPage := bg.vdev.Config.Pages[newPage]

		go bg.setKeyBackground(&currentPage)

		go bg.setLcdBackground(&currentPage)

		for i := range currentPage.Keys {
			key := &currentPage.Keys[i]

			go bg.setIndividualKeyBackground(key, i)
		}

		for i := range currentPage.Knobs {
			knob := &currentPage.Knobs[i]

			go bg.setIndividualLcdBackground(knob, i)
		}
	})
}
