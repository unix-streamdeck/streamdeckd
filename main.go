package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/driver"
	"github.com/unix-streamdeck/streamdeckd/handlers"
	"github.com/unix-streamdeck/streamdeckd/handlers/examples"
	"golang.org/x/sync/semaphore"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type VirtualDev struct {
	Deck streamdeck.Device
	Page int
	IsOpen bool
	Config []api.Page
}

var devs map[string]*VirtualDev
var config *api.Config
var configPath = os.Getenv("HOME") + string(os.PathSeparator) + ".streamdeck-config.json"
var disconnectSem = semaphore.NewWeighted(1)
var connectSem = semaphore.NewWeighted(1)
var basicConfig = api.Config{
	Modules: []string{},
	Decks: []api.Deck{
		{
		},
	},
}
var isRunning = true

func main() {
	configPtr := flag.String("config", configPath, "Path to config file")
	flag.Parse()
	if *configPtr != "" {
		configPath = *configPtr
	}
	cleanupHook()
	go InitDBUS()
	examples.RegisterBaseModules()
	loadConfig()
	devs = make(map[string]*VirtualDev)
	attemptConnection()
}

func attemptConnection() {
	for isRunning {
		dev := &VirtualDev{}
		dev, _ = openDevice()
		if dev.IsOpen {
			SetPage(dev, 0)
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
			go Listen(dev)
		}
	}
}

func disconnect(dev *VirtualDev) {
	ctx := context.Background()
	err := disconnectSem.Acquire(ctx, 1)
	if err != nil {
		return
	}
	defer disconnectSem.Release(1)
	if !dev.IsOpen {
		return
	}
	log.Println("Device (" + dev.Deck.Serial + ") disconnected")
	_ = dev.Deck.Close()
	dev.IsOpen = false
	unmountDevHandlers(dev)
}

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
		devNo = len(config.Decks) -1
	}
	dev := &VirtualDev{Deck: device, Page: 0, IsOpen: true, Config: config.Decks[devNo].Pages}
	devs[device.Serial] = dev
	log.Println("Device (" + device.Serial + ") connected")
	return dev, nil
}

func loadConfig() {
	var err error
	config, err = readConfig()
	if err != nil && !os.IsNotExist(err) {
		log.Println(err)
	} else if os.IsNotExist(err) {
		file, err := os.Create(configPath)
		if err != nil {
			log.Println(err)
		}
		err = file.Close()
		if err != nil {
			log.Println(err)
		}
		config = &basicConfig
		err = SaveConfig()
		if err != nil {
			log.Println(err)
		}
	}
	if len(config.Modules) > 0 {
		for _, module := range config.Modules {
			handlers.LoadModule(module)
		}
	}
}

func readConfig() (*api.Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return &api.Config{}, err
	}
	var config api.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return &api.Config{}, err
	}
	return &config, nil
}

func runCommand(command string) {
	//args := strings.Split(command, " ")
	c := exec.Command("/bin/sh", "-c", command)
	if err := c.Start(); err != nil {
		log.Println(err)
	}
	err := c.Wait()
	if err != nil {
		log.Printf("command failed: %s", err)
	}
}

func cleanupHook() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		<-sigs
		isRunning = false
		unmountHandlers()
		var err error
		for s := range devs {
			if devs[s].IsOpen {
				err = devs[s].Deck.Reset()
				if err != nil {
					log.Println(err)
				}
				err = devs[s].Deck.Close()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}

func SetConfig(configString string) error {
	unmountHandlers()
	var err error
	config = nil
	err = json.Unmarshal([]byte(configString), &config)
	if err != nil {
		return err
	}
	for s := range devs {
		dev := devs[s]
		for i := range config.Decks {
			if dev.Deck.Serial == config.Decks[i].Serial {
				dev.Config = config.Decks[i].Pages
			}
		}
		SetPage(dev, devs[s].Page)
	}
	return nil
}

func ReloadConfig() error {
	unmountHandlers()
	loadConfig()
	for s := range devs {
		dev := devs[s]
		for i := range config.Decks {
			if dev.Deck.Serial == config.Decks[i].Serial {
				dev.Config = config.Decks[i].Pages
			}
		}
		SetPage(dev, devs[s].Page)
	}
	return nil
}

func SaveConfig() error {
	f, err := os.OpenFile(configPath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	var configString []byte
	configString, err = json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = f.Write(configString)
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}
func unmountHandlers() {
	for s := range devs {
		dev := devs[s]
		unmountDevHandlers(dev)
	}
}

func unmountDevHandlers(dev *VirtualDev) {
	for i := range dev.Config {
		page := dev.Config[i]
		for i2 := 0; i2 < len(page); i2++ {
			key := &page[i2]
			if key.IconHandlerStruct != nil {
				key.IconHandlerStruct.Stop()
			}
		}
	}
}
