package gif

import (
	"context"
	"image"
	"image/draw"
	"image/gif"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/unix-streamdeck/api/v2"
	"golang.org/x/sync/semaphore"
)

type GifIconHandler struct {
	Running bool
	Lock    *semaphore.Weighted
	Quit    chan bool
	Gifs    []*image.Paletted
}

func (s *GifIconHandler) Start(key api.KeyConfigV3, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	if s.Quit == nil {
		s.Quit = make(chan bool)
	}
	if s.Lock == nil {
		s.Lock = semaphore.NewWeighted(1)
	}
	s.Running = true
	icon, ok := key.IconHandlerFields["icon"]
	if !ok {
		return
	}
	f, err := os.Open(icon.(string))
	if err != nil {
		log.Println(err)
		return
	}
	gifs, err := gif.DecodeAll(f)
	s.Gifs = gifs.Image
	if err != nil {
		log.Println(err)
		return
	}
	timeDelay := gifs.Delay[0]
	if timeDelay < 1 {
		timeDelay = 8
	}
	frames := make([]image.Image, len(gifs.Image))

	iconSize := info.IconSize

	overPaintImage := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))
	draw.Draw(overPaintImage, overPaintImage.Bounds(), api.ResizeImage(gifs.Image[0], iconSize), image.ZP, draw.Src)

	for i, frame := range gifs.Image {
		draw.Draw(overPaintImage, overPaintImage.Bounds(), api.ResizeImage(frame, iconSize), image.ZP, draw.Over)
		frame := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))
		draw.Draw(frame, frame.Bounds(), overPaintImage, image.ZP, draw.Over)
		img := frame.SubImage(frame.Rect)
		text, ok := key.IconHandlerFields["text"]
		if ok {
			text_size, ok := key.IconHandlerFields["text_size"]
			var size int64
			if ok {
				size, _ = strconv.ParseInt(text_size.(string), 10, 0)
			} else {
				size = 0
			}
			alignment, ok := key.IconHandlerFields["text_alignment"]
			if !ok {
				alignment = ""
			}
			fontFace, ok := key.IconHandlerFields["font_face"]
			if !ok {
				fontFace = ""
			}
			textColour, ok := key.IconHandlerFields["text_colour"]
			if !ok {
				textColour = ""
			}
			img, err = api.DrawText(img, text.(string), api.DrawTextOptions{
				FontSize:          size,
				VerticalAlignment: api.VerticalAlignment(alignment.(string)),
				FontFace:          fontFace.(string),
				Colour:            textColour.(string),
			})
			if err != nil {
				log.Println(err)
			}
		}
		frames[i] = img
	}
	go s.loop(frames, timeDelay, callback)
}

func (s *GifIconHandler) IsRunning() bool {
	return s.Running
}

func (s *GifIconHandler) SetRunning(running bool) {
	s.Running = running
}

func (s *GifIconHandler) Stop() {
	s.Running = false
	s.Quit <- true
}

func (s *GifIconHandler) loop(frames []image.Image, timeDelay int, callback func(image image.Image)) {

	ctx := context.Background()
	err := s.Lock.Acquire(ctx, 1)
	if err != nil {
		return
	}
	defer s.Lock.Release(1)
	gifIndex := 0
	for {
		select {
		case <-s.Quit:
			return
		default:
			img := frames[gifIndex]
			callback(img)
			gifIndex++
			if gifIndex >= len(frames) {
				gifIndex = 0
			}
			time.Sleep(time.Duration(timeDelay*10) * time.Millisecond)
		}
	}
}
