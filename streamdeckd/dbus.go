package streamdeckd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/godbus/dbus/v5"
	"github.com/unix-streamdeck/api"
	"image"
	"image/png"
	"log"
	"time"
)

var conn *dbus.Conn

var sDbus *StreamDeckDBus

type StreamDeckDBus struct {
}

func (s StreamDeckDBus) GetDeckInfo() (string, *dbus.Error) {
	var decks []api.StreamDeckInfoV1
	for _, dev := range Devs {
		decks = append(decks, dev.sdInfo)
	}
	infoString, err := json.Marshal(decks)
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
	for s := range Devs {
		if Devs[s].Deck.Serial == serial {
			dev := Devs[s]
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
	for _, module := range AvailableModules() {
		modules = append(modules, api.Module{Name: module.Name, IconFields: module.IconFields, KeyFields: module.KeyFields, IsIcon: module.NewIcon != nil, IsKey: module.NewKey != nil})
	}
	modulesString, err := json.Marshal(modules)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	return string(modulesString), nil
}

func (StreamDeckDBus) PressButton(serial string, keyIndex int) *dbus.Error {
	dev, ok := Devs[serial]
	if !ok || !dev.IsOpen{
		return dbus.MakeFailedError(errors.New("Can't find connected device: " + serial))
	}
	HandleKeyInput(dev, &dev.Config[dev.Page].Keys[keyIndex])
	return nil
}

func (StreamDeckDBus) GetHandlerExample(serial string, keyString string) (string, *dbus.Error) {
	var key *api.KeyConfigV3
	err := json.Unmarshal([]byte(keyString), &key)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	if key.IconHandler == "" || key.IconHandler == "Default" {
		return "", dbus.MakeFailedError(errors.New("Invalid icon handler"))
	}
	var handler api.IconHandler
	modules := AvailableModules()
	for _, module := range modules {
		if module.Name == key.IconHandler {
			handler = module.NewIcon()
			break
		}
	}
	if handler == nil {
		return "", dbus.MakeFailedError(errors.New("Invalid icon handler"))
	}
	var dev api.StreamDeckInfoV1
	sd, ok := Devs[serial]
	if !ok {
		return "", dbus.MakeFailedError(errors.New("could not find device"))
	}
	dev = sd.sdInfo
	var img image.Image
	log.Println("Created and running " + key.IconHandler + " for dbus")
	handler.Start(*key, dev, func(image image.Image) {
		if image.Bounds().Max.X != dev.IconSize || image.Bounds().Max.Y != dev.IconSize {
			image = api.ResizeImage(image, dev.IconSize)
		}
		img = image
		log.Println("Stopping " + key.IconHandler + " for dbus")
		handler.Stop()
		log.Println("Stopped " + key.IconHandler + " for dbus")
	})
	counter := 0
	for {
		if img != nil {
			buf := new(bytes.Buffer)
			err = png.Encode(buf, img)
			imageBits := buf.Bytes()
			return "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageBits), nil
		}
		if counter >= 100 {
			return "", dbus.MakeFailedError(errors.New("Handler did not respond in a timely fashion"))
		}
		counter += 1
		time.Sleep(50 * time.Millisecond)
	}
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

func EmitPage(dev *VirtualDev, page int) {
	if conn != nil {
		conn.Emit("/com/unixstreamdeck/streamdeckd", "com.unixstreamdeck.streamdeckd.Page", dev.Deck.Serial, page)
	}
}


type ScreensaverConnection struct {
	busobj dbus.BusObject
	conn *dbus.Conn
}

func ConnectScreensaver() (*ScreensaverConnection, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}
	return &ScreensaverConnection{
		conn: conn,
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