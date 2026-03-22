package examples

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

func RegisterGif() api.Module {
	return api.Module{
		Name: "Gif",
		NewForeground: func() api.ForegroundHandler {
			return &GifHandler{Running: true, Lock: semaphore.NewWeighted(1)}
		},
		ForegroundFields: []api.Field{
			{Title: "Icon", Name: "icon", Type: api.File, FileTypes: []string{".gif"}},
			{Title: "Text", Name: "text", Type: api.Text},
			{Title: "Text Size", Name: "text_size", Type: api.Number},
			{Title: "Text Alignment", Name: "text_alignment", Type: api.TextAlignment},
			{Title: "Font Face", Name: "font_face", Type: api.FontFace},
			{Title: "Text Colour", Name: "text_colour", Type: api.Colour},
		},
		NewBackground: func() api.BackgroundHandler {
			return &GifHandler{Running: true, Lock: semaphore.NewWeighted(1)}
		},
		BackgroundFields: []api.Field{
			{Title: "Icon", Name: "icon", Type: api.File, FileTypes: []string{".gif"}},
			{Title: "Text", Name: "text", Type: api.Text},
			{Title: "Text Size", Name: "text_size", Type: api.Number},
			{Title: "Text Alignment", Name: "text_alignment", Type: api.TextAlignment},
			{Title: "Font Face", Name: "font_face", Type: api.FontFace},
			{Title: "Text Colour", Name: "text_colour", Type: api.Colour},
		},
	}
}

type GifHandler struct {
	Running bool
	Lock    *semaphore.Weighted
	Quit    chan bool
	Gifs    []*image.Paletted
}

func (s *GifHandler) Start(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	if s.Quit == nil {
		s.Quit = make(chan bool)
	}
	if s.Lock == nil {
		s.Lock = semaphore.NewWeighted(1)
	}
	s.Running = true
	icon, ok := fields["icon"]
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
	width, height := info.GetDimensions(handlerType)

	frames := getFrames(fields, gifs, width, height, err)
	go loop(s, frames, timeDelay, callback)
}

func (s *GifHandler) StartGrid(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1, callback func(image []image.Image)) {
	if s.Quit == nil {
		s.Quit = make(chan bool)
	}
	if s.Lock == nil {
		s.Lock = semaphore.NewWeighted(1)
	}
	s.Running = true
	icon, ok := fields["icon"]
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
	width, height := info.GetGridDimensions(handlerType)

	frames := getFrames(fields, gifs, width, height, err)

	var gridFrames [][]image.Image

	for _, frame := range frames {

		gridFrame := info.SplitBackgroundImage(frame, handlerType)

		gridFrames = append(gridFrames, gridFrame)
	}

	go loop(s, gridFrames, timeDelay, callback)
}

func getFrames(fields map[string]any, gifs *gif.GIF, width int, height int, err error) []image.Image {
	frames := make([]image.Image, len(gifs.Image))

	overPaintImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(overPaintImage, overPaintImage.Bounds(), api.ResizeImageWH(gifs.Image[0], width, height), image.ZP, draw.Src)

	for i, frame := range gifs.Image {
		draw.Draw(overPaintImage, overPaintImage.Bounds(), api.ResizeImageWH(frame, width, height), image.ZP, draw.Over)
		frame := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.Draw(frame, frame.Bounds(), overPaintImage, image.ZP, draw.Over)
		img := frame.SubImage(frame.Rect)
		text, ok := fields["text"]
		if ok {
			text_size, ok := fields["text_size"]
			var size int64
			if ok {
				size, _ = strconv.ParseInt(text_size.(string), 10, 0)
			} else {
				size = 0
			}
			alignment, ok := fields["text_alignment"]
			if !ok {
				alignment = ""
			}
			fontFace, ok := fields["font_face"]
			if !ok {
				fontFace = ""
			}
			textColour, ok := fields["text_colour"]
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
	return frames
}

func (s *GifHandler) IsRunning() bool {
	return s.Running
}

func (s *GifHandler) SetRunning(running bool) {
	s.Running = running
}

func (s *GifHandler) Stop() {
	s.Running = false
	s.Quit <- true
}

func loop[T any](s *GifHandler, frames []T, timeDelay int, callback func(image T)) {

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
