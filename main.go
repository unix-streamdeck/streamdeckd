package main

import (
	"encoding/json"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/driver"
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

var dev streamdeck.Device
var config *api.Config
var configPath = os.Getenv("HOME") + "/.streamdeck-config.json"

var basicConfig = api.Config{
	Pages: []api.Page{
		{
			api.Key{
			},
		},
	},
}

func main() {
	d, err := streamdeck.Devices()
	if err != nil {
		log.Println(err)
	}
	if len(d) == 0 {
		log.Println("No Stream Deck devices found.")
	}
	dev = d[0]
	err = dev.Open()
	if err != nil {
		log.Println(err)
	}
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
	if len(config.Pages) == 0 {
		config.Pages = append(config.Pages, api.Page{})
	}
	cleanupHook()
	SetPage(config, 0, dev)
	go InitDBUS()
	Listen()
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
		_ = dev.Reset()
		os.Exit(0)
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
	if len(config.Pages) == 0 {
		config.Pages = append(config.Pages, api.Page{})
	}
	SetPage(config, p, dev)
	return nil
}

func ReloadConfig() error {
	unmountHandlers()
	var err error
	config, err = readConfig()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if len(config.Pages) == 0 {
		config.Pages = append(config.Pages, api.Page{})
	}
	SetPage(config, p, dev)
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
	for i := range config.Pages {
		page := config.Pages[i]
		for i2 := range page {
			key := page[i2]
			if key.IconHandlerStruct != nil {
				key.IconHandlerStruct.Stop()
			}
		}
	}
}