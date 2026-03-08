//go:build darwin

package streamdeckd

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#import <ApplicationServices/ApplicationServices.h>

CGEventRef CreateKeyEvent(int key, bool down) {
	CGEventRef event = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)key, down);
	return event;
}

void SendEvent(CGEventRef event) {
	CGEventPost(kCGAnnotatedSessionEventTap, event);
}

void ReleaseEvent(CGEventRef event) {
    CFRelease(event);
}

void AddActionKey(int type,CGEventRef event){
	CGEventSetFlags(event, (CGEventFlags) type);
}
*/
import "C"
import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/unix-streamdeck/api/v2"
)

var controlFlag = C.kCGEventFlagMaskControl
var shiftFlag = C.kCGEventFlagMaskShift
var commandFlag = C.kCGEventFlagMaskCommand

func UpdateApplication() {
	log.Println("Application based configuration is not currently supported on macOS, due to the difficulty in getting the current active application, sorry :(")
}

func EnableVirtualKeyboard() {
	log.Println("")
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

	control := strings.Contains(keybind, "ctrl")
	command := strings.Contains(keybind, "command")
	shift := strings.Contains(keybind, "shift")

	for _, key := range keys {
		downEvent := C.CreateKeyEvent(C.int(key), true)
		if command {
			C.AddActionKey(C.int(commandFlag), downEvent)
		}
		if control {
			C.AddActionKey(C.int(controlFlag), downEvent)
		}
		if shift {
			C.AddActionKey(C.int(shiftFlag), downEvent)
		}
		C.SendEvent(downEvent)
		C.ReleaseEvent(downEvent)
	}

	for i := len(keys) - 1; i >= 0; i-- {
		upEvent := C.CreateKeyEvent(C.int(keys[i]), false)
		C.SendEvent(upEvent)
		C.ReleaseEvent(upEvent)
	}

	return nil
}
