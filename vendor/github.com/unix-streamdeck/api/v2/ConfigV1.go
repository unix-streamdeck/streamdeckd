package api

import "image"

type ConfigV1 struct {
    Modules []string `json:"modules,omitempty"`
    Pages   []PageV1 `json:"pages"`
}

type PageV1 []KeyV1

type KeyV1 struct {
    Icon              string            `json:"icon,omitempty"`
    SwitchPage        int               `json:"switch_page,omitempty"`
    Text              string            `json:"text,omitempty"`
    TextSize          int               `json:"text_size,omitempty"`
    TextAlignment     string            `json:"text_alignment,omitempty"`
    Keybind           string            `json:"keybind,omitempty"`
    Command           string            `json:"command,omitempty"`
    Brightness        int               `json:"brightness,omitempty"`
    Url               string            `json:"url,omitempty"`
    ObsCommand        string            `json:"obs_command,omitempty"`
    ObsCommandParams  map[string]string `json:"obs_command_params,omitempty"`
    IconHandler       string            `json:"icon_handler,omitempty"`
    KeyHandler        string            `json:"key_handler,omitempty"`
    IconHandlerFields map[string]any    `json:"icon_handler_fields,omitempty"`
    KeyHandlerFields  map[string]any    `json:"key_handler_fields,omitempty"`
    Buff              image.Image       `json:"-"`
    IconHandlerStruct IconHandler       `json:"-"`
    KeyHandlerStruct  KeyHandler        `json:"-"`
}
