// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

import (
	"sync"
)

type AtomCache struct {
	atoms     map[string]Atom
	atomNames map[Atom]string
	mu        sync.RWMutex
}

func (ac *AtomCache) getVal(name string) (val Atom, ok bool) {
	ac.mu.RLock()
	val, ok = ac.atoms[name]
	ac.mu.RUnlock()
	return
}

func (ac *AtomCache) getName(val Atom) (name string, ok bool) {
	ac.mu.RLock()
	name, ok = ac.atomNames[val]
	ac.mu.RUnlock()
	return
}

func (ac *AtomCache) store(name string, val Atom) {
	ac.mu.Lock()
	ac.atoms[name] = val
	ac.atomNames[val] = name
	ac.mu.Unlock()
}

func NewAtomCache() *AtomCache {
	ac := &AtomCache{
		atoms:     make(map[string]Atom),
		atomNames: make(map[Atom]string),
	}
	return ac
}

var defaultAtomCache *AtomCache
var defaultAtomCacheMu sync.Mutex

func (c *Conn) getAtomCache() *AtomCache {
	c.atomCacheMu.Lock()
	if c.atomCache == nil {
		// try default atom cache
		defaultAtomCacheMu.Lock()
		if defaultAtomCache == nil {
			defaultAtomCache = NewAtomCache()
		}
		c.atomCache = defaultAtomCache
		defaultAtomCacheMu.Unlock()
	}
	v := c.atomCache
	c.atomCacheMu.Unlock()
	return v
}

func (c *Conn) SetAtomCache(ac *AtomCache) {
	c.atomCacheMu.Lock()
	c.atomCache = ac
	c.atomCacheMu.Unlock()
}

func (c *Conn) GetAtomCache() (ac *AtomCache) {
	c.atomCacheMu.Lock()
	ac = c.atomCache
	c.atomCacheMu.Unlock()
	return
}

func (c *Conn) GetAtom(name string) (Atom, error) {
	return c.getAtom(false, name)
}

func (c *Conn) GetAtomExisting(name string) (Atom, error) {
	return c.getAtom(true, name)
}

func (c *Conn) getAtom(onlyIfExists bool, name string) (Atom, error) {
	ac := c.getAtomCache()
	val, ok := ac.getVal(name)
	if ok {
		return val, nil
	}

	reply, err := InternAtom(c, onlyIfExists, name).Reply(c)
	if err != nil {
		return AtomNone, err
	}
	ac.store(name, reply.Atom)
	return reply.Atom, nil
}

func (c *Conn) GetAtomName(atom Atom) (string, error) {
	ac := c.getAtomCache()
	name, ok := ac.getName(atom)
	if ok {
		return name, nil
	}

	reply, err := GetAtomName(c, atom).Reply(c)
	if err != nil {
		return "", err
	}
	ac.store(reply.Name, atom)
	return reply.Name, nil
}
