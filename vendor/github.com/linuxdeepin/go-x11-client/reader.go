// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

import (
	"bytes"
	"errors"
)

var ErrDataLenShort = errors.New("data length is short")

type Reader struct {
	pos  int
	data []byte
}

func NewReaderFromData(data []byte) *Reader {
	return &Reader{
		data: data,
	}
}

func (r *Reader) Pos() int {
	return r.pos
}

func (r *Reader) Read1b() uint8 {
	v := r.data[r.pos]
	r.pos++
	return v
}

func (r *Reader) Read2b() uint16 {
	v := uint16(r.data[r.pos])
	v |= uint16(r.data[r.pos+1]) << 8
	r.pos += 2
	return v
}

func (r *Reader) Read4b() uint32 {
	v := uint32(r.data[r.pos])
	v |= uint32(r.data[r.pos+1]) << 8
	v |= uint32(r.data[r.pos+2]) << 16
	v |= uint32(r.data[r.pos+3]) << 24
	r.pos += 4
	return v
}

func (r *Reader) MustReadBytes(n int) []byte {
	v := r.data[r.pos : r.pos+n]
	r.pos += n
	return v
}

func (r *Reader) ReadBytes(n int) ([]byte, error) {
	if !r.RemainAtLeast(n) {
		return nil, ErrDataLenShort
	}

	v := r.data[r.pos : r.pos+n]
	r.pos += n
	return v, nil
}

func (r *Reader) ReadBytesWithPad(n int) ([]byte, error) {
	total := n + Pad(n)
	if !r.RemainAtLeast(total) {
		return nil, ErrDataLenShort
	}

	v := r.data[r.pos : r.pos+n]
	r.pos += total
	return v, nil
}

func (r *Reader) ReadString(n int) (string, error) {
	v, err := r.ReadBytes(n)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func (r *Reader) ReadStrWithPad(n int) (string, error) {
	v, err := r.ReadBytesWithPad(n)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func (r *Reader) ReadNulTermStr() string {
	idx := bytes.IndexByte(r.data[r.pos:], 0)
	var v []byte
	if idx == -1 {
		v = r.data[r.pos:]
		r.pos = len(r.data)
	} else {
		v = r.data[r.pos : r.pos+idx]
		r.pos += idx
	}
	return string(v)
}

func (r *Reader) ReadBool() bool {
	return Uint8ToBool(r.Read1b())
}

func (r *Reader) ReadPad(n int) {
	if r.pos+n > len(r.data) {
		panic("index out of range")
	}
	r.pos += n
}

func (r *Reader) Reset() {
	r.pos = 0
}

func (r *Reader) RemainAtLeast(n int) bool {
	return r.pos+n <= len(r.data)
}

func (r *Reader) RemainAtLeast4b(n int) bool {
	return r.RemainAtLeast(n << 2)
}

func (r *Reader) ReadReplyHeader() (data uint8, length uint32) {
	r.ReadPad(1)
	data = r.Read1b()
	// seq
	r.ReadPad(2)
	// length
	length = r.Read4b() // 2
	return
}

func (r *Reader) ReadEventHeader() (data uint8, seq uint16) {
	r.ReadPad(1)
	data = r.Read1b()
	// seq
	seq = r.Read2b() // 1
	return
}
