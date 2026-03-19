package gif

import (
	"context"
	"image"
	"image/draw"
	"image/gif"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/unix-streamdeck/api/v2"
	"golang.org/x/sync/semaphore"
)

type GifKeyGridBackgroundHandler struct {
	Running bool
	Lock    *semaphore.Weighted
	Quit    chan bool
	Gifs    []*image.Paletted
}

func (s *GifKeyGridBackgroundHandler) StartIndividual(fields map[string]any, info api.StreamDeckInfoV1, callback func(img image.Image)) {
	log.Println("Not Implemented")
	return
}

func (s *GifKeyGridBackgroundHandler) Start(fields map[string]any, info api.StreamDeckInfoV1, callback func(imgs []image.Image)) {
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

	w := (info.IconSize * info.Cols) + (info.PaddingX * (info.Cols - 1))
	h := (info.IconSize * info.Rows) + (info.PaddingX * (info.Rows - 1))

	frames := make([][]image.Image, len(gifs.Image))

	overPaintImage := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(overPaintImage, overPaintImage.Bounds(), api.ResizeImageWH(gifs.Image[0], w, h), image.ZP, draw.Src)

	for i, frame := range gifs.Image {
		draw.Draw(overPaintImage, overPaintImage.Bounds(), api.ResizeImageWH(frame, w, h), image.ZP, draw.Over)
		frame := image.NewRGBA(image.Rect(0, 0, w, h))
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
		frameArr := make([]image.Image, info.LcdCols)
		for keyIndex := range info.Cols * info.Rows {
			keyX := keyIndex % info.Cols
			keyY := int(math.Floor(float64(keyIndex) / float64(info.Cols)))

			x0, y0 := keyX*(info.IconSize+info.PaddingX), keyY*(info.IconSize+info.PaddingY)
			x1, y1 := keyX*(info.IconSize+info.PaddingX)+info.IconSize, keyY*(info.IconSize+info.PaddingY)+info.IconSize

			frameArr = append(frameArr, api.SubImage(img, x0, y0, x1, y1))
		}
		frames[i] = frameArr
	}
	go s.loop(frames, timeDelay, callback)
}

func (s *GifKeyGridBackgroundHandler) IsRunning() bool {
	return s.Running
}

func (s *GifKeyGridBackgroundHandler) SetRunning(running bool) {
	s.Running = running
}

func (s *GifKeyGridBackgroundHandler) Stop() {
	s.Running = false
	s.Quit <- true
}

func (s *GifKeyGridBackgroundHandler) loop(frames [][]image.Image, timeDelay int, callback func(image []image.Image)) {
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
