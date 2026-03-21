package streamdeckd

import (
	"errors"
	"fmt"
	"image"
	"log"
	"os"
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

		InitDevice(rawDev)
		log.Println(fmt.Sprintf("Device (%s) connected", rawDev.Serial))
	}
	return nil
}

func InitDevice(rawDev *streamdeck.Device) {
	dev, ok := Devs[rawDev.Serial]

	if !ok {
		// initial connect
		config := findConfig(rawDev)
		dev = &VirtualDev{
			Deck:           rawDev,
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

		dev.backgrounder = &Backgrounder{
			vdev:   dev,
			sdInfo: dev.sdInfo,
		}

		dev.pageManager = &PageManager{
			vdev: dev,
			page: 0,
		}

		dev.handlerPruner = &HandlerPruner{
			vdev: dev,
		}

		dev.inputManager = &InputManager{
			vdev: dev,
		}

		dev.foregrounder = &Foregrounder{
			vdev: dev,
		}

		dev.backgrounder.AttachPageChangeListener()

		dev.foregrounder.AttachPageChangeListener()
		dev.foregrounder.AttachAppChangeListener()

		dev.handlerPruner.OnPageChange()
		dev.handlerPruner.OnAppSwitch()

		dev.logger = log.New(os.Stdout, fmt.Sprintf("(%s) ", dev.sdInfo.Serial), log.Lshortfile|log.Ltime)

		Devs[rawDev.Serial] = dev
	} else {
		//reconnect
		dev.IsOpen = true
		dev.Deck = rawDev
		dev.sdInfo.LastConnected = time.Now()
		dev.sdInfo.Connected = true
	}

	dev.pageManager.SetPage(dev.pageManager.page)

	go dev.backgrounder.setKeyBackground(&dev.Config)

	go dev.backgrounder.setLcdBackground(&dev.Config)

	go dev.HandleInput()
	dev.Render()
}

type VirtualDev struct {
	Deck           *streamdeck.Device
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

	foregrounder  *Foregrounder
	backgrounder  *Backgrounder
	pageManager   *PageManager
	handlerPruner *HandlerPruner
	inputManager  *InputManager

	logger *log.Logger
}

func (dev *VirtualDev) SetKeyBackground(keyIndex int, page int) {
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
	if dev.pageManager.page != page {
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

func (dev *VirtualDev) SetPanelBackground(knobIndex int, page int) {
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
	if dev.pageManager.page != page {
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

func (dev *VirtualDev) SetBrightness(brightness uint8) error {
	return dev.Deck.SetBrightness(brightness)
}

func (dev *VirtualDev) SetSdInfo() {

	manufacturer, err := dev.Deck.Device.GetManufacturer()
	if err != nil {
		dev.logger.Println(err)
	}
	product, err := dev.Deck.Device.GetProduct()
	if err != nil {
		dev.logger.Println(err)
	}

	info := api.StreamDeckInfoV1{
		Cols:                    int(dev.Deck.Columns),
		Rows:                    int(dev.Deck.Rows),
		IconSize:                int(dev.Deck.Pixels),
		Page:                    0,
		Serial:                  dev.Deck.Serial,
		Name:                    manufacturer + " " + product,
		Connected:               true,
		LastConnected:           time.Now(),
		LcdWidth:                int(dev.Deck.LcdWidth),
		LcdHeight:               int(dev.Deck.LcdHeight),
		LcdCols:                 int(dev.Deck.LcdColumns),
		KnobCols:                int(dev.Deck.Knobs),
		PaddingX:                int(dev.Deck.PaddingX),
		PaddingY:                int(dev.Deck.PaddingY),
		KeyGridBackgroundWidth:  int((dev.Deck.Pixels * uint(dev.Deck.Columns)) + (dev.Deck.PaddingX * uint(dev.Deck.Columns-1))),
		KeyGridBackgroundHeight: int((dev.Deck.Pixels * uint(dev.Deck.Rows)) + (dev.Deck.PaddingX * uint(dev.Deck.Rows-1))),
		LcdBackgroundWidth:      int(dev.Deck.LcdWidth * uint(dev.Deck.LcdColumns)),
		LcdBackgroundHeight:     int(dev.Deck.LcdHeight),
	}

	dev.sdInfo = info
}

func (dev *VirtualDev) HandleScreenLockChange(locked bool) {
	if locked {
		dev.handlerPruner.StopAllHandlers()
		dev.Deck.Reset()
	} else {
		dev.pageManager.SetPage(dev.pageManager.page)
	}
}

func (dev *VirtualDev) Render() {
	err := dev.Deck.SetImage(0, image.NewRGBA(image.Rect(0, 0, int(dev.Deck.Pixels), int(dev.Deck.Pixels))))
	if err != nil {
		dev.logger.Println(err)
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
			dev.logger.Println(err)
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
				dev.logger.Println(fmt.Sprintf("%s provided: %d x %d", err.Error(), bounds.X, bounds.Y))
				return
			}

			dev.logger.Println(err)
		}
	}
}

func (dev *VirtualDev) RenderKnob() {
	for dev.IsOpen && !dev.shuttingDown {

		knobIndex := <-dev.knobUpdateChan

		mergedImage, err := api.LayerImages(int(dev.Deck.LcdWidth), int(dev.Deck.LcdHeight), dev.PanelBGBuffs[knobIndex], dev.PanelFGBuffs[knobIndex])

		if err != nil {
			dev.knobUpdateChan <- knobIndex
			dev.logger.Println(err)
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
				dev.logger.Println(fmt.Sprintf("%s provided: %d x %d", err.Error(), bounds.X, bounds.Y))
				return
			}

			dev.logger.Println(err)
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
				page := dev.Config.Pages[dev.pageManager.page]
				if uint8(len(page.Keys)) > event.Index {
					dev.inputManager.HandleKeyInput(&page.Keys[event.Index], event.EventType == streamdeck.KEY_PRESS)
				}
			} else if event.EventType == streamdeck.SCREEN_SWIPE {
				if event.ScreenEndX < event.ScreenX {
					if dev.pageManager.page < len(dev.Config.Pages)-1 {
						dev.pageManager.SetPage(dev.pageManager.page + 1)
					}
				} else {
					if dev.pageManager.page > 0 {
						dev.pageManager.SetPage(dev.pageManager.page - 1)
					}
				}
			} else if dev.Deck.HasLCD && dev.Deck.HasKnobs {
				page := dev.Config.Pages[dev.pageManager.page]
				if uint8(len(page.Knobs)) > event.Index {
					dev.inputManager.HandleKnobInput(&page.Knobs[event.Index], event)
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
	dev.logger.Println("Device (" + dev.Deck.Serial + ") disconnected")
	err := dev.Deck.Close()
	if err != nil {
		dev.logger.Println(err)
	}
	dev.IsOpen = false
	dev.sdInfo.Connected = false
	dev.sdInfo.LastDisconnected = time.Now()
	dev.handlerPruner.StopAllHandlers()
}

func (dev *VirtualDev) Stop() {
	dev.shuttingDown = true
	if dev.IsOpen {
		err := dev.Deck.Reset()
		if err != nil {
			dev.logger.Println(err)
		}
		err = dev.Deck.Close()
		if err != nil {
			dev.logger.Println(err)
		}
	}
}
