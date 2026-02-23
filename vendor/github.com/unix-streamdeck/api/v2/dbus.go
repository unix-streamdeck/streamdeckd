package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/godbus/dbus/v5"
	"image"
	"image/png"
	"strings"
)

type IConn interface {
	Close() error
	AddMatchSignal(options ...dbus.MatchOption) error
	Signal(ch chan<- *dbus.Signal)
}

type Conn struct {
	conn *dbus.Conn
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) AddMatchSignal(options ...dbus.MatchOption) error {
	return c.conn.AddMatchSignal(options...)
}

func (c *Conn) Signal(ch chan<- *dbus.Signal) {
	c.conn.Signal(ch)
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
		conn:   &Conn{conn: conn},
		busobj: conn.Object("com.unixstreamdeck.streamdeckd", "/com/unixstreamdeck/streamdeckd"),
	}, nil
}

func (c *Connection) Close() {
	c.conn.Close()
}

func (c *Connection) GetInfo() ([]*StreamDeckInfoV1, error) {
	var s string
	err := c.busobj.Call("com.unixstreamdeck.streamdeckd.GetDeckInfo", 0).Store(&s)
	if err != nil {
		return nil, err
	}
	var info []*StreamDeckInfoV1
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

func (c *Connection) GetConfig() (*ConfigV3, error) {
	var s string
	err := c.busobj.Call("com.unixstreamdeck.streamdeckd.GetConfig", 0).Store(&s)
	if err != nil {
		return nil, err
	}
	var config *ConfigV3
	err = json.Unmarshal([]byte(s), &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Connection) SetConfig(config *ConfigV3) error {
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

func (c *Connection) GetModules() ([]*Module, error) {
	var s string
	err := c.busobj.Call("com.unixstreamdeck.streamdeckd.GetModules", 0).Store(&s)
	if err != nil {
		return nil, err
	}
	var modules []*Module
	err = json.Unmarshal([]byte(s), &modules)
	if err != nil {
		return nil, err
	}
	return modules, nil
}

func (c *Connection) PressButton(serial string, keyIndex int) error {
	return c.busobj.Call("com.unixstreamdeck.streamdeckd.PressButton", 0, serial, keyIndex).Err
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

func (c *Connection) GetHandlerExample(serial string, keyConfig KeyConfigV3) (image.Image, error) {
	configString, err := json.Marshal(keyConfig)
	if err != nil {
		return nil, err
	}
	var s string
	err = c.busobj.Call("com.unixstreamdeck.streamdeckd.GetHandlerExample", 0, string(serial), string(configString)).Store(&s)
	if err != nil {
		return nil, err
	}
	s = strings.ReplaceAll(s, "data:image/png;base64,", "")
	bytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return png.Decode(strings.NewReader(string(bytes)))
}
func (c *Connection) GetKnobHandlerExample(serial string, knobConfig KnobConfigV3) (image.Image, error) {
	configString, err := json.Marshal(knobConfig)
	if err != nil {
		return nil, err
	}
	var s string
	err = c.busobj.Call("com.unixstreamdeck.streamdeckd.GetKnobHandlerExample", 0, string(serial), string(configString)).Store(&s)
	if err != nil {
		return nil, err
	}
	s = strings.ReplaceAll(s, "data:image/png;base64,", "")
	bytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return png.Decode(strings.NewReader(string(bytes)))
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
