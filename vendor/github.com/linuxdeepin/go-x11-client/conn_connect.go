// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func NewConn() (*Conn, error) {
	return NewConnDisplay("")
}

func NewConnDisplay(display string) (*Conn, error) {
	c := new(Conn)
	err := c.connectDisplay(display)
	if err != nil {
		return nil, err
	}

	c.bufReader = bufio.NewReader(c.conn)
	c.in.init(&c.ioMu)
	c.out.bw = bufio.NewWriter(c.conn)
	c.ridAllocator.init(c.setup.ResourceIdBase, c.setup.ResourceIdMask)
	go c.readLoop()
	go c.eventSendLoop()
	return c, nil
}

func (c *Conn) connectDisplay(display string) error {
	err := c.dial(display)
	if err != nil {
		return err
	}

	return c.postConnect()
}

// dial initializes the actual net connection with X.
func (c *Conn) dial(display string) error {
	if len(display) == 0 {
		display = os.Getenv("DISPLAY")
	}

	display0 := display
	if len(display) == 0 {
		return errors.New("empty display string")
	}

	colonIdx := strings.LastIndex(display, ":")
	if colonIdx < 0 {
		return errors.New("bad display string: " + display0)
	}

	var protocol, socket string

	if display[0] == '/' {
		socket = display[0:colonIdx]
	} else {
		slashIdx := strings.LastIndex(display, "/")
		if slashIdx >= 0 {
			protocol = display[0:slashIdx]
			c.host = display[slashIdx+1 : colonIdx]
		} else {
			c.host = display[0:colonIdx]
		}
	}

	display = display[colonIdx+1:]
	if len(display) == 0 {
		return errors.New("bad display string: " + display0)
	}

	var scr string
	dotIdx := strings.LastIndex(display, ".")
	if dotIdx < 0 {
		c.display = display[0:]
	} else {
		c.display = display[0:dotIdx]
		scr = display[dotIdx+1:]
	}

	var err error
	c.DisplayNumber, err = strconv.Atoi(c.display)
	if err != nil || c.DisplayNumber < 0 {
		return errors.New("bad display string: " + display0)
	}

	if scr != "" {
		c.ScreenNumber, err = strconv.Atoi(scr)
		if err != nil {
			return errors.New("bad display string: " + display0)
		}
	}
	logPrintf("socket: %q, protocol: %q\n", socket, protocol)

	// Connect to server
	if len(socket) != 0 {
		c.conn, err = net.Dial("unix", socket+":"+c.display)
	} else if len(c.host) != 0 {
		if protocol == "" {
			protocol = "tcp"
		}
		logPrintln("dial tcp")
		c.conn, err = net.Dial(protocol,
			c.host+":"+strconv.Itoa(6000+c.DisplayNumber))
	} else {
		logPrintln("dial unix")
		c.conn, err = net.Dial("unix", "/tmp/.X11-unix/X"+c.display)
	}

	if err != nil {
		return errors.New("cannot connect to " + display0 + ": " + err.Error())
	}
	return nil
}

// do the postConnect action after Conn get it's underly net.Conn
func (c *Conn) postConnect() error {
	// Get authentication data
	authName, authData, err := readAuthority(c.host, c.display)
	noAuth := false
	if err != nil {
		logPrintf("Could not get authority info: %v", err)
		logPrintln("Trying connection without authority info...")
		authName = ""
		authData = []byte{}
		noAuth = true
	}

	// Assume that the authentication protocol is "MIT-MAGIC-COOKIE-1".
	if !noAuth && (authName != "MIT-MAGIC-COOKIE-1" || len(authData) != 16) {
		return errors.New("unsupported auth protocol " + authName)
	}

	bufLen := 12 + len(authName) + Pad(len(authName)) + len(authData) + Pad(len(authData))
	buf := make([]byte, bufLen)
	buf[0] = 0x6c
	buf[1] = 0
	Put16(buf[2:], 11)
	Put16(buf[4:], 0)
	Put16(buf[6:], uint16(len(authName)))
	Put16(buf[8:], uint16(len(authData)))
	Put16(buf[10:], 0)
	copy(buf[12:], []byte(authName))
	copy(buf[12+len(authName)+Pad(len(authName)):], authData)
	if _, err = c.conn.Write(buf); err != nil {
		return err
	}

	head := make([]byte, 8)
	if _, err = io.ReadFull(c.conn, head[0:8]); err != nil {
		return err
	}
	respType := head[0]
	major := Get16(head[2:])
	minor := Get16(head[4:])
	dataLen := Get16(head[6:])

	if major != 11 || minor != 0 {
		return fmt.Errorf("x protocol version mismatch: %d.%d", major, minor)
	}

	buf = make([]byte, 8+dataLen*4)
	copy(buf, head)
	if _, err = io.ReadFull(c.conn, buf[8:]); err != nil {
		return err
	}

	if respType == ResponseTypeError {
		reasonLen := head[1]
		reason := buf[8 : 8+reasonLen]
		return fmt.Errorf("x protocol authentication refused: %s",
			string(reason))
	}

	var setup Setup
	r := NewReaderFromData(buf)
	err = readSetup(r, &setup)
	if err != nil {
		return err
	}
	c.setup = &setup

	/* Make sure requested screen number is in bounds for this server */
	if c.ScreenNumber >= len(setup.Roots) {
		return errors.New("invalid screen")
	}

	return nil
}
