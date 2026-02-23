// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

import (
	"sync"
)

type ReadErrorFunc func(*Reader) Error

type lazyReplyTag uint

const (
	lazyNone lazyReplyTag = iota
	lazyCookie
	lazyForced
)

type Extension struct {
	id               int
	name             string // extension-xname
	maxErrCode       uint8
	errCodeNameMap   map[uint8]string
	reqOpcodeNameMap map[uint]string
}

func (ext *Extension) Name() string {
	return ext.name
}

var nextExtId = 1

// only call it in init() func
func NewExtension(name string, maxErrCode uint8,
	errCodeNameMap map[uint8]string, reqOpcodeNameMap map[uint]string) *Extension {
	id := nextExtId
	nextExtId++
	return &Extension{
		id:               id,
		name:             name,
		maxErrCode:       maxErrCode,
		errCodeNameMap:   errCodeNameMap,
		reqOpcodeNameMap: reqOpcodeNameMap,
	}
}

type lazyReply struct {
	tag    lazyReplyTag
	reply  *QueryExtensionReply
	cookie QueryExtensionCookie
}

type extAndData struct {
	ext   *Extension
	reply *lazyReply
}

type ext struct {
	mu         sync.RWMutex
	extensions []extAndData
}

func (e *ext) grow(n int) {
	if n <= len(e.extensions) {
		return
	}

	logPrintf("ext grow %d -> %d\n", len(e.extensions), n)
	bigger := make([]extAndData, n)
	copy(bigger, e.extensions)
	e.extensions = bigger
}

func (e *ext) getById(id int) *extAndData {
	if id > len(e.extensions) {
		e.grow(id * 2)
	}
	return &e.extensions[id-1]
}

func (e *ext) getExtByMajorOpcode(majorOpcode uint8) *Extension {
	for _, extAndData := range e.extensions {
		ext := extAndData.ext
		lzr := extAndData.reply
		if ext != nil && lzr != nil && lzr.reply != nil {
			if majorOpcode == lzr.reply.MajorOpcode {
				return ext
			}
		}
	}
	return nil
}

func (e *ext) getExtErrName(errCode uint8) string {
	for _, extAndData := range e.extensions {
		ext := extAndData.ext
		lzr := extAndData.reply
		if ext != nil && lzr != nil && lzr.reply != nil {
			base := lzr.reply.FirstError
			max := lzr.reply.FirstError + ext.maxErrCode

			if base <= errCode && errCode <= max {
				return ext.errCodeNameMap[errCode-base]
			}
		}
	}
	return ""
}

func (e *ext) getLazyReply(conn *Conn, ext *Extension) (lzr *lazyReply) {
	extAndData := e.getById(ext.id)

	if extAndData.ext == nil {
		// init extAndData
		extAndData.ext = ext
		extAndData.reply = &lazyReply{tag: lazyNone}
	}
	lzr = extAndData.reply

	// lazyReply tag: lazyNone -> lazyCookie
	if lzr.tag == lazyNone {
		lzr.tag = lazyCookie
		lzr.cookie = QueryExtension(conn, ext.name)
	}
	return
}

func (e *ext) getExtData(c *Conn, ext *Extension) *QueryExtensionReply {
	lzr := e.getLazyReply(c, ext)

	// lazyReply tag: lazyCookie -> lazyForced
	if lzr.tag == lazyCookie {
		lzr.tag = lazyForced
		lzr.reply, _ = lzr.cookie.Reply(c)
	}
	return lzr.reply
}

func (c *Conn) GetExtensionData(ext *Extension) *QueryExtensionReply {
	c.ext.mu.Lock()
	data := c.ext.getExtData(c, ext)
	c.ext.mu.Unlock()
	return data
}

func (c *Conn) PrefetchExtensionData(ext *Extension) {
	c.ext.mu.Lock()
	c.ext.getLazyReply(c, ext)
	c.ext.mu.Unlock()
}
