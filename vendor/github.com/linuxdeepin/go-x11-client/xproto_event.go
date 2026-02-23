// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

type KeyEvent struct {
	Detail     Keycode
	Sequence   uint16
	Time       Timestamp
	Root       Window
	Event      Window
	Child      Window
	RootX      int16
	RootY      int16
	EventX     int16
	EventY     int16
	State      uint16
	SameScreen bool
}

func readKeyEvent(r *Reader, v *KeyEvent) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}

	var detail uint8
	detail, v.Sequence = r.ReadEventHeader()
	v.Detail = Keycode(detail)

	v.Time = Timestamp(r.Read4b())

	v.Root = Window(r.Read4b())

	v.Event = Window(r.Read4b())

	v.Child = Window(r.Read4b()) // 5

	v.RootX = int16(r.Read2b())
	v.RootY = int16(r.Read2b())

	v.EventX = int16(r.Read2b())
	v.EventY = int16(r.Read2b())

	v.State = r.Read2b()
	v.SameScreen = r.ReadBool() // 8

	return nil
}

type KeyPressEvent struct {
	KeyEvent
}

func readKeyPressEvent(r *Reader, v *KeyPressEvent) error {
	return readKeyEvent(r, &v.KeyEvent)
}

type KeyReleaseEvent struct {
	KeyEvent
}

func readKeyReleaseEvent(r *Reader, v *KeyReleaseEvent) error {
	return readKeyEvent(r, &v.KeyEvent)
}

type ButtonEvent struct {
	Detail     Button
	Sequence   uint16
	Time       Timestamp
	Root       Window
	Event      Window
	Child      Window
	RootX      int16
	RootY      int16
	EventX     int16
	EventY     int16
	State      uint16
	SameScreen bool
}

func readButtonEvent(r *Reader, v *ButtonEvent) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}

	var detail uint8
	detail, v.Sequence = r.ReadEventHeader()
	v.Detail = Button(detail)

	v.Time = Timestamp(r.Read4b())

	v.Root = Window(r.Read4b())

	v.Event = Window(r.Read4b())

	v.Child = Window(r.Read4b()) // 5

	v.RootX = int16(r.Read2b())
	v.RootY = int16(r.Read2b())

	v.EventX = int16(r.Read2b())
	v.EventY = int16(r.Read2b()) // 7

	v.State = r.Read2b()
	v.SameScreen = r.ReadBool() // 8

	return nil
}

type ButtonPressEvent struct {
	ButtonEvent
}

func readButtonPressEvent(r *Reader, v *ButtonPressEvent) error {
	return readButtonEvent(r, &v.ButtonEvent)
}

type ButtonReleaseEvent struct {
	ButtonEvent
}

func readButtonReleaseEvent(r *Reader, v *ButtonReleaseEvent) error {
	return readButtonEvent(r, &v.ButtonEvent)
}

type MotionNotifyEvent struct {
	Detail     uint8
	Sequence   uint16
	Time       Timestamp
	Root       Window
	Event      Window
	Child      Window
	RootX      int16
	RootY      int16
	EventX     int16
	EventY     int16
	State      uint16
	SameScreen bool
}

func readMotionNotifyEvent(r *Reader, v *MotionNotifyEvent) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}

	v.Detail, v.Sequence = r.ReadEventHeader()

	v.Time = Timestamp(r.Read4b())

	v.Root = Window(r.Read4b())

	v.Event = Window(r.Read4b())

	v.Child = Window(r.Read4b()) // 5

	v.RootX = int16(r.Read2b())
	v.RootY = int16(r.Read2b())

	v.EventX = int16(r.Read2b())
	v.EventY = int16(r.Read2b())

	v.State = r.Read2b()
	v.SameScreen = r.ReadBool() // 8

	return nil
}

type PointerWindowEvent struct {
	Detail          uint8
	Sequence        uint16
	Time            Timestamp
	Root            Window
	Event           Window
	Child           Window
	RootX           int16
	RootY           int16
	EventX          int16
	EventY          int16
	State           uint16
	Mode            uint8
	SameScreenFocus uint8
}

func readPointerWindowEvent(r *Reader, v *PointerWindowEvent) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}

	v.Detail, v.Sequence = r.ReadEventHeader()

	v.Time = Timestamp(r.Read4b())

	v.Root = Window(r.Read4b())

	v.Event = Window(r.Read4b())

	v.Child = Window(r.Read4b()) // 5

	v.RootX = int16(r.Read2b())
	v.RootY = int16(r.Read2b())

	v.EventX = int16(r.Read2b())
	v.EventY = int16(r.Read2b())

	v.State = r.Read2b()
	v.Mode = r.Read1b()
	v.SameScreenFocus = r.Read1b() // 8

	return nil
}

type EnterNotifyEvent struct {
	PointerWindowEvent
}

func readEnterNotifyEvent(r *Reader, v *EnterNotifyEvent) error {
	return readPointerWindowEvent(r, &v.PointerWindowEvent)
}

type LeaveNotifyEvent struct {
	PointerWindowEvent
}

func readLeaveNotifyEvent(r *Reader, v *LeaveNotifyEvent) error {
	return readPointerWindowEvent(r, &v.PointerWindowEvent)
}

type FocusEvent struct {
	Detail   uint8
	Sequence uint16
	Event    Window
	Mode     uint8
}

func readFocusEvent(r *Reader, v *FocusEvent) error {
	if !r.RemainAtLeast4b(3) {
		return ErrDataLenShort
	}

	v.Detail, v.Sequence = r.ReadEventHeader()

	v.Event = Window(r.Read4b())

	v.Mode = r.Read1b() // 3

	return nil
}

type FocusInEvent struct {
	FocusEvent
}

func readFocusInEvent(r *Reader, v *FocusInEvent) error {
	return readFocusEvent(r, &v.FocusEvent)
}

type FocusOutEvent struct {
	FocusEvent
}

func readFocusOutEvent(r *Reader, v *FocusOutEvent) error {
	return readFocusEvent(r, &v.FocusEvent)
}

type KeymapNotifyEvent struct {
	Keys []byte
}

func readKeymapNotifyEvent(r *Reader, v *KeymapNotifyEvent) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}

	r.ReadPad(1)
	v.Keys = r.MustReadBytes(31) // 8
	return nil
}

type ExposeEvent struct {
	Sequence uint16
	Window   Window
	X        uint16
	Y        uint16
	Width    uint16
	Height   uint16
	Count    uint16
}

func readExposeEvent(r *Reader, v *ExposeEvent) error {
	if !r.RemainAtLeast4b(5) {
		return ErrDataLenShort
	}

	_, v.Sequence = r.ReadEventHeader()

	v.Window = Window(r.Read4b())

	v.X = r.Read2b()
	v.Y = r.Read2b() // 3

	v.Width = r.Read2b()
	v.Height = r.Read2b()

	v.Count = r.Read2b() // 5

	return nil
}

type GraphicsExposureEvent struct {
	Sequence    uint16
	Drawable    Drawable
	X           uint16
	Y           uint16
	Width       uint16
	Height      uint16
	MinorOpcode uint16
	Count       uint16
	MajorOpcode uint8
}

func readGraphicsExposureEvent(r *Reader, v *GraphicsExposureEvent) error {
	if !r.RemainAtLeast4b(6) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Drawable = Drawable(r.Read4b())

	v.X = r.Read2b()
	v.Y = r.Read2b() // 3

	v.Width = r.Read2b()
	v.Height = r.Read2b()

	v.MinorOpcode = r.Read2b()
	v.Count = r.Read2b()

	v.MajorOpcode = r.Read1b() // 6

	return nil
}

type NoExposureEvent struct {
	Sequence    uint16
	Drawable    Drawable
	MinorOpcode uint16
	MajorOpcode uint8
}

func readNoExposureEvent(r *Reader, v *NoExposureEvent) error {
	if !r.RemainAtLeast4b(3) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Drawable = Drawable(r.Read4b())

	v.MinorOpcode = r.Read2b()
	v.MajorOpcode = r.Read1b() // 3

	return nil
}

type VisibilityNotifyEvent struct {
	Sequence uint16
	Window   Window
	State    uint8
}

func readVisibilityNotifyEvent(r *Reader, v *VisibilityNotifyEvent) error {
	if !r.RemainAtLeast4b(3) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Window = Window(r.Read4b())

	v.State = r.Read1b() // 3

	return nil
}

type CreateNotifyEvent struct {
	Sequence         uint16
	Parent           Window
	Window           Window
	X                int16
	Y                int16
	Width            uint16
	Height           uint16
	BorderWidth      uint16
	OverrideRedirect bool
}

func readCreateNotifyEvent(r *Reader, v *CreateNotifyEvent) error {
	if !r.RemainAtLeast4b(5) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Window = Window(r.Read4b())

	v.X = int16(r.Read2b())
	v.Y = int16(r.Read2b()) // 3

	v.Width = r.Read2b()
	v.Height = r.Read2b()

	v.BorderWidth = r.Read2b()
	v.OverrideRedirect = r.ReadBool() // 5

	return nil
}

type DestroyNotifyEvent struct {
	Sequence uint16
	Event    Window
	Window   Window
}

func readDestroyNotifyEvent(r *Reader, v *DestroyNotifyEvent) error {
	if !r.RemainAtLeast4b(3) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Event = Window(r.Read4b())

	v.Window = Window(r.Read4b()) // 3

	return nil
}

type UnmapNotifyEvent struct {
	Sequence      uint16
	Event         Window
	Window        Window
	FromConfigure bool
}

func readUnmapNotifyEvent(r *Reader, v *UnmapNotifyEvent) error {
	if !r.RemainAtLeast4b(4) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Event = Window(r.Read4b())

	v.Window = Window(r.Read4b())

	v.FromConfigure = r.ReadBool() // 4

	return nil
}

type MapNotifyEvent struct {
	Sequence         uint16
	Event            Window
	Window           Window
	OverrideRedirect bool
}

func readMapNotifyEvent(r *Reader, v *MapNotifyEvent) error {
	if !r.RemainAtLeast4b(4) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Event = Window(r.Read4b())

	v.Window = Window(r.Read4b())

	v.OverrideRedirect = r.ReadBool() // 4

	return nil
}

type MapRequestEvent struct {
	Sequence uint16
	Parent   Window
	Window   Window
}

func readMapRequestEvent(r *Reader, v *MapRequestEvent) error {
	if !r.RemainAtLeast4b(3) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Parent = Window(r.Read4b())

	v.Window = Window(r.Read4b()) // 3

	return nil
}

type ReparentNotifyEvent struct {
	Sequence         uint16
	Event            Window
	Window           Window
	Parent           Window
	X                int16
	Y                int16
	OverrideRedirect bool
}

func readReparentNotifyEvent(r *Reader, v *ReparentNotifyEvent) error {
	if !r.RemainAtLeast4b(6) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Event = Window(r.Read4b())

	v.Window = Window(r.Read4b())

	v.Parent = Window(r.Read4b()) // 4

	v.X = int16(r.Read2b())
	v.Y = int16(r.Read2b())

	v.OverrideRedirect = r.ReadBool() // 6

	return nil
}

type ConfigureNotifyEvent struct {
	Sequence         uint16
	Event            Window
	Window           Window
	AboveSibling     Window
	X                int16
	Y                int16
	Width            uint16
	Height           uint16
	BorderWidth      uint16
	OverrideRedirect bool
}

func readConfigureNotifyEvent(r *Reader, v *ConfigureNotifyEvent) error {
	if !r.RemainAtLeast4b(7) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Event = Window(r.Read4b())

	v.Window = Window(r.Read4b())

	v.AboveSibling = Window(r.Read4b()) // 4

	v.X = int16(r.Read2b())
	v.Y = int16(r.Read2b())

	v.Width = r.Read2b()
	v.Height = r.Read2b()

	v.BorderWidth = r.Read2b()
	v.OverrideRedirect = r.ReadBool() // 7

	return nil
}

type ConfigureRequestEvent struct {
	StackMode   uint8
	Sequence    uint16
	Parent      Window
	Window      Window
	Sibling     Window
	X           int16
	Y           int16
	Width       uint16
	Height      uint16
	BorderWidth uint16
	ValueMask   uint16
}

func readConfigureRequestEvent(r *Reader, v *ConfigureRequestEvent) error {
	if !r.RemainAtLeast4b(7) {
		return ErrDataLenShort
	}
	v.StackMode, v.Sequence = r.ReadEventHeader()

	v.Parent = Window(r.Read4b())

	v.Window = Window(r.Read4b())

	v.Sibling = Window(r.Read4b())

	v.X = int16(r.Read2b())
	v.Y = int16(r.Read2b()) // 5

	v.Width = r.Read2b()
	v.Height = r.Read2b()

	v.BorderWidth = r.Read2b()
	v.ValueMask = r.Read2b() // 7

	return nil
}

type GravityNotifyEvent struct {
	Sequence uint16
	Event    Window
	Window   Window
	X        int16
	Y        int16
}

func readGravityNotifyEvent(r *Reader, v *GravityNotifyEvent) error {
	if !r.RemainAtLeast4b(4) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Event = Window(r.Read4b())

	v.Window = Window(r.Read4b())

	v.X = int16(r.Read2b())
	v.Y = int16(r.Read2b()) // 4

	return nil
}

type ResizeRequestEvent struct {
	Sequence uint16
	Window   Window
	Width    uint16
	Height   uint16
}

func readResizeRequestEvent(r *Reader, v *ResizeRequestEvent) error {
	if !r.RemainAtLeast4b(3) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Window = Window(r.Read4b())

	v.Width = r.Read2b()
	v.Height = r.Read2b() // 3

	return nil
}

type CirculateEvent struct {
	Sequence uint16
	Event    Window
	Window   Window
	Place    uint8
}

func readCirculateEvent(r *Reader, v *CirculateEvent) error {
	if !r.RemainAtLeast4b(5) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Event = Window(r.Read4b())

	v.Window = Window(r.Read4b())

	// unused
	r.ReadPad(4)

	v.Place = r.Read1b() // 5

	return nil
}

type CirculateNotifyEvent struct {
	CirculateEvent
}

func readCirculateNotifyEvent(r *Reader, v *CirculateNotifyEvent) error {
	return readCirculateEvent(r, &v.CirculateEvent)
}

type CirculateRequestEvent struct {
	CirculateEvent
}

func readCirculateRequestEvent(r *Reader, v *CirculateRequestEvent) error {
	return readCirculateEvent(r, &v.CirculateEvent)
}

type PropertyNotifyEvent struct {
	Sequence uint16
	Window   Window
	Atom     Atom
	Time     Timestamp
	State    uint8
}

func readPropertyNotifyEvent(r *Reader, v *PropertyNotifyEvent) error {
	if !r.RemainAtLeast4b(5) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Window = Window(r.Read4b())

	v.Atom = Atom(r.Read4b())

	v.Time = Timestamp(r.Read4b())

	v.State = r.Read1b() // 5

	return nil
}

type SelectionClearEvent struct {
	Sequence  uint16
	Time      Timestamp
	Owner     Window
	Selection Atom
}

func readSelectionClearEvent(r *Reader, v *SelectionClearEvent) error {
	if !r.RemainAtLeast4b(4) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Time = Timestamp(r.Read4b())

	v.Owner = Window(r.Read4b())

	v.Selection = Atom(r.Read4b()) // 4

	return nil
}

type SelectionRequestEvent struct {
	Sequence  uint16
	Time      Timestamp
	Owner     Window
	Requestor Window
	Selection Atom
	Target    Atom
	Property  Atom
}

func readSelectionRequestEvent(r *Reader, v *SelectionRequestEvent) error {
	if !r.RemainAtLeast4b(7) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Time = Timestamp(r.Read4b())

	v.Owner = Window(r.Read4b())

	v.Requestor = Window(r.Read4b())

	v.Selection = Atom(r.Read4b()) // 5

	v.Target = Atom(r.Read4b())

	v.Property = Atom(r.Read4b()) // 7

	return nil
}

type SelectionNotifyEvent struct {
	Sequence  uint16
	Time      Timestamp
	Requestor Window
	Selection Atom
	Target    Atom
	Property  Atom
}

func readSelectionNotifyEvent(r *Reader, v *SelectionNotifyEvent) error {
	if !r.RemainAtLeast4b(6) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Time = Timestamp(r.Read4b())

	v.Requestor = Window(r.Read4b())

	v.Selection = Atom(r.Read4b())

	v.Target = Atom(r.Read4b())

	v.Property = Atom(r.Read4b()) // 6

	return nil
}

func WriteSelectionNotifyEvent(w *Writer, v *SelectionNotifyEvent) {
	w.Write1b(SelectionNotifyEventCode)
	w.WritePad(1)
	w.Write2b(v.Sequence)
	w.Write4b(uint32(v.Time))
	w.Write4b(uint32(v.Requestor))
	w.Write4b(uint32(v.Selection))
	w.Write4b(uint32(v.Target))
	w.Write4b(uint32(v.Property))
	w.WritePad(8)
}

type ColormapNotifyEvent struct {
	Sequence uint16
	Window   Window
	Colormap Colormap
	New      bool
	State    uint8
}

func readColormapNotifyEvent(r *Reader, v *ColormapNotifyEvent) error {
	if !r.RemainAtLeast4b(4) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Window = Window(r.Read4b())

	v.Colormap = Colormap(r.Read4b()) // 3

	v.New = r.ReadBool()
	v.State = r.Read1b() // 4

	return nil
}

type ClientMessageEvent struct {
	Format   uint8
	Sequence uint16
	Window   Window
	Type     Atom
	Data     ClientMessageData
}

func readClientMessageEvent(r *Reader, v *ClientMessageEvent) error {
	if !r.RemainAtLeast4b(8) {
		return ErrDataLenShort
	}
	v.Format, v.Sequence = r.ReadEventHeader()

	v.Window = Window(r.Read4b())

	v.Type = Atom(r.Read4b()) // 3

	data := r.MustReadBytes(20) // 8
	v.Data = ClientMessageData{
		data: data,
	}

	return nil
}

func WriteClientMessageEvent(w *Writer, v *ClientMessageEvent) {
	w.Write1b(ClientMessageEventCode)
	w.Write1b(v.Format)
	w.Write2b(v.Sequence)
	w.Write4b(uint32(v.Window))
	w.Write4b(uint32(v.Type))
	writeClientMessageData(w, &v.Data)
}

type MappingNotifyEvent struct {
	Sequence     uint16
	Request      uint8
	FirstKeycode Keycode
	Count        uint8
}

func readMappingNotifyEvent(r *Reader, v *MappingNotifyEvent) error {
	if !r.RemainAtLeast4b(2) {
		return ErrDataLenShort
	}
	_, v.Sequence = r.ReadEventHeader()

	v.Request = r.Read1b()
	v.FirstKeycode = Keycode(r.Read1b())
	v.Count = r.Read1b() // 2

	return nil
}

type GeGenericEvent struct {
	Extension uint8
	Sequence  uint16
	Length    uint32
	EventType uint16
	Data      []byte
}

func readGeGenericEvent(r *Reader, v *GeGenericEvent) error {
	if !r.RemainAtLeast4b(3) {
		return ErrDataLenShort
	}
	v.Extension, v.Sequence = r.ReadEventHeader()

	v.Length = r.Read4b()

	v.EventType = r.Read2b() // 3

	var err error
	v.Data, err = r.ReadBytes(22 + (int(v.Length) * 4))
	return err
}
