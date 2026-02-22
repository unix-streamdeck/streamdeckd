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
	//	kb, _ = uinput.CreateKeyboard("/dev/uinput", []byte("streamdeckd"))
	//	defer kb.Close()
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
