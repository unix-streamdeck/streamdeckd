package streamdeckd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/unix-streamdeck/api/v3"
	streamdeck "github.com/unix-streamdeck/driver"
	"github.com/unix-streamdeck/streamdeckd/streamdeckd/handlers/basic"
	"github.com/unix-streamdeck/streamdeckd/streamdeckd/handlers/examples"
)

var configPath string
var basicConfig = api.Config{
	Modules: []string{},
	Decks: []api.Deck{
		{},
	},
}
var config *api.Config

var configSem sync.Mutex

func LoadConfig() {
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
	RegisterModuleGroup(basic.GetModuleGroup())
	RegisterModuleGroup(examples.GetModuleGroup())
	if len(config.Modules) > 0 {
		for _, module := range config.Modules {
			LoadModule(module)
		}
	}
	tryConnectObs()
}

func readConfig() (*api.Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return &api.Config{}, err
	}
	var config api.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln("Could not parse config, shutting down", err)
		return &api.Config{}, err
	}
	return &config, nil
}

func SetConfig(configString string) error {
	configSem.Lock()
	defer configSem.Unlock()
	UnmountHandlers()
	var err error
	config = nil
	err = json.Unmarshal([]byte(configString), &config)
	if err != nil {
		return err
	}
	for s := range Devs {
		dev := Devs[s]
		for i := range config.Decks {
			if dev.Serial() == config.Decks[i].Serial {
				dev.SetConfig(config.Decks[i])
			}
		}
	}
	return nil
}

func ReloadConfig() error {
	configSem.Lock()
	defer configSem.Unlock()
	UnmountHandlers()
	LoadConfig()
	for s := range Devs {
		dev := Devs[s]
		for i := range config.Decks {
			if dev.Serial() == config.Decks[i].Serial {
				dev.SetConfig(config.Decks[i])
			}
		}
	}
	return nil
}

func SaveConfig() error {
	configSem.Lock()
	defer configSem.Unlock()
	return SaveFile(configPath, config)
}

func SaveFile(path string, value any) error {
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	var configString []byte
	configString, err = json.Marshal(value)
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

func SetConfigPath(path string) {
	if path != "" {
		configPath = path
	} else {
		basePath := os.Getenv("HOME") + string(os.PathSeparator) + ".config"
		if os.Getenv("XDG_CONFIG_HOME") != "" {
			basePath = os.Getenv("XDG_CONFIG_HOME")
		}
		configPath = basePath + string(os.PathSeparator) + ".streamdeck-config.json"
	}
}

func findConfig(device *streamdeck.Device) api.Deck {
	for _, deck := range config.Decks {
		if deck.Serial == device.Serial {
			return deck
		}
	}

	return makeEmptyDeckConfig(device)
}

func makeEmptyDeckConfig(device *streamdeck.Device) api.Deck {
	var pages []api.Page
	pages = append(pages, makeEmptyPageConfig(device))
	devConf := api.Deck{Serial: device.Serial, Pages: pages}
	config.Decks = append(config.Decks, devConf)
	_ = SaveConfig()
	return devConf
}

func makeEmptyPageConfig(device *streamdeck.Device) api.Page {
	page := api.Page{}
	for i := 0; i < int(device.Rows)*int(device.Columns); i++ {
		applications := make(map[string]*api.KeyConfig)
		applications[""] = &api.KeyConfig{}
		page.Keys = append(page.Keys, api.Key{
			Application: applications,
		})
	}
	for i := 0; i < int(device.Knobs); i++ {
		applications := make(map[string]*api.KnobConfig)
		applications[""] = &api.KnobConfig{}
		page.Knobs = append(page.Knobs, api.Knob{
			Application: applications,
		})
	}
	return page
}
