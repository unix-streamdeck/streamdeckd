package streamdeckd

import (
	"encoding/json"
	"github.com/unix-streamdeck/api/v2"
	"io/ioutil"
	"log"
	"os"
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
	if err != nil || config.Decks == nil || allKeysEmpty(config) {
		var config1 api.ConfigV1
		var config2 api.ConfigV2
		err = json.Unmarshal(data, &config1)
		if err == nil && config1.Pages != nil {
			SaveFile(configPath+".v1", config1)
			config2 = api.ConfigV2{Modules: config1.Modules, Decks: []api.DeckV2{{Pages: config1.Pages, Serial: ""}}}
			migrateConfigFromV1 = true
		} else {
			err = json.Unmarshal(data, &config2)
			if err != nil {
				log.Fatalln("Could not parse config, shutting down", err)
				return &api.ConfigV3{}, err
			}
		}
		SaveFile(configPath+".v2", config2)
		var decksV3 []api.DeckV3
		for _, deck := range config2.Decks {
			var pagesV3 []api.PageV3
			deckV3 := api.DeckV3{Serial: deck.Serial}
			for _, page := range deck.Pages {
				pageV3 := api.PageV3{}
				for _, key := range page {
					defaultMap := make(map[string]*api.KeyConfigV3)
					defaultMap[""] = &api.KeyConfigV3{
						Icon:              key.Icon,
						SwitchPage:        key.SwitchPage,
						Text:              key.Text,
						TextSize:          key.TextSize,
						TextAlignment:     key.TextAlignment,
						Keybind:           key.Keybind,
						Command:           key.Command,
						Brightness:        key.Brightness,
						Url:               key.Url,
						ObsCommand:        key.ObsCommand,
						ObsCommandParams:  key.ObsCommandParams,
						IconHandler:       key.IconHandler,
						KeyHandler:        key.KeyHandler,
						IconHandlerFields: key.IconHandlerFields,
						KeyHandlerFields:  key.KeyHandlerFields,
					}
					pageV3.Keys = append(pageV3.Keys, api.KeyV3{
						Application: defaultMap,
					})
				}
				pagesV3 = append(pagesV3, pageV3)
			}
			deckV3.Pages = pagesV3
			decksV3 = append(decksV3, deckV3)
		}
		config = api.ConfigV3{
			Modules:           config2.Modules,
			Decks:             decksV3,
			ObsConnectionInfo: config2.ObsConnectionInfo,
		}
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
				dev.Config = config.Decks[i].Pages
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
				dev.Config = config.Decks[i].Pages
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

func allKeysEmpty(config api.ConfigV3) bool {
	for _, deck := range config.Decks {
		for _, page := range deck.Pages {
			for _, key := range page.Keys {
				if !api.CompareKeys(key, api.KeyV3{}) {
					return false
				}
			}
		}
	}
	return true
}
