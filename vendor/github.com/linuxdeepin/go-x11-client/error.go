// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

import "fmt"

type Error struct {
	conn       *Conn
	Code       uint8  // ErrorCode
	Sequence   uint16 // Sequence number
	ResourceID uint32 // Resource ID for requests with side effects only
	MinorCode  uint16 // Minor opcode of the failed request
	MajorCode  uint8  // Major opcode of the failed request
}

func (err *Error) Error() string {
	var errDesc string
	var errName string
	if 1 <= err.Code && err.Code <= 127 {
		// core error code in range [1,127]
		errName = errorCodeNameMap[err.Code]
	} else {
		// is ext error
		errName = err.conn.ext.getExtErrName(err.Code)
	}
	if errName != "" {
		errDesc = " (" + errName + ")"
	}

	var majorCodeDesc, minorCodeDesc string

	if 1 <= err.MajorCode && err.MajorCode <= 127 {
		// is core request
		reqName := requestOpcodeNameMap[uint(err.MajorCode)]
		if reqName != "" {
			majorCodeDesc = " (" + reqName + ")"
		}
	} else {
		// is ext request
		ext := err.conn.ext.getExtByMajorOpcode(err.MajorCode)
		if ext != nil {
			reqName := ext.reqOpcodeNameMap[uint(err.MinorCode)]
			if reqName != "" {
				minorCodeDesc = " (" + reqName + ")"
			}
			majorCodeDesc = " (" + ext.name + ")"
		}
	}

	return fmt.Sprintf("x.Error: %d%s, sequence: %d, resource id: %d,"+
		" major code: %d%s, minor code: %d%s",
		err.Code, errDesc, err.Sequence, err.ResourceID,
		err.MajorCode, majorCodeDesc,
		err.MinorCode, minorCodeDesc)
}

func newError(data []byte) *Error {
	var v Error

	responseType := data[0]
	if responseType != ResponseTypeError {
		panic("not error")
	}
	b := 1

	v.Code = data[b]
	b += 1

	v.Sequence = Get16(data[b:])
	b += 2

	v.ResourceID = Get32(data[b:])
	b += 4

	v.MinorCode = Get16(data[b:])
	b += 2

	v.MajorCode = data[b]

	return &v
}

func (c *Conn) NewError(data []byte) *Error {
	err := newError(data)
	err.conn = c
	return err
}
