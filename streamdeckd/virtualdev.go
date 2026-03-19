package streamdeckd

import (
	"errors"
	"fmt"
	"image"
	"log"
	"math"
	"regexp"
	"sync"
	"time"

	"github.com/unix-streamdeck/api/v2"
	streamdeck "github.com/unix-streamdeck/driver"
)

var disconnectSem sync.Mutex
var connectSem sync.Mutex
var Devs map[string]*VirtualDev

func OpenDevice() error {
	connectSem.Lock()
	defer connectSem.Unlock()
	rawDevs, err := streamdeck.Devices()
	if err != nil {
		return err
	}
	if len(rawDevs) == 0 {
		return errors.New("No streamdeck devices found")
	}
	for _, rawDev := range rawDevs {
		if len(rawDev.Serial) != 12 {
			continue
		}
		dev, ok := Devs[rawDev.Serial]
		if ok && dev.IsOpen {
			continue
		}
		err = rawDev.Open()
		if err != nil {
			log.Println(err)
			continue
		}
		if !ok {
			// initial connect
			config := findConfig(rawDev)
			dev = &VirtualDev{
				Deck:           rawDev,
				Page:           0,
				IsOpen:         true,
				Config:         config,
				keyUpdateChan:  make(chan int),
				knobUpdateChan: make(chan int),
				KeyBGBuffs:     make([]image.Image, rawDev.Keys),
				KeyFGBuffs:     make([]image.Image, rawDev.Keys),
				PanelBGBuffs:   make([]image.Image, rawDev.LcdColumns),
				PanelFGBuffs:   make([]image.Image, rawDev.LcdColumns),
			}
			dev.SetSdInfo()
			Devs[rawDev.Serial] = dev
		} else {
			//reconnect
			dev.IsOpen = true
			dev.Deck = rawDev
			dev.sdInfo.LastConnected = time.Now()
			dev.sdInfo.Connected = true
		}

		w := (dev.Deck.Pixels * uint(dev.Deck.Columns)) + (dev.Deck.PaddingX * uint(dev.Deck.Columns-1))
		h := (dev.Deck.Pixels * uint(dev.Deck.Rows)) + (dev.Deck.PaddingX * uint(dev.Deck.Rows-1))
		go dev.setKeyBackground(&dev.Config, int(w), int(h))

		go dev.setLcdBackground(&dev.Config, dev.sdInfo.LcdWidth*dev.sdInfo.LcdCols, dev.sdInfo.LcdHeight)

		go dev.HandleInput()
		dev.Render()
		dev.SetPage(dev.Page)
		log.Println(fmt.Sprintf("Device (%s) connected", rawDev.Serial))
	}
	return nil
}

func findConfig(device *streamdeck.Device) api.DeckV3 {
	if migrateConfigFromV1 {
		config.Decks[0].Serial = device.Serial
		_ = SaveConfig()
		migrateConfigFromV1 = false
		return config.Decks[0]
	}
	for _, deck := range config.Decks {
		if deck.Serial == device.Serial {
			return deck
		}
	}

	return makeEmptyDeckConfig(device)
}

func makeEmptyDeckConfig(device *streamdeck.Device) api.DeckV3 {
	var pages []api.PageV3
	page := api.PageV3{}
	for i := 0; i < int(device.Rows)*int(device.Columns); i++ {
		applications := make(map[string]*api.KeyConfigV3)
		applications[""] = &api.KeyConfigV3{}
		page.Keys = append(page.Keys, api.KeyV3{
			Application: applications,
		})
	}
	pages = append(pages, page)
	devConf := api.DeckV3{Serial: device.Serial, Pages: pages}
	config.Decks = append(config.Decks, devConf)
	_ = SaveConfig()
	return devConf
}

type VirtualDev struct {
	Deck           *streamdeck.Device
	Page           int
	IsOpen         bool
	Config         api.DeckV3
	mu             sync.Mutex
	shuttingDown   bool
	sdInfo         api.StreamDeckInfoV1
	keyUpdateChan  chan int
	knobUpdateChan chan int
	KeyFGBuffs     []image.Image
	KeyBGBuffs     []image.Image
	PanelFGBuffs   []image.Image
	PanelBGBuffs   []image.Image
}

func (dev *VirtualDev) SetPage(page int) {
	if locked {
		return
	}
	if page != dev.Page {
		dev.unmountPageHandlersOnPageSwitch(dev.Config.Pages[dev.Page])
	}
	dev.Page = page
	currentPage := dev.Config.Pages[page]

	w := (dev.Deck.Pixels * uint(dev.Deck.Columns)) + (dev.Deck.PaddingX * uint(dev.Deck.Columns-1))
	h := (dev.Deck.Pixels * uint(dev.Deck.Rows)) + (dev.Deck.PaddingX * uint(dev.Deck.Rows-1))
	go dev.setKeyBackground(&dev.Config.Pages[page], int(w), int(h))

	go dev.setLcdBackground(&dev.Config.Pages[page], dev.sdInfo.LcdWidth*dev.sdInfo.LcdCols, dev.sdInfo.LcdHeight)

	if dev.Config.Pages[page].GetKeyGridBackgroundHandler() != nil {
		go dev.Config.GetKeyGridBackgroundHandler().Stop()
	} else {
		go dev.setKeyBackground(&dev.Config, int(w), int(h))
	}

	if dev.Config.Pages[page].GetTouchPanelBackgroundHandler() != nil {
		go dev.Config.GetTouchPanelBackgroundHandler().Stop()
	} else {
		go dev.setLcdBackground(&dev.Config, dev.sdInfo.LcdWidth*dev.sdInfo.LcdCols, dev.sdInfo.LcdHeight)
	}

	for i, _ := range currentPage.Keys {
		key := &currentPage.Keys[i]

		go dev.setIndividualKeyBackground(key, i, dev.sdInfo.IconSize, dev.sdInfo.IconSize)

		if key.Application == nil {
			key.Application = map[string]*api.KeyConfigV3{}
			key.Application[""] = &api.KeyConfigV3{}
			currentPage.Keys[i] = *key
			log.Println(fmt.Sprintf("Setting empty application on key: %d on page: %d", i, page))
			SaveConfig()
		}
		_, keyHasApp := key.Application[currentApplication]
		if key.ActiveApplication != "" && !keyHasApp {
			key.ActiveApplication = ""
		}
		if keyHasApp {
			key.ActiveApplication = currentApplication
		}
		go SetKey(dev, key.Application[key.ActiveApplication], i, page, key.ActiveApplication)
	}
	for i, _ := range currentPage.Knobs {
		knob := &currentPage.Knobs[i]

		go dev.setIndividualLcdBackground(knob, i, dev.sdInfo.LcdWidth, dev.sdInfo.LcdHeight)

		if knob.Application == nil {
			knob.Application = map[string]*api.KnobConfigV3{}
			knob.Application[""] = &api.KnobConfigV3{}
			currentPage.Knobs[i] = *knob
			log.Println(fmt.Sprintf("Setting empty application on knob: %d on page: %d", i, page))
			SaveConfig()
		}
		_, knobHasApp := knob.Application[currentApplication]
		if knob.ActiveApplication != "" && !knobHasApp {
			knob.ActiveApplication = ""
		}
		if knobHasApp {
			knob.ActiveApplication = currentApplication
		}
		go SetKnob(dev, knob.Application[knob.ActiveApplication], i, page, knob.ActiveApplication)
	}
	dev.sdInfo.Page = page
	EmitPage(dev, page)
}

func (dev *VirtualDev) CompositeKeyImage(keyIndex int, page int) {
	var background image.Image
	keyV3 := dev.Config.Pages[page].Keys[keyIndex]
	keyConfigV3 := keyV3.Application[keyV3.ActiveApplication]

	kcbg := keyConfigV3.GetKeyBackgroundBuff()

	if kcbg != nil {
		background = kcbg
	}

	if background == nil {
		kbg := keyV3.GetKeyBackgroundBuff()
		if kbg != nil {
			background = kbg
		}
	}

	if background == nil {
		pbg := dev.Config.Pages[page].GetTouchPanelBackgroundBuff()

		if pbg != nil {
			background = pbg[keyIndex]
		}
	}

	if background == nil {
		cbg := dev.Config.GetKeyGridBackgroundBuff()

		if cbg != nil {
			background = cbg[keyIndex]
		}
	}

	if dev.KeyBGBuffs[keyIndex] != background {
		dev.KeyBGBuffs[keyIndex] = background
		dev.keyUpdateChan <- keyIndex
	}
}

func (dev *VirtualDev) SetKeyForeground(img image.Image, keyIndex int, page int) {
	if dev.Page != page {
		return
	}

	bounds := img.Bounds().Max
	if bounds.X != dev.sdInfo.IconSize || bounds.Y != dev.sdInfo.IconSize {
		img = api.ResizeImage(img, dev.sdInfo.IconSize)
	}

	if dev.KeyFGBuffs[keyIndex] != img {
		dev.KeyFGBuffs[keyIndex] = img
		dev.keyUpdateChan <- keyIndex
	}
}

func (dev *VirtualDev) CompositePanelImage(knobIndex int, page int) {
	var background image.Image
	knobV3 := dev.Config.Pages[page].Knobs[knobIndex]
	knobConfigV3 := knobV3.Application[knobV3.ActiveApplication]

	kcbg := knobConfigV3.GetTouchPanelBackgroundBuff()

	if kcbg != nil {
		background = kcbg
	}

	if background == nil {
		kbg := knobV3.GetTouchPanelBackgroundBuff()
		if kbg != nil {
			background = kbg
		}
	}

	if background == nil {
		pbg := dev.Config.Pages[page].GetTouchPanelBackgroundBuff()
		if pbg != nil {
			background = pbg[knobIndex]
		}
	}

	if background == nil {
		cbg := dev.Config.GetTouchPanelBackgroundBuff()
		if cbg != nil {
			background = cbg[knobIndex]
		}
	}

	if dev.PanelBGBuffs[knobIndex] != background {
		dev.PanelBGBuffs[knobIndex] = background
		dev.knobUpdateChan <- knobIndex
	}
}

func (dev *VirtualDev) SetPanelForeground(img image.Image, knobIndex int, page int) {
	if dev.Page != page {
		return
	}

	bounds := img.Bounds().Max
	if bounds.X != dev.sdInfo.LcdWidth || bounds.Y != dev.sdInfo.LcdHeight {
		img = api.ResizeImageWH(img, dev.sdInfo.LcdWidth, dev.sdInfo.LcdHeight)
	}

	if dev.PanelFGBuffs[knobIndex] != img {
		dev.PanelFGBuffs[knobIndex] = img
		dev.knobUpdateChan <- knobIndex
	}
}

func (dev *VirtualDev) UnmountHandlers() {
	for i := range dev.Config.Pages {
		dev.unmountPageHandlersOnPageSwitch(dev.Config.Pages[i])
	}
}

func (dev *VirtualDev) SetBrightness(brightness uint8) error {
	return dev.Deck.SetBrightness(brightness)
}

func (dev *VirtualDev) SetSdInfo() {

	manufacturer, err := dev.Deck.Device.GetManufacturer()
	if err != nil {
		log.Println(err)
	}
	product, err := dev.Deck.Device.GetProduct()
	if err != nil {
		log.Println(err)
	}
	info := api.StreamDeckInfoV1{
		Cols:          int(dev.Deck.Columns),
		Rows:          int(dev.Deck.Rows),
		IconSize:      int(dev.Deck.Pixels),
		Page:          0,
		Serial:        dev.Deck.Serial,
		Name:          manufacturer + " " + product,
		Connected:     true,
		LastConnected: time.Now(),
		LcdWidth:      int(dev.Deck.LcdWidth),
		LcdHeight:     int(dev.Deck.LcdHeight),
		LcdCols:       int(dev.Deck.LcdColumns),
		KnobCols:      int(dev.Deck.Knobs),
		PaddingX:      int(dev.Deck.PaddingX),
		PaddingY:      int(dev.Deck.PaddingY),
	}

	dev.sdInfo = info
}

func (dev *VirtualDev) ApplicationUpdated() {
	if locked {
		return
	}
	page := dev.Config.Pages[dev.Page]
	for i := range page.Keys {
		key := &page.Keys[i]
		_, keyHasApp := key.Application[currentApplication]
		activeApp := key.ActiveApplication
		if key.Application[key.ActiveApplication].KeyHold != 0 && (keyHasApp || key.ActiveApplication != "") {
			kb.KeyUp(key.Application[key.ActiveApplication].KeyHold)
		}
		if key.ActiveApplication != "" && !keyHasApp {
			key.ActiveApplication = ""
		}
		if keyHasApp {
			key.ActiveApplication = currentApplication
		}
		if key.ActiveApplication != activeApp {
			go SetKey(dev, key.Application[key.ActiveApplication], i, dev.Page, key.ActiveApplication)
		}
	}
	for i := range page.Knobs {
		knob := &page.Knobs[i]
		_, keyHasApp := knob.Application[currentApplication]
		activeApp := knob.ActiveApplication
		if knob.ActiveApplication != "" && !keyHasApp {
			knob.ActiveApplication = ""
		}
		if keyHasApp {
			knob.ActiveApplication = currentApplication
		}
		if knob.ActiveApplication != activeApp {
			go SetKnob(dev, knob.Application[knob.ActiveApplication], i, dev.Page, knob.ActiveApplication)
		}
	}
	dev.unmountPageHandlersOnAppSwitch(page)
}

func (dev *VirtualDev) HandleScreenLockChange(locked bool) {
	if locked {
		dev.UnmountHandlers()
		dev.Deck.Reset()
	} else {
		dev.SetPage(dev.Page)
	}
}

func (dev *VirtualDev) Render() {
	err := dev.Deck.SetImage(0, image.NewRGBA(image.Rect(0, 0, int(dev.Deck.Pixels), int(dev.Deck.Pixels))))
	if err != nil {
		log.Println(err)
		return
	}

	go dev.RenderKey()

	go dev.RenderKnob()
}

func (dev *VirtualDev) RenderKey() {
	for dev.IsOpen && !dev.shuttingDown {

		keyIndex := <-dev.keyUpdateChan

		mergedImage, err := api.LayerImages(int(dev.Deck.Pixels), int(dev.Deck.Pixels), dev.KeyBGBuffs[keyIndex], dev.KeyFGBuffs[keyIndex])

		if err != nil {
			dev.keyUpdateChan <- keyIndex
			log.Println(err)
			continue
		}

		bounds := mergedImage.Bounds().Max

		dev.mu.Lock()

		err = dev.Deck.SetImage(uint8(keyIndex), mergedImage)

		dev.mu.Unlock()

		if err != nil {
			dev.keyUpdateChan <- keyIndex
			match, _ := regexp.MatchString(`.*hidapi.*`, err.Error())
			if match {
				dev.Disconnect()
				return
			}
			match, _ = regexp.MatchString(`.*dimensions.*`, err.Error())
			if match {
				log.Println(fmt.Sprintf("%s provided: %d x %d", err.Error(), bounds.X, bounds.Y))
				return
			}

			log.Println(err)
		}
	}
}

func (dev *VirtualDev) RenderKnob() {
	for dev.IsOpen && !dev.shuttingDown {

		knobIndex := <-dev.knobUpdateChan

		mergedImage, err := api.LayerImages(int(dev.Deck.LcdWidth), int(dev.Deck.LcdHeight), dev.PanelBGBuffs[knobIndex], dev.PanelFGBuffs[knobIndex])

		if err != nil {
			dev.knobUpdateChan <- knobIndex
			log.Println(err)
			continue
		}

		bounds := mergedImage.Bounds().Max

		dev.mu.Lock()

		err = dev.Deck.SetLcdImage(knobIndex, mergedImage)

		dev.mu.Unlock()

		if err != nil {
			dev.knobUpdateChan <- knobIndex
			match, _ := regexp.MatchString(`.*hidapi.*`, err.Error())
			if match {
				dev.Disconnect()
				return
			}
			match, _ = regexp.MatchString(`.*dimensions.*`, err.Error())
			if match {
				log.Println(fmt.Sprintf("%s provided: %d x %d", err.Error(), bounds.X, bounds.Y))
				return
			}

			log.Println(err)
		}
	}
}

func (dev *VirtualDev) HandleInput() {
	defer func() {
		if err := recover(); err != nil {
			dev.Disconnect()
		}
	}()
	dev.Deck.HandleInput(func(event streamdeck.InputEvent) {
		if !locked {
			if event.EventType == streamdeck.KEY_PRESS || event.EventType == streamdeck.KEY_RELEASE {
				page := dev.Config.Pages[dev.Page]
				if uint8(len(page.Keys)) > event.Index {
					HandleKeyInput(dev, &page.Keys[event.Index], event.EventType == streamdeck.KEY_PRESS)
				}
			} else if event.EventType == streamdeck.SCREEN_SWIPE {
				if event.ScreenEndX < event.ScreenX {
					if dev.Page < len(dev.Config.Pages)-1 {
						dev.SetPage(dev.Page + 1)
					}
				} else {
					if dev.Page > 0 {
						dev.SetPage(dev.Page - 1)
					}
				}
			} else if dev.Deck.HasLCD && dev.Deck.HasKnobs {
				page := dev.Config.Pages[dev.Page]
				if uint8(len(page.Knobs)) > event.Index {
					HandleKnobInput(dev, &page.Knobs[event.Index], event)
				}
			}
		}
	})
}

func (dev *VirtualDev) Disconnect() {
	disconnectSem.Lock()
	defer disconnectSem.Unlock()
	if !dev.IsOpen {
		return
	}
	log.Println("Device (" + dev.Deck.Serial + ") disconnected")
	err := dev.Deck.Close()
	if err != nil {
		log.Println(err)
	}
	dev.IsOpen = false
	dev.sdInfo.Connected = false
	dev.sdInfo.LastDisconnected = time.Now()
	dev.UnmountHandlers()
}

func (dev *VirtualDev) unmountPageHandlersOnPageSwitch(page api.PageV3) {

	if page.KeyGridBackgroundHandler != nil {
		page.KeyGridBackgroundHandler.Stop()
	}

	if page.TouchPanelBackgroundHandler != nil {
		page.TouchPanelBackgroundHandler.Stop()
	}

	for i2 := 0; i2 < len(page.Keys); i2++ {
		key := &page.Keys[i2]

		if key.KeyBackgroundHandler != nil {
			key.KeyBackgroundHandler.Stop()
		}

		for _, keyConfig := range key.Application {

			if keyConfig.KeyBackgroundHandler != nil {
				keyConfig.KeyBackgroundHandler.Stop()
			}

			if keyConfig.IconHandlerStruct != nil {
				log.Printf("Stopping %s\n", keyConfig.IconHandler)
				if keyConfig.IconHandlerStruct.IsRunning() {
					go UnmountKeyHandler(keyConfig)
				}
			}
		}

	}
	for i2 := 0; i2 < len(page.Knobs); i2++ {
		knob := &page.Knobs[i2]

		if knob.TouchPanelBackgroundHandler != nil {
			knob.TouchPanelBackgroundHandler.Stop()
		}

		for _, knobConfig := range knob.Application {

			if knobConfig.TouchPanelBackgroundHandler != nil {
				knobConfig.TouchPanelBackgroundHandler.Stop()
			}

			if knobConfig.LcdHandlerStruct != nil {
				log.Printf("Stopping %s\n", knobConfig.LcdHandler)
				if knobConfig.LcdHandlerStruct.IsRunning() {
					go UnmountKnobHandler(knobConfig)
				}
			}
		}
	}
}

func (dev *VirtualDev) unmountPageHandlersOnAppSwitch(page api.PageV3) {

	for i2 := 0; i2 < len(page.Keys); i2++ {
		key := &page.Keys[i2]

		_, keyHasApp := key.Application[currentApplication]
		for app := range key.Application {
			keyConfig := key.Application[app]

			if keyConfig.KeyBackgroundHandler != nil {
				keyConfig.KeyBackgroundHandler.Stop()
			}

			if (keyHasApp && app == currentApplication) || (!keyHasApp && app == "") {
				continue
			}
			if keyConfig.IconHandlerStruct != nil && keyConfig.IconHandlerStruct.IsRunning() {
				log.Printf("Stopping %s\n", keyConfig.IconHandler)
				if keyConfig.IconHandlerStruct.IsRunning() {
					go UnmountKeyHandler(keyConfig)
				}
			}
		}

	}
	for i2 := 0; i2 < len(page.Knobs); i2++ {
		knob := &page.Knobs[i2]

		_, keyHasApp := knob.Application[currentApplication]
		for app := range knob.Application {
			knobConfig := knob.Application[app]

			if knobConfig.TouchPanelBackgroundHandler != nil {
				go knobConfig.TouchPanelBackgroundHandler.Stop()
			}

			if (keyHasApp && app == currentApplication) || (!keyHasApp && app == "") {
				continue
			}
			if knobConfig.LcdHandlerStruct != nil && knobConfig.LcdHandlerStruct.IsRunning() {
				log.Printf("Stopping %s\n", knobConfig.LcdHandler)
				if knobConfig.LcdHandlerStruct.IsRunning() {
					go UnmountKnobHandler(knobConfig)
				}
			}
		}
	}
}

func (dev *VirtualDev) setLcdBackground(backgrounder api.LcdBackgrounder, w, h int) {
	if backgrounder.GetTouchPanelBackground() == "" {
		return
	}

	if backgrounder.GetTouchPanelBackgroundHandler() == nil {
		var handler api.TouchPanelBackgroundHandler

		for _, module := range modules {
			if module.Name == backgrounder.GetTouchPanelBackground() {
				handler = module.NewTouchPanelBackgroundHandler()
			}
		}

		backgrounder.SetTouchPanelBackgroundHandler(handler)
	}

	if backgrounder.GetTouchPanelBackgroundHandlerFields() != nil {
		go backgrounder.GetTouchPanelBackgroundHandler().Start(backgrounder.GetTouchPanelBackgroundHandlerFields(), dev.sdInfo, func(imgs []image.Image) {
			if len(imgs) == int(dev.Deck.Knobs) {
				backgrounder.SetTouchPanelBackgroundBuff(imgs)

				for u := range imgs {
					dev.CompositePanelImage(u, dev.Page)
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
		log.Println(err)
		return
	}

	img = api.ResizeImageWH(img, w, h)

	var imgs []image.Image

	for lcdIndex := range int(dev.Deck.LcdColumns) {
		x0, y0 := dev.sdInfo.LcdWidth*lcdIndex, 0
		x1, y1 := dev.sdInfo.LcdWidth*(lcdIndex+1), dev.sdInfo.LcdHeight

		imgs = append(imgs, api.SubImage(img, x0, y0, x1, y1))
	}

	backgrounder.SetTouchPanelBackgroundBuff(imgs)

	for index, _ := range imgs {
		dev.CompositePanelImage(index, dev.Page)
	}
}

func (dev *VirtualDev) setKeyBackground(backgrounder api.KeyGridBackgrounder, w, h int) {
	if backgrounder.GetKeyGridBackground() == "" {
		return
	}

	if backgrounder.GetKeyGridBackgroundHandler() == nil {
		var handler api.KeyGridBackgroundHandler

		for _, module := range modules {
			if module.Name == backgrounder.GetKeyGridBackground() {
				handler = module.NewKeyGridBackground()
			}
		}

		backgrounder.SetKeyGridBackgroundHandler(handler)
	}

	if backgrounder.GetKeyGridBackgroundHandler() != nil {
		go backgrounder.GetKeyGridBackgroundHandler().Start(backgrounder.GetKeyGridBackgroundHandlerFields(), dev.sdInfo, func(imgs []image.Image) {
			if len(imgs) == int(dev.Deck.Keys) {
				backgrounder.SetKeyGridBackgroundBuff(imgs)

				for u := range imgs {
					dev.CompositeKeyImage(u, dev.Page)
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
		log.Println(err)
		return
	}

	img = api.ResizeImageWH(img, w, h)

	var imgs []image.Image
	for keyIndex := range int(dev.Deck.Keys) {
		keyX := keyIndex % int(dev.Deck.Columns)
		keyY := int(math.Floor(float64(keyIndex) / float64(dev.Deck.Columns)))

		x0, y0 := keyX*int(dev.Deck.Pixels+dev.Deck.PaddingX), keyY*int(dev.Deck.Pixels+dev.Deck.PaddingY)
		x1, y1 := keyX*int(dev.Deck.Pixels+dev.Deck.PaddingX)+int(dev.Deck.Pixels), keyY*int(dev.Deck.Pixels+dev.Deck.PaddingY)+int(dev.Deck.Pixels)

		imgs = append(imgs, api.SubImage(img, x0, y0, x1, y1))
	}
	backgrounder.SetKeyGridBackgroundBuff(imgs)

	for index, _ := range imgs {
		dev.CompositeKeyImage(index, dev.Page)
	}
}

func (dev *VirtualDev) setIndividualLcdBackground(backgrounder api.LcdSegmentBackgrounder, index, w, h int) {
	if backgrounder.GetTouchPanelBackground() == "" {
		return
	}

	if backgrounder.GetTouchPanelBackgroundHandler() == nil {
		var handler api.TouchPanelBackgroundHandler

		for _, module := range modules {
			if module.Name == backgrounder.GetTouchPanelBackground() {
				handler = module.NewTouchPanelBackgroundHandler()
			}
		}

		backgrounder.SetTouchPanelBackgroundHandler(handler)
	}

	if backgrounder.GetTouchPanelBackgroundHandlerFields() != nil {
		go backgrounder.GetTouchPanelBackgroundHandler().StartIndividual(backgrounder.GetTouchPanelBackgroundHandlerFields(), dev.sdInfo, func(img image.Image) {
			backgrounder.SetTouchPanelBackgroundBuff(img)

			dev.CompositePanelImage(index, dev.Page)
		})
		return
	}

	if backgrounder.GetTouchPanelBackgroundBuff() != nil {
		return
	}

	img, err := LoadImage(backgrounder.GetTouchPanelBackground())
	if err != nil {
		log.Println(err)
		return
	}

	img = api.ResizeImageWH(img, w, h)

	backgrounder.SetTouchPanelBackgroundBuff(img)

	dev.CompositePanelImage(index, dev.Page)
}

func (dev *VirtualDev) setIndividualKeyBackground(backgrounder api.KeyBackgrounder, index, w, h int) {
	if backgrounder.GetKeyBackground() == "" {
		return
	}

	if backgrounder.GetKeyBackgroundHandler() == nil {
		var handler api.KeyGridBackgroundHandler

		for _, module := range modules {
			if module.Name == backgrounder.GetKeyBackground() {
				handler = module.NewKeyGridBackground()
			}
		}

		backgrounder.SetKeyBackgroundHandler(handler)
	}

	if backgrounder.GetKeyBackgroundHandler() != nil {
		go backgrounder.GetKeyBackgroundHandler().StartIndividual(backgrounder.GetKeyBackgroundHandlerFields(), dev.sdInfo, func(img image.Image) {
			backgrounder.SetKeyBackgroundBuff(img)

			dev.CompositeKeyImage(index, dev.Page)
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

	img = api.ResizeImageWH(img, w, h)

	dev.CompositeKeyImage(index, dev.Page)
}

func (dev *VirtualDev) Stop() {
	dev.shuttingDown = true
	if dev.IsOpen {
		err := dev.Deck.Reset()
		if err != nil {
			log.Println(err)
		}
		err = dev.Deck.Close()
		if err != nil {
			log.Println(err)
		}
	}
}
