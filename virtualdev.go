package main

import (
	"context"
	"errors"
	"github.com/unix-streamdeck/api"
	streamdeck "github.com/unix-streamdeck/driver"
	"golang.org/x/sync/semaphore"
	"image"
	"log"
	"strings"
)

func openDevice() (*VirtualDev, error) {
	ctx := context.Background()
	err := connectSem.Acquire(ctx, 1)
	if err != nil {
		return &VirtualDev{}, err
	}
	defer connectSem.Release(1)
	d, err := streamdeck.Devices()
	if err != nil {
		return &VirtualDev{}, err
	}
	if len(d) == 0 {
		return &VirtualDev{}, errors.New("No streamdeck devices found")
	}
	device := streamdeck.Device{Serial: ""}
	for i := range d {
		found := false
		for s := range devs {
			if d[i].ID == devs[s].Deck.ID && devs[s].IsOpen {
				found = true
				break
			} else if d[i].Serial == s && !devs[s].IsOpen {
				err = d[i].Open()
				if err != nil {
					return &VirtualDev{}, err
				}
				devs[s].Deck = d[i]
				devs[s].IsOpen = true
				return devs[s], nil
			}
		}
		if !found {
			device = d[i]
		}
	}
	if len(device.Serial) != 12 {
		return &VirtualDev{}, errors.New("No streamdeck devices found")
	}
	err = device.Open()
	if err != nil {
		return &VirtualDev{}, err
	}
	devNo := -1
	if migrateConfig {
		config.Decks[0].Serial = device.Serial
		_ = SaveConfig()
		migrateConfig = false
	}
	for i := range config.Decks {
		if config.Decks[i].Serial == device.Serial {
			devNo = i
		}
	}
	if devNo == -1 {
		var pages []api.Page
		page := api.Page{}
		for i := 0; i < int(device.Rows)*int(device.Columns); i++ {
			page = append(page, api.Key{})
		}
		pages = append(pages, page)
		config.Decks = append(config.Decks, api.Deck{Serial: device.Serial, Pages: pages})
		devNo = len(config.Decks) - 1
	}
	dev, ok := devs[device.Serial]
	if !ok {
		dev = &VirtualDev{Deck: device, Page: 0, IsOpen: true, Config: config.Decks[devNo].Pages, sem: semaphore.NewWeighted(int64(1))}
		devs[device.Serial] = dev
	} else {
		dev.Deck = device
	}
	log.Println("Device (" + device.Serial + ") connected")
	return dev, nil
}

type VirtualDev struct {
	Deck   streamdeck.Device
	Page   int
	IsOpen bool
	Config []api.Page
	sem *semaphore.Weighted
	shuttingDown bool
}

func (dev *VirtualDev) Listen() {
	kch, err := dev.Deck.ReadKeys()
	if err != nil {
		log.Println(err)
	}
	for dev.IsOpen {
		select {
		case k, ok := <-kch:
			if !ok {
				disconnect(dev)
				return
			}
			if k.Pressed == true {
				if len(dev.Config)-1 >= dev.Page && len(dev.Config[dev.Page])-1 >= int(k.Index) {
					HandleInput(dev, &dev.Config[dev.Page][k.Index], dev.Page)
				}
			}
		}
	}
}

func (dev *VirtualDev) SetPage(page int) {
	if page != dev.Page {
		unmountPageHandlers(dev.Config[dev.Page])
	}
	dev.Page = page
	currentPage := dev.Config[page]
	for i := 0; i < len(currentPage); i++ {
		currentKey := &currentPage[i]
		go SetKey(dev, currentKey, i, page)
	}
	EmitPage(dev, page)
}

func (dev *VirtualDev) SetImage(img image.Image, keyIndex int, page int) {
	ctx := context.Background()
	err := dev.sem.Acquire(ctx, 1)
	if err != nil {
		log.Println(err)
		return
	}
	defer dev.sem.Release(1)

	if dev.Page == page && dev.IsOpen && !dev.shuttingDown {
		err := dev.Deck.SetImage(uint8(keyIndex), img)
		if err != nil {
			if strings.Contains(err.Error(), "hidapi") {
				disconnect(dev)
			} else if strings.Contains(err.Error(), "dimensions") {
				log.Println(err)
			}else {
				log.Println(err)
			}
		}
	}
}

func (dev *VirtualDev) UnmountHandlers() {
	for i := range dev.Config {
		unmountPageHandlers(dev.Config[i])
	}
}

func (dev *VirtualDev) SetBrightness(brightness uint8) error {
	return dev.Deck.SetBrightness(brightness)
}

func (dev *VirtualDev) AppendSdInfo() {
	found := false
	for i := range sDInfo {
		if sDInfo[i].Serial == dev.Deck.Serial {
			found = true
		}
	}
	if !found {
		sDInfo = append(sDInfo, api.StreamDeckInfo{
			Cols:     int(dev.Deck.Columns),
			Rows:     int(dev.Deck.Rows),
			IconSize: int(dev.Deck.Pixels),
			Page:     0,
			Serial:   dev.Deck.Serial,
		})
	}
}