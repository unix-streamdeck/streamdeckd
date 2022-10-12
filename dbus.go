package main

import (
	"encoding/json"
	"errors"
	"github.com/godbus/dbus/v5"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/streamdeckd/handlers"
	"log"
)

var conn *dbus.Conn

var sDbus *StreamDeckDBus
var sDInfo []api.StreamDeckInfo

type StreamDeckDBus struct {
}

func (s StreamDeckDBus) GetDeckInfo() (string, *dbus.Error) {
	infoString, err := json.Marshal(sDInfo)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	return string(infoString), nil
}

func (StreamDeckDBus) GetConfig() (string, *dbus.Error) {
	configString, err := json.Marshal(config)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	return string(configString), nil
}

func (StreamDeckDBus) ReloadConfig() *dbus.Error {
	err := ReloadConfig()
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (StreamDeckDBus) SetPage(serial string, page int) *dbus.Error {
	for s := range devs {
		if devs[s].Deck.Serial == serial {
			dev := devs[s]
			dev.SetPage(page)
			return nil
		}
	}
	return dbus.MakeFailedError(errors.New("Device with Serial: " + serial + " could not be found"))
}

func (StreamDeckDBus) SetConfig(configString string) *dbus.Error {
	err := SetConfig(configString)
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (StreamDeckDBus) CommitConfig() *dbus.Error {
	err := SaveConfig()
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (StreamDeckDBus) GetModules() (string, *dbus.Error) {
	var modules []api.Module
	for _, module := range handlers.AvailableModules() {
		modules = append(modules, api.Module{Name: module.Name, IconFields: module.IconFields, KeyFields: module.KeyFields, IsIcon: module.NewIcon != nil, IsKey: module.NewKey != nil})
	}
	modulesString, err := json.Marshal(modules)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	return string(modulesString), nil
}

func (StreamDeckDBus) PressButton(serial string, keyIndex int) *dbus.Error {
	dev, ok := devs[serial]
	if !ok || !dev.IsOpen{
		return dbus.MakeFailedError(errors.New("Can't find connected device: " + serial))
	}
	HandleInput(dev, &dev.Config[dev.Page][keyIndex], dev.Page)
	return nil
}

func InitDBUS() error {
	var err error
	conn, err = dbus.SessionBus()
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close()

	sDbus = &StreamDeckDBus{}
	conn.ExportAll(sDbus, "/com/unixstreamdeck/streamdeckd", "com.unixstreamdeck.streamdeckd")
	reply, err := conn.RequestName("com.unixstreamdeck.streamdeckd",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Println(err)
		return err
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return errors.New("DBus: Name already taken")
	}
	select {}
}

func EmitPage(dev *VirtualDev, page int) {
	if conn != nil {
		conn.Emit("/com/unixstreamdeck/streamdeckd", "com.unixstreamdeck.streamdeckd.Page", dev.Deck.Serial, page)
	}
	for i := range sDInfo {
		if sDInfo[i].Serial == dev.Deck.Serial {
			sDInfo[i].Page = page
		}
	}
}
