//go:build linux

package streamdeckd

import (
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/bendahl/uinput"
	"github.com/godbus/dbus/v5"
	x "github.com/linuxdeepin/go-x11-client"
	"golang.org/x/sys/unix"
)

var c *x.Conn

var kb uinput.Keyboard

func UpdateApplication() {
	sessionType, found := unix.Getenv("XDG_SESSION_TYPE")

	if found && sessionType == "x11" {
		updateX11Application()
		return
	}

	desktopSession, found := unix.Getenv("XDG_CURRENT_DESKTOP")
	if !found {
		log.Println("[WARN] Could not find XDG_CURRENT_DESKTOP env var, application based configuration is unavailable")
	}
	// I get Wayland has security concerns at the core of it, but having to do stuff like this to get metadata about the focused application is insane
	switch desktopSession {
	case "Hyprland":
		updateHyprlandApplication()
	case "KDE":
		updateKDEApplication()
	default:
		// if your favourite desktop session is not here, feel free to submit a PR, I just don't have time to test on all sessions
		log.Println("[WARN] Desktop session:", desktopSession, "is unsupported, application based configuration is unavailable")
	}
}

func setApplication(activePs string) {
	activePs = strings.Trim(activePs, "\n")
	if currentApplication != activePs {
		currentApplication = activePs
		log.Println("Application updated to: " + currentApplication)
		for _, dev := range Devs {
			dev.ApplicationUpdated()
		}
	}
}

type HyprlandActiveWindow struct {
	Class string `json:"class"`
}

func updateHyprlandApplication() {
	errCounter := 0
	for {
		out, err := exec.Command("/bin/sh", "-c", "hyprctl activewindow -j").Output()
		if err != nil {
			log.Println(err)
			errCounter += 1
			if errCounter == 10 {
				log.Println("[WARN] Repeated error getting Hyprland active window, application based configuration is unavailable until streamdeckd restart")
				return
			}
		} else {
			errCounter = 0
		}
		var activeWindow HyprlandActiveWindow
		err = json.Unmarshal(out, &activeWindow)
		if err != nil {
			log.Println(err)
			errCounter += 1
			if errCounter == 10 {
				log.Println("[WARN] Repeated error getting Hyprland active window, application based configuration is unavailable until streamdeckd restart")
				return
			}
		} else {
			errCounter = 0
		}
		if activeWindow.Class != "" {
			setApplication(activeWindow.Class)
		}
		time.Sleep(time.Duration((errCounter+1)*50) * time.Millisecond)
	}
}
func updateKDEApplication() {
	errCounter := 0
	for {
		out, err := exec.Command("/bin/sh", "-c", "kdotool getwindowclassname $(kdotool getactivewindow)").Output()
		if err != nil {
			log.Println(err)
			errCounter += 1
			if errCounter == 10 {
				log.Println("[WARN] Repeated error getting KDE active window, application based configuration is unavailable until streamdeckd restart")
				return
			}
		} else {
			errCounter = 0
		}
		outString := string(out)
		outString = strings.Trim(outString, " \n")
		if outString != "" {
			setApplication(outString)
		}
		time.Sleep(time.Duration((errCounter+1)*50) * time.Millisecond)
	}
}

func updateX11Application() {
	c, _ = x.NewConn()
	defer c.Close()
	var activeWindow x.Window
	for {
		win, _ := x.GetInputFocus(c).Reply(c)
		if activeWindow != win.Focus {
			activeWindow = win.Focus
			active, err := exec.Command("xdotool", "getactivewindow", "getwindowclassname").Output()
			if err != nil {
				active = []byte("")
			}
			activePs := string(active)
			log.Println(activePs)
			activePs = strings.Trim(activePs, "\n")
			setApplication(activePs)
		}
		time.Sleep(50 * time.Millisecond)
	}
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
