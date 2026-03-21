package streamdeckd

import (
	"image"
	"log"
	"math"

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

	var imgs []image.Image

	for lcdIndex := range bg.sdInfo.LcdCols {
		x0, y0 := bg.sdInfo.LcdWidth*lcdIndex, 0
		x1, y1 := bg.sdInfo.LcdWidth*(lcdIndex+1), bg.sdInfo.LcdHeight

		imgs = append(imgs, api.SubImage(img, x0, y0, x1, y1))
	}

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

	var imgs []image.Image
	for keyIndex := range bg.sdInfo.Cols * bg.sdInfo.Rows {
		keyX := keyIndex % bg.sdInfo.Cols
		keyY := int(math.Floor(float64(keyIndex) / float64(bg.sdInfo.Cols)))

		x0, y0 := keyX*(bg.sdInfo.IconSize+bg.sdInfo.PaddingX), keyY*(bg.sdInfo.IconSize+bg.sdInfo.PaddingY)
		x1, y1 := keyX*(bg.sdInfo.IconSize+bg.sdInfo.PaddingX)+bg.sdInfo.IconSize, keyY*(bg.sdInfo.IconSize+bg.sdInfo.PaddingY)+bg.sdInfo.IconSize

		imgs = append(imgs, api.SubImage(img, x0, y0, x1, y1))
	}
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
