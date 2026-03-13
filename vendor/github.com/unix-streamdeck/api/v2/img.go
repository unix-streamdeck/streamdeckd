package api

import (
	"errors"
	"image"
	"image/color"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/unix-streamdeck/gg"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gobolditalic"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/gofont/gomediumitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/gomonobolditalic"
	"golang.org/x/image/font/gofont/gomonoitalic"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/gofont/gosmallcaps"
	"golang.org/x/image/font/gofont/gosmallcapsitalic"
)

const BorderClearance = 10

// TODO replace use of gg with native font.Drawer
type VerticalAlignment string

const (
	Top    VerticalAlignment = "TOP"
	Center VerticalAlignment = "CENTER"
	Bottom VerticalAlignment = "BOTTOM"
)

type DrawTextOptions struct {
	FontSize          int64
	VerticalAlignment VerticalAlignment
	FontFace          string
	Colour            string
}

type IContext interface {
	SetRGB(r, g, b float64)
	SetHexColor(color string)
	SetFontFace(font font.Face)
	WordWrap(text string, width float64) []string
	DrawStringWrapped(s string, x, y, ax, ay, width, lineSpacing float64, align gg.Align)
	Image() image.Image
	MeasureMultilineString(text string, lineSpacing float64) (float64, float64)
	Width() int
	Height() int
}

func DrawText(img image.Image, text string, options DrawTextOptions) (image.Image, error) {
	ggImg := gg.NewContextForImage(img)
	return drawText(ggImg, text, options)
}

func drawText(img IContext, text string, options DrawTextOptions) (image.Image, error) {
	width, height := img.Width(), img.Height()
	img.SetRGB(1, 1, 1)
	matched, _ := regexp.MatchString(`#?([0-9a-fA-F]{8}|[0-9a-fA-F]{6}|[0-9a-fA-F]{3})`, options.Colour)
	if matched {
		img.SetHexColor(options.Colour)
	}
	f, err := truetype.Parse(loadFontFace(options.FontFace))
	if err != nil {
		return nil, err
	}
	fSize := calculateFontSize(f, text, img)

	if options.FontSize != 0 {
		fSize = float64(options.FontSize)
	}

	face := truetype.NewFace(f, &truetype.Options{Size: fSize})
	defer face.Close()
	img.SetFontFace(face)

	lines := img.WordWrap(text, float64(width-BorderClearance))
	lineCount := float64(len(lines))

	if strings.Contains(text, "\n") {
		lineCount += float64(strings.Count(text, "\n") + 1)
	}

	valign, y := calculateVerticalAlignment(options.VerticalAlignment, height)
	img.DrawStringWrapped(text, float64(width/2), y, 0.5, valign, float64(width-BorderClearance), 1, gg.AlignCenter)
	return img.Image(), nil
}

// TODO Support loading fonts via fontconfig on linux and whatever the equivalent is on darwin
func loadFontFace(fontName string) []byte {
	switch fontName {
	case "bold":
		return gobold.TTF
	case "bolditalic":
		return gobolditalic.TTF
	case "italic":
		return goitalic.TTF
	case "medium":
		return gomedium.TTF
	case "mediumitalic":
		return gomediumitalic.TTF
	case "mono":
		return gomono.TTF
	case "monobold":
		return gomonobold.TTF
	case "monobolditalic":
		return gomonobolditalic.TTF
	case "monoitalic":
		return gomonoitalic.TTF
	case "smallcaps":
		return gosmallcaps.TTF
	case "smallcapsitalic":
		return gosmallcapsitalic.TTF
	case "regular":
		fallthrough
	default:
		return goregular.TTF
	}
}

func calculateVerticalAlignment(alignment VerticalAlignment, height int) (float64, float64) {
	if alignment == Top {
		return 0.0, BorderClearance / 2
	}
	if alignment == Bottom {
		return 1.0, float64(height) - (BorderClearance / 2)
	}
	return 0.5, float64(height) / 2
}

func calculateFontSize(f *truetype.Font, text string, img IContext) float64 {
	width := img.Width()
	fontSize := float64(width) / 3
	face := truetype.NewFace(f, &truetype.Options{Size: fontSize})
	defer face.Close()
	img.SetFontFace(face)
	textWidth, _ := img.MeasureMultilineString(text, 1.0)
	fSize := fontSize
	if textWidth >= float64(width-BorderClearance) {
		oversizeRatio := float64(width-BorderClearance) / textWidth
		scaledFontSize := math.Min(oversizeRatio*fontSize, 12)
		for size := fontSize; size >= scaledFontSize; size -= 0.5 {
			if attemptFontSize(f, text, img, size) {
				return size
			}
		}
		return scaledFontSize
	}
	return fSize
}

func attemptFontSize(f *truetype.Font, text string, img IContext, fSize float64) bool {
	width := img.Width()
	height := img.Height()
	face := truetype.NewFace(f, &truetype.Options{Size: fSize})
	defer face.Close()
	img.SetFontFace(face)
	wrappedGroups := img.WordWrap(text, float64(width-BorderClearance))
	wrappedText := strings.Join(wrappedGroups, "\n")
	textWidth, textHeight := img.MeasureMultilineString(wrappedText, 1.0)
	return textHeight < float64(height-BorderClearance) && textWidth < float64(width-BorderClearance)
}

func ResizeImage(img image.Image, keySize int) image.Image {
	return ResizeImageWH(img, keySize, keySize)
}

func ResizeImageWH(img image.Image, width int, height int) image.Image {

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.BiLinear.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

	return dst
}

func DrawProgressBar(img image.Image, label string, x, y, h, w, progress float64) (image.Image, error) {
	return DrawProgressBarWithAccent(img, label, x, y, h, w, progress, "#777777")
}

func DrawProgressBarWithAccent(img image.Image, label string, x, y, h, w, progress float64, hex string) (image.Image, error) {
	ggImg := gg.NewContextForImage(img)

	f, err := truetype.Parse(goregular.TTF)

	if err != nil {
		return nil, err
	}

	face := truetype.NewFace(f, &truetype.Options{Size: h / 2})
	defer face.Close()

	ggImg.SetFillRule(gg.FillRuleEvenOdd)

	ggImg.SetFillStyle(gg.NewSolidPattern(HexColor("#333333")))

	ggImg.DrawRoundedRectangle(x, y, w, h, 5)

	ggImg.Fill()

	ggImg.SetFillStyle(gg.NewSolidPattern(HexColor(hex)))

	ggImg.DrawRoundedRectangle(x, y, w/100*progress, h, 5)

	ggImg.Fill()

	ggImg.SetHexColor("#FFFFFF")

	ggImg.DrawStringAnchored(label, (x+w)/2, y+(h/2), 0.5, 0.5)

	return ggImg.Image(), nil
}

func HexColor(hex string) color.RGBA {
	values, _ := strconv.ParseUint(hex[1:], 16, 32)
	return color.RGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}

func LayerImages(background, foreground image.Image) (image.Image, error) {
	if background.Bounds().Size() != foreground.Bounds().Size() {
		return nil, errors.New("images must be same size")
	}

	ggBG := gg.NewContextForImage(background)

	ggBG.DrawImage(foreground, background.Bounds().Min.X, background.Bounds().Min.Y)

	return ggBG.Image(), nil
}
