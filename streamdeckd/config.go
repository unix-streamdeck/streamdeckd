package streamdeckd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/unix-streamdeck/api/v2"
)

var configPath string
var basicConfig = api.ConfigV3{
	Modules: []string{},
	Decks: []api.DeckV3{
		{},
	},
}
var config *api.ConfigV3
var migrateConfigFromV1 = false

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
	if len(config.Modules) > 0 {
		for _, module := range config.Modules {
			LoadModule(module)
		}
	}
	tryConnectObs()
}

func readConfig() (*api.ConfigV3, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return &api.ConfigV3{}, err
	}
	var config api.ConfigV3
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln("Could not parse config, shutting down", err)
		return &api.ConfigV3{}, err
	}
	return &config, nil
}

func SetConfig(configString string) error {
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
			if dev.Deck.Serial == config.Decks[i].Serial {
				dev.Config = config.Decks[i]
			}
		}
		dev.SetPage(Devs[s].Page)
	}
	return nil
}

func ReloadConfig() error {
	UnmountHandlers()
	LoadConfig()
	for s := range Devs {
		dev := Devs[s]
		for i := range config.Decks {
			if dev.Deck.Serial == config.Decks[i].Serial {
				dev.Config = config.Decks[i]
			}
		}
		dev.SetPage(Devs[s].Page)
	}
	return nil
}

func SaveConfig() error {
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
