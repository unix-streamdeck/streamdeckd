package streamdeckd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"log"
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

    <method name="SendKnobEvent">
      <arg name="serial" type="s" direction="in"/>
      <arg name="knobIndex" type="i" direction="in"/>
      <arg name="event" type="i" direction="in"/>
    </method>

    <method name="GetLayersExample">
      <arg name="serial" type="s" direction="in"/>
      <arg name="layersString" type="s" direction="in"/>
      <arg name="handlerType" type="y" direction="in"/>
      <arg name="image" type="s" direction="out"/>
    </method>

    <method name="GetBackgroundExample">
      <arg name="serial" type="s" direction="in"/>
      <arg name="handler" type="s" direction="in"/>
      <arg name="fieldsString" type="s" direction="in"/>
      <arg name="handlerType" type="y" direction="in"/>
      <arg name="grid" type="b" direction="in"/>
      <arg name="images" type="s" direction="out"/>
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
	SendKnobEvent(serial string, knobIndex int, event int) *dbus.Error
	GetLayersExample(serial string, layersString string, handlerType api.HandlerType) (string, *dbus.Error)
	GetBackgroundExample(serial string, handler string, fieldsString string, handlerType api.HandlerType, grid bool) (string, *dbus.Error)
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
	time.Sleep(50 * time.Millisecond)

	dev.InputManager().HandleKeyInput(&dev.Config().Pages[dev.PageManager().GetPage()].Keys[keyIndex], streamdeck.InputEvent{
		EventType: streamdeck.KEY_RELEASE,
	})
	return nil
}

func (StreamDeckDBus) SendKnobEvent(serial string, knobIndex int, event int) *dbus.Error {
	dev, ok := Devs[serial]
	if !ok || !dev.IsOpen() {
		return dbus.MakeFailedError(errors.New("Can't find connected device: " + serial))
	}

	knobs := dev.Config().Pages[dev.PageManager().GetPage()].Knobs
	if knobIndex < 0 || knobIndex >= len(knobs) {
		return dbus.MakeFailedError(errors.New("knob index out of range"))
	}

	info := dev.SdInfo()

	inputEvent := streamdeck.InputEvent{
		EventType: streamdeck.InputEventType(event),
		Index:     uint8(knobIndex),
	}
	switch inputEvent.EventType {
	case streamdeck.KNOB_CW, streamdeck.KNOB_CCW:
		inputEvent.RotateNotches = 1
	case streamdeck.SCREEN_SHORT_TAP, streamdeck.SCREEN_LONG_TAP:
		inputEvent.ScreenX = uint16(knobIndex*info.LcdWidth + info.LcdWidth/2)
		inputEvent.ScreenY = uint16(info.LcdHeight / 2)
	case streamdeck.KNOB_PRESS:
	case streamdeck.SCREEN_SWIPE:
		return dbus.MakeFailedError(errors.New("screen swipe events cannot be triggered"))
	default:
		return dbus.MakeFailedError(errors.New("not a knob or screen event"))
	}
	dev.InputManager().HandleKnobInput(&knobs[knobIndex], inputEvent)
	return nil
}

var (
	exampleTimeout  = 5 * time.Second
	examplePollTick = 50 * time.Millisecond
)

func awaitFrames(start func(emit func([]image.Image)), stop func(), want, width, height int) ([]image.Image, error) {
	var (
		mu     sync.Mutex
		frames []image.Image
	)

	go start(func(imgs []image.Image) {
		if len(imgs) != want {
			return
		}

		resized := make([]image.Image, want)
		for i, img := range imgs {
			if img == nil {
				return
			}
			if img.Bounds().Dx() != width || img.Bounds().Dy() != height {
				img = api.ResizeImageWH(img, width, height)
			}
			resized[i] = img
		}

		mu.Lock()
		first := frames == nil
		if first {
			frames = resized
		}
		mu.Unlock()

		if first {
			go stop()
		}
	})

	deadline := time.Now().Add(exampleTimeout)
	for time.Now().Before(deadline) {
		mu.Lock()
		captured := frames
		mu.Unlock()

		if captured != nil {
			return captured, nil
		}
		time.Sleep(examplePollTick)
	}

	go stop()
	return nil, errors.New("handler did not respond in a timely fashion")
}

func pngDataURL(img image.Image) (string, error) {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

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

	rendered := make([]image.Image, len(visual))
	failures := make([]error, len(visual))
	handlers := make([]api.ForegroundHandler, len(visual))

	for i, layer := range visual {
		handler := newExampleHandler(layer.Handler)
		if handler == nil {
			return "", dbus.MakeFailedError(errors.New("invalid handler: " + layer.Handler))
		}
		handlers[i] = handler
	}

	var wg sync.WaitGroup
	for i := range visual {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			layer := visual[index]
			handler := handlers[index]
			log.Println("Created and running " + layer.Handler + " for dbus")

			frames, err := awaitFrames(func(emit func([]image.Image)) {
				handler.Start(layer.Fields, handlerType, *info, func(img image.Image) {
					emit([]image.Image{img})
				})
			}, func() {
				log.Println("Stopping " + layer.Handler + " for dbus")
				handler.Stop()
			}, 1, width, height)

			if err != nil {
				failures[index] = err
				return
			}
			rendered[index] = frames[0]
		}(i)
	}
	wg.Wait()

	for _, err := range failures {
		if err != nil {
			return "", dbus.MakeFailedError(err)
		}
	}

	img, err := api.LayerImages(width, height, rendered...)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}

	dataURL, err := pngDataURL(img)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	return dataURL, nil
}

func backgroundGridSize(info api.StreamDeckInfo, handlerType api.HandlerType) int {
	if handlerType == api.LCD {
		return info.KnobCols
	}
	return info.Cols * info.Rows
}

func (StreamDeckDBus) GetBackgroundExample(serial string, handler string, fieldsString string, handlerType api.HandlerType, grid bool) (string, *dbus.Error) {
	var fields map[string]any
	if err := json.Unmarshal([]byte(fieldsString), &fields); err != nil {
		return "", dbus.MakeFailedError(err)
	}

	dev, ok := Devs[serial]
	if !ok {
		return "", dbus.MakeFailedError(errors.New("could not find device: " + serial))
	}
	info := dev.SdInfo()

	module, err := GetModule(handler, BACKGROUND)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	if module.NewBackground == nil {
		return "", dbus.MakeFailedError(errors.New("invalid background handler: " + handler))
	}
	background := module.NewBackground()
	if background == nil {
		return "", dbus.MakeFailedError(errors.New("invalid background handler: " + handler))
	}

	width, height := info.GetDimensions(handlerType)

	want := 1
	start := func(emit func([]image.Image)) {
		background.Start(fields, handlerType, *info, func(img image.Image) {
			emit([]image.Image{img})
		})
	}

	if grid {
		want = backgroundGridSize(*info, handlerType)
		if want <= 0 {
			return "", dbus.MakeFailedError(errors.New("device has no grid for the requested handler type"))
		}
		start = func(emit func([]image.Image)) {
			background.StartGrid(fields, handlerType, *info, emit)
		}
	}

	log.Println("Created and running background " + handler + " for dbus")

	frames, err := awaitFrames(start, func() {
		log.Println("Stopping background " + handler + " for dbus")
		background.Stop()
	}, want, width, height)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}

	dataURLs := make([]string, 0, len(frames))
	for _, frame := range frames {
		dataURL, err := pngDataURL(frame)
		if err != nil {
			return "", dbus.MakeFailedError(err)
		}
		dataURLs = append(dataURLs, dataURL)
	}

	payload, err := json.Marshal(dataURLs)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	return string(payload), nil
}

func newExampleHandler(name string) api.ForegroundHandler {
	module, err := GetModule(name, FOREGROUND)
	if err != nil {
		return nil
	}
	handler, ok := module.NewCombined().(api.ForegroundHandler)
	if !ok {
		return nil
	}
	return handler
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
