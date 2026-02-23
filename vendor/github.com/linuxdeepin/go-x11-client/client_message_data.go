// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

type ClientMessageData struct {
	data []byte
}

func ClientMessageDataRead(r *Reader, v *ClientMessageData) int {
	v.data = r.MustReadBytes(20)
	return 20
}

func writeClientMessageData(w *Writer, v *ClientMessageData) int {
	w.WriteNBytes(20, v.data)
	return 20
}

func (v *ClientMessageData) GetData8() []byte {
	return v.data
}

func (v *ClientMessageData) GetData16() []uint16 {
	ret := make([]uint16, 10)
	idx := 0
	for i := 0; i < 20; i += 2 {
		ret[idx] = Get16(v.data[i : i+2])
		idx++
	}
	return ret
}

func (v *ClientMessageData) GetData32() []uint32 {
	ret := make([]uint32, 5)
	idx := 0
	for i := 0; i < 20; i += 4 {
		ret[idx] = Get32(v.data[i : i+4])
		idx++
	}
	return ret
}

func (v *ClientMessageData) SetData8(p *[20]byte) {
	v.data = p[:]
}

func (v *ClientMessageData) SetData16(p *[10]uint16) {
	w := NewWriter()
	for i := 0; i < 10; i++ {
		w.Write2b(p[i])
	}
	v.data = w.Bytes()
}

func (v *ClientMessageData) SetData32(p *[5]uint32) {
	w := NewWriter()
	for i := 0; i < 5; i++ {
		w.Write4b(p[i])
	}
	v.data = w.Bytes()
}
