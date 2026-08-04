//go:build linux

package streamdeckd

import (
	"bufio"
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/bendahl/uinput"
	"github.com/godbus/dbus/v5"
	"golang.org/x/sys/unix"
)

var kb uinput.Keyboard

func UpdateApplication() {

	desktopSession, found := unix.Getenv("XDG_CURRENT_DESKTOP")
	if !found {
		log.Println("[WARN] Could not find XDG_CURRENT_DESKTOP env var, application based configuration is unavailable")
	}
	// I get Wayland has security concerns at the core of it, but having to do stuff like this to get metadata about the focused application is insane
	errCount := 0
	for {
		var err error
		switch desktopSession {
		case "Hyprland":
			err = updateHyprlandApplication()
		case "KDE":
			err = updateKDEApplication()
		case "niri":
			err = updateNiriApplication()
		default:
			// if your favourite desktop session is not here, feel free to submit a PR, I just don't have time to test on all sessions
			log.Println("[WARN] Desktop session:", desktopSession, "is unsupported, application based configuration is unavailable")
			return
		}
		if err != nil {
			log.Println(err)
			if errCount >= 10 {
				log.Println("[WARN] Repeated error getting " + desktopSession + " active window, application based configuration is unavailable until streamdeckd restart")
				return
			}
			errCount++
		} else {
			errCount = 0
		}
		time.Sleep(time.Duration((errCount+1)*50) * time.Millisecond)
	}

}

func setApplication(activePs string) {
	activePs = strings.Trim(activePs, "\n")
	applicationManager.SetApplication(activePs)
}

type HyprlandActiveWindow struct {
	Class string `json:"class"`
}

func updateHyprlandApplication() error {
	out, err := exec.Command("/bin/sh", "-c", "hyprctl activewindow -j").Output()
	if err != nil {
		return err
	}
	var activeWindow HyprlandActiveWindow
	err = json.Unmarshal(out, &activeWindow)
	if err != nil {
		return err
	}
	if activeWindow.Class != "" {
		setApplication(activeWindow.Class)
	}
	return nil
}

type NiriActiveWindow struct {
	AppId string `json:"app_id"`
}

type NiriEvent struct {
	WindowFocusChanged struct {
		Id *int `json:"id,omitempty"`
	} `json:"WindowFocusChanged"`
}

func updateNiriApplication() error {
	setWindow := func() error {
		out, err := exec.Command("/bin/sh", "-c", "niri msg -j focused-window").Output()
		if err != nil {
			return err
		}
		var activeWindow NiriActiveWindow
		err = json.Unmarshal(out, &activeWindow)
		if err != nil {
			return err
		}
		if activeWindow.AppId != "" {
			setApplication(activeWindow.AppId)
		}
		return nil
	}

	err := setWindow()
	if err != nil {
		return err
	}

	cmd := exec.Command("/bin/sh", "-c", "niri msg -j event-stream")
	outputPipe, _ := cmd.StdoutPipe()
	done := make(chan error)

	go func() {
		scanner := bufio.NewScanner(outputPipe)
		for scanner.Scan() {
			line := scanner.Text()
			var focusedWindow NiriEvent
			err := json.Unmarshal([]byte(line), &focusedWindow)
			if err != nil {
				done <- err
			}
			if focusedWindow.WindowFocusChanged.Id != nil {
				err = setWindow()
				if err != nil {
					done <- err
				}
			}
		}
		done <- nil
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = <-done
	if err != nil {
		return err
	}

	defer func(cmd *exec.Cmd) {
		_ = cmd.Wait()
	}(cmd)

	return nil
}

func updateKDEApplication() error {
	out, err := exec.Command("/bin/sh", "-c", "kdotool getwindowclassname $(kdotool getactivewindow)").Output()
	if err != nil {
		return err
	}
	outString := string(out)
	outString = strings.Trim(outString, " \n")
	if outString != "" {
		setApplication(outString)
	}
	return nil
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
				if deck.IsOpen() {
					deck.HandleScreenLockChange(locked)
				}
			}
		}
	}
}
