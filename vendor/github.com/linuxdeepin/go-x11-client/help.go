// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

/*
help.go is meant to contain a rough hodge podge of functions that are mainly
used in the auto generated code. Indeed, several functions here are simple
wrappers so that the sub-packages don't need to be smart about which stdlib
packages to import.

Also, the 'Get..' and 'Put..' functions are used through the core xgb package
too. (xgbutil uses them too.)
*/

import (
	"fmt"
	//"github.com/gavv/monotime"
	"strings"
)

//func GetServerCurrentTime() Timestamp {
//	now := monotime.Now()
//	ns := now.Nanoseconds()
//	return Timestamp(ns / 1e6)
//}

// StringsJoin is an alias to strings.Join. It allows us to avoid having to
// import 'strings' in each of the generated Go files.
func StringsJoin(ss []string, sep string) string {
	return strings.Join(ss, sep)
}

// Sprintf is so we don't need to import 'fmt' in the generated Go files.
func Sprintf(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}

// Errorf is just a wrapper for fmt.Errorf. Exists for the same reason
// that 'stringsJoin' and 'sprintf' exists.
func Errorf(format string, v ...interface{}) error {
	return fmt.Errorf(format, v...)
}

func Pad(e int) int {
	// pad(E) = (4 - (E mod 4)) mod 4
	return (4 - (e % 4)) % 4
}

// PopCount counts the number of bits set in a value list mask.
func PopCount(mask0 int) int {
	mask := uint32(mask0)
	n := 0
	for i := uint32(0); i < 32; i++ {
		if mask&(1<<i) != 0 {
			n++
		}
	}
	return n
}

func SizeIn4bWithPad(length int) int {
	return (length + Pad(length)) / 4
}

// Put16 takes a 16 bit integer and copies it into a byte slice.
func Put16(buf []byte, v uint16) {
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
}

// Put32 takes a 32 bit integer and copies it into a byte slice.
func Put32(buf []byte, v uint32) {
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
	buf[2] = byte(v >> 16)
	buf[3] = byte(v >> 24)
}

// Put64 takes a 64 bit integer and copies it into a byte slice.
func Put64(buf []byte, v uint64) {
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
	buf[2] = byte(v >> 16)
	buf[3] = byte(v >> 24)
	buf[4] = byte(v >> 32)
	buf[5] = byte(v >> 40)
	buf[6] = byte(v >> 48)
	buf[7] = byte(v >> 56)
}

// Get16 constructs a 16 bit integer from the beginning of a byte slice.
func Get16(buf []byte) uint16 {
	v := uint16(buf[0])
	v |= uint16(buf[1]) << 8
	return v
}

// Get32 constructs a 32 bit integer from the beginning of a byte slice.
func Get32(buf []byte) uint32 {
	v := uint32(buf[0])
	v |= uint32(buf[1]) << 8
	v |= uint32(buf[2]) << 16
	v |= uint32(buf[3]) << 24
	return v
}

// Get64 constructs a 64 bit integer from the beginning of a byte slice.
func Get64(buf []byte) uint64 {
	v := uint64(buf[0])
	v |= uint64(buf[1]) << 8
	v |= uint64(buf[2]) << 16
	v |= uint64(buf[3]) << 24
	v |= uint64(buf[4]) << 32
	v |= uint64(buf[5]) << 40
	v |= uint64(buf[6]) << 48
	v |= uint64(buf[7]) << 56
	return v
}

func BoolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func Uint8ToBool(v uint8) bool {
	return v != 0
}

func TruncateStr(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
