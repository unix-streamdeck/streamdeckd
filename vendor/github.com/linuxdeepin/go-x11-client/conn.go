// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
)

var Logger *log.Logger
var debugEnabled bool

func init() {
	if os.Getenv("DEBUG_X11_CLIENT") == "1" {
		debugEnabled = true
		Logger = log.New(os.Stderr, "[x] ", log.Lshortfile)
	}
}

func logPrintln(v ...interface{}) {
	if debugEnabled {
		_ = Logger.Output(2, fmt.Sprintln(v...))
	}
}
func logPrintf(format string, v ...interface{}) {
	if debugEnabled {
		_ = Logger.Output(2, fmt.Sprintf(format, v...))
	}
}

type Conn struct {
	conn      net.Conn
	closed    int32
	bufReader *bufio.Reader

	host          string
	display       string
	DisplayNumber int
	ScreenNumber  int
	setup         *Setup

	ioMu sync.Mutex
	in   in
	out  out

	ext          ext
	ridAllocator resourceIdAllocator
	atomCache    *AtomCache
	atomCacheMu  sync.Mutex
	errorCb      func(err *Error)
}

func (c *Conn) GetSetup() *Setup {
	return c.setup
}

func (c *Conn) GetDefaultScreen() *Screen {
	return &c.setup.Roots[c.ScreenNumber]
}

func (c *Conn) isClosed() bool {
	return 1 == atomic.LoadInt32(&c.closed)
}

func (c *Conn) markClosed() {
	atomic.StoreInt32(&c.closed, 1)
}

func (c *Conn) Close() {
	if c.isClosed() {
		return
	}
	c.markClosed()
	c.conn.Close()

	c.in.eventsCond.Signal()
	go func() {
		c.ioMu.Lock()
		c.in.wakeUpAllReaders()
		c.ioMu.Unlock()
	}()
}

var errConnClosed = errors.New("conn closed")

func (c *Conn) AddEventChan(eventChan chan<- GenericEvent) {
	if c.isClosed() || eventChan == nil {
		return
	}
	c.in.addEventChan(eventChan)
}

func (c *Conn) RemoveEventChan(eventChan chan<- GenericEvent) {
	if eventChan == nil {
		return
	}
	c.in.removeEventChan(eventChan)
}

func (c *Conn) SetErrorCallback(fn func(err *Error)) {
	c.errorCb = fn
}

func (c *Conn) MakeAndAddEventChan(bufSize int) chan GenericEvent {
	ch := make(chan GenericEvent, bufSize)
	c.AddEventChan(ch)

	return ch
}
