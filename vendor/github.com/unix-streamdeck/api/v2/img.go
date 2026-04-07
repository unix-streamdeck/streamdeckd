package api

import (
	"errors"
	"image"
	"image/color"
	"log"
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
	"golang.org/x/image/math/fixed"
)

const BorderClearance = 10

type VerticalAlignment string

const (
	Top    VerticalAlignment = "TOP"
	Center VerticalAlignment = "CENTER"
	Bottom VerticalAlignment = "BOTTOM"
)

type HorizontalAlignment string

const (
	Left   HorizontalAlignment = "LEFT"
	Middle HorizontalAlignment = "MIDDLE"
	Right  HorizontalAlignment = "RIGHT"
)

type Overflow string

const (
	Wrap Overflow = "WRAP"
	Fade Overflow = "FADE"
)

type DrawTextOptions struct {
	FontSize            int64
	VerticalAlignment   VerticalAlignment
	HorizontalAlignment HorizontalAlignment
	FontFace            string
	Colour              string
	Overflow            Overflow
	// With anchor set, the alignments indicate the position of the text the anchor will be in, e.g BOTTOM & RIGHT will put the anchor in the bottom right of the text, so all text will be up and left of it
	Anchor *image.Point
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

	if options.Overflow == "" {
		options.Overflow = Wrap
	}

	drawImg, ok := img.(draw.Image)

	if !ok {
		return img, errors.New("cannot convert")
	}

	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	availableWidth, availableHeight := width, height

	if options.Anchor != nil {
		availableWidth, availableHeight = width-options.Anchor.X, height-options.Anchor.Y

		if options.VerticalAlignment == Center {
			availableHeight = height - options.Anchor.Y
		}

		if options.HorizontalAlignment == Middle {
			availableWidth = width - options.Anchor.X
		}
	}

	col := color.RGBA{0xff, 0xff, 0xff, 0xff}

	//img.SetRGB(1, 1, 1)
	matched, _ := regexp.MatchString(`#?([0-9a-fA-F]{8}|[0-9a-fA-F]{6}|[0-9a-fA-F]{3})`, options.Colour)
	if matched {
		col = HexColor(options.Colour)
	}
	f, err := truetype.Parse(loadFontFace(options.FontFace))
	if err != nil {
		return nil, err
	}
	fSize := calculateFontSize(f, text, availableWidth, availableHeight, options.Overflow)

	if options.FontSize != 0 {
		fSize = float64(options.FontSize)
	}

	face := truetype.NewFace(f, &truetype.Options{Size: fSize})
	defer face.Close()

	lines := text

	if options.Overflow == Wrap {
		lines = wrapString(text, availableWidth, face)
	}

	lineCount := strings.Count(lines, "\n") + 1

	var y float64

	var x int

	w, _ := getTextBounds(lines, face)

	if options.Anchor != nil {
		x = calculateHorizonalAlignment(options.HorizontalAlignment, int(w), options.Anchor.X, true)
		y = calculateVerticalAlignment(options.VerticalAlignment, options.Anchor.Y, lineCount, fSize, true)
	} else {
		x = calculateHorizonalAlignment(options.HorizontalAlignment, int(w), width, false)
		y = calculateVerticalAlignment(options.VerticalAlignment, height, lineCount, fSize, false)
	}

	d := &font.Drawer{
		Dst:  drawImg,
		Src:  image.NewUniform(col),
		Face: face,
	}

	d.Dot = fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(int(y) + (int(fSize) / 4)),
	}

	linesSplit := strings.Split(lines, "\n")

	if len(linesSplit) == 1 {
		d.DrawString(lines)
		return img, nil
	}

	linesAbove := float64(lineCount) / 2

	linesAbove = linesAbove - 1

	startingLineY := y - (linesAbove * fSize)

	for i, line := range linesSplit {
		w, _ := getTextBounds(line, face)

		if options.Anchor != nil {
			x = calculateHorizonalAlignment(options.HorizontalAlignment, int(w), options.Anchor.X, true)
		} else {
			x = calculateHorizonalAlignment(options.HorizontalAlignment, int(w), width, false)
		}

		d.Dot = fixed.Point26_6{
			X: fixed.I(x),
			Y: fixed.I(int(startingLineY + (float64(i) * fSize))),
		}
		d.DrawString(line)
	}
	return img, nil
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

func calculateVerticalAlignment(alignment VerticalAlignment, height, lines int, fSize float64, anchor bool) float64 {
	textMidPoint := (float64(lines) / 2.0) * fSize
	if !anchor {
		if alignment == Top {
			return (BorderClearance / 2) + (textMidPoint)
		}
		if alignment == Bottom {
			return float64(height) - (BorderClearance / 2) - textMidPoint
		}
		return float64(height) / 2
	}

	if alignment == Top {
		return float64(height) + textMidPoint
	}
	if alignment == Bottom {
		return float64(height) - textMidPoint
	}
	return float64(height)
}

func calculateHorizonalAlignment(alignment HorizontalAlignment, textWidth, width int, anchor bool) int {
	if !anchor {
		if alignment == Left {
			return BorderClearance / 2
		}
		if alignment == Right {
			return width - (BorderClearance / 2) - textWidth
		}
		return ((width) / 2) - (int(textWidth) / 2)
	}

	if alignment == Left {
		return width
	}

	if alignment == Right {
		return width - textWidth
	}

	return width - (textWidth / 2)
}

func calculateFontSize(f *truetype.Font, text string, width, height int, overflow Overflow) float64 {
	fontSize := float64(width) / 3
	face := truetype.NewFace(f, &truetype.Options{Size: fontSize})
	defer face.Close()
	w, _ := getTextBounds(text, face)
	fSize := fontSize
	if w >= float64(width-BorderClearance) {
		oversizeRatio := float64(width-BorderClearance) / w
		scaledFontSize := math.Min(oversizeRatio*fontSize, 12)
		for size := fontSize; size >= scaledFontSize; size -= 0.5 {
			if attemptFontSize(f, text, width, height, size, overflow) {
				return size
			}
		}
		return scaledFontSize
	}
	return fSize
}

func attemptFontSize(f *truetype.Font, text string, width, height int, fSize float64, overflow Overflow) bool {
	face := truetype.NewFace(f, &truetype.Options{Size: fSize})
	defer face.Close()
	w, h := getTextBounds(text, face)
	if w <= float64(width-BorderClearance) {
		return true
	}
	if h > float64(height) {
		return false
	}
	if overflow != Wrap {
		return false
	}
	lines := wrapString(text, width, face)
	if lines == "" {
		return false
	}
	maxTextWidth := 0.0
	for _, s := range strings.Split(lines, "\n") {
		textWidth, _ := getTextBounds(s, face)
		if textWidth > maxTextWidth {
			maxTextWidth = textWidth
		}
	}
	textHeight := float64(strings.Count(lines, "\n")+1) * fSize
	if textHeight < float64(height-BorderClearance) && maxTextWidth < float64(width-BorderClearance) {
		return true
	}
	return false
}

func wrapString(text string, width int, face font.Face) string {
	splitMessage := strings.Split(text, " ")
	if len(splitMessage) == 1 {
		return text
	}

	var lines []string
	nextWordIndex := 0
	for nextWordIndex < len(splitMessage) {
		lineLength := 0.0
		var line []string
		for lineLength <= float64(width-BorderClearance) && nextWordIndex < len(splitMessage) {
			w, _ := getTextBounds(splitMessage[nextWordIndex], face)
			if w > float64(width-BorderClearance) {
				return ""
			}
			if w+lineLength > float64(width-BorderClearance) {
				break
			}
			lineLength += w
			line = append(line, splitMessage[nextWordIndex])
			nextWordIndex += 1
		}
		lines = append(lines, strings.Join(line, " "))
	}
	return strings.Join(lines, "\n")
}

func getTextBounds(text string, face font.Face) (float64, float64) {
	bounds, _ := font.BoundString(face, text)

	return (float64(bounds.Max.X.Round()) - float64(bounds.Min.X.Round())), (float64(bounds.Max.Y.Round()) - float64(bounds.Min.Y.Round()))
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
	var err error

	ggImg := gg.NewContextForImage(img)

	ggImg.SetFillRule(gg.FillRuleEvenOdd)

	ggImg.SetFillStyle(gg.NewSolidPattern(HexColor("#333333")))

	ggImg.DrawRoundedRectangle(x, y, w, h, 5)

	ggImg.Fill()

	ggImg.SetFillStyle(gg.NewSolidPattern(HexColor(hex)))

	ggImg.DrawRoundedRectangle(x, y, w/100*progress, h, 5)

	ggImg.Fill()

	ggImg.SetHexColor("#FFFFFF")

	img = ggImg.Image()

	if label != "" {
		img, err = DrawText(img, label, DrawTextOptions{
			Anchor: &image.Point{
				X: int(x) + int(w/2),
				Y: int(y) + int((h/2)*1.2),
			},
			FontSize: int64(math.Max(h/2, h-5)),
		})
	}

	return img, err
}

func HexColor(hex string) color.RGBA {
	values, _ := strconv.ParseUint(hex[1:], 16, 32)
	return color.RGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}

func LayerImages(x, y int, images ...image.Image) (image.Image, error) {

	if len(images) == 0 {
		return nil, errors.New("no images supplied")
	}

	layers := 0

	dst := image.NewRGBA(image.Rect(0, 0, x, y))

	for _, img := range images {
		if img == nil {
			continue
		}
		if img.Bounds().Dx() != x || img.Bounds().Dy() != y {
			continue
		}
		layers += 1
		draw.Copy(dst, dst.Bounds().Min, img, img.Bounds(), draw.Over, &draw.Options{})
	}

	if layers == 0 {
		return nil, errors.New("no valid images supplied")
	}

	return dst, nil
}

func SubImage(img image.Image, x0, y0, x1, y1 int) image.Image {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	simg, ok := img.(subImager)

	if !ok {
		log.Println("Couldn't resize")
		return nil
	}

	rect := image.Rect(x0, y0, x1, y1)

	img = simg.SubImage(rect)

	return img
}
