// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

import (
	"container/list"
	"io"
	"sync"
)

const (
	ResponseTypeError = 0
	ResponseTypeReply = 1
)

//enum xcb_send_request_flags_t {
//XCB_REQUEST_CHECKED = 1 << 0,
//XCB_REQUEST_RAW = 1 << 1,
//XCB_REQUEST_DISCARD_REPLY = 1 << 2,
//XCB_REQUEST_REPLY_FDS = 1 << 3
//};
const (
	RequestChecked      = 1 << 0
	RequestRaw          = 1 << 1 // ?
	RequestDiscardReply = 1 << 2
	RequestReplyFds     = 1 << 3 // ?
)

//typedef struct {
//uint8_t   response_type;  [>*< Type of the response <]
//uint8_t  pad0;           [>*< Padding <]
//uint16_t sequence;       [>*< Sequence number <]
//uint32_t length;         [>*< Length of the response <]
//} xcb_generic_reply_t;

type GenericReply struct {
	responseType uint8
	sequence     uint16
	length       uint32
}

func NewGenericReply(buf []byte) *GenericReply {
	r := &GenericReply{}
	r.responseType = buf[0]
	// skip pad buf[1]
	r.sequence = Get16(buf[2:])
	r.length = Get32(buf[4:])
	return r
}

// go c.readLoop
func (c *Conn) readLoop() {
	for {
		err := c.readPacket()
		if err != nil {
			c.Close()
			return
		}
	}
}

func (c *Conn) readPacket() error {
	length := 32
	buf := make([]byte, length)
	_, err := io.ReadFull(c.bufReader, buf)
	if err != nil {
		return err
	}
	genReply := NewGenericReply(buf)
	logPrintf("genReply: %#v\n", genReply)

	if (genReply.responseType == ResponseTypeReply ||
		genReply.responseType == GeGenericEventCode) && genReply.length > 0 {
		length += int(genReply.length) * 4
		// grow buf
		biggerBuf := make([]byte, length)
		copy(biggerBuf[:32], buf)
		_, err = io.ReadFull(c.bufReader, biggerBuf[32:])
		if err != nil {
			return err
		}
		buf = biggerBuf
	}

	c.ioMu.Lock()
	defer c.ioMu.Unlock()
	// update c.in.request*

	// KeymapNotifyEvent 没有 sequence
	if genReply.responseType != KeymapNotifyEventCode {
		lastRead := c.in.requestRead
		c.in.requestRead = (lastRead & 0xffffffffffff0000) | SeqNum(genReply.sequence)
		if c.in.requestRead < lastRead {
			c.in.requestRead += 0x10000
		}

		if c.in.requestRead > c.in.requestExpected {
			c.in.requestExpected = c.in.requestRead
		}

		if c.in.requestRead != lastRead {
			// 这个 reply 是对一个新的请求的 reply
			curReply := c.in.currentReply
			if curReply != nil && curReply.Len() != 0 {
				c.in.replies[lastRead] = curReply
				c.in.currentReply = nil
			}
			c.in.requestCompleted = c.in.requestRead - 1
		}

		c.in.removeFinishedPendingReplies()

		if genReply.responseType == ResponseTypeError {
			c.in.requestCompleted = c.in.requestRead
		}
		c.in.removeFinishedReaders()
	}

	var pend *PendingReply
	if genReply.responseType == ResponseTypeReply ||
		genReply.responseType == ResponseTypeError {

		if prFront := c.in.pendingReplies.Front(); prFront != nil {
			pend = prFront.Value.(*PendingReply)

			if !(pend.firstRequest <= c.in.requestRead &&
				c.in.requestRead <= pend.lastRequest) {
				pend = nil
			}
		}
	}

	logPrintf("pend %#v\n", pend)
	if pend != nil && pend.flags&RequestDiscardReply != 0 {
		// discard reply
		logPrintln("discard reply", c.in.requestRead)
		return nil
	}

	/* reply, or checked error */
	if genReply.responseType == ResponseTypeReply ||
		(genReply.responseType == ResponseTypeError &&
			pend != nil && pend.flags&RequestChecked != 0) {

		if c.in.currentReply == nil {
			c.in.currentReply = list.New()
		}
		c.in.currentReply.PushBack(buf)
		logPrintf("pushBack buf %d len=%d\n", c.in.requestRead, len(buf))

		front := c.in.readers.Front()
		if front != nil {
			reader := front.Value.(*ReplyReader)
			if reader.request == c.in.requestRead {
				logPrintf("readPacket reader %d signal\n", reader.request)
				reader.cond.Signal()
			}
		}
		return nil
	}

	/* event, or unchecked error */
	// not special event
	if genReply.responseType == ResponseTypeError {
		// is unchecked error
		if c.errorCb != nil {
			c.errorCb(c.NewError(buf))
		} else {
			c.in.addEvent(GenericEvent(buf))
		}
	} else {
		// is event
		c.in.addEvent(GenericEvent(buf))
	}

	return nil
}

func (c *Conn) pollForReply(request SeqNum) (replyBuf []byte, isErr, stop bool) {
	if c.isClosed() {
		stop = true
		return
	}

	logPrintln("pollForReply", request)
	var front *list.Element
	if request < c.in.requestRead {
		/* We've read requests past the one we want, so if it has replies we have
		 * them all and they're in the replies map. */
		l := c.in.replies[request]
		logPrintln("reply in replies map")
		if l != nil {
			// pop front
			front = l.Front()
			if front != nil {
				l.Remove(front)
			}
			if l.Len() == 0 {
				delete(c.in.replies, request)
			}
		}
	} else if request == c.in.requestRead && c.in.currentReply != nil && c.in.currentReply.Front() != nil {
		/* We're currently processing the responses to the request we want, and we
		 * have a reply ready to return. So just return it without blocking. */
		logPrintln("reply in currentReply")
		front = c.in.currentReply.Front()
		c.in.currentReply.Remove(front)

	} else if request == c.in.requestCompleted {
		/* We know this request can't have any more replies, and we've already
		 * established it doesn't have a reply now. Don't bother blocking. */
		logPrintln("request completed")
		stop = true
		return
	} else {
		/* We may have more replies on the way for this request: block until we're sure. */
		// stop = false
		return
	}

	if front != nil {
		replyBuf = front.Value.([]byte)
		respType := replyBuf[0]
		if respType == ResponseTypeError {
			isErr = true
		}
	}
	stop = true
	return
}

func (c *Conn) waitForReply(request SeqNum) (replyBuf []byte, err error) {
	err = c.out.flushTo(request)
	if err != nil {
		c.Close()
		return nil, err
	}

	cond := sync.NewCond(&c.ioMu)
	r := c.in.insertNewReader(request, cond)
	var isErr bool
	var stop bool
	for {
		replyBuf, isErr, stop = c.pollForReply(request)
		if stop {
			break
		}

		logPrintf("waitForReply reader %d wait\n", request)
		cond.Wait()
	}
	if c.isClosed() {
		return nil, errConnClosed
	}
	c.in.removeReader(r)
	c.in.wakeUpNextReader()

	if isErr {
		return nil, c.NewError(replyBuf)
	}
	return replyBuf, nil
}

func (c *Conn) WaitForReply(request SeqNum) (replyBuf []byte, err error) {
	if c.isClosed() {
		return nil, errConnClosed
	}
	err = request.err()
	if err != nil {
		return nil, err
	}

	c.ioMu.Lock()
	replyBuf, err = c.waitForReply(request)
	c.ioMu.Unlock()
	return
}

type VoidCookie SeqNum

func (cookie VoidCookie) Check(c *Conn) error {
	if c.isClosed() {
		return errConnClosed
	}
	err := SeqNum(cookie).err()
	if err != nil {
		return err
	}

	return c.requestCheck(SeqNum(cookie))
}

func (c *Conn) requestCheck(request SeqNum) error {
	c.ioMu.Lock()
	if request >= c.in.requestExpected &&
		request > c.in.requestCompleted {

		// send sync request
		c.sendSync()
		err := c.out.flushTo(c.out.request)
		if err != nil {
			c.ioMu.Unlock()
			return err
		}
	}
	replyBuf, err := c.waitForReply(request)
	c.ioMu.Unlock()

	if err != nil {
		return err
	} else {
		// if not err, replyBuf must be nil
		if replyBuf != nil {
			panic("replyBuf is not nil")
		}
	}
	return nil
}

func (c *Conn) pollForEvent() (ev GenericEvent) {
	f := c.in.events.Front()
	if f != nil {
		ev = f.Value.(GenericEvent)
		c.in.events.Remove(f)
	}
	return
}

func (c *Conn) waitForEvent() (ev GenericEvent) {
	for {
		if c.isClosed() {
			return nil
		}
		ev = c.pollForEvent()
		if ev != nil {
			break
		}
		c.in.eventsCond.Wait()
	}
	return ev
}

func (c *Conn) eventSendLoop() {
	for {
		c.ioMu.Lock()
		ev := c.waitForEvent()
		c.ioMu.Unlock()
		if ev == nil {
			c.in.closeEventChans()
			return
		}
		c.in.sendEvent(ev)
	}
}
