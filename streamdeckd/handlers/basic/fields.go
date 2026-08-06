package basic

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/unix-streamdeck/api/v3"
)

func fieldString(fields map[string]any, name string) (string, bool) {
	value, ok := fields[name]
	if !ok {
		return "", false
	}

	str, ok := value.(string)
	if !ok {
		log.Printf("[WARN] Field `%s` should be a string, got %T\n", name, value)
		return "", false
	}

	return str, str != ""
}

func fieldInt(fields map[string]any, name string) (int, bool) {
	value, ok := fields[name]
	if !ok {
		return 0, false
	}

	switch typed := value.(type) {
	case int:
		return typed, true
	case int32:
		return int(typed), true
	case int64:
		return int(typed), true
	case float32:
		return int(typed), true
	case float64:
		return int(typed), true
	case json.Number:
		parsed, err := typed.Int64()
		if err != nil {
			log.Printf("[WARN] Field `%s` is not a whole number: %s\n", name, err)
			return 0, false
		}
		return int(parsed), true
	case string:
		parsed, err := strconv.Atoi(typed)
		if err != nil {
			log.Printf("[WARN] Field `%s` is not a whole number: %s\n", name, err)
			return 0, false
		}
		return parsed, true
	}

	log.Printf("[WARN] Field `%s` should be a number, got %T\n", name, value)
	return 0, false
}

func shouldAct(event api.InputEvent) bool {
	return event.EventType != api.KEY_RELEASE
}
