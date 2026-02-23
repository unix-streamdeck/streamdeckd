// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

func CreateWindow(conn *Conn, depth uint8, wid, parent Window, x, y int16, width, height, borderWidth uint16, class uint16, visual VisualID, valueMask uint32, valueList []uint32) {
	headerData, body := encodeCreateWindow(depth, wid, parent, x, y, width, height, borderWidth, class, visual, valueMask, valueList)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CreateWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func CreateWindowChecked(conn *Conn, depth uint8, wid, parent Window, x, y int16, width, height, borderWidth uint16, class uint16, visual VisualID, valueMask uint32, valueList []uint32) VoidCookie {
	headerData, body := encodeCreateWindow(depth, wid, parent, x, y, width, height, borderWidth, class, visual, valueMask, valueList)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CreateWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func ChangeWindowAttributes(conn *Conn, window Window, valueMask uint32, valueList []uint32) {
	headerData, body := encodeChangeWindowAttributes(window, valueMask, valueList)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ChangeWindowAttributesOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func ChangeWindowAttributesChecked(conn *Conn, window Window, valueMask uint32, valueList []uint32) VoidCookie {
	headerData, body := encodeChangeWindowAttributes(window, valueMask, valueList)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ChangeWindowAttributesOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GetWindowAttributes(conn *Conn, window Window) GetWindowAttributesCookie {
	headerData, body := encodeGetWindowAttributes(window)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: GetWindowAttributesOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return GetWindowAttributesCookie(seq)
}

func (cookie GetWindowAttributesCookie) Reply(conn *Conn) (*GetWindowAttributesReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply GetWindowAttributesReply
	err = readGetWindowAttributesReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func DestroyWindow(conn *Conn, window Window) {
	headerData, body := encodeDestroyWindow(window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: DestroyWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func DestroyWindowChecked(conn *Conn, window Window) VoidCookie {
	headerData, body := encodeDestroyWindow(window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: DestroyWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func DestroySubwindows(conn *Conn, window Window) {
	headerData, body := encodeDestroySubwindows(window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: DestroySubwindowsOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func DestroySubwindowsChecked(conn *Conn, window Window) VoidCookie {
	headerData, body := encodeDestroySubwindows(window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: DestroySubwindowsOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func ChangeSaveSet(conn *Conn, mode uint8, window Window) {
	headerData, body := encodeChangeSaveSet(mode, window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ChangeSaveSetOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func ChangeSaveSetChecked(conn *Conn, mode uint8, window Window) VoidCookie {
	headerData, body := encodeChangeSaveSet(mode, window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ChangeSaveSetOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func ReparentWindow(conn *Conn, window, parent Window, x, y int16) {
	headerData, body := encodeReparentWindow(window, parent, x, y)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ReparentWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func ReparentWindowChecked(conn *Conn, window, parent Window, x, y int16) VoidCookie {
	headerData, body := encodeReparentWindow(window, parent, x, y)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ReparentWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func MapWindow(conn *Conn, window Window) {
	headerData, body := encodeMapWindow(window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: MapWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func MapWindowChecked(conn *Conn, window Window) VoidCookie {
	headerData, body := encodeMapWindow(window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: MapWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func MapSubwindows(conn *Conn, window Window) {
	headerData, body := encodeMapSubwindows(window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: MapSubwindowsOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func MapSubwindowsChecked(conn *Conn, window Window) VoidCookie {
	headerData, body := encodeMapSubwindows(window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: MapSubwindowsOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func UnmapWindow(conn *Conn, window Window) {
	headerData, body := encodeUnmapWindow(window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UnmapWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func UnmapWindowChecked(conn *Conn, window Window) VoidCookie {
	headerData, body := encodeUnmapWindow(window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UnmapWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func UnmapSubwindows(conn *Conn, window Window) {
	headerData, body := encodeUnmapSubwindows(window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UnmapSubwindowsOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func UnmapSubwindowsChecked(conn *Conn, window Window) VoidCookie {
	headerData, body := encodeUnmapSubwindows(window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UnmapSubwindowsOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func ConfigureWindow(conn *Conn, window Window, valueMask uint16, valueList []uint32) {
	headerData, body := encodeConfigureWindow(window, valueMask, valueList)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ConfigureWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func ConfigureWindowChecked(conn *Conn, window Window, valueMask uint16, valueList []uint32) VoidCookie {
	headerData, body := encodeConfigureWindow(window, valueMask, valueList)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ConfigureWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func CirculateWindow(conn *Conn, direction uint8, window Window) {
	headerData, body := encodeCirculateWindow(direction, window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CirculateWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func CirculateWindowChecked(conn *Conn, direction uint8, window Window) VoidCookie {
	headerData, body := encodeCirculateWindow(direction, window)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CirculateWindowOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GetGeometry(conn *Conn, drawable Drawable) GetGeometryCookie {
	headerData, body := encodeGetGeometry(drawable)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: GetGeometryOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return GetGeometryCookie(seq)
}

func (cookie GetGeometryCookie) Reply(conn *Conn) (*GetGeometryReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply GetGeometryReply
	err = readGetGeometryReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func QueryTree(conn *Conn, window Window) QueryTreeCookie {
	headerData, body := encodeQueryTree(window)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: QueryTreeOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return QueryTreeCookie(seq)
}

func (cookie QueryTreeCookie) Reply(conn *Conn) (*QueryTreeReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply QueryTreeReply
	err = readQueryTreeReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func InternAtom(conn *Conn, onlyIfExists bool, name string) InternAtomCookie {
	headerData, body := encodeInternAtom(onlyIfExists, name)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: InternAtomOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return InternAtomCookie(seq)
}

func (cookie InternAtomCookie) Reply(conn *Conn) (*InternAtomReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply InternAtomReply
	err = readInternAtomReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func GetAtomName(conn *Conn, atom Atom) GetAtomNameCookie {
	headerData, body := encodeGetAtomName(atom)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: GetAtomNameOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return GetAtomNameCookie(seq)
}

func (cookie GetAtomNameCookie) Reply(conn *Conn) (*GetAtomNameReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply GetAtomNameReply
	err = readGetAtomNameReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func ChangeProperty(conn *Conn, mode uint8, window Window, property, type0 Atom, format uint8, data []byte) {
	headerData, body := encodeChangeProperty(mode, window, property, type0, format, data)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ChangePropertyOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func ChangePropertyChecked(conn *Conn, mode uint8, window Window, property, type0 Atom, format uint8, data []byte) VoidCookie {
	headerData, body := encodeChangeProperty(mode, window, property, type0, format, data)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ChangePropertyOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func DeleteProperty(conn *Conn, window Window, property Atom) {
	headerData, body := encodeDeleteProperty(window, property)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: DeletePropertyOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func DeletePropertyChecked(conn *Conn, window Window, property Atom) VoidCookie {
	headerData, body := encodeDeleteProperty(window, property)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: DeletePropertyOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GetProperty(conn *Conn, delete bool, window Window, property, type0 Atom, longOffset, longLength uint32) GetPropertyCookie {
	headerData, body := encodeGetProperty(delete, window, property, type0, longOffset, longLength)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: GetPropertyOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return GetPropertyCookie(seq)
}

func (cookie GetPropertyCookie) Reply(conn *Conn) (*GetPropertyReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply GetPropertyReply
	err = readGetPropertyReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func ListProperties(conn *Conn, window Window) ListPropertiesCookie {
	headerData, body := encodeListProperties(window)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: ListPropertiesOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return ListPropertiesCookie(seq)
}

func (cookie ListPropertiesCookie) Reply(conn *Conn) (*ListPropertiesReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply ListPropertiesReply
	err = readListPropertiesReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func SetSelectionOwner(conn *Conn, owner Window, selection Atom, time Timestamp) {
	headerData, body := encodeSetSelectionOwner(owner, selection, time)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SetSelectionOwnerOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func SetSelectionOwnerChecked(conn *Conn, owner Window, selection Atom, time Timestamp) VoidCookie {
	headerData, body := encodeSetSelectionOwner(owner, selection, time)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SetSelectionOwnerOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GetSelectionOwner(conn *Conn, selection Atom) GetSelectionOwnerCookie {
	headerData, body := encodeGetSelectionOwner(selection)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: GetSelectionOwnerOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return GetSelectionOwnerCookie(seq)
}

func (cookie GetSelectionOwnerCookie) Reply(conn *Conn) (*GetSelectionOwnerReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply GetSelectionOwnerReply
	err = readGetSelectionOwnerReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func ConvertSelection(conn *Conn, requestor Window, selection, target, property Atom, time Timestamp) {
	headerData, body := encodeConvertSelection(requestor, selection, target, property, time)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ConvertSelectionOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func ConvertSelectionChecked(conn *Conn, requestor Window, selection, target, property Atom, time Timestamp) VoidCookie {
	headerData, body := encodeConvertSelection(requestor, selection, target, property, time)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ConvertSelectionOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func SendEvent(conn *Conn, propagate bool, destination Window, eventMask uint32, event []byte) {
	headerData, body := encodeSendEvent(propagate, destination, eventMask, event)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SendEventOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func SendEventChecked(conn *Conn, propagate bool, destination Window, eventMask uint32, event []byte) VoidCookie {
	headerData, body := encodeSendEvent(propagate, destination, eventMask, event)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SendEventOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GrabPointer(conn *Conn, ownerEvents bool, grabWindow Window, eventMask uint16, pointerMode, keyboardMode uint8, confineTo Window, cursor Cursor, time Timestamp) GrabPointerCookie {
	headerData, body := encodeGrabPointer(ownerEvents, grabWindow, eventMask, pointerMode, keyboardMode, confineTo, cursor, time)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: GrabPointerOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return GrabPointerCookie(seq)
}

func (cookie GrabPointerCookie) Reply(conn *Conn) (*GrabPointerReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply GrabPointerReply
	err = readGrabPointerReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func UngrabPointer(conn *Conn, time Timestamp) {
	headerData, body := encodeUngrabPointer(time)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UngrabPointerOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func UngrabPointerChecked(conn *Conn, time Timestamp) VoidCookie {
	headerData, body := encodeUngrabPointer(time)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UngrabPointerOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GrabButton(conn *Conn, ownerEvents bool, grabWindow Window, eventMask uint16, pointerMode, keyboardMode uint8, confineTo Window, cursor Cursor, button uint8, modifiers uint16) {
	headerData, body := encodeGrabButton(ownerEvents, grabWindow, eventMask, pointerMode, keyboardMode, confineTo, cursor, button, modifiers)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: GrabButtonOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func GrabButtonChecked(conn *Conn, ownerEvents bool, grabWindow Window, eventMask uint16, pointerMode, keyboardMode uint8, confineTo Window, cursor Cursor, button uint8, modifiers uint16) VoidCookie {
	headerData, body := encodeGrabButton(ownerEvents, grabWindow, eventMask, pointerMode, keyboardMode, confineTo, cursor, button, modifiers)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: GrabButtonOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func UngrabButton(conn *Conn, button uint8, grabWindow Window, modifiers uint16) {
	headerData, body := encodeUngrabButton(button, grabWindow, modifiers)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UngrabButtonOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func UngrabButtonChecked(conn *Conn, button uint8, grabWindow Window, modifiers uint16) VoidCookie {
	headerData, body := encodeUngrabButton(button, grabWindow, modifiers)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UngrabButtonOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func ChangeActivePointerGrab(conn *Conn, cursor Cursor, time Timestamp, eventMask uint16) {
	headerData, body := encodeChangeActivePointerGrab(cursor, time, eventMask)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ChangeActivePointerGrabOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func ChangeActivePointerGrabChecked(conn *Conn, cursor Cursor, time Timestamp, eventMask uint16) VoidCookie {
	headerData, body := encodeChangeActivePointerGrab(cursor, time, eventMask)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ChangeActivePointerGrabOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GrabKeyboard(conn *Conn, ownerEvents bool, grabWindow Window, time Timestamp, pointerMode, keyboardMode uint8) GrabKeyboardCookie {
	headerData, body := encodeGrabKeyboard(ownerEvents, grabWindow, time, pointerMode, keyboardMode)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: GrabKeyboardOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return GrabKeyboardCookie(seq)
}

func (cookie GrabKeyboardCookie) Reply(conn *Conn) (*GrabKeyboardReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply GrabKeyboardReply
	err = readGrabKeyboardReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func UngrabKeyboard(conn *Conn, time Timestamp) {
	headerData, body := encodeUngrabKeyboard(time)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UngrabKeyboardOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func UngrabKeyboardChecked(conn *Conn, time Timestamp) VoidCookie {
	headerData, body := encodeUngrabKeyboard(time)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UngrabKeyboardOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GrabKey(conn *Conn, ownerEvents bool, grabWindow Window, modifiers uint16, key Keycode, pointerMode, keyboardMode uint8) {
	headerData, body := encodeGrabKey(ownerEvents, grabWindow, modifiers, key, pointerMode, keyboardMode)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: GrabKeyOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func GrabKeyChecked(conn *Conn, ownerEvents bool, grabWindow Window, modifiers uint16, key Keycode, pointerMode, keyboardMode uint8) VoidCookie {
	headerData, body := encodeGrabKey(ownerEvents, grabWindow, modifiers, key, pointerMode, keyboardMode)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: GrabKeyOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func UngrabKey(conn *Conn, key Keycode, grabWindow Window, modifiers uint16) {
	headerData, body := encodeUngrabKey(key, grabWindow, modifiers)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UngrabKeyOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func UngrabKeyChecked(conn *Conn, key Keycode, grabWindow Window, modifiers uint16) VoidCookie {
	headerData, body := encodeUngrabKey(key, grabWindow, modifiers)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UngrabKeyOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func AllowEvents(conn *Conn, mode uint8, time Timestamp) {
	headerData, body := encodeAllowEvents(mode, time)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: AllowEventsOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func AllowEventsChecked(conn *Conn, mode uint8, time Timestamp) VoidCookie {
	headerData, body := encodeAllowEvents(mode, time)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: AllowEventsOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GrabServer(conn *Conn) {
	headerData, body := encodeGrabServer()
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: GrabServerOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func GrabServerChecked(conn *Conn) VoidCookie {
	headerData, body := encodeGrabServer()
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: GrabServerOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func UngrabServer(conn *Conn) {
	headerData, body := encodeUngrabServer()
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UngrabServerOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func UngrabServerChecked(conn *Conn) VoidCookie {
	headerData, body := encodeUngrabServer()
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: UngrabServerOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func QueryPointer(conn *Conn, window Window) QueryPointerCookie {
	headerData, body := encodeQueryPointer(window)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: QueryPointerOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return QueryPointerCookie(seq)
}

func (cookie QueryPointerCookie) Reply(conn *Conn) (*QueryPointerReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply QueryPointerReply
	err = readQueryPointerReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func GetMotionEvents(conn *Conn, window Window, start, stop Timestamp) GetMotionEventsCookie {
	headerData, body := encodeGetMotionEvents(window, start, stop)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: GetMotionEventsOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return GetMotionEventsCookie(seq)
}

func (cookie GetMotionEventsCookie) Reply(conn *Conn) (*GetMotionEventsReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply GetMotionEventsReply
	err = readGetMotionEventsReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func TranslateCoordinates(conn *Conn, srcWindow, dstWindow Window, srcX, srcY int16) TranslateCoordinatesCookie {
	headerData, body := encodeTranslateCoordinates(srcWindow, dstWindow, srcX, srcY)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: TranslateCoordinatesOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return TranslateCoordinatesCookie(seq)
}

func (cookie TranslateCoordinatesCookie) Reply(conn *Conn) (*TranslateCoordinatesReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply TranslateCoordinatesReply
	err = readTranslateCoordinatesReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func WarpPointer(conn *Conn, srcWindow, dstWindow Window, srcX, srcY int16, srcWidth, srcHeight uint16, dstX, dstY int16) {
	headerData, body := encodeWarpPointer(srcWindow, dstWindow, srcX, srcY, srcWidth, srcHeight, dstX, dstY)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: WarpPointerOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func WarpPointerChecked(conn *Conn, srcWindow, dstWindow Window, srcX, srcY int16, srcWidth, srcHeight uint16, dstX, dstY int16) VoidCookie {
	headerData, body := encodeWarpPointer(srcWindow, dstWindow, srcX, srcY, srcWidth, srcHeight, dstX, dstY)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: WarpPointerOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func SetInputFocus(conn *Conn, revertTo uint8, focus Window, time Timestamp) {
	headerData, body := encodeSetInputFocus(revertTo, focus, time)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SetInputFocusOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func SetInputFocusChecked(conn *Conn, revertTo uint8, focus Window, time Timestamp) VoidCookie {
	headerData, body := encodeSetInputFocus(revertTo, focus, time)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SetInputFocusOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GetInputFocus(conn *Conn) GetInputFocusCookie {
	headerData, body := encodeGetInputFocus()
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: GetInputFocusOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return GetInputFocusCookie(seq)
}

func (cookie GetInputFocusCookie) Reply(conn *Conn) (*GetInputFocusReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply GetInputFocusReply
	err = readGetInputFocusReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func QueryKeymap(conn *Conn) QueryKeymapCookie {
	headerData, body := encodeQueryKeymap()
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: QueryKeymapOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return QueryKeymapCookie(seq)
}

func (cookie QueryKeymapCookie) Reply(conn *Conn) (*QueryKeymapReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply QueryKeymapReply
	err = readQueryKeymapReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func OpenFont(conn *Conn, fid Font, name string) {
	headerData, body := encodeOpenFont(fid, name)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: OpenFontOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func OpenFontChecked(conn *Conn, fid Font, name string) VoidCookie {
	headerData, body := encodeOpenFont(fid, name)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: OpenFontOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func CloseFont(conn *Conn, font Font) {
	headerData, body := encodeCloseFont(font)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CloseFontOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func CloseFontChecked(conn *Conn, font Font) VoidCookie {
	headerData, body := encodeCloseFont(font)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CloseFontOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func QueryFont(conn *Conn, font Fontable) QueryFontCookie {
	headerData, body := encodeQueryFont(font)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: QueryFontOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return QueryFontCookie(seq)
}

func (cookie QueryFontCookie) Reply(conn *Conn) (*QueryFontReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply QueryFontReply
	err = readQueryFontReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func ListFonts(conn *Conn, maxNames uint16, pattern string) ListFontsCookie {
	headerData, body := encodeListFonts(maxNames, pattern)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: ListFontsOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return ListFontsCookie(seq)
}

func (cookie ListFontsCookie) Reply(conn *Conn) (*ListFontsReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply ListFontsReply
	err = readListFontsReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func ListFontsWithInfo(conn *Conn, maxNames uint16, pattern string) ListFontsWithInfoCookie {
	headerData, body := encodeListFontsWithInfo(maxNames, pattern)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: ListFontsWithInfoOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return ListFontsWithInfoCookie(seq)
}

func (cookie ListFontsWithInfoCookie) Reply(conn *Conn) (*ListFontsWithInfoReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply ListFontsWithInfoReply
	err = readListFontsWithInfoReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func SetFontPath(conn *Conn, paths []string) {
	headerData, body := encodeSetFontPath(paths)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SetFontPathOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func SetFontPathChecked(conn *Conn, paths []string) VoidCookie {
	headerData, body := encodeSetFontPath(paths)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SetFontPathOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GetFontPath(conn *Conn) GetFontPathCookie {
	headerData, body := encodeGetFontPath()
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: GetFontPathOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return GetFontPathCookie(seq)
}

func (cookie GetFontPathCookie) Reply(conn *Conn) (*GetFontPathReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply GetFontPathReply
	err = readGetFontPathReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func CreatePixmap(conn *Conn, depth uint8, pid Pixmap, drawable Drawable, width, height uint16) {
	headerData, body := encodeCreatePixmap(depth, pid, drawable, width, height)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CreatePixmapOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func CreatePixmapChecked(conn *Conn, depth uint8, pid Pixmap, drawable Drawable, width, height uint16) VoidCookie {
	headerData, body := encodeCreatePixmap(depth, pid, drawable, width, height)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CreatePixmapOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func FreePixmap(conn *Conn, pixmap Pixmap) {
	headerData, body := encodeFreePixmap(pixmap)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: FreePixmapOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func FreePixmapChecked(conn *Conn, pixmap Pixmap) VoidCookie {
	headerData, body := encodeFreePixmap(pixmap)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: FreePixmapOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func CreateGC(conn *Conn, cid GContext, drawable Drawable, valueMask uint32, valueList []uint32) {
	headerData, body := encodeCreateGC(cid, drawable, valueMask, valueList)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CreateGCOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func CreateGCChecked(conn *Conn, cid GContext, drawable Drawable, valueMask uint32, valueList []uint32) VoidCookie {
	headerData, body := encodeCreateGC(cid, drawable, valueMask, valueList)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CreateGCOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func ChangeGC(conn *Conn, gc GContext, valueMask uint32, valueList []uint32) {
	headerData, body := encodeChangeGC(gc, valueMask, valueList)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ChangeGCOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func ChangeGCChecked(conn *Conn, gc GContext, valueMask uint32, valueList []uint32) VoidCookie {
	headerData, body := encodeChangeGC(gc, valueMask, valueList)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ChangeGCOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func CopyGC(conn *Conn, srcGC, dstGC GContext, valueMask uint32) {
	headerData, body := encodeCopyGC(srcGC, dstGC, valueMask)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CopyGCOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func CopyGCChecked(conn *Conn, srcGC, dstGC GContext, valueMask uint32) VoidCookie {
	headerData, body := encodeCopyGC(srcGC, dstGC, valueMask)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CopyGCOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func SetDashes(conn *Conn, gc GContext, dashOffset uint16, dashes []uint8) {
	headerData, body := encodeSetDashes(gc, dashOffset, dashes)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SetDashesOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func SetDashesChecked(conn *Conn, gc GContext, dashOffset uint16, dashes []uint8) VoidCookie {
	headerData, body := encodeSetDashes(gc, dashOffset, dashes)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SetDashesOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func SetClipRectangles(conn *Conn, ordering uint8, gc GContext, clipXOrigin, clipYOrigin int16, rectangles []Rectangle) {
	headerData, body := encodeSetClipRectangles(ordering, gc, clipXOrigin, clipYOrigin, rectangles)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SetClipRectanglesOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func SetClipRectanglesChecked(conn *Conn, ordering uint8, gc GContext, clipXOrigin, clipYOrigin int16, rectangles []Rectangle) VoidCookie {
	headerData, body := encodeSetClipRectangles(ordering, gc, clipXOrigin, clipYOrigin, rectangles)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SetClipRectanglesOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func FreeGC(conn *Conn, gc GContext) {
	headerData, body := encodeFreeGC(gc)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: FreeGCOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func FreeGCChecked(conn *Conn, gc GContext) VoidCookie {
	headerData, body := encodeFreeGC(gc)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: FreeGCOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func ClearArea(conn *Conn, exposures bool, window Window, x, y int16, width, height uint16) {
	headerData, body := encodeClearArea(exposures, window, x, y, width, height)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ClearAreaOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func ClearAreaChecked(conn *Conn, exposures bool, window Window, x, y int16, width, height uint16) VoidCookie {
	headerData, body := encodeClearArea(exposures, window, x, y, width, height)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ClearAreaOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func CopyArea(conn *Conn, srcDrawable, dstDrawable Drawable, gc GContext, srcX, srcY, dstX, dstY int16, width, height uint16) {
	headerData, body := encodeCopyArea(srcDrawable, dstDrawable, gc, srcX, srcY, dstX, dstY, width, height)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CopyAreaOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func CopyAreaChecked(conn *Conn, srcDrawable, dstDrawable Drawable, gc GContext, srcX, srcY, dstX, dstY int16, width, height uint16) VoidCookie {
	headerData, body := encodeCopyArea(srcDrawable, dstDrawable, gc, srcX, srcY, dstX, dstY, width, height)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CopyAreaOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func CopyPlane(conn *Conn, srcDrawable, dstDrawable Drawable, gc GContext, srcX, srcY, dstX, dstY int16, width, height uint16, bitPlane uint32) {
	headerData, body := encodeCopyPlane(srcDrawable, dstDrawable, gc, srcX, srcY, dstX, dstY, width, height, bitPlane)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CopyPlaneOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func CopyPlaneChecked(conn *Conn, srcDrawable, dstDrawable Drawable, gc GContext, srcX, srcY, dstX, dstY int16, width, height uint16, bitPlane uint32) VoidCookie {
	headerData, body := encodeCopyPlane(srcDrawable, dstDrawable, gc, srcX, srcY, dstX, dstY, width, height, bitPlane)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: CopyPlaneOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func PolyPoint(conn *Conn, coordinateMode uint8, drawable Drawable, gc GContext, points []Point) {
	headerData, body := encodePolyPoint(coordinateMode, drawable, gc, points)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolyPointOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func PolyPointChecked(conn *Conn, coordinateMode uint8, drawable Drawable, gc GContext, points []Point) VoidCookie {
	headerData, body := encodePolyPoint(coordinateMode, drawable, gc, points)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolyPointOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func PolyLine(conn *Conn, coordinateMode uint8, drawable Drawable, gc GContext, points []Point) {
	headerData, body := encodePolyLine(coordinateMode, drawable, gc, points)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolyLineOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func PolyLineChecked(conn *Conn, coordinateMode uint8, drawable Drawable, gc GContext, points []Point) VoidCookie {
	headerData, body := encodePolyLine(coordinateMode, drawable, gc, points)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolyLineOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func PolySegment(conn *Conn, drawable Drawable, gc GContext, segments []Segment) {
	headerData, body := encodePolySegment(drawable, gc, segments)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolySegmentOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func PolySegmentChecked(conn *Conn, drawable Drawable, gc GContext, segments []Segment) VoidCookie {
	headerData, body := encodePolySegment(drawable, gc, segments)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolySegmentOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func PolyRectangle(conn *Conn, drawable Drawable, gc GContext, rectangles []Rectangle) {
	headerData, body := encodePolyRectangle(drawable, gc, rectangles)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolyRectangleOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func PolyRectangleChecked(conn *Conn, drawable Drawable, gc GContext, rectangles []Rectangle) VoidCookie {
	headerData, body := encodePolyRectangle(drawable, gc, rectangles)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolyRectangleOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func PolyArc(conn *Conn, drawable Drawable, gc GContext, arcs []Arc) {
	headerData, body := encodePolyArc(drawable, gc, arcs)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolyArcOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func PolyArcChecked(conn *Conn, drawable Drawable, gc GContext, arcs []Arc) VoidCookie {
	headerData, body := encodePolyArc(drawable, gc, arcs)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolyArcOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func FillPoly(conn *Conn, drawable Drawable, gc GContext, shape uint8, coordinateMode uint8, points []Point) {
	headerData, body := encodeFillPoly(drawable, gc, shape, coordinateMode, points)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: FillPolyOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func FillPolyChecked(conn *Conn, drawable Drawable, gc GContext, shape uint8, coordinateMode uint8, points []Point) VoidCookie {
	headerData, body := encodeFillPoly(drawable, gc, shape, coordinateMode, points)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: FillPolyOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func PolyFillRectangle(conn *Conn, drawable Drawable, gc GContext, rectangles []Rectangle) {
	headerData, body := encodePolyFillRectangle(drawable, gc, rectangles)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolyFillRectangleOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func PolyFillRectangleChecked(conn *Conn, drawable Drawable, gc GContext, rectangles []Rectangle) VoidCookie {
	headerData, body := encodePolyFillRectangle(drawable, gc, rectangles)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolyFillRectangleOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func PolyFillArc(conn *Conn, drawable Drawable, gc GContext, arcs []Arc) {
	headerData, body := encodePolyFillArc(drawable, gc, arcs)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolyFillArcOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func PolyFillArcChecked(conn *Conn, drawable Drawable, gc GContext, arcs []Arc) VoidCookie {
	headerData, body := encodePolyFillArc(drawable, gc, arcs)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PolyFillArcOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func PutImage(conn *Conn, format uint8, drawable Drawable, gc GContext, width, height uint16, dstX, dstY int16, leftPad, depth uint8, data []byte) {
	headerData, body := encodePutImage(format, drawable, gc, width, height, dstX, dstY, leftPad, depth, data)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PutImageOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func PutImageChecked(conn *Conn, format uint8, drawable Drawable, gc GContext, width, height uint16, dstX, dstY int16, leftPad, depth uint8, data []byte) VoidCookie {
	headerData, body := encodePutImage(format, drawable, gc, width, height, dstX, dstY, leftPad, depth, data)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: PutImageOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GetImage(conn *Conn, format uint8, drawable Drawable, x, y int16, width, height uint16, planeMask uint32) GetImageCookie {
	headerData, body := encodeGetImage(format, drawable, x, y, width, height, planeMask)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: GetImageOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return GetImageCookie(seq)
}

func (cookie GetImageCookie) Reply(conn *Conn) (*GetImageReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply GetImageReply
	err = readGetImageReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func QueryExtension(conn *Conn, name string) QueryExtensionCookie {
	headerData, body := encodeQueryExtension(name)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: QueryExtensionOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return QueryExtensionCookie(seq)
}

func (cookie QueryExtensionCookie) Reply(conn *Conn) (*QueryExtensionReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply QueryExtensionReply
	err = readQueryExtensionReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func ListExtensions(conn *Conn) ListExtensionsCookie {
	headerData, body := encodeListExtensions()
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: ListExtensionsOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return ListExtensionsCookie(seq)
}

func (cookie ListExtensionsCookie) Reply(conn *Conn) (*ListExtensionsReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply ListExtensionsReply
	err = readListExtensionsReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func KillClient(conn *Conn, resource uint32) {
	headerData, body := encodeKillClient(resource)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: KillClientOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func KillClientChecked(conn *Conn, resource uint32) VoidCookie {
	headerData, body := encodeKillClient(resource)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: KillClientOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GetKeyboardMapping(conn *Conn, firstKeycode Keycode, count uint8) GetKeyboardMappingCookie {
	headerData, body := encodeGetKeyboardMapping(firstKeycode, count)
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: GetKeyboardMappingOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return GetKeyboardMappingCookie(seq)
}

func (cookie GetKeyboardMappingCookie) Reply(conn *Conn) (*GetKeyboardMappingReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply GetKeyboardMappingReply
	err = readGetKeyboardMappingReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func SetScreenSaver(conn *Conn, timeout, interval int16, preferBlanking, allowExposures uint8) {
	headerData, body := encodeSetScreenSaver(timeout, interval, preferBlanking, allowExposures)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SetScreenSaverOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func SetScreenSaverChecked(conn *Conn, timeout, interval int16, preferBlanking, allowExposures uint8) VoidCookie {
	headerData, body := encodeSetScreenSaver(timeout, interval, preferBlanking, allowExposures)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: SetScreenSaverOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func GetScreenSaver(conn *Conn) GetScreenSaverCookie {
	headerData, body := encodeGetScreenSaver()
	req := &ProtocolRequest{
		Header: RequestHeader{
			Opcode: GetScreenSaverOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return GetScreenSaverCookie(seq)
}

func (cookie GetScreenSaverCookie) Reply(conn *Conn) (*GetScreenSaverReply, error) {
	replyBuf, err := conn.WaitForReply(SeqNum(cookie))
	if err != nil {
		return nil, err
	}
	r := NewReaderFromData(replyBuf)
	var reply GetScreenSaverReply
	err = readGetScreenSaverReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func ForceScreenSaver(conn *Conn, mode uint8) {
	headerData, body := encodeForceScreenSaver(mode)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ForceScreenSaverOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func ForceScreenSaverChecked(conn *Conn, mode uint8) VoidCookie {
	headerData, body := encodeForceScreenSaver(mode)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ForceScreenSaverOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func NoOperation(conn *Conn, n int) {
	headerData, body := encodeNoOperation(n)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: NoOperationOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func NoOperationChecked(conn *Conn, n int) VoidCookie {
	headerData, body := encodeNoOperation(n)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: NoOperationOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func FreeCursor(conn *Conn, cursor Cursor) {
	headerData, body := encodeFreeCursor(cursor)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: FreeCursorOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func FreeCursorChecked(conn *Conn, cursor Cursor) VoidCookie {
	headerData, body := encodeFreeCursor(cursor)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: FreeCursorOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}

func ChangeHosts(conn *Conn, mode, family uint8, address string) {
	headerData, body := encodeChangeHosts(mode, family, address)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ChangeHostsOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	conn.SendRequest(0, req)
}

func ChangeHostsChecked(conn *Conn, mode, family uint8, address string) VoidCookie {
	headerData, body := encodeChangeHosts(mode, family, address)
	req := &ProtocolRequest{
		NoReply: true,
		Header: RequestHeader{
			Opcode: ChangeHostsOpcode,
			Data:   headerData,
		},
		Body: body,
	}
	seq := conn.SendRequest(RequestChecked, req)
	return VoidCookie(seq)
}
