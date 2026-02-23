// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

// simple ('xcb', 'WINDOW')
type Window uint32

// simple ('xcb', 'PIXMAP')
type Pixmap uint32

// simple ('xcb', 'CURSOR')
type Cursor uint32

// simple ('xcb', 'FONT')
type Font uint32

// simple ('xcb', 'GCONTEXT')
type GContext uint32

// simple ('xcb', 'COLORMAP')
type Colormap uint32

// simple ('xcb', 'ATOM')
type Atom uint32

// simple ('xcb', 'DRAWABLE')
type Drawable uint32

// simple ('xcb', 'FONTABLE')
type Fontable uint32

// simple ('xcb', 'BOOL32')
type Bool32 uint32

// simple ('xcb', 'VISUALID')
type VisualID uint32

// simple ('xcb', 'TIMESTAMP')
type Timestamp uint32

// simple ('xcb', 'KEYSYM')
type Keysym uint32

// simple ('xcb', 'KEYCODE')
type Keycode uint8

// simple ('xcb', 'KEYCODE32')
type Keycode32 uint32

// simple ('xcb', 'BUTTON')
type Button uint8

// enum VisualClass
const (
	VisualClassStaticGray  = 0
	VisualClassGrayScale   = 1
	VisualClassStaticColor = 2
	VisualClassPseudoColor = 3
	VisualClassTrueColor   = 4
	VisualClassDirectColor = 5
)

// enum EventMask
const (
	EventMaskNoEvent              = 0
	EventMaskKeyPress             = 1
	EventMaskKeyRelease           = 2
	EventMaskButtonPress          = 4
	EventMaskButtonRelease        = 8
	EventMaskEnterWindow          = 16
	EventMaskLeaveWindow          = 32
	EventMaskPointerMotion        = 64
	EventMaskPointerMotionHint    = 128
	EventMaskButton1Motion        = 256
	EventMaskButton2Motion        = 512
	EventMaskButton3Motion        = 1024
	EventMaskButton4Motion        = 2048
	EventMaskButton5Motion        = 4096
	EventMaskButtonMotion         = 8192
	EventMaskKeymapState          = 16384
	EventMaskExposure             = 32768
	EventMaskVisibilityChange     = 65536
	EventMaskStructureNotify      = 131072
	EventMaskResizeRedirect       = 262144
	EventMaskSubstructureNotify   = 524288
	EventMaskSubstructureRedirect = 1048576
	EventMaskFocusChange          = 2097152
	EventMaskPropertyChange       = 4194304
	EventMaskColorMapChange       = 8388608
	EventMaskOwnerGrabButton      = 16777216
)

// enum BackingStore
const (
	BackingStoreNotUseful  = 0
	BackingStoreWhenMapped = 1
	BackingStoreAlways     = 2
)

// enum ImageOrder
const (
	ImageOrderLSBFirst = 0
	ImageOrderMSBFirst = 1
)

// enum ModMask
const (
	ModMaskShift   = 1
	ModMaskLock    = 2
	ModMaskControl = 4
	ModMask1       = 8
	ModMask2       = 16
	ModMask3       = 32
	ModMask4       = 64
	ModMask5       = 128
	ModMaskAny     = 32768
)

// enum KeyButMask
const (
	KeyButMaskShift   = 1
	KeyButMaskLock    = 2
	KeyButMaskControl = 4
	KeyButMaskMod1    = 8
	KeyButMaskMod2    = 16
	KeyButMaskMod3    = 32
	KeyButMaskMod4    = 64
	KeyButMaskMod5    = 128
	KeyButMaskButton1 = 256
	KeyButMaskButton2 = 512
	KeyButMaskButton3 = 1024
	KeyButMaskButton4 = 2048
	KeyButMaskButton5 = 4096
)

// enum Window
const (
	WindowNone = 0
)

const KeyPressEventCode = 2

func NewKeyPressEvent(data []byte) (*KeyPressEvent, error) {
	var ev KeyPressEvent
	r := NewReaderFromData(data)
	err := readKeyPressEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const KeyReleaseEventCode = 3

func NewKeyReleaseEvent(data []byte) (*KeyReleaseEvent, error) {
	var ev KeyReleaseEvent
	r := NewReaderFromData(data)
	err := readKeyReleaseEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

// enum ButtonMask
const (
	ButtonMask1   = 256
	ButtonMask2   = 512
	ButtonMask3   = 1024
	ButtonMask4   = 2048
	ButtonMask5   = 4096
	ButtonMaskAny = 32768
)

const ButtonPressEventCode = 4

func NewButtonPressEvent(data []byte) (*ButtonPressEvent, error) {
	var ev ButtonPressEvent
	r := NewReaderFromData(data)
	err := readButtonPressEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const ButtonReleaseEventCode = 5

func NewButtonReleaseEvent(data []byte) (*ButtonReleaseEvent, error) {
	var ev ButtonReleaseEvent
	r := NewReaderFromData(data)
	err := readButtonReleaseEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

// enum Motion
const (
	MotionNormal = 0
	MotionHint   = 1
)

const MotionNotifyEventCode = 6

func NewMotionNotifyEvent(data []byte) (*MotionNotifyEvent, error) {
	var ev MotionNotifyEvent
	r := NewReaderFromData(data)
	err := readMotionNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

// enum NotifyDetail
const (
	NotifyDetailAncestor         = 0
	NotifyDetailVirtual          = 1
	NotifyDetailInferior         = 2
	NotifyDetailNonlinear        = 3
	NotifyDetailNonlinearVirtual = 4
	NotifyDetailPointer          = 5
	NotifyDetailPointerRoot      = 6
	NotifyDetailNone             = 7
)

// enum NotifyMode
const (
	NotifyModeNormal       = 0
	NotifyModeGrab         = 1
	NotifyModeUngrab       = 2
	NotifyModeWhileGrabbed = 3
)

const EnterNotifyEventCode = 7

func NewEnterNotifyEvent(data []byte) (*EnterNotifyEvent, error) {
	var ev EnterNotifyEvent
	r := NewReaderFromData(data)
	err := readEnterNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const LeaveNotifyEventCode = 8

func NewLeaveNotifyEvent(data []byte) (*LeaveNotifyEvent, error) {
	var ev LeaveNotifyEvent
	r := NewReaderFromData(data)
	err := readLeaveNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const FocusInEventCode = 9

func NewFocusInEvent(data []byte) (*FocusInEvent, error) {
	var ev FocusInEvent
	r := NewReaderFromData(data)
	err := readFocusInEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const FocusOutEventCode = 10

func NewFocusOutEvent(data []byte) (*FocusOutEvent, error) {
	var ev FocusOutEvent
	r := NewReaderFromData(data)
	err := readFocusOutEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const KeymapNotifyEventCode = 11

func NewKeymapNotifyEvent(data []byte) (*KeymapNotifyEvent, error) {
	var ev KeymapNotifyEvent
	r := NewReaderFromData(data)
	err := readKeymapNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const ExposeEventCode = 12

func NewExposeEvent(data []byte) (*ExposeEvent, error) {
	var ev ExposeEvent
	r := NewReaderFromData(data)
	err := readExposeEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const GraphicsExposureEventCode = 13

func NewGraphicsExposureEvent(data []byte) (*GraphicsExposureEvent, error) {
	var ev GraphicsExposureEvent
	r := NewReaderFromData(data)
	err := readGraphicsExposureEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const NoExposureEventCode = 14

func NewNoExposureEvent(data []byte) (*NoExposureEvent, error) {
	var ev NoExposureEvent
	r := NewReaderFromData(data)
	err := readNoExposureEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

// enum Visibility
const (
	VisibilityUnobscured        = 0
	VisibilityPartiallyObscured = 1
	VisibilityFullyObscured     = 2
)

const VisibilityNotifyEventCode = 15

func NewVisibilityNotifyEvent(data []byte) (*VisibilityNotifyEvent, error) {
	var ev VisibilityNotifyEvent
	r := NewReaderFromData(data)
	err := readVisibilityNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const CreateNotifyEventCode = 16

func NewCreateNotifyEvent(data []byte) (*CreateNotifyEvent, error) {
	var ev CreateNotifyEvent
	r := NewReaderFromData(data)
	err := readCreateNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const DestroyNotifyEventCode = 17

func NewDestroyNotifyEvent(data []byte) (*DestroyNotifyEvent, error) {
	var ev DestroyNotifyEvent
	r := NewReaderFromData(data)
	err := readDestroyNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const UnmapNotifyEventCode = 18

func NewUnmapNotifyEvent(data []byte) (*UnmapNotifyEvent, error) {
	var ev UnmapNotifyEvent
	r := NewReaderFromData(data)
	err := readUnmapNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const MapNotifyEventCode = 19

func NewMapNotifyEvent(data []byte) (*MapNotifyEvent, error) {
	var ev MapNotifyEvent
	r := NewReaderFromData(data)
	err := readMapNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const MapRequestEventCode = 20

func NewMapRequestEvent(data []byte) (*MapRequestEvent, error) {
	var ev MapRequestEvent
	r := NewReaderFromData(data)
	err := readMapRequestEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const ReparentNotifyEventCode = 21

func NewReparentNotifyEvent(data []byte) (*ReparentNotifyEvent, error) {
	var ev ReparentNotifyEvent
	r := NewReaderFromData(data)
	err := readReparentNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const ConfigureNotifyEventCode = 22

func NewConfigureNotifyEvent(data []byte) (*ConfigureNotifyEvent, error) {
	var ev ConfigureNotifyEvent
	r := NewReaderFromData(data)
	err := readConfigureNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const ConfigureRequestEventCode = 23

func NewConfigureRequestEvent(data []byte) (*ConfigureRequestEvent, error) {
	var ev ConfigureRequestEvent
	r := NewReaderFromData(data)
	err := readConfigureRequestEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const GravityNotifyEventCode = 24

func NewGravityNotifyEvent(data []byte) (*GravityNotifyEvent, error) {
	var ev GravityNotifyEvent
	r := NewReaderFromData(data)
	err := readGravityNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const ResizeRequestEventCode = 25

func NewResizeRequestEvent(data []byte) (*ResizeRequestEvent, error) {
	var ev ResizeRequestEvent
	r := NewReaderFromData(data)
	err := readResizeRequestEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

// enum Place
const (
	PlaceOnTop    = 0
	PlaceOnBottom = 1
)

const CirculateNotifyEventCode = 26

func NewCirculateNotifyEvent(data []byte) (*CirculateNotifyEvent, error) {
	var ev CirculateNotifyEvent
	r := NewReaderFromData(data)
	err := readCirculateNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const CirculateRequestEventCode = 27

func NewCirculateRequestEvent(data []byte) (*CirculateRequestEvent, error) {
	var ev CirculateRequestEvent
	r := NewReaderFromData(data)
	err := readCirculateRequestEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

// enum Property
const (
	PropertyNewValue = 0
	PropertyDelete   = 1
)

const PropertyNotifyEventCode = 28

func NewPropertyNotifyEvent(data []byte) (*PropertyNotifyEvent, error) {
	var ev PropertyNotifyEvent
	r := NewReaderFromData(data)
	err := readPropertyNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const SelectionClearEventCode = 29

func NewSelectionClearEvent(data []byte) (*SelectionClearEvent, error) {
	var ev SelectionClearEvent
	r := NewReaderFromData(data)
	err := readSelectionClearEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

// enum Time
const (
	TimeCurrentTime = 0
)

// enum Atom
const (
	AtomNone               = 0
	AtomAny                = 0
	AtomPrimary            = 1
	AtomSecondary          = 2
	AtomArc                = 3
	AtomAtom               = 4
	AtomBitmap             = 5
	AtomCardinal           = 6
	AtomColormap           = 7
	AtomCursor             = 8
	AtomCutBuffer0         = 9
	AtomCutBuffer1         = 10
	AtomCutBuffer2         = 11
	AtomCutBuffer3         = 12
	AtomCutBuffer4         = 13
	AtomCutBuffer5         = 14
	AtomCutBuffer6         = 15
	AtomCutBuffer7         = 16
	AtomDrawable           = 17
	AtomFont               = 18
	AtomInteger            = 19
	AtomPixmap             = 20
	AtomPoint              = 21
	AtomRectangle          = 22
	AtomResourceManager    = 23
	AtomRGBColorMap        = 24
	AtomRGBBestMap         = 25
	AtomRGBBlueMap         = 26
	AtomRGBDefaultMap      = 27
	AtomRGBGrayMap         = 28
	AtomRGBGreenMap        = 29
	AtomRGBRedMap          = 30
	AtomString             = 31
	AtomVisualID           = 32
	AtomWindow             = 33
	AtomWMCommand          = 34
	AtomWMHints            = 35
	AtomWMClientMachine    = 36
	AtomWMIconName         = 37
	AtomWMIconSize         = 38
	AtomWMName             = 39
	AtomWMNormalHints      = 40
	AtomWMSizeHints        = 41
	AtomWMZoomHints        = 42
	AtomMinSpace           = 43
	AtomNormSpace          = 44
	AtomMaxSpace           = 45
	AtomEndSpace           = 46
	AtomSuperscriptX       = 47
	AtomSuperscriptY       = 48
	AtomSubscriptX         = 49
	AtomSubscriptY         = 50
	AtomUnderlinePosition  = 51
	AtomUnderlineThickness = 52
	AtomStrikeoutAscent    = 53
	AtomStrikeoutDescent   = 54
	AtomItalicAngle        = 55
	AtomXHeight            = 56
	AtomQuadWidth          = 57
	AtomWeight             = 58
	AtomPointSize          = 59
	AtomResolution         = 60
	AtomCopyright          = 61
	AtomNotice             = 62
	AtomFontName           = 63
	AtomFamilyName         = 64
	AtomFullName           = 65
	AtomCapHeight          = 66
	AtomWMClass            = 67
	AtomWMTransientFor     = 68
)

const SelectionRequestEventCode = 30

func NewSelectionRequestEvent(data []byte) (*SelectionRequestEvent, error) {
	var ev SelectionRequestEvent
	r := NewReaderFromData(data)
	err := readSelectionRequestEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const SelectionNotifyEventCode = 31

func NewSelectionNotifyEvent(data []byte) (*SelectionNotifyEvent, error) {
	var ev SelectionNotifyEvent
	r := NewReaderFromData(data)
	err := readSelectionNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

// enum ColormapState
const (
	ColormapStateUninstalled = 0
	ColormapStateInstalled   = 1
)

// enum Colormap
const (
	ColormapNone = 0
)

const ColormapNotifyEventCode = 32

func NewColormapNotifyEvent(data []byte) (*ColormapNotifyEvent, error) {
	var ev ColormapNotifyEvent
	r := NewReaderFromData(data)
	err := readColormapNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const ClientMessageEventCode = 33

func NewClientMessageEvent(data []byte) (*ClientMessageEvent, error) {
	var ev ClientMessageEvent
	r := NewReaderFromData(data)
	err := readClientMessageEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

// enum Mapping
const (
	MappingModifier = 0
	MappingKeyboard = 1
	MappingPointer  = 2
)

const MappingNotifyEventCode = 34

func NewMappingNotifyEvent(data []byte) (*MappingNotifyEvent, error) {
	var ev MappingNotifyEvent
	r := NewReaderFromData(data)
	err := readMappingNotifyEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const GeGenericEventCode = 35

func NewGeGenericEvent(data []byte) (*GeGenericEvent, error) {
	var ev GeGenericEvent
	r := NewReaderFromData(data)
	err := readGeGenericEvent(r, &ev)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

const RequestErrorCode = 1
const ValueErrorCode = 2
const WindowErrorCode = 3
const PixmapErrorCode = 4
const AtomErrorCode = 5
const CursorErrorCode = 6
const FontErrorCode = 7
const MatchErrorCode = 8
const DrawableErrorCode = 9
const AccessErrorCode = 10
const AllocErrorCode = 11
const ColormapErrorCode = 12
const GContextErrorCode = 13
const IdChoiceErrorCode = 14
const NameErrorCode = 15
const LengthErrorCode = 16
const ImplementationErrorCode = 17

// enum WindowClass
const (
	WindowClassCopyFromParent = 0
	WindowClassInputOutput    = 1
	WindowClassInputOnly      = 2
)

// enum CW
const (
	CWBackPixmap       = 1
	CWBackPixel        = 2
	CWBorderPixmap     = 4
	CWBorderPixel      = 8
	CWBitGravity       = 16
	CWWinGravity       = 32
	CWBackingStore     = 64
	CWBackingPlanes    = 128
	CWBackingPixel     = 256
	CWOverrideRedirect = 512
	CWSaveUnder        = 1024
	CWEventMask        = 2048
	CWDontPropagate    = 4096
	CWColormap         = 8192
	CWCursor           = 16384
)

// enum BackPixmap
const (
	BackPixmapNone           = 0
	BackPixmapParentRelative = 1
)

// enum Gravity
const (
	GravityBitForget = 0
	GravityWinUnmap  = 0
	GravityNorthWest = 1
	GravityNorth     = 2
	GravityNorthEast = 3
	GravityWest      = 4
	GravityCenter    = 5
	GravityEast      = 6
	GravitySouthWest = 7
	GravitySouth     = 8
	GravitySouthEast = 9
	GravityStatic    = 10
)

const CreateWindowOpcode = 1
const ChangeWindowAttributesOpcode = 2

// enum MapState
const (
	MapStateUnmapped   = 0
	MapStateUnviewable = 1
	MapStateViewable   = 2
)

const GetWindowAttributesOpcode = 3

type GetWindowAttributesCookie SeqNum

const DestroyWindowOpcode = 4
const DestroySubwindowsOpcode = 5

// enum SetMode
const (
	SetModeInsert = 0
	SetModeDelete = 1
)

const ChangeSaveSetOpcode = 6
const ReparentWindowOpcode = 7
const MapWindowOpcode = 8
const MapSubwindowsOpcode = 9
const UnmapWindowOpcode = 10
const UnmapSubwindowsOpcode = 11

// enum ConfigWindow
const (
	ConfigWindowX           = 1
	ConfigWindowY           = 2
	ConfigWindowWidth       = 4
	ConfigWindowHeight      = 8
	ConfigWindowBorderWidth = 16
	ConfigWindowSibling     = 32
	ConfigWindowStackMode   = 64
)

// enum StackMode
const (
	StackModeAbove    = 0
	StackModeBelow    = 1
	StackModeTopIf    = 2
	StackModeBottomIf = 3
	StackModeOpposite = 4
)

const ConfigureWindowOpcode = 12

// enum Circulate
const (
	CirculateRaiseLowest  = 0
	CirculateLowerHighest = 1
)

const CirculateWindowOpcode = 13
const GetGeometryOpcode = 14

type GetGeometryCookie SeqNum

const QueryTreeOpcode = 15

type QueryTreeCookie SeqNum

const InternAtomOpcode = 16

type InternAtomCookie SeqNum

const GetAtomNameOpcode = 17

type GetAtomNameCookie SeqNum

// enum PropMode
const (
	PropModeReplace = 0
	PropModePrepend = 1
	PropModeAppend  = 2
)

const ChangePropertyOpcode = 18
const DeletePropertyOpcode = 19

// enum GetPropertyType
const (
	GetPropertyTypeAny = 0
)

const GetPropertyOpcode = 20

type GetPropertyCookie SeqNum

const ListPropertiesOpcode = 21

type ListPropertiesCookie SeqNum

const SetSelectionOwnerOpcode = 22
const GetSelectionOwnerOpcode = 23

type GetSelectionOwnerCookie SeqNum

const ConvertSelectionOpcode = 24

// enum SendEventDest
const (
	SendEventDestPointerWindow = 0
	SendEventDestItemFocus     = 1
)

const SendEventOpcode = 25

// enum GrabMode
const (
	GrabModeSync  = 0
	GrabModeAsync = 1
)

// enum GrabStatus
const (
	GrabStatusSuccess        = 0
	GrabStatusAlreadyGrabbed = 1
	GrabStatusInvalidTime    = 2
	GrabStatusNotViewable    = 3
	GrabStatusFrozen         = 4
)

// enum Cursor
const (
	CursorNone = 0
)

const GrabPointerOpcode = 26

type GrabPointerCookie SeqNum

const UngrabPointerOpcode = 27

// enum ButtonIndex
const (
	ButtonIndexAny = 0
	ButtonIndex1   = 1
	ButtonIndex2   = 2
	ButtonIndex3   = 3
	ButtonIndex4   = 4
	ButtonIndex5   = 5
)

const GrabButtonOpcode = 28
const UngrabButtonOpcode = 29
const ChangeActivePointerGrabOpcode = 30
const GrabKeyboardOpcode = 31

type GrabKeyboardCookie SeqNum

const UngrabKeyboardOpcode = 32

// enum Grab
const (
	GrabAny = 0
)

const GrabKeyOpcode = 33
const UngrabKeyOpcode = 34

// enum Allow
const (
	AllowAsyncPointer   = 0
	AllowSyncPointer    = 1
	AllowReplayPointer  = 2
	AllowAsyncKeyboard  = 3
	AllowSyncKeyboard   = 4
	AllowReplayKeyboard = 5
	AllowAsyncBoth      = 6
	AllowSyncBoth       = 7
)

const AllowEventsOpcode = 35
const GrabServerOpcode = 36
const UngrabServerOpcode = 37
const QueryPointerOpcode = 38

type QueryPointerCookie SeqNum

const GetMotionEventsOpcode = 39

type GetMotionEventsCookie SeqNum

const TranslateCoordinatesOpcode = 40

type TranslateCoordinatesCookie SeqNum

const WarpPointerOpcode = 41

// enum InputFocus
const (
	InputFocusNone           = 0
	InputFocusPointerRoot    = 1
	InputFocusParent         = 2
	InputFocusFollowKeyboard = 3
)

const SetInputFocusOpcode = 42
const GetInputFocusOpcode = 43

type GetInputFocusCookie SeqNum

const QueryKeymapOpcode = 44

type QueryKeymapCookie SeqNum

const OpenFontOpcode = 45
const CloseFontOpcode = 46

// enum FontDraw
const (
	FontDrawLeftToRight = 0
	FontDrawRightToLeft = 1
)

const QueryFontOpcode = 47

type QueryFontCookie SeqNum

const QueryTextExtentsOpcode = 48

type QueryTextExtentsCookie SeqNum

const ListFontsOpcode = 49

type ListFontsCookie SeqNum

const ListFontsWithInfoOpcode = 50

type ListFontsWithInfoCookie SeqNum

const SetFontPathOpcode = 51
const GetFontPathOpcode = 52

type GetFontPathCookie SeqNum

const CreatePixmapOpcode = 53
const FreePixmapOpcode = 54

// enum GC
const (
	GCFunction           = 1
	GCPlaneMask          = 2
	GCForeground         = 4
	GCBackground         = 8
	GCLineWidth          = 16
	GCLineStyle          = 32
	GCCapStyle           = 64
	GCJoinStyle          = 128
	GCFillStyle          = 256
	GCFillRule           = 512
	GCTile               = 1024
	GCStipple            = 2048
	GCTileStippleOriginX = 4096
	GCTileStippleOriginY = 8192
	GCFont               = 16384
	GCSubwindowMode      = 32768
	GCGraphicsExposures  = 65536
	GCClipOriginX        = 131072
	GCClipOriginY        = 262144
	GCClipMask           = 524288
	GCDashOffset         = 1048576
	GCDashList           = 2097152
	GCArcMode            = 4194304
)

// enum GX
const (
	GXClear        = 0
	GXAnd          = 1
	GXAndReverse   = 2
	GXCopy         = 3
	GXAndInverted  = 4
	GXNoop         = 5
	GXXor          = 6
	GXOr           = 7
	GXNor          = 8
	GXEquiv        = 9
	GXInvert       = 10
	GXOrReverse    = 11
	GXCopyInverted = 12
	GXOrInverted   = 13
	GXNand         = 14
	GXSet          = 15
)

// enum LineStyle
const (
	LineStyleSolid      = 0
	LineStyleOnOffDash  = 1
	LineStyleDoubleDash = 2
)

// enum CapStyle
const (
	CapStyleNotLast    = 0
	CapStyleButt       = 1
	CapStyleRound      = 2
	CapStyleProjecting = 3
)

// enum JoinStyle
const (
	JoinStyleMiter = 0
	JoinStyleRound = 1
	JoinStyleBevel = 2
)

// enum FillStyle
const (
	FillStyleSolid          = 0
	FillStyleTiled          = 1
	FillStyleStippled       = 2
	FillStyleOpaqueStippled = 3
)

// enum FillRule
const (
	FillRuleEvenOdd = 0
	FillRuleWinding = 1
)

// enum SubwindowMode
const (
	SubwindowModeClipByChildren   = 0
	SubwindowModeIncludeInferiors = 1
)

// enum ArcMode
const (
	ArcModeChord    = 0
	ArcModePieSlice = 1
)

const CreateGCOpcode = 55
const ChangeGCOpcode = 56
const CopyGCOpcode = 57
const SetDashesOpcode = 58

// enum ClipOrdering
const (
	ClipOrderingUnsorted = 0
	ClipOrderingYSorted  = 1
	ClipOrderingYXSorted = 2
	ClipOrderingYXBanded = 3
)

const SetClipRectanglesOpcode = 59
const FreeGCOpcode = 60
const ClearAreaOpcode = 61
const CopyAreaOpcode = 62
const CopyPlaneOpcode = 63

// enum CoordMode
const (
	CoordModeOrigin   = 0
	CoordModePrevious = 1
)

const PolyPointOpcode = 64
const PolyLineOpcode = 65
const PolySegmentOpcode = 66
const PolyRectangleOpcode = 67
const PolyArcOpcode = 68

// enum PolyShape
const (
	PolyShapeComplex   = 0
	PolyShapeNonconvex = 1
	PolyShapeConvex    = 2
)

const FillPolyOpcode = 69
const PolyFillRectangleOpcode = 70
const PolyFillArcOpcode = 71

// enum ImageFormat
const (
	ImageFormatXYBitmap = 0
	ImageFormatXYPixmap = 1
	ImageFormatZPixmap  = 2
)

const PutImageOpcode = 72
const GetImageOpcode = 73

type GetImageCookie SeqNum

const PolyText8Opcode = 74
const PolyText16Opcode = 75
const ImageText8Opcode = 76
const ImageText16Opcode = 77

// enum ColormapAlloc
const (
	ColormapAllocNone = 0
	ColormapAllocAll  = 1
)

const CreateColormapOpcode = 78
const FreeColormapOpcode = 79
const CopyColormapAndFreeOpcode = 80
const InstallColormapOpcode = 81
const UninstallColormapOpcode = 82
const ListInstalledColormapsOpcode = 83

type ListInstalledColormapsCookie SeqNum

const AllocColorOpcode = 84

type AllocColorCookie SeqNum

const AllocNamedColorOpcode = 85

type AllocNamedColorCookie SeqNum

const AllocColorCellsOpcode = 86

type AllocColorCellsCookie SeqNum

const AllocColorPlanesOpcode = 87

type AllocColorPlanesCookie SeqNum

const FreeColorsOpcode = 88

// enum ColorFlag
const (
	ColorFlagRed   = 1
	ColorFlagGreen = 2
	ColorFlagBlue  = 4
)

const StoreColorsOpcode = 89
const StoreNamedColorOpcode = 90
const QueryColorsOpcode = 91

type QueryColorsCookie SeqNum

const LookupColorOpcode = 92

type LookupColorCookie SeqNum

// enum Pixmap
const (
	PixmapNone = 0
)

const CreateCursorOpcode = 93

// enum Font
const (
	FontNone = 0
)

const CreateGlyphCursorOpcode = 94
const FreeCursorOpcode = 95
const RecolorCursorOpcode = 96

// enum QueryShapeOf
const (
	QueryShapeOfLargestCursor  = 0
	QueryShapeOfFastestTile    = 1
	QueryShapeOfFastestStipple = 2
)

const QueryBestSizeOpcode = 97

type QueryBestSizeCookie SeqNum

const QueryExtensionOpcode = 98

type QueryExtensionCookie SeqNum

const ListExtensionsOpcode = 99

type ListExtensionsCookie SeqNum

const ChangeKeyboardMappingOpcode = 100
const GetKeyboardMappingOpcode = 101

type GetKeyboardMappingCookie SeqNum

// enum KB
const (
	KBKeyClickPercent = 1
	KBBellPercent     = 2
	KBBellPitch       = 4
	KBBellDuration    = 8
	KBLed             = 16
	KBLedMode         = 32
	KBKey             = 64
	KBAutoRepeatMode  = 128
)

// enum LedMode
const (
	LedModeOff = 0
	LedModeOn  = 1
)

// enum AutoRepeatMode
const (
	AutoRepeatModeOff     = 0
	AutoRepeatModeOn      = 1
	AutoRepeatModeDefault = 2
)

const ChangeKeyboardControlOpcode = 102
const GetKeyboardControlOpcode = 103

type GetKeyboardControlCookie SeqNum

const BellOpcode = 104
const ChangePointerControlOpcode = 105
const GetPointerControlOpcode = 106

type GetPointerControlCookie SeqNum

// enum Blanking
const (
	BlankingNotPreferred = 0
	BlankingPreferred    = 1
	BlankingDefault      = 2
)

// enum Exposures
const (
	ExposuresNotAllowed = 0
	ExposuresAllowed    = 1
	ExposuresDefault    = 2
)

const SetScreenSaverOpcode = 107
const GetScreenSaverOpcode = 108

type GetScreenSaverCookie SeqNum

// enum HostMode
const (
	HostModeInsert = 0
	HostModeDelete = 1
)

// enum Family
const (
	FamilyInternet          = 0
	FamilyDECnet            = 1
	FamilyChaos             = 2
	FamilyServerInterpreted = 5
	FamilyInternet6         = 6
)

const ChangeHostsOpcode = 109
const ListHostsOpcode = 110

type ListHostsCookie SeqNum

// enum AccessControl
const (
	AccessControlDisable = 0
	AccessControlEnable  = 1
)

const SetAccessControlOpcode = 111

// enum CloseDown
const (
	CloseDownDestroyAll      = 0
	CloseDownRetainPermanent = 1
	CloseDownRetainTemporary = 2
)

const SetCloseDownModeOpcode = 112

// enum Kill
const (
	KillAllTemporary = 0
)

const KillClientOpcode = 113
const RotatePropertiesOpcode = 114

// enum ScreenSaver
const (
	ScreenSaverReset  = 0
	ScreenSaverActive = 1
)

const ForceScreenSaverOpcode = 115

// enum MappingStatus
const (
	MappingStatusSuccess = 0
	MappingStatusBusy    = 1
	MappingStatusFailure = 2
)

const SetPointerMappingOpcode = 116

type SetPointerMappingCookie SeqNum

const GetPointerMappingOpcode = 117

type GetPointerMappingCookie SeqNum

// enum MapIndex
const (
	MapIndexShift   = 0
	MapIndexLock    = 1
	MapIndexControl = 2
	MapIndex1       = 3
	MapIndex2       = 4
	MapIndex3       = 5
	MapIndex4       = 6
	MapIndex5       = 7
)

const SetModifierMappingOpcode = 118

type SetModifierMappingCookie SeqNum

const GetModifierMappingOpcode = 119

type GetModifierMappingCookie SeqNum

const NoOperationOpcode = 127

var errorCodeNameMap = map[uint8]string{
	RequestErrorCode:        "BadRequest",
	ValueErrorCode:          "BadValue",
	WindowErrorCode:         "BadWindow",
	PixmapErrorCode:         "BadPixmap",
	AtomErrorCode:           "BadAtom",
	CursorErrorCode:         "BadCursor",
	FontErrorCode:           "BadFont",
	MatchErrorCode:          "BadMatch",
	DrawableErrorCode:       "BadDrawable",
	AccessErrorCode:         "BadAccess",
	AllocErrorCode:          "BadAlloc",
	ColormapErrorCode:       "BadColormap",
	GContextErrorCode:       "BadGContext",
	IdChoiceErrorCode:       "BadIdChoice",
	NameErrorCode:           "BadName",
	LengthErrorCode:         "BadLength",
	ImplementationErrorCode: "BadImplementation",
}
var requestOpcodeNameMap = map[uint]string{
	CreateWindowOpcode:            "CreateWindow",
	ChangeWindowAttributesOpcode:  "ChangeWindowAttributes",
	GetWindowAttributesOpcode:     "GetWindowAttributes",
	DestroyWindowOpcode:           "DestroyWindow",
	DestroySubwindowsOpcode:       "DestroySubwindows",
	ChangeSaveSetOpcode:           "ChangeSaveSet",
	ReparentWindowOpcode:          "ReparentWindow",
	MapWindowOpcode:               "MapWindow",
	MapSubwindowsOpcode:           "MapSubwindows",
	UnmapWindowOpcode:             "UnmapWindow",
	UnmapSubwindowsOpcode:         "UnmapSubwindows",
	ConfigureWindowOpcode:         "ConfigureWindow",
	CirculateWindowOpcode:         "CirculateWindow",
	GetGeometryOpcode:             "GetGeometry",
	QueryTreeOpcode:               "QueryTree",
	InternAtomOpcode:              "InternAtom",
	GetAtomNameOpcode:             "GetAtomName",
	ChangePropertyOpcode:          "ChangeProperty",
	DeletePropertyOpcode:          "DeleteProperty",
	GetPropertyOpcode:             "GetProperty",
	ListPropertiesOpcode:          "ListProperties",
	SetSelectionOwnerOpcode:       "SetSelectionOwner",
	GetSelectionOwnerOpcode:       "GetSelectionOwner",
	ConvertSelectionOpcode:        "ConvertSelection",
	SendEventOpcode:               "SendEvent",
	GrabPointerOpcode:             "GrabPointer",
	UngrabPointerOpcode:           "UngrabPointer",
	GrabButtonOpcode:              "GrabButton",
	UngrabButtonOpcode:            "UngrabButton",
	ChangeActivePointerGrabOpcode: "ChangeActivePointerGrab",
	GrabKeyboardOpcode:            "GrabKeyboard",
	UngrabKeyboardOpcode:          "UngrabKeyboard",
	GrabKeyOpcode:                 "GrabKey",
	UngrabKeyOpcode:               "UngrabKey",
	AllowEventsOpcode:             "AllowEvents",
	GrabServerOpcode:              "GrabServer",
	UngrabServerOpcode:            "UngrabServer",
	QueryPointerOpcode:            "QueryPointer",
	GetMotionEventsOpcode:         "GetMotionEvents",
	TranslateCoordinatesOpcode:    "TranslateCoordinates",
	WarpPointerOpcode:             "WarpPointer",
	SetInputFocusOpcode:           "SetInputFocus",
	GetInputFocusOpcode:           "GetInputFocus",
	QueryKeymapOpcode:             "QueryKeymap",
	OpenFontOpcode:                "OpenFont",
	CloseFontOpcode:               "CloseFont",
	QueryFontOpcode:               "QueryFont",
	QueryTextExtentsOpcode:        "QueryTextExtents",
	ListFontsOpcode:               "ListFonts",
	ListFontsWithInfoOpcode:       "ListFontsWithInfo",
	SetFontPathOpcode:             "SetFontPath",
	GetFontPathOpcode:             "GetFontPath",
	CreatePixmapOpcode:            "CreatePixmap",
	FreePixmapOpcode:              "FreePixmap",
	CreateGCOpcode:                "CreateGC",
	ChangeGCOpcode:                "ChangeGC",
	CopyGCOpcode:                  "CopyGC",
	SetDashesOpcode:               "SetDashes",
	SetClipRectanglesOpcode:       "SetClipRectangles",
	FreeGCOpcode:                  "FreeGC",
	ClearAreaOpcode:               "ClearArea",
	CopyAreaOpcode:                "CopyArea",
	CopyPlaneOpcode:               "CopyPlane",
	PolyPointOpcode:               "PolyPoint",
	PolyLineOpcode:                "PolyLine",
	PolySegmentOpcode:             "PolySegment",
	PolyRectangleOpcode:           "PolyRectangle",
	PolyArcOpcode:                 "PolyArc",
	FillPolyOpcode:                "FillPoly",
	PolyFillRectangleOpcode:       "PolyFillRectangle",
	PolyFillArcOpcode:             "PolyFillArc",
	PutImageOpcode:                "PutImage",
	GetImageOpcode:                "GetImage",
	PolyText8Opcode:               "PolyText8",
	PolyText16Opcode:              "PolyText16",
	ImageText8Opcode:              "ImageText8",
	ImageText16Opcode:             "ImageText16",
	CreateColormapOpcode:          "CreateColormap",
	FreeColormapOpcode:            "FreeColormap",
	CopyColormapAndFreeOpcode:     "CopyColormapAndFree",
	InstallColormapOpcode:         "InstallColormap",
	UninstallColormapOpcode:       "UninstallColormap",
	ListInstalledColormapsOpcode:  "ListInstalledColormaps",
	AllocColorOpcode:              "AllocColor",
	AllocNamedColorOpcode:         "AllocNamedColor",
	AllocColorCellsOpcode:         "AllocColorCells",
	AllocColorPlanesOpcode:        "AllocColorPlanes",
	FreeColorsOpcode:              "FreeColors",
	StoreColorsOpcode:             "StoreColors",
	StoreNamedColorOpcode:         "StoreNamedColor",
	QueryColorsOpcode:             "QueryColors",
	LookupColorOpcode:             "LookupColor",
	CreateCursorOpcode:            "CreateCursor",
	CreateGlyphCursorOpcode:       "CreateGlyphCursor",
	FreeCursorOpcode:              "FreeCursor",
	RecolorCursorOpcode:           "RecolorCursor",
	QueryBestSizeOpcode:           "QueryBestSize",
	QueryExtensionOpcode:          "QueryExtension",
	ListExtensionsOpcode:          "ListExtensions",
	ChangeKeyboardMappingOpcode:   "ChangeKeyboardMapping",
	GetKeyboardMappingOpcode:      "GetKeyboardMapping",
	ChangeKeyboardControlOpcode:   "ChangeKeyboardControl",
	GetKeyboardControlOpcode:      "GetKeyboardControl",
	BellOpcode:                    "Bell",
	ChangePointerControlOpcode:    "ChangePointerControl",
	GetPointerControlOpcode:       "GetPointerControl",
	SetScreenSaverOpcode:          "SetScreenSaver",
	GetScreenSaverOpcode:          "GetScreenSaver",
	ChangeHostsOpcode:             "ChangeHosts",
	ListHostsOpcode:               "ListHosts",
	SetAccessControlOpcode:        "SetAccessControl",
	SetCloseDownModeOpcode:        "SetCloseDownMode",
	KillClientOpcode:              "KillClient",
	RotatePropertiesOpcode:        "RotateProperties",
	ForceScreenSaverOpcode:        "ForceScreenSaver",
	SetPointerMappingOpcode:       "SetPointerMapping",
	GetPointerMappingOpcode:       "GetPointerMapping",
	SetModifierMappingOpcode:      "SetModifierMapping",
	GetModifierMappingOpcode:      "GetModifierMapping",
	NoOperationOpcode:             "NoOperation",
}
