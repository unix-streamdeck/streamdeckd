package main

type Pages [][] struct {
	Icon        string `json:"icon,omitempty"`
	SwitchPage  *int    `json:"switch_page,omitempty"`
	Text        string `json:"text,omitempty"`
	Keybind     string `json:"keybind,omitempty"`
	Command     string `json:"command,omitempty"`
	Brightness  *int    `json:"brightness,omitempty"`
}

type Config struct {
	Pages Pages `json:"pages"`
}


type Key struct {
	Icon        string `json:"icon,omitempty"`
	SwitchPage  *int    `json:"switch_page,omitempty"`
	Text        string `json:"text,omitempty"`
	Keybind     string `json:"keybind,omitempty"`
	Command     string `json:"command,omitempty"`
	Brightness  *int    `json:"brightness,omitempty"`
}