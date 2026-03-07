//go:build darwin

package streamdeckd

/*
 #cgo CFLAGS: -x objective-c
 #cgo LDFLAGS: -framework Cocoa
 #import <Cocoa/Cocoa.h>
 #import <ApplicationServices/ApplicationServices.h>
 CGEventRef CreateDown(int k){
	CGEventRef event = CGEventCreateKeyboardEvent (NULL, (CGKeyCode)k, true);
	return event;
 }
 CGEventRef CreateUp(int k){
	CGEventRef event = CGEventCreateKeyboardEvent (NULL, (CGKeyCode)k, false);
	return event;
 }
 void KeyTap(CGEventRef event){
	CGEventPost(kCGAnnotatedSessionEventTap, event);
	CFRelease(event);
 }
 void AddActionKey(CGEventFlags type,CGEventRef event){
 	CGEventSetFlags(event, type);
 }
*/
import "C"
import (
	"errors"
	"log"

	"github.com/bendahl/uinput"
	"github.com/godbus/dbus/v5"
	"github.com/unix-streamdeck/api/v2"
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

func ExecuteKeybind(keybind string) error {
	keys, err := api.ParseXDoToolKeybindString(keybind)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to parse keybind: %s", err))
	}

	for _, key := range keys {
		downEvent := C.CreateDown(C.int(key))
		C.KeyTap(downEvent)
	}

	for i := len(keys) - 1; i >= 0; i-- {
		upEvent := C.CreateUp(C.int(key))
		C.KeyTap(upEvent)
	}

	return nil
}
