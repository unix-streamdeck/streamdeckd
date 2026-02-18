//go:build linux
// +build linux

package streamdeckd

import (
	"errors"
	"log"

	"github.com/bendahl/uinput"
	"github.com/godbus/dbus/v5"
)

var kb uinput.Keyboard

//var c *x.Conn

//func UpdateApplication() {
//	go EnableVirtualKeyboard()
//	c, _ = x.NewConn()
//	defer c.Close()
//	var activeWindow x.Window
//	defer HandlePanic(func() {
//		log.Println("Restarting UpdateApplication")
//		go UpdateApplication()
//	})
//	for {
//		win, _ := x.GetInputFocus(c).Reply(c)
//		if activeWindow != win.Focus {
//			activeWindow = win.Focus
//			active, err := exec.Command("kdotool", "getactivewindow", "getwindowclassname").Output()
//			if err != nil {
//				active = []byte("")
//			}
//			activePs := string(active)
//			log.Println(activePs)
//			activePs = strings.Trim(activePs, "\n")
//			if currentApplication != activePs {
//				currentApplication = activePs
//				log.Println("Application updated to: " + currentApplication)
//				for _, dev := range Devs {
//					dev.ApplicationUpdated()
//				}
//			}
//		}
//		time.Sleep(50 * time.Millisecond)
//	}
//}

//func EnableVirtualKeyboard() {
//	defer HandlePanic(func() {
//		log.Println("VirtualKeyboard crash")
//	})
//	kb, _ = uinput.CreateKeyboard("/dev/uinput", []byte("streamdeckd"))
//	defer kb.Close()
//	select {}
//}

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
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}
	return &ScreensaverConnection{
		conn:   conn,
		busobj: conn.Object("org.freedesktop.ScreenSaver", "/org/freedesktop/ScreenSaver"),
	}, nil
}

func (c *ScreensaverConnection) RegisterScreensaverActiveListener() {
	defer HandlePanic(func() {
		log.Println("Restarting Screensaver Listener")
		go c.RegisterScreensaverActiveListener()
	})
	err := c.conn.AddMatchSignal(dbus.WithMatchObjectPath("/org/freedesktop/ScreenSaver"), dbus.WithMatchInterface("org.freedesktop.ScreenSaver"), dbus.WithMatchMember("ActiveChanged"))
	if err != nil {
		log.Println(err)
	}
	ch := make(chan *dbus.Signal, 10)
	c.conn.Signal(ch)
	for v := range ch {
		if locked != v.Body[0].(bool) {
			locked = v.Body[0].(bool)
			for _, deck := range Devs {
				if deck.IsOpen {
					deck.HandleScreenLockChange(locked)
				}
			}
		}
	}
}
