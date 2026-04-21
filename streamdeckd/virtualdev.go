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
	"golang.org/x/image/draw"
)

var disconnectSem sync.Mutex
var connectSem sync.Mutex
var Devs map[string]IVirtualDev

type IVirtualDev interface {
	IsOpen() bool
	Config() api.DeckV3
	SetConfig(v3 api.DeckV3)
	SdInfo() *api.StreamDeckInfoV1
	Serial() string
	Foregrounder() IForegrounder
	Backgrounder() IBackgrounder
	PageManager() IPageManager
	HandlerPruner() IHandlerPruner
	InputManager() IInputManager
	Logger() *log.Logger

	Open(rawDev *streamdeck.Device) error
	SetKeyBackground(keyIndex int, page int)
	SetKeyForeground(img image.Image, keyIndex int, page int)
	SetPanelBackground(knobIndex int, page int)
	SetPanelForeground(img image.Image, knobIndex int, page int)
	RedrawKey(keyIndex int)
	SetBrightness(brightness uint8) error
	HandleScreenLockChange(locked bool)
	Close()
}

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
		if len(rawDev.Serial) != 12 && len(rawDev.Serial) != 14 {
			continue
		}
		dev, ok := Devs[rawDev.Serial]
		if ok && dev.IsOpen() {
			continue
		}

		if dev == nil {
			dev = &VirtualDev{}
		}

		err := dev.Open(rawDev)
		if err == nil {
			log.Println(fmt.Sprintf("Device (%s) connected", rawDev.Serial))
		}
	}
	return nil
}

type VirtualDev struct {

	//Internal Properties
	mu             sync.Mutex
	shuttingDown   bool
	keyUpdateChan  chan int
	knobUpdateChan chan int
	keyFGBuffs     []image.Image
	keyBGBuffs     []image.Image
	panelFGBuffs   []image.Image
	panelBGBuffs   []image.Image
	roundedCorners image.Image

	//External Properties
	isOpen        bool
	config        api.DeckV3
	sdInfo        *api.StreamDeckInfoV1
	deck          *streamdeck.Device
	foregrounder  IForegrounder
	backgrounder  IBackgrounder
	pageManager   IPageManager
	handlerPruner IHandlerPruner
	inputManager  IInputManager
	logger        *log.Logger
}

func (dev *VirtualDev) Open(rawDev *streamdeck.Device) error {

	err := rawDev.Open()
	if err != nil {
		log.Println(err)
		return err
	}

	if dev.deck == nil {
		// initial connect
		config := findConfig(rawDev)
		dev = &VirtualDev{
			deck:           rawDev,
			isOpen:         true,
			config:         config,
			keyUpdateChan:  make(chan int),
			knobUpdateChan: make(chan int),
			keyBGBuffs:     make([]image.Image, rawDev.Keys),
			keyFGBuffs:     make([]image.Image, rawDev.Keys),
			panelBGBuffs:   make([]image.Image, rawDev.LcdColumns),
			panelFGBuffs:   make([]image.Image, rawDev.LcdColumns),
		}
		dev.setSdInfo()

		img, err := getRoundedCornersImage(int(rawDev.Pixels))

		if err == nil {
			dev.roundedCorners = img
		}

		dev.backgrounder = &Backgrounder{
			vdev: dev,
		}

		dev.pageManager = &PageManager{
			vdev: dev,
			page: 0,
		}

		dev.handlerPruner = &HandlerPruner{
			vdev: dev,
		}

		dev.inputManager = &InputManager{
			vdev:      dev,
			KeyStates: make([]bool, dev.deck.Keys),
		}

		dev.foregrounder = &Foregrounder{
			vdev: dev,
		}

		dev.backgrounder.AttachPageChangeListener()

		dev.pageManager.AttachListener(func(_, _ int) {
			dev.keyFGBuffs = make([]image.Image, rawDev.Keys)
			dev.panelFGBuffs = make([]image.Image, rawDev.LcdColumns)
		})

		dev.foregrounder.AttachPageChangeListener()
		dev.foregrounder.AttachAppChangeListener()

		dev.handlerPruner.OnPageChange()
		dev.handlerPruner.OnAppSwitch()

		dev.logger = log.New(os.Stdout, fmt.Sprintf("(%s) ", dev.sdInfo.Serial), log.Lshortfile|log.Ltime)

		Devs[rawDev.Serial] = dev
	} else {
		//reconnect
		dev.isOpen = true
		dev.deck = rawDev
		dev.sdInfo.LastConnected = time.Now()
		dev.sdInfo.Connected = true
	}

	dev.pageManager.SetPage(dev.pageManager.GetPage())

	go dev.backgrounder.SetKeyBackground(&dev.config)

	go dev.backgrounder.SetLcdBackground(&dev.config)

	go dev.handleInput()
	dev.render()

	return nil
}

func (dev *VirtualDev) IsOpen() bool {
	return dev.isOpen
}

func (dev *VirtualDev) Config() api.DeckV3 {
	return dev.config
}

func (dev *VirtualDev) SetConfig(config api.DeckV3) {
	dev.config = config
}

func (dev *VirtualDev) SdInfo() *api.StreamDeckInfoV1 {
	return dev.sdInfo
}

func (dev *VirtualDev) Serial() string {
	return dev.deck.Serial
}

func (dev *VirtualDev) Foregrounder() IForegrounder {
	return dev.foregrounder
}

func (dev *VirtualDev) Backgrounder() IBackgrounder {
	return dev.backgrounder
}

func (dev *VirtualDev) PageManager() IPageManager {
	return dev.pageManager
}

func (dev *VirtualDev) HandlerPruner() IHandlerPruner {
	return dev.handlerPruner
}

func (dev *VirtualDev) InputManager() IInputManager {
	return dev.inputManager
}

func (dev *VirtualDev) Logger() *log.Logger {
	return dev.logger
}

func (dev *VirtualDev) SetKeyBackground(keyIndex int, page int) {
	var background image.Image
	var keyV3 *api.KeyV3
	var keyConfigV3 *api.KeyConfigV3
	if len(dev.config.Pages[page].Keys) > keyIndex {
		keyV3 = &dev.config.Pages[page].Keys[keyIndex]
		kcv3, ok := keyV3.Application[keyV3.ActiveApplication]
		if ok {
			keyConfigV3 = kcv3
		}
	}

	if keyConfigV3 != nil {
		kcbg := keyConfigV3.GetKeyBackgroundBuff()

		if kcbg != nil {
			background = kcbg
		}

	}

	if background == nil && keyV3 != nil {
		kbg := keyV3.GetKeyBackgroundBuff()
		if kbg != nil {
			background = kbg
		}
	}

	if background == nil {
		pbg := dev.config.Pages[page].GetTouchPanelBackgroundBuff()

		if pbg != nil {
			background = pbg[keyIndex]
		}
	}

	if background == nil {
		cbg := dev.config.GetKeyGridBackgroundBuff()

		if cbg != nil {
			background = cbg[keyIndex]
		}
	}

	if dev.keyBGBuffs[keyIndex] != background {
		dev.keyBGBuffs[keyIndex] = background
		dev.keyUpdateChan <- keyIndex
	}
}

func (dev *VirtualDev) RedrawKey(keyIndex int) {
	dev.keyUpdateChan <- keyIndex
}

func (dev *VirtualDev) SetKeyForeground(img image.Image, keyIndex int, page int) {
	if dev.pageManager.GetPage() != page {
		return
	}

	bounds := img.Bounds().Max
	if bounds.X != dev.sdInfo.IconSize || bounds.Y != dev.sdInfo.IconSize {
		img = api.ResizeImage(img, dev.sdInfo.IconSize)
	}

	if dev.keyFGBuffs[keyIndex] != img {
		dev.keyFGBuffs[keyIndex] = img
		dev.keyUpdateChan <- keyIndex
	}
}

func (dev *VirtualDev) SetPanelBackground(knobIndex int, page int) {
	var background image.Image
	var knobV3 *api.KnobV3
	var knobConfigV3 *api.KnobConfigV3
	if len(dev.config.Pages[page].Knobs) > knobIndex {
		knobV3 = &dev.config.Pages[page].Knobs[knobIndex]
		kcv3, ok := knobV3.Application[knobV3.ActiveApplication]
		if ok {
			knobConfigV3 = kcv3
		}
	}

	if knobConfigV3 != nil {
		kcbg := knobConfigV3.GetTouchPanelBackgroundBuff()

		if kcbg != nil {
			background = kcbg
		}

	}

	if background == nil && knobV3 != nil {
		kbg := knobV3.GetTouchPanelBackgroundBuff()
		if kbg != nil {
			background = kbg
		}
	}

	if background == nil {
		pbg := dev.config.Pages[page].GetTouchPanelBackgroundBuff()
		if pbg != nil {
			background = pbg[knobIndex]
		}
	}

	if background == nil {
		cbg := dev.config.GetTouchPanelBackgroundBuff()
		if cbg != nil {
			background = cbg[knobIndex]
		}
	}

	if dev.panelBGBuffs[knobIndex] != background {
		dev.panelBGBuffs[knobIndex] = background
		dev.knobUpdateChan <- knobIndex
	}
}

func (dev *VirtualDev) SetPanelForeground(img image.Image, knobIndex int, page int) {
	if dev.pageManager.GetPage() != page {
		return
	}

	bounds := img.Bounds().Max
	if bounds.X != dev.sdInfo.LcdWidth || bounds.Y != dev.sdInfo.LcdHeight {
		img = api.ResizeImageWH(img, dev.sdInfo.LcdWidth, dev.sdInfo.LcdHeight)
	}

	if dev.panelFGBuffs[knobIndex] != img {
		dev.panelFGBuffs[knobIndex] = img
		dev.knobUpdateChan <- knobIndex
	}
}

func (dev *VirtualDev) SetBrightness(brightness uint8) error {
	return dev.deck.SetBrightness(brightness)
}

func (dev *VirtualDev) setSdInfo() {

	manufacturer, err := dev.deck.Device.GetManufacturer()
	if err != nil {
		dev.logger.Println(err)
	}
	product, err := dev.deck.Device.GetProduct()
	if err != nil {
		dev.logger.Println(err)
	}

	info := api.StreamDeckInfoV1{
		Cols:                    int(dev.deck.Columns),
		Rows:                    int(dev.deck.Rows),
		IconSize:                int(dev.deck.Pixels),
		Page:                    0,
		Serial:                  dev.deck.Serial,
		Name:                    manufacturer + " " + product,
		Connected:               true,
		LastConnected:           time.Now(),
		LcdWidth:                int(dev.deck.LcdWidth),
		LcdHeight:               int(dev.deck.LcdHeight),
		LcdCols:                 int(dev.deck.LcdColumns),
		KnobCols:                int(dev.deck.Knobs),
		PaddingX:                int(dev.deck.PaddingX),
		PaddingY:                int(dev.deck.PaddingY),
		KeyGridBackgroundWidth:  int((dev.deck.Pixels * uint(dev.deck.Columns)) + (dev.deck.PaddingX * uint(dev.deck.Columns-1))),
		KeyGridBackgroundHeight: int((dev.deck.Pixels * uint(dev.deck.Rows)) + (dev.deck.PaddingX * uint(dev.deck.Rows-1))),
		LcdBackgroundWidth:      int(dev.deck.LcdWidth * uint(dev.deck.LcdColumns)),
		LcdBackgroundHeight:     int(dev.deck.LcdHeight),
	}

	dev.sdInfo = &info
}

func (dev *VirtualDev) HandleScreenLockChange(locked bool) {
	if locked {
		dev.handlerPruner.StopAllHandlers()
		dev.deck.Reset()
	} else {
		dev.pageManager.SetPage(dev.pageManager.GetPage())
	}
}

func (dev *VirtualDev) render() {
	err := dev.deck.SetImage(0, image.NewRGBA(image.Rect(0, 0, dev.sdInfo.IconSize, dev.sdInfo.IconSize)))
	if err != nil {
		dev.logger.Println(err)
		return
	}

	go dev.renderKey()

	go dev.renderKnob()
}

func (dev *VirtualDev) renderKey() {
	for dev.isOpen && !dev.shuttingDown {

		keyIndex := <-dev.keyUpdateChan

		mergedImage, err := api.LayerImages(dev.sdInfo.IconSize, dev.sdInfo.IconSize, dev.keyBGBuffs[keyIndex], dev.keyFGBuffs[keyIndex], dev.roundedCorners)

		if err != nil {
			dev.keyUpdateChan <- keyIndex
			dev.logger.Println(err)
			continue
		}

		if dev.inputManager.GetKeyState(keyIndex) {

			bg := image.NewRGBA(image.Rect(0, 0, dev.sdInfo.IconSize, dev.sdInfo.IconSize))

			mergedImage = api.ResizeImage(mergedImage, int(float64(dev.sdInfo.IconSize)*.9))

			draw.Copy(bg, image.Pt(int(float64(dev.sdInfo.IconSize)*.05), int(float64(dev.sdInfo.IconSize)*.05)), mergedImage, mergedImage.Bounds(), draw.Over, &draw.Options{})

			mergedImage = bg
		}

		bounds := mergedImage.Bounds().Max

		dev.mu.Lock()

		err = dev.deck.SetImage(uint8(keyIndex), mergedImage)

		dev.mu.Unlock()

		if err != nil {
			dev.keyUpdateChan <- keyIndex
			match, _ := regexp.MatchString(`.*hidapi.*`, err.Error())
			if match {
				dev.disconnect()
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

func (dev *VirtualDev) renderKnob() {
	for dev.isOpen && !dev.shuttingDown {

		knobIndex := <-dev.knobUpdateChan

		mergedImage, err := api.LayerImages(dev.sdInfo.LcdWidth, dev.sdInfo.LcdHeight, dev.panelBGBuffs[knobIndex], dev.panelFGBuffs[knobIndex])

		if err != nil {
			dev.knobUpdateChan <- knobIndex
			dev.logger.Println(err)
			continue
		}

		bounds := mergedImage.Bounds().Max

		dev.mu.Lock()

		err = dev.deck.SetLcdImage(knobIndex, mergedImage)

		dev.mu.Unlock()

		if err != nil {
			dev.knobUpdateChan <- knobIndex
			match, _ := regexp.MatchString(`.*hidapi.*`, err.Error())
			if match {
				dev.disconnect()
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

func (dev *VirtualDev) handleInput() {
	defer func() {
		if err := recover(); err != nil {
			dev.disconnect()
		}
	}()
	dev.deck.HandleInput(func(event streamdeck.InputEvent) {
		if !locked {
			if event.EventType == streamdeck.KEY_PRESS || event.EventType == streamdeck.KEY_RELEASE {
				page := dev.config.Pages[dev.pageManager.GetPage()]
				if uint8(len(page.Keys)) > event.Index {
					dev.inputManager.HandleKeyInput(&page.Keys[event.Index], event)
				}
			} else if event.EventType == streamdeck.SCREEN_SWIPE {
				if event.ScreenEndX < event.ScreenX {
					if dev.pageManager.GetPage() < len(dev.config.Pages)-1 {
						dev.pageManager.SetPage(dev.pageManager.GetPage() + 1)
					}
				} else {
					if dev.pageManager.GetPage() > 0 {
						dev.pageManager.SetPage(dev.pageManager.GetPage() - 1)
					}
				}
			} else if dev.deck.HasLCD && dev.deck.HasKnobs {
				page := dev.config.Pages[dev.pageManager.GetPage()]
				if uint8(len(page.Knobs)) > event.Index {
					dev.inputManager.HandleKnobInput(&page.Knobs[event.Index], event)
				}
			}
		}
	})
}

func (dev *VirtualDev) disconnect() {
	disconnectSem.Lock()
	defer disconnectSem.Unlock()
	if !dev.isOpen {
		return
	}
	dev.logger.Println("Device (" + dev.deck.Serial + ") disconnected")
	err := dev.deck.Close()
	if err != nil {
		dev.logger.Println(err)
	}
	dev.isOpen = false
	dev.sdInfo.Connected = false
	dev.sdInfo.LastDisconnected = time.Now()
	dev.handlerPruner.StopAllHandlers()
}

func (dev *VirtualDev) Close() {
	dev.shuttingDown = true
	if dev.isOpen {
		err := dev.deck.Reset()
		if err != nil {
			dev.logger.Println(err)
		}
		dev.disconnect()
	}
}
