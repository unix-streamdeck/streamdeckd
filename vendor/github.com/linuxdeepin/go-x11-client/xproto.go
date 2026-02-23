// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

import (
	"bytes"
	"errors"
	"math"
)

type Setup struct {
	ProtocolMajorVersion     uint16
	ProtocolMinorVersion     uint16
	ReleaseNumber            uint32
	ResourceIdBase           uint32
	ResourceIdMask           uint32
	MonitorBufferSize        uint32
	MaximumRequestLength     uint16
	ImageByteOrder           uint8
	BitmapFormatBitOrder     uint8
	BitmapFormatScanlineUint uint8
	BitmapFormatScanlinePad  uint8
	MinKeycode               Keycode
	MaxKeycode               Keycode
	Vendor                   string
	PixmapFormats            []Format
	Roots                    []Screen
}

func readSetup(r *Reader, v *Setup) error {
	if !r.RemainAtLeast4b(10) {
		return ErrDataLenShort
	}

	// status
	status := r.Read1b()
	if status != 1 {
		return errors.New("status != 1")
	}
	r.ReadPad(1)
	v.ProtocolMajorVersion = r.Read2b() // 1
	v.ProtocolMinorVersion = r.Read2b()
	// length in 4-byte units of additional data
	r.ReadPad(2) // 2

	v.ReleaseNumber = r.Read4b() // 3

	v.ResourceIdBase = r.Read4b()

	v.ResourceIdMask = r.Read4b()

	v.MonitorBufferSize = r.Read4b()

	vendorLen := r.Read2b()
	v.MaximumRequestLength = r.Read2b() // 7

	screensLen := r.Read1b()
	formatsLen := r.Read1b()
	v.ImageByteOrder = r.Read1b()
	v.BitmapFormatBitOrder = r.Read1b() // 8

	v.BitmapFormatScanlineUint = r.Read1b()
	v.BitmapFormatScanlinePad = r.Read1b()
	v.MinKeycode = Keycode(r.Read1b())
	v.MaxKeycode = Keycode(r.Read1b()) // 9

	// unused
	r.Read4b() // 10

	var err error
	v.Vendor, err = r.ReadStrWithPad(int(vendorLen))
	if err != nil {
		return err
	}

	// formats
	if formatsLen > 0 {
		if !r.RemainAtLeast4b(int(formatsLen) * 2) {
			return ErrDataLenShort
		}
		v.PixmapFormats = make([]Format, int(formatsLen))
		for i := 0; i < int(formatsLen); i++ {
			readFormat(r, &v.PixmapFormats[i])
		}
	}

	// screens
	if screensLen > 0 {
		v.Roots = make([]Screen, int(screensLen))
		for i := 0; i < int(screensLen); i++ {
			err := readScreen(r, &v.Roots[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// size: 2 * 4b
type Format struct {
	Depth        uint8
	BitsPerPixel uint8
	ScanlinePad  uint8
}

func readFormat(r *Reader, v *Format) {
	v.Depth = r.Read1b()
	v.BitsPerPixel = r.Read1b()
	v.ScanlinePad = r.Read1b()

	// unused
	r.ReadPad(5)
}

type Screen struct {
	Root                Window
	DefaultColorMap     Colormap
	WhitePixel          uint32
	BlackPixel          uint32
	CurrentInputMask    uint32
	WidthInPixels       uint16
	HeightInPixels      uint16
	WidthInMillimeters  uint16
	HeightInMillimeters uint16
	MinInstalledMaps    uint16
	MaxInstalledMaps    uint16
	RootVisual          VisualID
	BackingStores       uint8
	SaveUnders          bool
	RootDepth           uint8
	AllowedDepths       []Depth
}

func readScreen(r *Reader, v *Screen) error {
	if !r.RemainAtLeast4b(10) {
		return ErrDataLenShort
	}

	v.Root = Window(r.Read4b())

	v.DefaultColorMap = Colormap(r.Read4b())

	v.WhitePixel = r.Read4b()

	v.BlackPixel = r.Read4b()

	v.CurrentInputMask = r.Read4b() // 5

	v.WidthInPixels = r.Read2b()
	v.HeightInPixels = r.Read2b()

	v.WidthInMillimeters = r.Read2b()
	v.HeightInMillimeters = r.Read2b()

	v.MinInstalledMaps = r.Read2b()
	v.MaxInstalledMaps = r.Read2b() // 8

	v.RootVisual = VisualID(r.Read4b()) // 9

	v.BackingStores = r.Read1b()
	v.SaveUnders = r.ReadBool()
	v.RootDepth = r.Read1b()
	depthsLen := r.Read1b() // 10

	// depths
	if depthsLen > 0 {
		v.AllowedDepths = make([]Depth, int(depthsLen))
		for i := 0; i < int(depthsLen); i++ {
			err := readDepth(r, &v.AllowedDepths[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type Depth struct {
	Depth   uint8
	Visuals []VisualType
}

func readDepth(r *Reader, v *Depth) error {
	if !r.RemainAtLeast4b(2) {
		return ErrDataLenShort
	}
	v.Depth = r.Read1b()
	r.ReadPad(1)
	visualsLen := r.Read2b()

	// unused
	r.ReadPad(4) // 2

	// visuals
	if visualsLen > 0 {
		if !r.RemainAtLeast4b(int(visualsLen) * 6) {
			return ErrDataLenShort
		}
		v.Visuals = make([]VisualType, int(visualsLen))
		for i := 0; i < int(visualsLen); i++ {
			readVisualType(r, &v.Visuals[i])
		}
	}

	return nil
}

// size: 6 * 4b
type VisualType struct {
	Id              VisualID
	Class           uint8
	BitsPerRGBValue uint8
	ColorMapEntries uint16
	RedMask         uint32
	GreenMask       uint32
	BlueMask        uint32
}

func readVisualType(r *Reader, v *VisualType) {
	v.Id = VisualID(r.Read4b())

	v.Class = r.Read1b()
	v.BitsPerRGBValue = r.Read1b()
	v.ColorMapEntries = r.Read2b()

	v.RedMask = r.Read4b()

	v.GreenMask = r.Read4b()

	v.BlueMask = r.Read4b()

	// unused
	r.ReadPad(4)
}

// size: 2 * 4b
type Rectangle struct {
	X, Y          int16
	Width, Height uint16
}

func ReadRectangle(r *Reader) Rectangle {
	var v Rectangle
	v.X = int16(r.Read2b())
	v.Y = int16(r.Read2b())

	v.Width = r.Read2b()
	v.Height = r.Read2b() // 2
	return v
}

func WriteRectangle(b *FixedSizeBuf, v Rectangle) {
	b.Write2b(uint16(v.X))
	b.Write2b(uint16(v.Y))
	b.Write2b(v.Width)
	b.Write2b(v.Height)
}

// #WREQ
func encodeCreateWindow(depth uint8, wid, parent Window, x, y int16,
	width, height, borderWidth uint16, class uint16, visual VisualID,
	valueMask uint32, valueList []uint32) (hd uint8, b RequestBody) {

	hd = depth
	b0 := b.AddBlock(7 + len(valueList)).
		Write4b(uint32(wid)).
		Write4b(uint32(parent)).
		Write2b(uint16(x)).
		Write2b(uint16(y)).
		Write2b(width).
		Write2b(height).
		Write2b(borderWidth).
		Write2b(class).
		Write4b(uint32(visual)).
		Write4b(valueMask)
	for _, value := range valueList {
		b0.Write4b(value)
	}
	b0.End()
	return
}

// #WREQ
func encodeChangeWindowAttributes(window Window, valueMask uint32,
	valueList []uint32) (hd uint8, b RequestBody) {

	b0 := b.AddBlock(2 + len(valueList)).
		Write4b(uint32(window)).
		Write4b(valueMask)
	for _, value := range valueList {
		b0.Write4b(value)
	}
	b0.End()
	return
}

// #WREQ
func encodeGetWindowAttributes(window Window) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(window)).
		End()
	return
}

type GetWindowAttributesReply struct {
	BackingStore       uint8
	Visual             VisualID
	Class              uint16
	BitGravity         uint8
	WinGravity         uint8
	BackingPlanes      uint32
	BackingPixel       uint32
	SaveUnder          bool
	MapIsInstalled     bool
	MapState           uint8
	OverrideRedirect   bool
	Colormap           Colormap
	AllEventMasks      uint32
	YourEventMask      uint32
	DoNotPropagateMask uint16
}

func readGetWindowAttributesReply(r *Reader, v *GetWindowAttributesReply) error {
	if !r.RemainAtLeast4b(11) {
		return ErrDataLenShort
	}

	v.BackingStore, _ = r.ReadReplyHeader()

	v.Visual = VisualID(r.Read4b())

	v.Class = r.Read2b()
	v.BitGravity = r.Read1b()
	v.WinGravity = r.Read1b() // 4

	v.BackingPlanes = r.Read4b()

	v.BackingPixel = r.Read4b()

	v.SaveUnder = r.ReadBool()
	v.MapIsInstalled = r.ReadBool()
	v.MapState = r.Read1b()
	v.OverrideRedirect = r.ReadBool() // 7

	v.Colormap = Colormap(r.Read4b())

	v.AllEventMasks = r.Read4b()

	v.YourEventMask = r.Read4b()

	v.DoNotPropagateMask = r.Read2b() // 11

	return nil
}

// #WREQ
func encodeDestroyWindow(window Window) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(window)).
		End()
	return
}

// #WREQ
func encodeDestroySubwindows(window Window) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(window)).
		End()
	return
}

// #WREQ
func encodeChangeSaveSet(mode uint8, window Window) (hd uint8, b RequestBody) {
	hd = mode
	b.AddBlock(1).
		Write4b(uint32(window)).
		End()
	return
}

// #WREQ
func encodeReparentWindow(window, parent Window, x, y int16) (hd uint8, b RequestBody) {
	b.AddBlock(3).
		Write4b(uint32(window)).
		Write4b(uint32(parent)).
		Write2b(uint16(x)).
		Write2b(uint16(y)).
		End()
	return
}

// #WREQ
func encodeMapWindow(window Window) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(window)).
		End()
	return
}

// #WREQ
func encodeMapSubwindows(window Window) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(window)).
		End()
	return
}

// #WREQ
func encodeUnmapWindow(window Window) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(window)).
		End()
	return
}

// #WREQ
func encodeUnmapSubwindows(window Window) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(window)).
		End()
	return
}

// #WREQ
func encodeConfigureWindow(window Window, valueMask uint16, valueList []uint32) (hd uint8, b RequestBody) {
	b0 := b.AddBlock(2 + len(valueList)).
		Write4b(uint32(window)).
		Write2b(valueMask).
		WritePad(2)
	for _, value := range valueList {
		b0.Write4b(value)
	}
	b0.End()
	return
}

// #WREQ
func encodeCirculateWindow(direction uint8, window Window) (hd uint8, b RequestBody) {
	hd = direction
	b.AddBlock(1).
		Write4b(uint32(window)).
		End()
	return
}

// #WREQ
func encodeGetGeometry(drawable Drawable) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(drawable)).
		End()
	return
}

type GetGeometryReply struct {
	Depth       uint8
	Root        Window
	X           int16
	Y           int16
	Width       uint16
	Height      uint16
	BorderWidth uint16
}

func readGetGeometryReply(r *Reader, v *GetGeometryReply) error {
	if !r.RemainAtLeast4b(6) {
		return ErrDataLenShort
	}
	v.Depth, _ = r.ReadReplyHeader() // 2

	v.Root = Window(r.Read4b())

	v.X = int16(r.Read2b())
	v.Y = int16(r.Read2b()) // 4

	v.Width = r.Read2b()
	v.Height = r.Read2b() // 5

	v.BorderWidth = r.Read2b() // 6

	return nil
}

// #WREQ
func encodeQueryTree(window Window) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(window)).
		End()
	return
}

type QueryTreeReply struct {
	Root, Parent Window
	Children     []Window
}

func readQueryTreeReply(r *Reader, v *QueryTreeReply) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}
	r.ReadPad(8) // 2

	v.Root = Window(r.Read4b())

	v.Parent = Window(r.Read4b()) // 4

	childrenLen := int(r.Read2b())

	// unused
	r.ReadPad(14) // 8

	if childrenLen > 0 {
		if !r.RemainAtLeast4b(childrenLen) {
			return ErrDataLenShort
		}
		v.Children = make([]Window, childrenLen)
		for i := 0; i < childrenLen; i++ {
			v.Children[i] = Window(r.Read4b())
		}
	}

	return nil
}

// #WREQ
func encodeInternAtom(onlyIfExists bool, name string) (hd uint8, b RequestBody) {
	hd = BoolToUint8(onlyIfExists)
	name = TruncateStr(name, math.MaxUint16)
	nameLen := len(name)
	b.AddBlock(1 + SizeIn4bWithPad(nameLen)).
		Write2b(uint16(nameLen)).
		WritePad(2).
		WriteString(name).
		WritePad(Pad(nameLen)).
		End()
	return
}

type InternAtomReply struct {
	Atom Atom
}

func readInternAtomReply(r *Reader, v *InternAtomReply) error {
	if !r.RemainAtLeast4b(3) {
		return ErrDataLenShort
	}
	r.ReadPad(8) // 2

	v.Atom = Atom(r.Read4b()) // 3
	return nil
}

// #WREQ
func encodeGetAtomName(atom Atom) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(atom)).
		End()
	return
}

type GetAtomNameReply struct {
	Name string
}

func readGetAtomNameReply(r *Reader, v *GetAtomNameReply) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}
	r.ReadPad(8) // 2

	// name len
	nameLen := r.Read2b()

	// unused
	r.ReadPad(22) // 8

	var err error
	v.Name, err = r.ReadString(int(nameLen))
	if err != nil {
		return err
	}

	return nil
}

// #WREQ
func encodeChangeProperty(mode uint8, window Window, property, type0 Atom,
	format uint8, data []byte) (hd uint8, b RequestBody) {
	hd = mode

	dataLen := len(data)
	b.AddBlock(5).
		Write4b(uint32(window)).
		Write4b(uint32(property)).
		Write4b(uint32(type0)).
		Write1b(format).
		WritePad(3).
		Write4b(uint32(dataLen / (int(format) / 8))).
		End()
	b.AddBytes(data)
	return
}

// #WREQ
func encodeDeleteProperty(window Window, property Atom) (hd uint8, b RequestBody) {
	b.AddBlock(2).
		Write4b(uint32(window)).
		Write4b(uint32(property)).
		End()
	return
}

// #WREQ
func encodeGetProperty(delete bool, window Window, property, type0 Atom,
	longOffset, longLength uint32) (hd uint8, b RequestBody) {

	hd = BoolToUint8(delete)
	b.AddBlock(5).
		Write4b(uint32(window)).
		Write4b(uint32(property)).
		Write4b(uint32(type0)).
		Write4b(longOffset).
		Write4b(longLength).
		End()
	return
}

type GetPropertyReply struct {
	Format     uint8
	Type       Atom
	BytesAfter uint32
	ValueLen   uint32
	Value      []byte
}

func readGetPropertyReply(r *Reader, v *GetPropertyReply) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}

	v.Format, _ = r.ReadReplyHeader() // 2

	v.Type = Atom(r.Read4b())

	v.BytesAfter = r.Read4b()

	v.ValueLen = r.Read4b() // 5

	// unused
	r.ReadPad(12) // 8

	n := int(v.ValueLen) * int(v.Format/8)
	var err error
	v.Value, err = r.ReadBytes(n)
	if err != nil {
		return err
	}

	return nil
}

// #WREQ
func encodeListProperties(window Window) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(window)).
		End()
	return
}

type ListPropertiesReply struct {
	Atoms []Atom
}

func readListPropertiesReply(r *Reader, v *ListPropertiesReply) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}
	r.ReadPad(8) // 2

	atomsLen := int(r.Read2b())

	// unused
	r.ReadPad(22) // 8

	if atomsLen > 0 {
		if !r.RemainAtLeast4b(atomsLen) {
			return ErrDataLenShort
		}
		v.Atoms = make([]Atom, atomsLen)
		for i := 0; i < atomsLen; i++ {
			v.Atoms[i] = Atom(r.Read4b())
		}
	}

	return nil
}

// #WREQ
func encodeSetSelectionOwner(owner Window, selection Atom, time Timestamp) (hd uint8, b RequestBody) {
	b.AddBlock(3).
		Write4b(uint32(owner)).
		Write4b(uint32(selection)).
		Write4b(uint32(time)).
		End()
	return
}

// #WREQ
func encodeGetSelectionOwner(selection Atom) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(selection)).
		End()
	return
}

type GetSelectionOwnerReply struct {
	Owner Window
}

func readGetSelectionOwnerReply(r *Reader, v *GetSelectionOwnerReply) error {
	if !r.RemainAtLeast4b(3) {
		return ErrDataLenShort
	}
	r.ReadPad(8) // 2

	v.Owner = Window(r.Read4b()) // 3
	return nil
}

// #WREQ
func encodeConvertSelection(requestor Window, selection, target, property Atom,
	time Timestamp) (hd uint8, b RequestBody) {
	b.AddBlock(5).
		Write4b(uint32(requestor)).
		Write4b(uint32(selection)).
		Write4b(uint32(target)).
		Write4b(uint32(property)).
		Write4b(uint32(time)).
		End()
	return
}

// #WREQ
func encodeSendEvent(propagate bool, destination Window, eventMask uint32,
	event []byte) (hd uint8, b RequestBody) {
	hd = BoolToUint8(propagate)
	b.AddBlock(2).
		Write4b(uint32(destination)).
		Write4b(eventMask)
	b.AddBytes(event)
	return
}

// #WREQ
func encodeGrabPointer(ownerEvents bool, grabWindow Window,
	eventMask uint16, pointerMode, keyboardMode uint8, confineTo Window,
	cursor Cursor, time Timestamp) (hd uint8, b RequestBody) {
	hd = BoolToUint8(ownerEvents)
	b.AddBlock(5).
		Write4b(uint32(grabWindow)).
		Write2b(eventMask).
		Write1b(pointerMode).
		Write1b(keyboardMode).
		Write4b(uint32(confineTo)).
		Write4b(uint32(cursor)).
		Write4b(uint32(time)).
		End()
	return
}

type GrabPointerReply struct {
	Status uint8
}

func readGrabPointerReply(r *Reader, v *GrabPointerReply) error {
	if !r.RemainAtLeast4b(2) {
		return ErrDataLenShort
	}

	v.Status, _ = r.ReadReplyHeader() // 2
	return nil
}

// #WREQ
func encodeUngrabPointer(time Timestamp) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(time)).
		End()
	return
}

// #WREQ
func encodeGrabButton(ownerEvents bool, grabWindow Window, eventMask uint16,
	pointerMode, keyboardMode uint8, confineTo Window, cursor Cursor, button uint8,
	modifiers uint16) (hd uint8, b RequestBody) {
	hd = BoolToUint8(ownerEvents)
	b.AddBlock(5).
		Write4b(uint32(grabWindow)).
		Write2b(eventMask).
		Write1b(pointerMode).
		Write1b(keyboardMode).
		Write4b(uint32(confineTo)).
		Write4b(uint32(cursor)).
		Write1b(button).
		WritePad(1).
		Write2b(modifiers).
		End()
	return
}

// #WREQ
func encodeUngrabButton(button uint8, grabWindow Window,
	modifiers uint16) (hd uint8, b RequestBody) {
	hd = button
	b.AddBlock(2).
		Write4b(uint32(grabWindow)).
		Write2b(modifiers).
		WritePad(2).
		End()
	return
}

// #WREQ
func encodeChangeActivePointerGrab(cursor Cursor, time Timestamp,
	eventMask uint16) (hd uint8, b RequestBody) {
	b.AddBlock(3).
		Write4b(uint32(cursor)).
		Write4b(uint32(time)).
		Write2b(eventMask).
		WritePad(2).
		End()
	return
}

// #WREQ
func encodeGrabKeyboard(ownerEvents bool, grabWindow Window, time Timestamp,
	pointerMode, keyboardMode uint8) (hd uint8, b RequestBody) {
	hd = BoolToUint8(ownerEvents)
	b.AddBlock(3).
		Write4b(uint32(grabWindow)).
		Write4b(uint32(time)).
		Write1b(pointerMode).
		Write1b(keyboardMode).
		WritePad(2).
		End()
	return
}

type GrabKeyboardReply struct {
	Status uint8
}

func readGrabKeyboardReply(r *Reader, v *GrabKeyboardReply) error {
	if !r.RemainAtLeast4b(2) {
		return ErrDataLenShort
	}
	v.Status, _ = r.ReadReplyHeader() // 2
	return nil
}

// #WREQ
func encodeUngrabKeyboard(time Timestamp) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(time)).
		End()
	return
}

// #WREQ
func encodeGrabKey(ownerEvents bool, grabWindow Window, modifiers uint16,
	key Keycode, pointerMode, keyboardMode uint8) (hd uint8, b RequestBody) {
	hd = BoolToUint8(ownerEvents)
	b.AddBlock(3).
		Write4b(uint32(grabWindow)).
		Write2b(modifiers).
		Write1b(uint8(key)).
		Write1b(pointerMode).
		Write1b(keyboardMode).
		WritePad(3).
		End()
	return
}

// #WREQ
func encodeUngrabKey(key Keycode, grabWindow Window, modifiers uint16) (hd uint8, b RequestBody) {
	hd = uint8(key)
	b.AddBlock(2).
		Write4b(uint32(grabWindow)).
		Write2b(modifiers).
		WritePad(2).
		End()
	return
}

// #WREQ
func encodeAllowEvents(mode uint8, time Timestamp) (hd uint8, b RequestBody) {
	hd = mode
	b.AddBlock(1).
		Write4b(uint32(time)).
		End()
	return
}

// #WREQ
func encodeGrabServer() (hd uint8, b RequestBody) {
	return
}

// #WREQ
func encodeUngrabServer() (hd uint8, b RequestBody) {
	return
}

// #WREQ
func encodeQueryPointer(window Window) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(window)).
		End()
	return
}

type QueryPointerReply struct {
	SameScreen bool
	Root       Window
	Child      Window
	RootX      int16
	RootY      int16
	WinX       int16
	WinY       int16
	Mask       uint16
}

func readQueryPointerReply(r *Reader, v *QueryPointerReply) error {
	if !r.RemainAtLeast4b(7) {
		return ErrDataLenShort
	}

	sameScreen, _ := r.ReadReplyHeader() // 2
	v.SameScreen = Uint8ToBool(sameScreen)

	v.Root = Window(r.Read4b())

	v.Child = Window(r.Read4b())

	v.RootX = int16(r.Read2b())
	v.RootY = int16(r.Read2b()) // 5

	v.WinX = int16(r.Read2b())
	v.WinY = int16(r.Read2b()) // 6

	v.Mask = r.Read2b() // 7
	return nil
}

// #WREQ
func encodeGetMotionEvents(window Window, start, stop Timestamp) (hd uint8, b RequestBody) {
	b.AddBlock(3).
		Write4b(uint32(window)).
		Write4b(uint32(start)).
		Write4b(uint32(stop)).
		End()
	return
}

type GetMotionEventsReply struct {
	Events []TimeCoord
}

// size: 2 * 4b
type TimeCoord struct {
	Time Timestamp
	X, Y int16
}

func readTimeCoord(r *Reader, v *TimeCoord) {
	v.Time = Timestamp(r.Read4b())

	v.X = int16(r.Read2b())
	v.Y = int16(r.Read2b())
}

func readGetMotionEventsReply(r *Reader, v *GetMotionEventsReply) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}

	r.ReadPad(8) // 2

	eventsLen := int(r.Read4b()) // 3

	// unused
	r.ReadPad(20) // 8

	if eventsLen > 0 {
		if !r.RemainAtLeast4b(eventsLen * 2) {
			return ErrDataLenShort
		}
		v.Events = make([]TimeCoord, eventsLen)
		for i := 0; i < eventsLen; i++ {
			readTimeCoord(r, &v.Events[i])
		}
	}

	return nil
}

// #WREQ
func encodeTranslateCoordinates(srcWindow, dstWindow Window,
	srcX, srcY int16) (hd uint8, b RequestBody) {
	b.AddBlock(3).
		Write4b(uint32(srcWindow)).
		Write4b(uint32(dstWindow)).
		Write2b(uint16(srcX)).
		Write2b(uint16(srcY)).
		End()
	return
}

type TranslateCoordinatesReply struct {
	SameScreen bool
	Child      Window
	DstX       int16
	DstY       int16
}

func readTranslateCoordinatesReply(r *Reader, v *TranslateCoordinatesReply) error {
	if !r.RemainAtLeast4b(4) {
		return ErrDataLenShort
	}
	sameScreen, _ := r.ReadReplyHeader() // 2
	v.SameScreen = Uint8ToBool(sameScreen)

	v.Child = Window(r.Read4b())

	v.DstX = int16(r.Read2b())
	v.DstY = int16(r.Read2b()) // 4
	return nil
}

// #WREQ
func encodeWarpPointer(srcWindow, dstWindow Window, srcX, srcY int16,
	srcWidth, srcHeight uint16, dstX, dstY int16) (hd uint8, b RequestBody) {
	b.AddBlock(5).
		Write4b(uint32(srcWindow)).
		Write4b(uint32(dstWindow)).
		Write2b(uint16(srcX)).
		Write2b(uint16(srcY)).
		Write2b(srcWidth).
		Write2b(srcHeight).
		Write2b(uint16(dstX)).
		Write2b(uint16(dstY)).
		End()
	return
}

// #WREQ
func encodeSetInputFocus(revertTo uint8, focus Window, time Timestamp) (hd uint8, b RequestBody) {
	hd = revertTo
	b.AddBlock(2).
		Write4b(uint32(focus)).
		Write4b(uint32(time)).
		End()
	return
}

// #WREQ
func encodeGetInputFocus() (hd uint8, b RequestBody) {
	return
}

type GetInputFocusReply struct {
	RevertTo uint8
	Focus    Window
}

func readGetInputFocusReply(r *Reader, v *GetInputFocusReply) error {
	if !r.RemainAtLeast4b(3) {
		return ErrDataLenShort
	}
	v.RevertTo, _ = r.ReadReplyHeader()

	v.Focus = Window(r.Read4b()) // 3
	return nil
}

// #WREQ
func encodeQueryKeymap() (hd uint8, b RequestBody) {
	return
}

type QueryKeymapReply struct {
	Keys []byte
}

func readQueryKeymapReply(r *Reader, v *QueryKeymapReply) error {
	if !r.RemainAtLeast4b(10) {
		return ErrDataLenShort
	}

	r.ReadPad(8) // 2

	// keys
	v.Keys = r.MustReadBytes(32) // 10

	return nil
}

// #WREQ
func encodeOpenFont(fid Font, name string) (hd uint8, b RequestBody) {
	name = TruncateStr(name, math.MaxUint16)
	nameLen := len(name)
	b.AddBlock(2 + SizeIn4bWithPad(nameLen)).
		Write4b(uint32(fid)).
		Write2b(uint16(nameLen)).
		WritePad(2).
		WriteString(name).
		WritePad(Pad(nameLen)).
		End()
	return
}

// #WREQ
func encodeCloseFont(font Font) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(font)).
		End()
	return
}

// #WREQ
func encodeQueryFont(font Fontable) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(font)).
		End()
	return
}

type QueryFontReply struct {
	MinBounds      CharInfo
	MaxBounds      CharInfo
	MinCharOrByte2 uint16
	MaxCharOrByte2 uint16
	DefaultChar    uint16
	PropertiesLen  uint16
	DrawDirection  uint8
	MinByte1       uint8
	MaxByte1       uint8
	AllCharsExist  bool
	FontAscent     int16
	FontDescent    int16
	CharInfosLen   uint32
	Properties     []FontProp
	CharInfos      []CharInfo
}

func readQueryFontReply(r *Reader, v *QueryFontReply) error {
	if !r.RemainAtLeast4b(15) {
		return ErrDataLenShort
	}
	r.ReadPad(8) // 2

	readCharInfo(r, &v.MinBounds) // 5

	// unused
	r.ReadPad(4) // 6

	readCharInfo(r, &v.MaxBounds) // 9

	// unused
	r.ReadPad(4) // 10

	v.MinCharOrByte2 = r.Read2b()
	v.MaxCharOrByte2 = r.Read2b()

	v.DefaultChar = r.Read2b()
	propsLen := int(r.Read2b()) // 12

	v.DrawDirection = r.Read1b()
	v.MinByte1 = r.Read1b()
	v.MaxByte1 = r.Read1b()
	v.AllCharsExist = r.ReadBool()

	v.FontAscent = int16(r.Read2b())
	v.FontDescent = int16(r.Read2b())

	charInfosLen := int(r.Read4b()) // 15

	if propsLen > 0 {
		if !r.RemainAtLeast4b(2 * propsLen) {
			return ErrDataLenShort
		}
		v.Properties = make([]FontProp, propsLen)
		for i := 0; i < propsLen; i++ {
			readFontProp(r, &v.Properties[i])
		}
	}

	if charInfosLen > 0 {
		if !r.RemainAtLeast4b(3 * charInfosLen) {
			return ErrDataLenShort
		}
		v.CharInfos = make([]CharInfo, charInfosLen)
		for i := 0; i < charInfosLen; i++ {
			readCharInfo(r, &v.CharInfos[i])
		}
	}

	return nil
}

// size: 3 * 4b
type CharInfo struct {
	LeftSideBearing  int16
	RightSideBearing int16
	CharacterWidth   int16
	Ascent           int16
	Descent          int16
	Attributes       uint16
}

func readCharInfo(r *Reader, v *CharInfo) {
	v.LeftSideBearing = int16(r.Read2b())
	v.RightSideBearing = int16(r.Read2b())

	v.CharacterWidth = int16(r.Read2b())
	v.Ascent = int16(r.Read2b())

	v.Descent = int16(r.Read2b())
	v.Attributes = r.Read2b()
}

// size: 2 * 4b
type FontProp struct {
	Name  Atom
	Value uint32
}

func readFontProp(r *Reader, v *FontProp) {
	v.Name = Atom(r.Read4b())

	v.Value = r.Read4b()
}

// TODO: QueryTextExtents

// #WREQ
func encodeListFonts(maxNames uint16, pattern string) (hd uint8, b RequestBody) {
	pattern = TruncateStr(pattern, math.MaxUint16)
	patternLen := len(pattern)
	b.AddBlock(1 + SizeIn4bWithPad(patternLen)).
		Write2b(maxNames).
		Write2b(uint16(patternLen)).
		WriteString(pattern).
		WritePad(Pad(patternLen)).
		End()
	return
}

type ListFontsReply struct {
	Names []string
}

func readListFontsReply(r *Reader, v *ListFontsReply) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}
	r.ReadPad(8) // 2

	namesLen := int(r.Read2b())

	// unused
	r.ReadPad(22) // 8

	// names
	if namesLen > 0 {
		v.Names = make([]string, namesLen)
		for i := 0; i < namesLen; i++ {
			var str Str
			_, err := readStr(r, &str)
			if err != nil {
				return err
			}
			v.Names[i] = str.Value
		}
	}

	return nil
}

// #WREQ
func encodeListFontsWithInfo(maxNames uint16, pattern string) (hd uint8, b RequestBody) {
	pattern = TruncateStr(pattern, math.MaxUint16)
	patternLen := len(pattern)
	b.AddBlock(1 + SizeIn4bWithPad(patternLen)).
		Write2b(maxNames).
		Write2b(uint16(patternLen)).
		WriteString(pattern).
		WritePad(Pad(patternLen)).
		End()
	return
}

type ListFontsWithInfoReply struct {
	LastReply      bool
	MinBounds      CharInfo
	MaxBounds      CharInfo
	MinCharOrByte2 uint16
	MaxCharOrByte2 uint16
	DefaultChar    uint16
	PropertiesLen  uint16
	DrawDirection  uint8
	MinByte1       uint8
	MaxByte1       uint8
	AllCharsExist  bool
	FontAscent     int16
	FontDescent    int16
	RepliesHint    uint32
	Properties     []FontProp
	Name           string
}

func readListFontsWithInfoReply(r *Reader, v *ListFontsWithInfoReply) error {
	if !r.RemainAtLeast4b(15) {
		return ErrDataLenShort
	}
	lastReplyIndicator, replyLen := r.ReadReplyHeader() // 2

	if lastReplyIndicator == 0 {
		v.LastReply = true
		return nil
	} else {
		v.LastReply = false
	}

	readCharInfo(r, &v.MinBounds) // 5

	// unused
	r.ReadPad(4) // 6

	readCharInfo(r, &v.MaxBounds) // 9

	// unused
	r.ReadPad(4)

	v.MinCharOrByte2 = r.Read2b()
	v.MaxCharOrByte2 = r.Read2b() // 11

	v.DefaultChar = r.Read2b()
	propsLen := int(r.Read2b())

	v.DrawDirection = r.Read1b()
	v.MinByte1 = r.Read1b()
	v.MaxByte1 = r.Read1b()
	v.AllCharsExist = r.ReadBool() // 13

	v.FontAscent = int16(r.Read2b())
	v.FontDescent = int16(r.Read2b())

	v.RepliesHint = r.Read4b() // 15

	if propsLen > 0 {
		if !r.RemainAtLeast4b(2 * propsLen) {
			return ErrDataLenShort
		}
		v.Properties = make([]FontProp, propsLen)
		for i := 0; i < propsLen; i++ {
			readFontProp(r, &v.Properties[i])
		}
	}

	// TODO: use r.ReadNulTermStr
	nameLen := (int(replyLen) - 7 - (2 * propsLen)) * 4
	nameBytes, err := r.ReadBytes(nameLen)
	if err != nil {
		return err
	}

	zeroIdx := bytes.IndexByte(nameBytes, 0)
	if zeroIdx != -1 {
		nameBytes = nameBytes[:zeroIdx]
	}
	v.Name = string(nameBytes)
	return nil
}

// #WREQ
func encodeSetFontPath(paths []string) (hd uint8, b RequestBody) {
	pathsLen := uint16(len(paths))

	var n int
	for _, p := range paths {
		pLen := len(p)
		if pLen > math.MaxUint8 {
			pLen = math.MaxUint8
		}
		n += 1 + pLen
	}

	b0 := b.AddBlock(1 + SizeIn4bWithPad(n)).
		Write2b(pathsLen).
		WritePad(2)

	for _, p := range paths {
		writeStr(b0, p)
	}
	b0.WritePad(Pad(n)).End()
	return
}

// #WREQ
func encodeGetFontPath() (hd uint8, b RequestBody) {
	return
}

type GetFontPathReply struct {
	Paths []string
}

func readGetFontPathReply(r *Reader, v *GetFontPathReply) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}
	r.ReadPad(8) // 2

	pathsLen := int(r.Read2b())

	// unused
	r.ReadPad(22) // 8

	if pathsLen > 0 {
		v.Paths = make([]string, pathsLen)
		for i := 0; i < pathsLen; i++ {
			var str Str
			_, err := readStr(r, &str)
			if err != nil {
				return err
			}
			v.Paths[i] = str.Value
		}
	}

	return nil
}

// #WREQ
func encodeCreatePixmap(depth uint8, pid Pixmap, drawable Drawable,
	width, height uint16) (hd uint8, b RequestBody) {
	hd = depth
	b.AddBlock(3).
		Write4b(uint32(pid)).
		Write4b(uint32(drawable)).
		Write2b(width).
		Write2b(height).
		End()
	return
}

// #WREQ
func encodeFreePixmap(pixmap Pixmap) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(pixmap)).
		End()
	return
}

// #WREQ
func encodeCreateGC(cid GContext, drawable Drawable, valueMask uint32,
	valueList []uint32) (hd uint8, b RequestBody) {

	b0 := b.AddBlock(3 + len(valueList)).
		Write4b(uint32(cid)).
		Write4b(uint32(drawable)).
		Write4b(valueMask)

	for _, value := range valueList {
		b0.Write4b(value)
	}
	b0.End()
	return
}

// #WREQ
func encodeChangeGC(gc GContext, valueMask uint32, valueList []uint32) (hd uint8, b RequestBody) {
	b0 := b.AddBlock(2 + len(valueList)).
		Write4b(uint32(gc)).
		Write4b(valueMask)
	for _, value := range valueList {
		b0.Write4b(value)
	}
	b0.End()
	return
}

// #WREQ
func encodeCopyGC(srcGC, dstGC GContext, valueMask uint32) (hd uint8, b RequestBody) {
	b.AddBlock(3).
		Write4b(uint32(srcGC)).
		Write4b(uint32(dstGC)).
		Write4b(valueMask).
		End()
	return
}

// #WREQ
func encodeSetDashes(gc GContext, dashOffset uint16, dashes []uint8) (hd uint8, b RequestBody) {
	dashesLen := uint16(len(dashes))
	b.AddBlock(2).
		Write4b(uint32(gc)).
		Write2b(dashOffset).
		Write2b(dashesLen).
		End()
	b.AddBytes(dashes)
	return
}

// #WREQ
func encodeSetClipRectangles(ordering uint8, gc GContext,
	clipXOrigin, clipYOrigin int16, rectangles []Rectangle) (hd uint8, b RequestBody) {
	hd = ordering
	b0 := b.AddBlock(2 + 2*len(rectangles)).
		Write4b(uint32(gc)).
		Write2b(uint16(clipXOrigin)).
		Write2b(uint16(clipYOrigin))

	for _, rect := range rectangles {
		WriteRectangle(b0, rect)
	}
	b0.End()
	return
}

// #WREQ
func encodeFreeGC(gc GContext) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(gc)).
		End()
	return
}

// #WREQ
func encodeClearArea(exposures bool, window Window, x, y int16,
	width, height uint16) (hd uint8, b RequestBody) {
	hd = BoolToUint8(exposures)
	b.AddBlock(3).
		Write4b(uint32(window)).
		Write2b(uint16(x)).
		Write2b(uint16(y)).
		Write2b(width).
		Write2b(height).
		End()
	return
}

// #WREQ
func encodeCopyArea(srcDrawable, dstDrawable Drawable, gc GContext,
	srcX, srcY, dstX, dstY int16, width, height uint16) (hd uint8, b RequestBody) {
	b.AddBlock(6).
		Write4b(uint32(srcDrawable)).
		Write4b(uint32(dstDrawable)).
		Write4b(uint32(gc)).
		Write2b(uint16(srcX)).
		Write2b(uint16(srcY)).
		Write2b(uint16(dstX)).
		Write2b(uint16(dstY)).
		Write2b(width).
		Write2b(height).
		End()
	return
}

// #WREQ
func encodeCopyPlane(srcDrawable, dstDrawable Drawable, gc GContext,
	srcX, srcY, dstX, dstY int16, width, height uint16, bitPlane uint32) (hd uint8, b RequestBody) {

	b.AddBlock(7).
		Write4b(uint32(srcDrawable)).
		Write4b(uint32(dstDrawable)).
		Write4b(uint32(gc)).
		Write2b(uint16(srcX)).
		Write2b(uint16(srcY)).
		Write2b(uint16(dstX)).
		Write2b(uint16(dstY)).
		Write2b(width).
		Write2b(height).
		Write4b(bitPlane).
		End()
	return
}

type Point struct {
	X, Y int16
}

func writePoint(b *FixedSizeBuf, v Point) {
	b.Write2b(uint16(v.X))
	b.Write2b(uint16(v.Y))
}

// #WREQ
func encodePolyPoint(coordinateMode uint8, drawable Drawable, gc GContext,
	points []Point) (hd uint8, b RequestBody) {

	hd = coordinateMode
	b0 := b.AddBlock(2 + len(points)).
		Write4b(uint32(drawable)).
		Write4b(uint32(gc))
	for _, p := range points {
		writePoint(b0, p)
	}
	b0.End()
	return
}

// #WREQ
func encodePolyLine(coordinateMode uint8, drawable Drawable, gc GContext,
	points []Point) (hd uint8, b RequestBody) {
	hd = coordinateMode
	b0 := b.AddBlock(2 + len(points)).
		Write4b(uint32(drawable)).
		Write4b(uint32(gc))

	for _, p := range points {
		writePoint(b0, p)
	}
	b0.End()
	return
}

// #WREQ
func encodePolySegment(drawable Drawable, gc GContext,
	segments []Segment) (hd uint8, b RequestBody) {

	b0 := b.AddBlock(2 + 2*len(segments)).
		Write4b(uint32(drawable)).
		Write4b(uint32(gc))
	for _, seg := range segments {
		writeSegment(b0, seg)
	}
	b0.End()
	return
}

type Segment struct {
	X1, Y1, X2, Y2 int16
}

func writeSegment(w *FixedSizeBuf, v Segment) {
	w.Write2b(uint16(v.X1))
	w.Write2b(uint16(v.Y1))
	w.Write2b(uint16(v.X2))
	w.Write2b(uint16(v.Y2))
}

// #WREQ
func encodePolyRectangle(drawable Drawable, gc GContext,
	rectangles []Rectangle) (hd uint8, b RequestBody) {

	b0 := b.AddBlock(2 + 2*len(rectangles)).
		Write4b(uint32(drawable)).
		Write4b(uint32(gc))
	for _, rect := range rectangles {
		WriteRectangle(b0, rect)
	}
	b0.End()
	return
}

// #WREQ
func encodePolyArc(drawable Drawable, gc GContext, arcs []Arc) (hd uint8, b RequestBody) {
	b0 := b.AddBlock(2 + 3*len(arcs)).
		Write4b(uint32(drawable)).
		Write4b(uint32(gc))
	for _, arc := range arcs {
		writeArc(b0, arc)
	}
	b0.End()
	return
}

type Arc struct {
	X, Y           int16
	Width, Height  uint16
	Angle1, Angle2 int16
}

func writeArc(w *FixedSizeBuf, v Arc) {
	w.Write2b(uint16(v.X))
	w.Write2b(uint16(v.Y))
	w.Write2b(v.Width)
	w.Write2b(v.Height)
	w.Write2b(uint16(v.Angle1))
	w.Write2b(uint16(v.Angle2))
}

// #WREQ
func encodeFillPoly(drawable Drawable, gc GContext, shape uint8,
	coordinateMode uint8, points []Point) (hd uint8, b RequestBody) {
	b0 := b.AddBlock(3 + len(points)).
		Write4b(uint32(drawable)).
		Write4b(uint32(gc)).
		Write1b(shape).
		Write1b(coordinateMode).
		WritePad(2)
	for _, p := range points {
		writePoint(b0, p)
	}
	b0.End()
	return
}

// #WREQ
func encodePolyFillRectangle(drawable Drawable, gc GContext,
	rectangles []Rectangle) (hd uint8, b RequestBody) {

	b0 := b.AddBlock(2 + 2*len(rectangles)).
		Write4b(uint32(drawable)).
		Write4b(uint32(gc))
	for _, rect := range rectangles {
		WriteRectangle(b0, rect)
	}
	b0.End()
	return
}

// #WREQ
func encodePolyFillArc(drawable Drawable, gc GContext, arcs []Arc) (hd uint8, b RequestBody) {
	b0 := b.AddBlock(2 + 3*len(arcs)).
		Write4b(uint32(drawable)).
		Write4b(uint32(gc))

	for _, arc := range arcs {
		writeArc(b0, arc)
	}
	b0.End()
	return
}

// #WREQ
func encodePutImage(format uint8, drawable Drawable, gc GContext, width, height uint16,
	dstX, dstY int16, leftPad, depth uint8, data []byte) (hd uint8, b RequestBody) {
	hd = format
	b.AddBlock(5).
		Write4b(uint32(drawable)).
		Write4b(uint32(gc)).
		Write2b(width).
		Write2b(height).
		Write2b(uint16(dstX)).
		Write2b(uint16(dstY)).
		Write1b(leftPad).
		Write1b(depth).
		WritePad(2).
		End()
	b.AddBytes(data)
	return
}

// #WREQ
func encodeGetImage(format uint8, drawable Drawable, x, y int16,
	width, height uint16, planeMask uint32) (hd uint8, b RequestBody) {
	hd = format
	b.AddBlock(4).
		Write4b(uint32(drawable)).
		Write2b(uint16(x)).
		Write2b(uint16(y)).
		Write2b(width).
		Write2b(height).
		Write4b(planeMask).
		End()
	return
}

type GetImageReply struct {
	Depth  uint8
	Visual VisualID
	Data   []byte
}

func readGetImageReply(r *Reader, v *GetImageReply) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}
	var replyLen uint32
	v.Depth, replyLen = r.ReadReplyHeader()

	v.Visual = VisualID(r.Read4b()) // 3

	// unused
	r.ReadPad(20) // 8

	dataLen := int(replyLen) * 4
	var err error
	v.Data, err = r.ReadBytes(dataLen)
	if err != nil {
		return err
	}

	return nil
}

// #WREQ
func encodeQueryExtension(name string) (hd uint8, b RequestBody) {
	name = TruncateStr(name, math.MaxUint16)
	nameLen := len(name)
	b.AddBlock(1 + SizeIn4bWithPad(nameLen)).
		Write2b(uint16(nameLen)).
		WritePad(2).
		WriteString(name).
		WritePad(Pad(nameLen)).
		End()
	return
}

type QueryExtensionReply struct {
	Present     bool
	MajorOpcode uint8
	FirstEvent  uint8
	FirstError  uint8
}

func readQueryExtensionReply(r *Reader, v *QueryExtensionReply) error {
	if !r.RemainAtLeast4b(3) {
		return ErrDataLenShort
	}

	r.ReadPad(8) // 2

	v.Present = r.ReadBool()
	v.MajorOpcode = r.Read1b()
	v.FirstEvent = r.Read1b()
	v.FirstError = r.Read1b() // 3

	return nil
}

// #WREQ
func encodeListExtensions() (hd uint8, b RequestBody) {
	return
}

type ListExtensionsReply struct {
	Names []string
}

type Str struct {
	Value string
}

func ReadStr(r *Reader) (string, error) {
	var v Str
	_, err := readStr(r, &v)
	if err != nil {
		return "", err
	}
	return v.Value, nil
}

func readStr(r *Reader, v *Str) (int, error) {
	if !r.RemainAtLeast(1) {
		return 0, ErrDataLenShort
	}

	nameLen := int(r.Read1b())

	var err error
	v.Value, err = r.ReadString(nameLen)
	if err != nil {
		return 0, err
	}

	return 1 + nameLen, nil
}

func writeStr(b *FixedSizeBuf, str string) {
	str = TruncateStr(str, math.MaxUint8)
	b.Write1b(uint8(len(str)))
	b.WriteString(str)
}

func readListExtensionsReply(r *Reader, v *ListExtensionsReply) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}
	strsLen, _ := r.ReadReplyHeader() // 2

	// unused
	r.ReadPad(24) // 8

	if strsLen > 0 {
		names := make([]string, strsLen)
		for i := 0; i < int(strsLen); i++ {
			var str Str
			_, err := readStr(r, &str)
			if err != nil {
				return err
			}
			names[i] = str.Value
		}
		v.Names = names
	}

	return nil
}

// #WREQ
func encodeKillClient(resource uint32) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(resource).
		End()
	return
}

// #WREQ
func encodeGetKeyboardMapping(firstKeycode Keycode, count uint8) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write1b(uint8(firstKeycode)).
		Write1b(count).
		WritePad(2).
		End()
	return
}

type GetKeyboardMappingReply struct {
	KeysymsPerKeycode uint8
	Keysyms           []Keysym
}

func readGetKeyboardMappingReply(r *Reader, v *GetKeyboardMappingReply) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}
	var keysymsLen uint32
	v.KeysymsPerKeycode, keysymsLen = r.ReadReplyHeader() // 2

	// unused
	r.ReadPad(24) // 8

	if !r.RemainAtLeast4b(int(keysymsLen)) {
		return ErrDataLenShort
	}
	v.Keysyms = make([]Keysym, keysymsLen)
	for i := 0; i < int(keysymsLen); i++ {
		v.Keysyms[i] = Keysym(r.Read4b())
	}

	return nil
}

// #WREQ
func encodeSetScreenSaver(timeout, interval int16, preferBlanking,
	allowExposures uint8) (hd uint8, b RequestBody) {
	b.AddBlock(2).
		Write2b(uint16(timeout)).
		Write2b(uint16(interval)).
		Write1b(preferBlanking).
		Write1b(allowExposures).
		WritePad(2).
		End()
	return
}

// #WREQ
func encodeGetScreenSaver() (hd uint8, b RequestBody) {
	return
}

type GetScreenSaverReply struct {
	Timeout        uint16
	Interval       uint16
	PreferBlanking uint8
	AllowExposures uint8
}

func readGetScreenSaverReply(r *Reader, v *GetScreenSaverReply) error {
	if !r.RemainAtLeast4b(4) {
		return ErrDataLenShort
	}
	r.ReadPad(8) // 2

	v.Timeout = r.Read2b()
	v.Interval = r.Read2b()

	v.PreferBlanking = r.Read1b()
	v.AllowExposures = r.Read1b() // 4

	return nil
}

// #WREQ
func encodeForceScreenSaver(mode uint8) (hd uint8, b RequestBody) {
	hd = mode
	return
}

// #WREQ
func encodeNoOperation(n int) (hd uint8, b RequestBody) {
	b.AddBlock(n)
	return
}

// #WREQ
func encodeFreeCursor(cursor Cursor) (hd uint8, b RequestBody) {
	b.AddBlock(1).
		Write4b(uint32(cursor)).
		End()
	return
}

// #WREQ
func encodeChangeHosts(mode, family uint8, address string)  (hd uint8, b RequestBody) {
	address = TruncateStr(address, math.MaxUint16)
	hd = mode
	b.AddBlock(1).
		Write1b(family).
		Write1b(0).
		Write2b(uint16(len(address))).
		End()
	b.AddBytes([]byte(address))
	return
}