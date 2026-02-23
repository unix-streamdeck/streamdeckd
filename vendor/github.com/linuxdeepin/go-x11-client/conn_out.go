// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

import (
	"bufio"
	"errors"
	"io"
	"math"
)

func (c *Conn) Flush() (err error) {
	c.ioMu.Lock()
	err = c.out.flushTo(c.out.request)
	c.ioMu.Unlock()
	if err != nil {
		c.Close()
	}
	return
}

func (c *Conn) sendSync() {
	header := RequestHeader{
		Opcode: GetInputFocusOpcode,
		Length: 4,
	}
	logPrintln("sendSync")
	c.sendRequest(false, 0, RequestDiscardReply,
		header.toBytes(), nil)
}

func (c *Conn) SendSync() {
	c.ioMu.Lock()
	c.sendSync()
	c.ioMu.Unlock()
}

// sequence number
type SeqNum uint64

const (
	seqNumConnClosed    SeqNum = 0
	seqNumExtNotPresent SeqNum = math.MaxUint64 - iota
	seqNumQueryExtErr
)

var errExtNotPresent = errors.New("extension not present")
var errQueryExtErr = errors.New("query extension error")

func (n SeqNum) err() error {
	switch n {
	case seqNumConnClosed:
		return errConnClosed
	case seqNumExtNotPresent:
		return errExtNotPresent
	case seqNumQueryExtErr:
		return errQueryExtErr
	default:
		return nil
	}
}

type out struct {
	request        SeqNum
	requestWritten SeqNum
	bw             *bufio.Writer
}

func (o *out) flushTo(request SeqNum) error {
	if !(request <= o.request) {
		panic("assert request < o.request failed")
	}

	if o.requestWritten >= request {
		return nil
	}

	logPrintln("flushTo", request)
	err := o.bw.Flush()
	if err != nil {
		return err
	}
	o.requestWritten = o.request
	return nil
}

func (c *Conn) writeRequest(header []byte, body RequestBody) error {
	_, err := c.out.bw.Write(header)
	if err != nil {
		return err
	}
	var emptyBuf [3]byte
	for _, data := range body {
		_, err = c.out.bw.Write(data)
		if err != nil {
			return err
		}

		pad := Pad(len(data))
		if pad > 0 {
			_, err = c.out.bw.Write(emptyBuf[:pad])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Conn) sendRequest(noReply bool, workaround uint, flags uint, header []byte, body RequestBody) {
	if c.isClosed() {
		return
	}
	c.out.request++
	if !noReply {
		// has reply
		c.in.requestExpected = c.out.request
	}
	logPrintln("sendRequest seq:", c.out.request)

	if workaround != 0 || flags != 0 {
		c.in.expectReply(c.out.request, workaround, flags)
	}

	err := c.writeRequest(header, body)
	if err != nil {
		logPrintln("write error:", err)
		c.Close()
		return
	}

	if c.out.bw.Buffered() == 0 {
		// write all the data of request to c.conn
		c.out.requestWritten = c.out.request
	}
}

type ProtocolRequest struct {
	Ext     *Extension
	NoReply bool
	Header  RequestHeader
	Body    RequestBody
}

// return sequence id
func (c *Conn) SendRequest(flags uint, req *ProtocolRequest) SeqNum {
	if c.isClosed() {
		return seqNumConnClosed
	}

	// process data auto field
	// set the major opcode, and the minor opcode for extensions
	if req.Ext != nil {
		extension := c.GetExtensionData(req.Ext)
		if extension == nil {
			return seqNumQueryExtErr
		}
		if !extension.Present {
			return seqNumExtNotPresent
		}

		req.Header.Opcode = extension.MajorOpcode
	}

	var requestLen uint64
	requestLen = 4
	for _, data := range req.Body {
		requestLen += uint64(len(data))
		requestLen += uint64(Pad(len(data)))
	}
	req.Header.Length = requestLen

	header := req.Header.toBytes()

	c.ioMu.Lock()
	c.sendRequest(req.NoReply, 0, flags, header, req.Body)
	request := c.out.request
	c.ioMu.Unlock()
	return request
}

type RequestHeader struct {
	Opcode uint8  // major opcode
	Data   uint8  // data or minor opcode
	Length uint64 // unit is byte
}

func (rh RequestHeader) toBytes() []byte {
	b := make([]byte, 4)
	b[0] = rh.Opcode
	b[1] = rh.Data

	if rh.Length&3 != 0 {
		panic("length is not a multiple of 4")
	}

	length := uint16(rh.Length)
	length >>= 2
	b[2] = byte(length)
	b[3] = byte(length >> 8)
	return b
}

type FixedSizeBuf struct {
	data   []byte
	offset int
}

func (b *FixedSizeBuf) Write1b(v uint8) *FixedSizeBuf {
	b.data[b.offset] = v
	b.offset++
	return b
}

// byte order: least significant byte first

func (b *FixedSizeBuf) Write2b(v uint16) *FixedSizeBuf {
	b.data[b.offset] = byte(v)
	b.data[b.offset+1] = byte(v >> 8)
	b.offset += 2
	return b
}

func (b *FixedSizeBuf) Write4b(v uint32) *FixedSizeBuf {
	b.data[b.offset] = byte(v)
	b.data[b.offset+1] = byte(v >> 8)
	b.data[b.offset+2] = byte(v >> 16)
	b.data[b.offset+3] = byte(v >> 24)
	b.offset += 4
	return b
}

func (b *FixedSizeBuf) Write8b(v uint64) *FixedSizeBuf {
	b.data[b.offset] = byte(v)
	b.data[b.offset+1] = byte(v >> 8)
	b.data[b.offset+2] = byte(v >> 16)
	b.data[b.offset+3] = byte(v >> 24)
	b.data[b.offset+4] = byte(v >> 32)
	b.data[b.offset+5] = byte(v >> 40)
	b.data[b.offset+6] = byte(v >> 48)
	b.data[b.offset+7] = byte(v >> 56)
	b.offset += 8
	return b
}

func (b *FixedSizeBuf) WritePad(n int) *FixedSizeBuf {
	b.offset += n
	return b
}

func (b *FixedSizeBuf) WriteString(s string) *FixedSizeBuf {
	n := copy(b.data[b.offset:], s)
	b.offset += n
	if n != len(s) {
		panic(io.ErrShortWrite)
	}
	return b
}

func (b *FixedSizeBuf) WriteBytes(p []byte) *FixedSizeBuf {
	n := copy(b.data[b.offset:], p)
	b.offset += n
	if n != len(p) {
		panic(io.ErrShortWrite)
	}
	return b
}

func (b *FixedSizeBuf) WriteBool(v bool) *FixedSizeBuf {
	return b.Write1b(BoolToUint8(v))
}

func (b *FixedSizeBuf) End() {
	if len(b.data) != b.offset {
		panic("not end")
	}
}

func (b *FixedSizeBuf) Bytes() []byte {
	return b.data
}

func NewFixedSizeBuf(size int) *FixedSizeBuf {
	return &FixedSizeBuf{
		data: make([]byte, size),
	}
}

type RequestBody [][]byte

func (rb *RequestBody) AddBlock(n int) *FixedSizeBuf {
	b := NewFixedSizeBuf(n * 4)
	*rb = append(*rb, b.data)
	return b
}

func (rb *RequestBody) AddBytes(data []byte) {
	*rb = append(*rb, data)
}
