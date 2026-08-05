package streamdeckd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/unix-streamdeck/api/v3"
	streamdeck "github.com/unix-streamdeck/driver"
)

var conn *dbus.Conn

var sDbus IStreamDeckDBus

const intro = `<!DOCTYPE node PUBLIC "-//freedesktop//DTD D-BUS Object Introspection 1.0//EN"
 "http://www.freedesktop.org/standards/dbus/1.0/introspect.dtd">
<node name="/com/unixstreamdeck/streamdeckd">
  <interface name="com.unixstreamdeck.streamdeckd">

    <method name="GetDeckInfo">
      <arg name="info" type="s" direction="out"/>
    </method>

    <method name="GetConfig">
      <arg name="config" type="s" direction="out"/>
    </method>

    <method name="ReloadConfig">
    </method>

    <method name="SetPage">
      <arg name="serial" type="s" direction="in"/>
      <arg name="page" type="i" direction="in"/>
    </method>

    <method name="SetConfig">
      <arg name="configString" type="s" direction="in"/>
    </method>

    <method name="CommitConfig">
    </method>

    <method name="GetModules">
      <arg name="modules" type="s" direction="out"/>
    </method>

    <method name="PressButton">
      <arg name="serial" type="s" direction="in"/>
      <arg name="keyIndex" type="i" direction="in"/>
    </method>

    <method name="GetLayersExample">
      <arg name="serial" type="s" direction="in"/>
      <arg name="layersString" type="s" direction="in"/>
      <arg name="handlerType" type="y" direction="in"/>
      <arg name="image" type="s" direction="out"/>
    </method>

    <signal name="Page">
      <arg name="serial" type="s"/>
      <arg name="page" type="i"/>
    </signal>

  </interface>

  <interface name="org.freedesktop.DBus.Introspectable">
    <method name="Introspect">
      <arg name="xml_data" type="s" direction="out"/>
    </method>
  </interface>

  <interface name="org.freedesktop.DBus.Peer">
    <method name="Ping">
    </method>
    <method name="GetMachineId">
      <arg name="machine_uuid" type="s" direction="out"/>
    </method>
  </interface>` + introspect.IntrospectDataString + `</node>`

type IStreamDeckDBus interface {
	GetDeckInfo() (string, *dbus.Error)
	GetConfig() (string, *dbus.Error)
	ReloadConfig() *dbus.Error
	SetPage(serial string, page int) *dbus.Error
	SetConfig(configString string) *dbus.Error
	CommitConfig() *dbus.Error
	GetModules() (string, *dbus.Error)
	PressButton(serial string, keyIndex int) *dbus.Error
	GetLayersExample(serial string, layersString string, handlerType api.HandlerType) (string, *dbus.Error)
}

type StreamDeckDBus struct {
}

func (s StreamDeckDBus) GetDeckInfo() (string, *dbus.Error) {
	var decks []api.StreamDeckInfo
	for _, dev := range Devs {
		decks = append(decks, *dev.SdInfo())
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
		if Devs[s].Serial() == serial {
			dev := Devs[s]
			dev.PageManager().SetPage(page)
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
	modulesString, err := json.Marshal(AvailableModules())
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	return string(modulesString), nil
}

func (StreamDeckDBus) PressButton(serial string, keyIndex int) *dbus.Error {
	dev, ok := Devs[serial]
	if !ok || !dev.IsOpen() {
		return dbus.MakeFailedError(errors.New("Can't find connected device: " + serial))
	}
	dev.InputManager().HandleKeyInput(&dev.Config().Pages[dev.PageManager().GetPage()].Keys[keyIndex], streamdeck.InputEvent{
		EventType: streamdeck.KEY_PRESS,
	})
	return nil
}

var (
	exampleTimeout  = 5 * time.Second
	examplePollTick = 50 * time.Millisecond
)

func (StreamDeckDBus) GetLayersExample(serial string, layersString string, handlerType api.HandlerType) (string, *dbus.Error) {
	var layers []*api.Layer
	if err := json.Unmarshal([]byte(layersString), &layers); err != nil {
		return "", dbus.MakeFailedError(err)
	}

	dev, ok := Devs[serial]
	if !ok {
		return "", dbus.MakeFailedError(errors.New("could not find device: " + serial))
	}
	info := dev.SdInfo()
	width, height := info.GetDimensions(handlerType)

	var visual []*api.Layer
	for _, layer := range layers {
		if layer == nil || layer.Global || !layer.Visual || layer.Handler == "" {
			continue
		}
		visual = append(visual, layer)
	}
	if len(visual) == 0 {
		return "", dbus.MakeFailedError(errors.New("no visual layers to render"))
	}

	var mu sync.Mutex
	buffs := make([]image.Image, len(visual))
	handlers := make([]api.ForegroundHandler, len(visual))

	for i, layer := range visual {
		handler := newExampleHandler(layer.Handler)
		if handler == nil {
			return "", dbus.MakeFailedError(errors.New("invalid handler: " + layer.Handler))
		}
		handlers[i] = handler

		index := i
		log.Println("Created and running " + layer.Handler + " for dbus")

		go handler.Start(layer.Fields, handlerType, *info, func(img image.Image) {
			if img == nil {
				return
			}
			if img.Bounds().Dx() != width || img.Bounds().Dy() != height {
				img = api.ResizeImageWH(img, width, height)
			}

			mu.Lock()
			first := buffs[index] == nil
			if first {
				buffs[index] = img
			}
			mu.Unlock()

			if first {
				log.Println("Stopping " + visual[index].Handler + " for dbus")
				go handler.Stop()
			}
		})
	}

	defer func() {
		for i, handler := range handlers {
			mu.Lock()
			produced := buffs[i] != nil
			mu.Unlock()
			if !produced {
				go handler.Stop()
			}
		}
	}()

	deadline := time.Now().Add(exampleTimeout)
	for time.Now().Before(deadline) {
		mu.Lock()
		missing := slices.Contains(buffs, nil)
		mu.Unlock()

		if !missing {
			break
		}
		time.Sleep(examplePollTick)
	}

	mu.Lock()
	rendered := slices.Clone(buffs)
	mu.Unlock()

	img, err := api.LayerImages(width, height, rendered...)
	if err != nil {
		return "", dbus.MakeFailedError(errors.New("handlers did not respond in a timely fashion"))
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return "", dbus.MakeFailedError(err)
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func newExampleHandler(name string) api.ForegroundHandler {
	for _, module := range AvailableModules() {
		if module.Name != name || module.NewCombined == nil {
			continue
		}
		handler, ok := module.NewCombined().(api.ForegroundHandler)
		if !ok {
			return nil
		}
		return handler
	}
	return nil
}

func EmitPage(dev IVirtualDev, page int) {
	if conn != nil {
		conn.Emit("/com/unixstreamdeck/streamdeckd", "com.unixstreamdeck.streamdeckd.Page", dev.Serial(), page)
	}
}

type ScreensaverConnection struct {
	busobj dbus.BusObject
	conn   *dbus.Conn
}
