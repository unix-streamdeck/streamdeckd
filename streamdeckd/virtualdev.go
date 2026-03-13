package streamdeckd

import (
	"errors"
	"fmt"
	"image"
	"log"
	"regexp"
	"strings"
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
			dev = &VirtualDev{Deck: rawDev, Page: 0, IsOpen: true, Config: config}
			dev.SetSdInfo()
			Devs[rawDev.Serial] = dev
		} else {
			//reconnect
			dev.IsOpen = true
			dev.Deck = rawDev
			dev.sdInfo.LastConnected = time.Now()
			dev.sdInfo.Connected = true
		}
		go dev.ReadDevKey()
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
	Deck         *streamdeck.Device
	Page         int
	IsOpen       bool
	Config       api.DeckV3
	mu           sync.Mutex
	shuttingDown bool
	sdInfo       api.StreamDeckInfoV1
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
	for i := 0; i < len(currentPage.Keys); i++ {
		currentKey := &currentPage.Keys[i]
		if currentKey.Application == nil {
			currentKey.Application = map[string]*api.KeyConfigV3{}
			currentKey.Application[""] = &api.KeyConfigV3{}
			currentPage.Keys[i] = *currentKey
			log.Println(fmt.Sprintf("Setting empty application on key: %d on page: %d", i, page))
			SaveConfig()
		}
		_, keyHasApp := currentKey.Application[currentApplication]
		if currentKey.ActiveApplication != "" && !keyHasApp {
			currentKey.ActiveApplication = ""
		}
		if keyHasApp {
			currentKey.ActiveApplication = currentApplication
		}
		go SetKey(dev, currentKey.Application[currentKey.ActiveApplication], i, page, currentKey.ActiveApplication)
	}
	for i, knob := range currentPage.Knobs {
		go SetKnob(dev, knob.Application[knob.ActiveApplication], i, page, knob.ActiveApplication)
	}
	dev.sdInfo.Page = page
	EmitPage(dev, page)
}

func (dev *VirtualDev) SetImage(img image.Image, keyIndex int, page int) {
	defer func() {
		if err := recover(); err != nil {
			dev.Disconnect()
		}
	}()
	if locked {
		return
	}
	dev.mu.Lock()
	defer dev.mu.Unlock()
	bounds := img.Bounds().Max
	if bounds.X != dev.sdInfo.IconSize || bounds.Y != dev.sdInfo.IconSize {
		img = api.ResizeImage(img, dev.sdInfo.IconSize)
	}
	if bounds.X == 0 && bounds.Y == 0 {
		log.Println("Empty image received")
		return
	}

	if dev.Page == page && dev.IsOpen && !dev.shuttingDown {

		err := dev.Deck.SetImage(uint8(keyIndex), img)
		if err != nil {
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

func GetSubImage(img image.Image, x0, y0, x1, y1 int) image.Image {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	simg, ok := img.(subImager)

	if !ok {
		log.Println("Couldn't resize")
		return nil
	}

	rect := image.Rect(x0, y0, x1, y1)

	img = simg.SubImage(rect)

	return img
}

func (dev *VirtualDev) GetPanelImagePart(knobIndex int, page int) image.Image {
	pageV3 := dev.Config.Pages[page]

	x0, y0 := dev.sdInfo.LcdWidth*knobIndex, 0
	x1, y1 := dev.sdInfo.LcdWidth*(knobIndex+1), dev.sdInfo.LcdHeight

	pbg := dev.GetBackgroundBuff(&pageV3, dev.sdInfo.LcdWidth*dev.sdInfo.LcdCols, dev.sdInfo.LcdHeight)

	if pbg != nil {
		return GetSubImage(pageV3.TouchPanelBackgroundBuff, x0, y0, x1, y1)
	}

	cbg := dev.GetBackgroundBuff(&dev.Config, dev.sdInfo.LcdWidth*dev.sdInfo.LcdCols, dev.sdInfo.LcdHeight)

	if cbg != nil {
		return GetSubImage(dev.Config.TouchPanelBackgroundBuff, x0, y0, x1, y1)
	}

	return nil
}

func (dev *VirtualDev) GetBackgroundBuff(backgrounder api.Backgrounder, w, h int) image.Image {
	if backgrounder.GetTouchPanelBackground() != "" {
		if backgrounder.GetTouchPanelBackgroundBuff() == nil {
			img, err := LoadImage(backgrounder.GetTouchPanelBackground())
			if err != nil {
				log.Println(err)
				return nil
			}
			img = api.ResizeImageWH(img, w, h)
			backgrounder.SetTouchPanelBackgroundBuff(img)
		}
		return backgrounder.GetTouchPanelBackgroundBuff()
	}
	return nil
}

func (dev *VirtualDev) CompositePanelImage(foreground image.Image, knobIndex int, page int) image.Image {
	knobV3 := dev.Config.Pages[page].Knobs[knobIndex]
	knobConfigV3 := knobV3.Application[knobV3.ActiveApplication]

	kcbg := dev.GetBackgroundBuff(knobConfigV3, dev.sdInfo.LcdWidth, dev.sdInfo.LcdHeight)

	if kcbg != nil {
		mergedImage, err := api.LayerImages(knobConfigV3.TouchPanelBackgroundBuff, foreground)
		if err != nil {
			log.Println(err)
			return foreground
		}
		return mergedImage
	}

	kbg := dev.GetBackgroundBuff(knobConfigV3, dev.sdInfo.LcdWidth, dev.sdInfo.LcdHeight)

	if kbg != nil {
		mergedImage, err := api.LayerImages(knobV3.TouchPanelBackgroundBuff, foreground)
		if err != nil {
			log.Println(err)
			return foreground
		}
		return mergedImage
	}

	background := dev.GetPanelImagePart(knobIndex, page)

	if background == nil {
		return foreground
	}

	mergedImage, err := api.LayerImages(background, foreground)

	if err != nil {
		log.Println(err)
		return foreground
	}

	return mergedImage
}

func (dev *VirtualDev) SetPanelImage(img image.Image, knobIndex int, page int) {
	defer func() {
		if err := recover(); err != nil {
			dev.Disconnect()
		}
	}()
	if locked {
		return
	}
	dev.mu.Lock()
	defer dev.mu.Unlock()
	bounds := img.Bounds().Max
	if bounds.X != dev.sdInfo.LcdWidth || bounds.Y != dev.sdInfo.LcdHeight {
		img = api.ResizeImageWH(img, dev.sdInfo.LcdWidth, dev.sdInfo.LcdHeight)
	}
	img = dev.CompositePanelImage(img, knobIndex, page)
	if bounds.X == 0 && bounds.Y == 0 {
		log.Println("Empty image received")
		return
	}

	if dev.Page == page && dev.IsOpen && !dev.shuttingDown {

		err := dev.Deck.SetLcdImage(knobIndex, img)
		if err != nil {
			if strings.Contains(err.Error(), "hidapi") {
				dev.Disconnect()
			} else if strings.Contains(err.Error(), "dimensions") {
				log.Println(fmt.Sprintf("%s provided: %d x %d", err.Error(), bounds.X, bounds.Y))
			} else {
				log.Println(err)
			}
		}
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

func (dev *VirtualDev) ReadDevKey() {
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
	for i2 := 0; i2 < len(page.Keys); i2++ {
		key := &page.Keys[i2]
		for _, keyConfig := range key.Application {
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
		for _, knobConfig := range knob.Application {
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
			if (keyHasApp && app == currentApplication) || (!keyHasApp && app == "") {
				continue
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
		_, keyHasApp := knob.Application[currentApplication]
		for app := range knob.Application {
			knobConfig := knob.Application[app]
			if (keyHasApp && app == currentApplication) || (!keyHasApp && app == "") {
				continue
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
