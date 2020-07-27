package main

import "image"

type Page []Key

type Config struct {
	Pages []Page `json:"pages"`
}

type Key struct {
	Icon       string `json:"icon,omitempty"`
	SwitchPage *int   `json:"switch_page,omitempty"`
	Text       string `json:"text,omitempty"`
	Keybind    string `json:"keybind,omitempty"`
	Command    string `json:"command,omitempty"`
	Brightness *int   `json:"brightness,omitempty"`
	Url        string `json:"url,omitempty"`
	buff       image.Image
}
