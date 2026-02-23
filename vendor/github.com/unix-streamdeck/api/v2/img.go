package api

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"log"
	"strings"
)

func DrawText(currentImage image.Image, text string, fontSize int, fontAlignment string) (image.Image, error) {
	width, height := currentImage.Bounds().Max.X, currentImage.Bounds().Max.Y
	img := gg.NewContextForImage(currentImage)
	img.SetRGB(1, 1, 1)
	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	fSize := calculateFontSize(f, text, img)
	if fontSize != 0 {
		fSize = float64(fontSize)
	}
	face := truetype.NewFace(f, &truetype.Options{Size: fSize})
	img.SetFontFace(face)

	lineCount := 1.0

	if strings.Contains(text, "\n") {
		lineCount = float64(strings.Count(text, "\n") + 1)
	} else {
		lines := img.WordWrap(text, float64(width-10))
		lineCount = float64(len(lines))
	}

	verticalAlignment := 0.5
	y := float64(height-5) / 2
	if strings.ToUpper(fontAlignment) == "TOP" {
		verticalAlignment = 1.0
		y = (fSize/2)*lineCount + 10*lineCount
	} else if strings.ToUpper(fontAlignment) == "BOTTOM" {
		verticalAlignment = 0.0
		y = float64(height-5) - (fSize * lineCount)
	}
	img.DrawStringWrapped(text, float64(width-5)/2, y, 0.5, verticalAlignment, float64(width-10), 1, gg.AlignCenter)
	return img.Image(), nil
}

func calculateFontSize(f *truetype.Font, text string, img *gg.Context) float64 {
	width, height := img.Image().Bounds().Max.X, img.Image().Bounds().Max.Y
	fontSize := float64(img.Image().Bounds().Max.Y) / 3
	face := truetype.NewFace(f, &truetype.Options{Size: fontSize})
	img.SetFontFace(face)
	textWidth, _ := img.MeasureMultilineString(text, 1.0)
	fSize := fontSize
	if textWidth >= float64(width-10) {
		s := (float64(width-10) / textWidth) * fontSize
		if s > 12 || !strings.Contains(text, " ") {
			fSize = s
		} else {
			words := img.WordWrap(text, float64(width-10))
			t := ""
			for i := 0; i < strings.Count(text, "\n"); i++ {
				t += "\n"
			}
			if strings.Contains(text, "\n") {
				log.Println(t)
			}
			for i, word := range words {
				t += word
				if i < len(words)-1 {
					t += "\n"
				}
			}
			textWidth, textHeight := img.MeasureMultilineString(t, 1.0)
			if textHeight > textWidth && textHeight > float64(height-10) {
				fSize = (float64(height-10) / textHeight) * fontSize

			} else if textWidth > textHeight && textWidth > float64(width-10) {
				fSize = (float64(height-10) / textWidth) * fontSize
			}
		}
	}
	return fSize
}

func ResizeImage(img image.Image, keySize int) image.Image {
	return resize.Resize(uint(keySize), uint(keySize), img, resize.Lanczos3)
}

func ResizeImageWH(img image.Image, width int, height int) image.Image {
	return resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
}
