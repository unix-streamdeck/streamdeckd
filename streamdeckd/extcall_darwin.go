//go:build darwin
// +build darwin

package streamdeckd

import (
	"errors"
	"log"

	"github.com/bendahl/uinput"
	"github.com/godbus/dbus/v5"
)

var kb uinput.Keyboard

func UpdateApplication() {
	log.Println("Application based configuration is not currently supported on macOS, due to the difficulty in getting the current active application, sorry :(")
}

func EnableVirtualKeyboard() {
	defer HandlePanic(func() {
		log.Println("VirtualKeyboard crash")
	})
	var err error
	kb, err = uinput.CreateKeyboard("/dev/uinput", []byte("streamdeckd"))
	if err != nil {
		log.Println(err)
	}
	defer kb.Close()
	select {}
}

func reInitDBus() {
	log.Println("Restarting DBUS")
	go InitDBUS()
}

func InitDBUS() {
	log.Println("DBUS Started")
	var err error
	conn, err = dbus.SessionBus()
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	defer HandlePanic(reInitDBus)

	sDbus = &StreamDeckDBus{}
	conn.ExportAll(sDbus, "/com/unixstreamdeck/streamdeckd", "com.unixstreamdeck.streamdeckd")
	reply, err := conn.RequestName("com.unixstreamdeck.streamdeckd",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("DBUS Started 2")
	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Println(errors.New("DBus: Name already taken"))
		return
	}
	select {}
}

func ConnectScreensaver() (*ScreensaverConnection, error) {
	return &ScreensaverConnection{}, nil
}

func (c *ScreensaverConnection) RegisterScreensaverActiveListener() {

}
