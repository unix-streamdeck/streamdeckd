package api

import (
	"errors"
	"strings"

	"github.com/bendahl/uinput"
)

var keyNameToCode = map[string]int{
	"escape":    uinput.KeyEsc,
	"esc":       uinput.KeyEsc,
	"backspace": uinput.KeyBackspace,
	"tab":       uinput.KeyTab,
	"enter":     uinput.KeyEnter,
	"return":    uinput.KeyEnter,
	"space":     uinput.KeySpace,
	"delete":    uinput.KeyDelete,
	"insert":    uinput.KeyInsert,
	"home":      uinput.KeyHome,
	"end":       uinput.KeyEnd,
	"pageup":    uinput.KeyPageup,
	"page_up":   uinput.KeyPageup,
	"pagedown":  uinput.KeyPagedown,
	"page_down": uinput.KeyPagedown,
	"pause":     uinput.KeyPause,

	"up":    uinput.KeyUp,
	"down":  uinput.KeyDown,
	"left":  uinput.KeyLeft,
	"right": uinput.KeyRight,

	"ctrl":       uinput.KeyLeftctrl,
	"control":    uinput.KeyLeftctrl,
	"ctrl_l":     uinput.KeyLeftctrl,
	"control_l":  uinput.KeyLeftctrl,
	"ctrl_r":     uinput.KeyRightctrl,
	"control_r":  uinput.KeyRightctrl,
	"shift":      uinput.KeyLeftshift,
	"shift_l":    uinput.KeyLeftshift,
	"shift_r":    uinput.KeyRightshift,
	"alt":        uinput.KeyLeftalt,
	"alt_l":      uinput.KeyLeftalt,
	"alt_r":      uinput.KeyRightalt,
	"meta":       uinput.KeyLeftmeta,
	"super":      uinput.KeyLeftmeta,
	"super_l":    uinput.KeyLeftmeta,
	"super_r":    uinput.KeyRightmeta,
	"compose":    uinput.KeyCompose,
	"menu":       uinput.KeyMenu,
	"capslock":   uinput.KeyCapslock,
	"caps_lock":  uinput.KeyCapslock,
	"numlock":    uinput.KeyNumlock,
	"num_lock":   uinput.KeyNumlock,
	"scrolllock": uinput.KeyScrolllock,

	"a": uinput.KeyA,
	"b": uinput.KeyB,
	"c": uinput.KeyC,
	"d": uinput.KeyD,
	"e": uinput.KeyE,
	"f": uinput.KeyF,
	"g": uinput.KeyG,
	"h": uinput.KeyH,
	"i": uinput.KeyI,
	"j": uinput.KeyJ,
	"k": uinput.KeyK,
	"l": uinput.KeyL,
	"m": uinput.KeyM,
	"n": uinput.KeyN,
	"o": uinput.KeyO,
	"p": uinput.KeyP,
	"q": uinput.KeyQ,
	"r": uinput.KeyR,
	"s": uinput.KeyS,
	"t": uinput.KeyT,
	"u": uinput.KeyU,
	"v": uinput.KeyV,
	"w": uinput.KeyW,
	"x": uinput.KeyX,
	"y": uinput.KeyY,
	"z": uinput.KeyZ,

	"0": uinput.Key0,
	"1": uinput.Key1,
	"2": uinput.Key2,
	"3": uinput.Key3,
	"4": uinput.Key4,
	"5": uinput.Key5,
	"6": uinput.Key6,
	"7": uinput.Key7,
	"8": uinput.Key8,
	"9": uinput.Key9,

	"f1":  uinput.KeyF1,
	"f2":  uinput.KeyF2,
	"f3":  uinput.KeyF3,
	"f4":  uinput.KeyF4,
	"f5":  uinput.KeyF5,
	"f6":  uinput.KeyF6,
	"f7":  uinput.KeyF7,
	"f8":  uinput.KeyF8,
	"f9":  uinput.KeyF9,
	"f10": uinput.KeyF10,
	"f11": uinput.KeyF11,
	"f12": uinput.KeyF12,
	"f13": uinput.KeyF13,
	"f14": uinput.KeyF14,
	"f15": uinput.KeyF15,
	"f16": uinput.KeyF16,
	"f17": uinput.KeyF17,
	"f18": uinput.KeyF18,
	"f19": uinput.KeyF19,
	"f20": uinput.KeyF20,
	"f21": uinput.KeyF21,
	"f22": uinput.KeyF22,
	"f23": uinput.KeyF23,
	"f24": uinput.KeyF24,

	"kp_0":         uinput.KeyKp0,
	"kp_1":         uinput.KeyKp1,
	"kp_2":         uinput.KeyKp2,
	"kp_3":         uinput.KeyKp3,
	"kp_4":         uinput.KeyKp4,
	"kp_5":         uinput.KeyKp5,
	"kp_6":         uinput.KeyKp6,
	"kp_7":         uinput.KeyKp7,
	"kp_8":         uinput.KeyKp8,
	"kp_9":         uinput.KeyKp9,
	"kp_enter":     uinput.KeyKpenter,
	"kp_plus":      uinput.KeyKpplus,
	"kp_minus":     uinput.KeyKpminus,
	"kp_multiply":  uinput.KeyKpasterisk,
	"kp_divide":    uinput.KeyKpslash,
	"kp_decimal":   uinput.KeyKpdot,
	"kp_separator": uinput.KeyKpcomma,
	"kp_equal":     uinput.KeyKpequal,

	"minus":        uinput.KeyMinus,
	"equal":        uinput.KeyEqual,
	"bracketleft":  uinput.KeyLeftbrace,
	"bracketright": uinput.KeyRightbrace,
	"semicolon":    uinput.KeySemicolon,
	"apostrophe":   uinput.KeyApostrophe,
	"grave":        uinput.KeyGrave,
	"backslash":    uinput.KeyBackslash,
	"comma":        uinput.KeyComma,
	"period":       uinput.KeyDot,
	"dot":          uinput.KeyDot,
	"slash":        uinput.KeySlash,
	"underscore":   uinput.KeyMinus,
	"plus":         uinput.KeyEqual,
	"braceleft":    uinput.KeyLeftbrace,
	"braceright":   uinput.KeyRightbrace,
	"colon":        uinput.KeySemicolon,
	"quotedbl":     uinput.KeyApostrophe,
	"asciitilde":   uinput.KeyGrave,
	"bar":          uinput.KeyBackslash,
	"less":         uinput.KeyComma,
	"greater":      uinput.KeyDot,
	"question":     uinput.KeySlash,
	"exclam":       uinput.Key1,
	"at":           uinput.Key2,
	"numbersign":   uinput.Key3,
	"dollar":       uinput.Key4,
	"percent":      uinput.Key5,
	"asciicircum":  uinput.Key6,
	"ampersand":    uinput.Key7,
	"asterisk":     uinput.Key8,
	"parenleft":    uinput.Key9,
	"parenright":   uinput.Key0,

	"xf86audiolowervolume": uinput.KeyVolumedown,
	"xf86audioraisevolume": uinput.KeyVolumeup,
	"xf86audiomute":        uinput.KeyMute,
	"xf86audioplay":        uinput.KeyPlaypause,
	"xf86audiostop":        uinput.KeyStopcd,
	"xf86audioprev":        uinput.KeyPrevioussong,
	"xf86audionext":        uinput.KeyNextsong,
	"xf86audiopause":       uinput.KeyPausecd,
	"xf86audiorecord":      uinput.KeyRecord,
	"xf86audiorewind":      uinput.KeyRewind,
	"xf86audioforward":     uinput.KeyFastforward,
	"xf86audiomicmute":     uinput.KeyMicmute,

	"xf86monbrightnessup":   uinput.KeyBrightnessup,
	"xf86monbrightnessdown": uinput.KeyBrightnessdown,
	"xf86display":           uinput.KeySwitchvideomode,

	"xf86sleep":    uinput.KeySleep,
	"xf86wakeup":   uinput.KeyWakeup,
	"xf86poweroff": uinput.KeyPower,
	"xf86suspend":  uinput.KeySuspend,

	"xf86calculator": uinput.KeyCalc,
	"xf86mail":       uinput.KeyMail,
	"xf86www":        uinput.KeyWww,
	"xf86homepage":   uinput.KeyHomepage,
	"xf86back":       uinput.KeyBack,
	"xf86forward":    uinput.KeyForward,
	"xf86search":     uinput.KeySearch,
	"xf86refresh":    uinput.KeyRefresh,
	"xf86stop":       uinput.KeyStop,
	"xf86favorites":  uinput.KeyBookmarks,
	"xf86mycomputer": uinput.KeyComputer,
	"xf86documents":  uinput.KeyDocuments,
	"xf86explorer":   uinput.KeyFile,

	"xf86kbdbrightnessup":   uinput.KeyKbdillumup,
	"xf86kbdbrightnessdown": uinput.KeyKbdillumdown,

	"xf86eject": uinput.KeyEjectcd,

	"xf86touchpadtoggle": uinput.KeyF21, // Often mapped to F21

	"xf86wlan": uinput.KeyWlan,

	"xf86bluetooth": uinput.KeyBluetooth,

	"xf86tools": uinput.KeyConfig,

	"xf86battery": uinput.KeyBattery,

	"xf86launch0": uinput.KeyProg1,
	"xf86launch1": uinput.KeyProg2,
	"xf86launch2": uinput.KeyProg3,
	"xf86launch3": uinput.KeyProg4,

	"print":    uinput.KeySysrq,
	"sysrq":    uinput.KeySysrq,
	"help":     uinput.KeyHelp,
	"undo":     uinput.KeyUndo,
	"redo":     uinput.KeyRedo,
	"cut":      uinput.KeyCut,
	"copy":     uinput.KeyCopy,
	"paste":    uinput.KeyPaste,
	"find":     uinput.KeyFind,
	"open":     uinput.KeyOpen,
	"close":    uinput.KeyClose,
	"save":     uinput.KeySave,
	"cancel":   uinput.KeyCancel,
	"exit":     uinput.KeyExit,
	"linefeed": uinput.KeyLinefeed,
}

func ParseXDoToolKeybindString(keybind string) ([]int, error) {
	if keybind == "" {
		return nil, errors.New("empty keybind")
	}

	parts := strings.Split(keybind, "+")
	if len(parts) == 0 {
		return nil, errors.New("invalid keybind format")
	}

	keys := make([]int, 0, len(parts))
	for _, part := range parts {
		keyName := strings.ToLower(strings.TrimSpace(part))
		if keyName == "" {
			continue
		}

		keyCode, ok := keyNameToCode[keyName]
		if !ok {
			return nil, errors.New("unknown key: " + part)
		}
		keys = append(keys, keyCode)
	}

	if len(keys) == 0 {
		return nil, errors.New("no valid keys in keybind")
	}

	return keys, nil
}

func FindXDoToolKeybindString(code int) string {
	for keybind, keyCode := range keyNameToCode {
		if code == keyCode {
			return keybind
		}
	}
	return ""
}
