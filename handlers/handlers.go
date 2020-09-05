package handlers

import "image"

type CounterIconHandler struct {
	Count    int
	Running  bool
	Callback func(image image.Image)
}

type GifIconHandler struct {
	Running bool
}

type TimeIconHandler struct {
	Running bool
}
type SpotifyIconHandler struct {
	Running bool
	oldUrl string
}