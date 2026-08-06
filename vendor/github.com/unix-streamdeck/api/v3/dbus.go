package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/png"
	"strings"

	"github.com/godbus/dbus/v5"
)

type IConn interface {
	Close() error
	AddMatchSignal(options ...dbus.MatchOption) error
	Signal(ch chan<- *dbus.Signal)
	Object(dest string, path dbus.ObjectPath) dbus.BusObject
}

type Connection struct {
	busobj dbus.BusObject
	conn   IConn
}

func Connect() (*Connection, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}
	return &Connection{
		conn:   conn,
		busobj: conn.Object("com.unixstreamdeck.streamdeckd", "/com/unixstreamdeck/streamdeckd"),
	}, nil
}

func (c *Connection) Close() {
	c.conn.Close()
}

func (c *Connection) GetInfo() ([]*StreamDeckInfo, error) {
	var s string
	err := c.busobj.Call("com.unixstreamdeck.streamdeckd.GetDeckInfo", 0).Store(&s)
	if err != nil {
		return nil, err
	}
	var info []*StreamDeckInfo
	err = json.Unmarshal([]byte(s), &info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (c *Connection) SetPage(serial string, page int) error {
	call := c.busobj.Call("com.unixstreamdeck.streamdeckd.SetPage", 0, serial, page)
	if call.Err != nil {
		return call.Err
	}
	return nil
}

func (c *Connection) GetConfig() (*Config, error) {
	var s string
	err := c.busobj.Call("com.unixstreamdeck.streamdeckd.GetConfig", 0).Store(&s)
	if err != nil {
		return nil, err
	}
	var config *Config
	err = json.Unmarshal([]byte(s), &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Connection) SetConfig(config *Config) error {
	configString, err := json.Marshal(config)
	if err != nil {
		return err
	}
	call := c.busobj.Call("com.unixstreamdeck.streamdeckd.SetConfig", 0, string(configString))
	if call.Err != nil {
		return call.Err
	}
	return nil
}

func (c *Connection) ReloadConfig() error {
	call := c.busobj.Call("com.unixstreamdeck.streamdeckd.ReloadConfig", 0)
	if call.Err != nil {
		return call.Err
	}
	return nil
}

func (c *Connection) CommitConfig() error {
	call := c.busobj.Call("com.unixstreamdeck.streamdeckd.CommitConfig", 0)
	if call.Err != nil {
		return call.Err
	}
	return nil
}

func (c *Connection) GetModules() (*AvailableModules, error) {
	var s string
	err := c.busobj.Call("com.unixstreamdeck.streamdeckd.GetModules", 0).Store(&s)
	if err != nil {
		return nil, err
	}
	var modules *AvailableModules
	err = json.Unmarshal([]byte(s), &modules)
	if err != nil {
		return nil, err
	}
	return modules, nil
}

func (c *Connection) PressButton(serial string, keyIndex int) error {
	return c.busobj.Call("com.unixstreamdeck.streamdeckd.PressButton", 0, serial, keyIndex).Err
}

func (c *Connection) SendKnobEvent(serial string, knobIndex int, event InputEventType) error {
	return c.busobj.Call("com.unixstreamdeck.streamdeckd.SendKnobEvent", 0, serial, knobIndex, int(event)).Err
}

func (c *Connection) GetObsFields() ([]*Field, error) {
	var s string
	err := c.busobj.Call("com.unixstreamdeck.streamdeckd.GetObsFields", 0).Store(&s)
	if err != nil {
		return nil, err
	}
	var fields []*Field
	err = json.Unmarshal([]byte(s), &fields)
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func decodePNGDataURL(dataURL string) (image.Image, error) {
	raw, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(dataURL, "data:image/png;base64,", ""))
	if err != nil {
		return nil, err
	}
	return png.Decode(bytes.NewReader(raw))
}

func (c *Connection) GetLayersExample(serial string, layers []*Layer, handlerType HandlerType) (image.Image, error) {
	layersString, err := json.Marshal(layers)
	if err != nil {
		return nil, err
	}
	var s string
	err = c.busobj.Call("com.unixstreamdeck.streamdeckd.GetLayersExample", 0, serial, string(layersString), handlerType).Store(&s)
	if err != nil {
		return nil, err
	}
	return decodePNGDataURL(s)
}

func (c *Connection) GetBackgroundExample(serial string, handler string, fields map[string]any, handlerType HandlerType, grid bool) ([]image.Image, error) {
	fieldsString, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}
	var s string
	err = c.busobj.Call("com.unixstreamdeck.streamdeckd.GetBackgroundExample", 0, serial, handler, string(fieldsString), handlerType, grid).Store(&s)
	if err != nil {
		return nil, err
	}

	var dataURLs []string
	if err := json.Unmarshal([]byte(s), &dataURLs); err != nil {
		return nil, err
	}

	images := make([]image.Image, 0, len(dataURLs))
	for _, dataURL := range dataURLs {
		img, err := decodePNGDataURL(dataURL)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

func (c *Connection) RegisterPageListener(cback func(string, int32)) error {
	err := c.conn.AddMatchSignal(dbus.WithMatchObjectPath("/com/unixstreamdeck/streamdeckd"), dbus.WithMatchInterface("com.unixstreamdeck.streamdeckd"), dbus.WithMatchMember("Page"))
	if err != nil {
		return err
	}
	ch := make(chan *dbus.Signal, 10)
	c.conn.Signal(ch)
	for v := range ch {
		cback(v.Body[0].(string), v.Body[1].(int32))
	}
	return nil
}
