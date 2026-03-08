package api

import (
	"errors"
	"strings"
)

func ParseXDoToolKeybindString(keybind string) ([]int, error) {
	if keybind == "" {
		return nil, errors.New("empty keybind")
	}

	parts := strings.Split(keybind, "+")

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
