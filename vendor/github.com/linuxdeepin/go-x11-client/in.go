// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

import (
	"container/list"
	"fmt"
	"os"
	"sync"
)

type in struct {
	requestExpected  SeqNum
	requestRead      SeqNum
	requestCompleted SeqNum

	currentReply   *list.List
	replies        map[SeqNum]*list.List
	pendingReplies *list.List

	readers    *list.List
	events     *list.List
	eventsCond *sync.Cond

	eventChans []chan<- GenericEvent
	chansMu    sync.Mutex
}

func (in *in) init(ioMu *sync.Mutex) {
	in.replies = make(map[SeqNum]*list.List)
	in.pendingReplies = list.New()
	in.readers = list.New()
	in.events = list.New()
	in.eventsCond = sync.NewCond(ioMu)
}

func (in *in) addEventChan(eventChan chan<- GenericEvent) {
	in.chansMu.Lock()

	for _, ch := range in.eventChans {
		if ch == eventChan {
			// exist
			in.chansMu.Unlock()
			return
		}
	}

	// COW
	newEventChans := make([]chan<- GenericEvent, len(in.eventChans)+1)
	copy(newEventChans, in.eventChans)
	newEventChans[len(newEventChans)-1] = eventChan
	in.eventChans = newEventChans
	in.chansMu.Unlock()
}

func (in *in) removeEventChan(eventChan chan<- GenericEvent) {
	in.chansMu.Lock()

	chans := in.eventChans
	var found bool
	for _, ch := range chans {
		if ch == eventChan {
			found = true
			break
		}
	}

	if !found {
		// not found
		in.chansMu.Unlock()
		return
	}

	// COW
	newEventChans := make([]chan<- GenericEvent, 0, len(in.eventChans)-1)
	for _, ch := range chans {
		if ch != eventChan {
			newEventChans = append(newEventChans, ch)
		}
	}
	in.eventChans = newEventChans
	in.chansMu.Unlock()
}

func (in *in) addEvent(e GenericEvent) {
	logPrintln("add event", e)
	if in.events.Len() > 0xfff {
		_, _ = fmt.Fprintf(os.Stderr, "<warning> too many events are not processed, len: %d\n",
			in.events.Len())
	}
	in.events.PushBack(e)
	in.eventsCond.Signal()
}

func (in *in) sendEvent(e GenericEvent) {
	in.chansMu.Lock()
	eventChans := in.eventChans
	in.chansMu.Unlock()

	for _, ch := range eventChans {
		ch <- e
	}
}

func (in *in) closeEventChans() {
	in.chansMu.Lock()

	for _, ch := range in.eventChans {
		close(ch)
	}

	in.chansMu.Unlock()
}

type ReplyReader struct {
	request SeqNum
	cond    *sync.Cond
}

func (in *in) insertNewReader(request SeqNum, cond *sync.Cond) *ReplyReader {
	r := &ReplyReader{
		request: request,
		cond:    cond,
	}

	var mark *list.Element
	l := in.readers
	for e := l.Front(); e != nil; e = e.Next() {
		reader := e.Value.(*ReplyReader)
		if reader.request >= request {
			mark = e
			break
		}
	}

	if mark != nil {
		logPrintf("insertNewReader %d before %d\n", request, mark.Value.(*ReplyReader).request)
		l.InsertBefore(r, mark)
	} else {
		logPrintf("insertNewReader %d at end\n", request)
		l.PushBack(r)
	}
	return r
}

func (in *in) removeReader(r *ReplyReader) {
	l := in.readers
	for e := l.Front(); e != nil; e = e.Next() {
		reader := e.Value.(*ReplyReader)
		if reader.request == r.request {
			logPrintln("remove reader", reader.request)
			l.Remove(e)
			break
		}
	}
}

func (in *in) removeFinishedReaders() {
	l := in.readers
	e := l.Front()
	for e != nil {
		reader := e.Value.(*ReplyReader)
		if reader.request <= in.requestCompleted {
			reader.cond.Signal()
			logPrintln("remove finished reader", reader.request)
			tmp := e
			e = e.Next()
			l.Remove(tmp)
		} else {
			break
		}
	}
}

func (in *in) wakeUpAllReaders() {
	l := in.readers
	for e := l.Front(); e != nil; e = e.Next() {
		reader := e.Value.(*ReplyReader)
		reader.cond.Signal()
	}
}

func (in *in) wakeUpNextReader() {
	if in.readers.Front() != nil {
		reader := in.readers.Front().Value.(*ReplyReader)
		logPrintln("wake up next reader", reader.request)
		reader.cond.Signal()
	}
}

type PendingReply struct {
	firstRequest SeqNum
	lastRequest  SeqNum
	workaround   uint
	flags        uint
}

func (in *in) expectReply(request SeqNum, workaround uint, flags uint) {
	pend := &PendingReply{
		firstRequest: request,
		lastRequest:  request,
		workaround:   workaround,
		flags:        flags,
	}
	in.pendingReplies.PushBack(pend)
}

func (in *in) removeFinishedPendingReplies() {
	l := in.pendingReplies
	e := l.Front()
	for e != nil {
		pend := e.Value.(*PendingReply)
		if pend.lastRequest <= in.requestCompleted {
			// remove pend from list
			tmp := e
			e = e.Next()
			l.Remove(tmp)
		} else {
			break
		}
	}
}
