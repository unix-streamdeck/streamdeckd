package handlers

import (
	"golang.org/x/sync/semaphore"
	"image"
)

type CounterIconHandler struct {
	Count    int
	Running  bool
	Callback func(image image.Image)
}

type GifIconHandler struct {
	Running bool
	Lock *semaphore.Weighted
}

type TimeIconHandler struct {
	Running bool
}