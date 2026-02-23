package api

import "time"

type StreamDeckInfoV1 struct {
    Cols             int       `json:"cols,omitempty"`
    Rows             int       `json:"rows,omitempty"`
    IconSize         int       `json:"icon_size,omitempty"`
    Page             int       `json:"page"`
    Serial           string    `json:"serial,omitempty"`
    Name             string    `json:"name,omitempty"`
    Connected        bool      `json:"connected"`
    LastConnected    time.Time `json:"last_connected,omitempty"`
    LastDisconnected time.Time `json:"last_disconnected,omitempty"`
    LcdWidth         int       `json:"lcd_width,omitempty"`
    LcdHeight        int       `json:"lcd_height,omitempty"`
    LcdCols          int       `json:"lcd_cols,omitempty"`
    KnobCols         int       `json:"knob_cols,omitempty"`
}

type DeckV2 struct {
    Serial string   `json:"serial"`
    Pages  []PageV1 `json:"pages"`
}

type ConfigV2 struct {
    Modules           []string            `json:"modules,omitempty"`
    Decks             []DeckV2            `json:"decks"`
    ObsConnectionInfo ObsConnectionInfoV2 `json:"obs_connection_info,omitempty"`
}

type ObsConnectionInfoV2 struct {
    Host string `json:"host,omitempty"`
    Port int    `json:"port,omitempty"`
}
