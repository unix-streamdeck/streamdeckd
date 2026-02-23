// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

import (
	"bytes"
)

type Writer struct {
	buf bytes.Buffer
}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) Reset() {
	w.buf.Reset()
}

func (w *Writer) WritePad(n int) {
	for i := 0; i < n; i++ {
		w.buf.WriteByte(0)
	}
}

func (w *Writer) Write1b(v uint8) {
	w.buf.WriteByte(v)
}

func (w *Writer) Write2b(v uint16) {
	w.buf.WriteByte(byte(v))
	w.buf.WriteByte(byte(v >> 8))
}

func (w *Writer) Write4b(v uint32) {
	w.buf.WriteByte(byte(v))
	w.buf.WriteByte(byte(v >> 8))
	w.buf.WriteByte(byte(v >> 16))
	w.buf.WriteByte(byte(v >> 24))
}

func (w *Writer) Write8b(v uint64) {
	w.buf.WriteByte(byte(v))
	w.buf.WriteByte(byte(v >> 8))
	w.buf.WriteByte(byte(v >> 16))
	w.buf.WriteByte(byte(v >> 24))
	w.buf.WriteByte(byte(v >> 32))
	w.buf.WriteByte(byte(v >> 40))
	w.buf.WriteByte(byte(v >> 48))
	w.buf.WriteByte(byte(v >> 56))
}

func (w *Writer) WriteBytes(p []byte) {
	w.buf.Write(p)
}

func (w *Writer) WriteNBytes(n int, p []byte) {
	if len(p) < n {
		w.buf.Write(p)
		w.WritePad(n - len(p))
	} else {
		w.buf.Write(p[:n])
	}
}

func (w *Writer) WriteString(s string) {
	w.buf.WriteString(s)
}

func (w *Writer) Bytes() []byte {
	return w.buf.Bytes()
}
