package main

import (
	"encoding/json"
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"log"
	"os"
)

const intro = `
<node>
	<interface name="com.thejonsey.streamdeckd">
		<method name="GetDeckInfo">
			<arg direction="out" type="s"/>
		</method>
		<method name="GetConfig">
			<arg direction="out" type="s"/>
		</method>
		<method name="ReloadConfig">
		</method>
		<method name="SetPage">
			<arg direction="in" type"i"/>
		</method>
		<method name="SetConfig">
			<arg direction="in" type="s"/>
		</method>
		<method name="CommitConfig">
		</method>
	</interface>` + introspect.IntrospectDataString + `</node> `

type StreamDeckDBus struct {
	Cols int `json:"cols,omitempty"`
	Rows int `json:"rows,omitempty"`
	IconSize uint `json:"icon_size,omitempty"`
}

func (s StreamDeckDBus) GetDeckInfo() (string, *dbus.Error) {
	infoString, err := json.Marshal(s)
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

func (StreamDeckDBus) SetPage(page int) *dbus.Error  {
	SetPage(config, page, dev)
	return nil
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

func InitDBUS() {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	s := StreamDeckDBus{
		Cols: int(dev.Columns),
		Rows: int(dev.Rows),
		IconSize: dev.Pixels,
	}
	conn.Export(s, "/com/thejonsey/streamdeckd", "com.thejonsey.streamdeckd")
	conn.Export(introspect.Introspectable(intro), "/com/thejonsey/streamdeckd",
		"org.freedesktop.DBus.Introspectable")
	reply, err := conn.RequestName("com.thejonsey.streamdeckd",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}
		select {}
}