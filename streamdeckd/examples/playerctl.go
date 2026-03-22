package examples

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"maps"
	"math"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/godbus/dbus/v5"
	"github.com/unix-streamdeck/api/v2"
	"golang.org/x/sync/semaphore"

	"github.com/Endg4meZer0/go-mpris"
)

type KeypressOperation string

const (
	PlayPause  KeypressOperation = "PlayPause"
	Play       KeypressOperation = "Play"
	Pause      KeypressOperation = "Pause"
	Previous   KeypressOperation = "Previous"
	Next       KeypressOperation = "Next"
	Shuffle    KeypressOperation = "Shuffle"
	LoopStatus KeypressOperation = "LoopStatus"
)

var operationsMap = map[string]KeypressOperation{
	"PlayPause":  PlayPause,
	"Play":       Play,
	"Pause":      Pause,
	"Previous":   Previous,
	"Next":       Next,
	"Shuffle":    Shuffle,
	"LoopStatus": LoopStatus,
}

type PlayerctlHandlerType string

const (
	Playback PlayerctlHandlerType = "Playback"
	Volume   PlayerctlHandlerType = "Volume"
)

var playerFilters = map[PlayerctlHandlerType]func(player *mpris.Player) bool{
	Playback: func(player *mpris.Player) bool {
		return true
	},
	Volume: func(player *mpris.Player) bool {
		_, err := player.GetVolume()
		return err == nil
	},
}

var calculateTextAndPercentage = map[PlayerctlHandlerType]func(player *mpris.Player) (string, float64){
	Playback: func(player *mpris.Player) (string, float64) {
		position, err := player.GetPosition()
		if err != nil {
			return "", 0.0
		}
		meta, err := player.GetMetadata()
		if err != nil {
			return "", 0.0
		}
		length, err := meta.Length()
		if err != nil {
			return "", 0.0
		}
		percentage := math.Round(float64(position)*1_000_000) / math.Round(float64(length)*1_000_000) * 100.0
		text := formatDuration(position) + "/" + formatDuration(length)
		return text, percentage
	},
	Volume: func(player *mpris.Player) (string, float64) {
		vol, err := player.GetVolume()
		if err != nil {
			return "", 0.0
		}
		volume := math.Round(vol * 100.0)
		text := strconv.Itoa(int(volume)) + "%"
		return text, volume
	},
}

type PlayerCtlHandler struct {
	Running                  bool
	Quit                     chan bool
	Lock                     *semaphore.Weighted
	Client                   *dbus.Conn
	AccentColour             string
	CurrentPlayerImage       image.Image
	StaticImage              bool
	PlayerName               string
	CurrentPlayerImageSource string
	Percentage               float64
	Text                     string
	FinalImage               image.Image
	PreviousPlayer           string
	Type                     PlayerctlHandlerType
	ActivePlayer             *mpris.Player
}

func (v *PlayerCtlHandler) Start(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	if v.Quit == nil {
		v.Quit = make(chan bool)
	}
	if v.Lock == nil {
		v.Lock = semaphore.NewWeighted(1)
	}
	if v.CurrentPlayerImage == nil {
		v.CurrentPlayerImage = v.GetImage("icon", fields, info, handlerType)
		if v.CurrentPlayerImage != nil {
			v.StaticImage = true
		}
	}

	accentColour, ok := fields["colour"]

	if ok {
		v.AccentColour = accentColour.(string)
	}

	playerctlHandlerType, ok := fields["type"]

	if !ok {
		log.Println("Type not specified")
		return
	}

	v.Type = PlayerctlHandlerType(playerctlHandlerType.(string))

	playerName, ok := fields["player_name"]
	if ok {
		v.PlayerName = playerName.(string)
	}
	v.Running = true
	v.Run(info, handlerType, callback)
}

func (v *PlayerCtlHandler) IsRunning() bool {
	return v.Running
}

func (v *PlayerCtlHandler) SetRunning(running bool) {
	v.Running = running
}

func (v *PlayerCtlHandler) Stop() {
	v.Running = false
	v.Quit <- true
	v.CurrentPlayerImage = nil
	v.CurrentPlayerImageSource = ""
	v.PlayerName = ""
	v.CurrentPlayerImageSource = ""
	v.Percentage = 0
	v.Text = ""
	v.FinalImage = nil
	v.PreviousPlayer = ""
	v.ActivePlayer = nil
}

func (v *PlayerCtlHandler) GetImage(index string, fields map[string]any, info api.StreamDeckInfoV1, handlerType api.HandlerType) image.Image {
	path, ok := fields[index]
	if !ok {
		return nil
	}
	f, err := os.Open(path.(string))
	defer f.Close()
	if err != nil {
		log.Println(err)
		return nil
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Println(err)
		return nil
	}
	return resizeThumbnail(img, info, handlerType)
}

func (v *PlayerCtlHandler) Run(info api.StreamDeckInfoV1, handlerType api.HandlerType, callback func(image image.Image)) {
	ctx := context.Background()
	err := v.Lock.Acquire(ctx, 1)
	defer v.Lock.Release(1)
	if err != nil {
		return
	}
	for {
		select {
		case <-v.Quit:
			return
		default:
			if playerNeedsRefreshing(v.ActivePlayer) {
				v.ActivePlayer = choosePlayer(v.Client, v.PlayerName, v.PreviousPlayer, playerFilters[v.Type])
				if v.ActivePlayer == nil {
					break
				}
				v.PreviousPlayer = v.ActivePlayer.GetShortName()
			}
			var img image.Image
			img = v.CurrentPlayerImage
			previousImage := v.CurrentPlayerImage
			if !v.StaticImage {
				img, err = v.FindImage(v.ActivePlayer, info, handlerType)
				if img == nil {
					if v.CurrentPlayerImageSource == v.PreviousPlayer {
						img = v.CurrentPlayerImage
					} else {
						w, h := info.GetDimensions(handlerType)
						img = image.NewNRGBA(image.Rect(0, 0, w, h))
						img = resizeThumbnail(img, info, handlerType)
						img, err = api.DrawText(img, v.PreviousPlayer, api.DrawTextOptions{
							VerticalAlignment: api.Center,
						})
						v.CurrentPlayerImage = img
						v.CurrentPlayerImageSource = v.PreviousPlayer
					}
				}
			}
			imgNeedsRefreshing := previousImage != img || v.FinalImage == nil
			finalImage := v.FinalImage
			if imgNeedsRefreshing {
				finalImage = overlayImage(img, info, handlerType)
				v.FinalImage = finalImage
			}
			text, percentage := calculateTextAndPercentage[v.Type](v.ActivePlayer)
			if math.IsNaN(percentage) || percentage < 0 || percentage > 100 {
				percentage = 100.0
			}
			infoNeedsRefreshing := false
			if percentage != v.Percentage {
				infoNeedsRefreshing = true
			}
			if text != v.Text {
				infoNeedsRefreshing = true
			}
			if !imgNeedsRefreshing && !infoNeedsRefreshing {
				break
			}
			v.Percentage = percentage
			v.Text = text
			if v.AccentColour == "" {
				v.AccentColour = getAverageColour(img)
			}
			w, h := info.GetDimensions(handlerType)
			imgParsed, err := api.DrawProgressBarWithAccent(finalImage, text, 5, float64(h-25), 20, float64(w-10), percentage, v.AccentColour)
			if err != nil {
				log.Println(err)
			} else {
				callback(imgParsed)
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (v *PlayerCtlHandler) FindImage(player *mpris.Player, info api.StreamDeckInfoV1, handlerType api.HandlerType) (image.Image, error) {
	metadata, err := player.GetMetadata()
	if err == nil && metadata != nil {
		artUrl, err := metadata.ArtURL()
		if err == nil && artUrl != "" {
			if artUrl == v.CurrentPlayerImageSource && v.CurrentPlayerImage != nil {
				return v.CurrentPlayerImage, nil
			}
			v.AccentColour = ""
			img, err := ExtractImage(artUrl)
			if err != nil {
				log.Println(err)
				err = nil
			}
			if img != nil {
				img = resizeThumbnail(img, info, handlerType)
				v.CurrentPlayerImage = img
				v.CurrentPlayerImageSource = artUrl
				return img, nil
			}
		}
		if err != nil {
			log.Println(err)
		}
		err = nil
	}
	if err != nil && err.Error() != "No player is being controlled by playerctld" {
		log.Println(err)
	}
	return nil, errors.New("couldn't find image")
}

func ExtractImage(icon string) (image.Image, error) {
	match, err := regexp.MatchString(`(https?://)[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&/=]*)`, icon)
	if match {
		return getHttpImage(icon)
	}
	if err != nil {
		log.Println(err)
		err = nil
	}
	match, err = regexp.MatchString(`(file://)?(/)+[a-zA-Z0-9\\\-_/ .]*\.+[a-z0-9A-Z]+`, icon)
	if match {
		icon = strings.ReplaceAll(icon, "file://", "")
		return loadImage(icon)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return nil, errors.New("couldn't find image")
}

func getHttpImage(url string) (image.Image, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("Couldn't get Image from URL")
	}
	defer response.Body.Close()
	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (v *PlayerCtlHandler) Input(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1, event api.InputEvent) {

	handlerTypeString, ok := fields["type"]

	if !ok && handlerType == api.LCD {
		log.Println("Type not specified")
		return
	}

	if handlerType == api.KEY {
		handlerTypeString = "Playback"
	}

	playerctlHandlerType := PlayerctlHandlerType(handlerTypeString.(string))

	playerName, ok := fields["player_name"]
	var playerNameString string
	if ok {
		playerNameString = playerName.(string)
	}
	player := choosePlayer(v.Client, playerNameString, v.PreviousPlayer, playerFilters[playerctlHandlerType])
	if player == nil {
		return
	}
	v.PreviousPlayer = player.GetName()

	if playerctlHandlerType == Volume {
		volume, err := player.GetVolume()
		if err != nil {
			log.Println(err)
			return
		}
		volume = math.Round(volume * 100.0)

		if event.EventType == api.KNOB_CCW {
			volume -= 1.0 * float64(event.RotateNotches)
		} else if event.EventType == api.KNOB_CW {
			volume += 1.0 * float64(event.RotateNotches)
		}
		volume /= 100.0
		err = player.SetVolume(volume)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		if event.EventType == api.KNOB_CCW {
			canGoPrevious, err := player.CanGoPrevious()
			if err != nil {
				return
			}
			if canGoPrevious {
				player.Previous()
			}
		} else if event.EventType == api.KNOB_CW {
			canGoNext, err := player.CanGoNext()
			if err != nil {
				return
			}
			if canGoNext {
				player.Next()
			}
		} else if event.EventType == api.KNOB_PRESS {
			player.PlayPause()
		} else if event.EventType == api.SCREEN_SHORT_TAP {
			status, err := player.GetLoopStatus()
			if err != nil {
				return
			}
			err = player.SetLoopStatus(getNextLoopStatus(status))
		} else if event.EventType == api.SCREEN_LONG_TAP {
			shuffle, err := player.GetShuffle()
			if err != nil {
				return
			}
			err = player.SetShuffle(!shuffle)
		} else if event.EventType == api.KEY_PRESS {
			operation, ok := fields["operation"]
			if !ok {
				log.Println("No MPRIS player operation specified")
				return
			}
			op, ok := operationsMap[operation.(string)]
			if !ok {
				log.Println("Invalid MPRIS player operation specified")
				return
			}
			var err error
			switch op {
			case PlayPause:
				err = player.PlayPause()
			case Play:
				err = player.Play()
			case Pause:
				err = player.Pause()
			case Previous:
				err = player.Previous()
			case Next:
				err = player.Next()
			case Shuffle:
				shuffle, err := player.GetShuffle()
				if err != nil {
					log.Println(err)
					return
				}
				err = player.SetShuffle(!shuffle)
				break
			case LoopStatus:
				status, err := player.GetLoopStatus()
				if err != nil {
					log.Println(err)
					return
				}
				err = player.SetLoopStatus(getNextLoopStatus(status))
				break
			}
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func getNextLoopStatus(status mpris.LoopStatus) mpris.LoopStatus {
	switch status {
	case "None":
		return "Track"
	case "Track":
		return "Playlist"
	case "Playlist":
		return "None"
	}
	return "None"
}
func choosePlayer(client *dbus.Conn, playerName, previousPlayerName string, filter func(player *mpris.Player) bool) *mpris.Player {
	var player *mpris.Player
	if !client.Connected() {
		client, _ = dbus.SessionBus()
	}
	players, err := mpris.List(client)
	if err != nil {
		log.Println(err)
		return nil
	}
	if playerName == "" {
		for _, p := range players {
			if strings.Contains(p, "playerctld") {
				log.Println("Using playerctld")
				player = mpris.New(client, p)
			}
		}
	}
	var previousPlayer *mpris.Player
	var pausedOption *mpris.Player
	for _, p := range players {
		pl := mpris.New(client, p)
		if pl.GetName() == previousPlayerName {
			previousPlayer = pl
		}
		if playerName != "" {
			if pl.GetShortName() == playerName {
				player = mpris.New(client, p)
				break
			}
		} else {
			if !filter(pl) {
				continue
			}
			status, err := pl.GetPlaybackStatus()
			if err != nil {
				log.Println(err)
				continue
			}
			if status == mpris.PlaybackPlaying {
				player = pl
				break
			} else {
				pausedOption = pl
			}
		}
	}
	if player == nil {
		if previousPlayer != nil {
			player = previousPlayer
		} else if pausedOption != nil {
			player = pausedOption
		}
	}
	return player
}

func formatDuration(microseconds int64) string {
	seconds := int(microseconds / 1000000)
	if seconds < 60 {
		return strconv.Itoa(seconds)
	}
	minutes := int(math.Floor(float64(seconds) / 60.0))
	seconds = seconds % 60
	if minutes < 60 {
		return strconv.Itoa(minutes) + ":" + pad(strconv.Itoa(seconds))
	}
	hours := int(math.Floor(float64(minutes) / 60.0))
	minutes = minutes % 60
	return strconv.Itoa(hours) + ":" + pad(strconv.Itoa(minutes)) + ":" + pad(strconv.Itoa(seconds))
}

func pad(timeSegment string) string {
	if len(timeSegment) == 1 {
		return "0" + timeSegment
	}
	return timeSegment
}

func resizeThumbnail(img image.Image, info api.StreamDeckInfoV1, handlerType api.HandlerType) image.Image {
	_, height := info.GetDimensions(handlerType)
	newSize := float64(height - 30)
	scalingFactor := newSize / float64(img.Bounds().Max.Y)
	x := float64(img.Bounds().Max.X) * scalingFactor
	y := float64(img.Bounds().Max.Y) * scalingFactor
	img = api.ResizeImageWH(img, int(math.Round(x)), int(math.Round(y)))
	return img
}

func overlayImage(img image.Image, info api.StreamDeckInfoV1, handlerType api.HandlerType) image.Image {
	width, height := info.GetDimensions(handlerType)
	mprisImg := img
	img = image.NewNRGBA(image.Rect(0, 0, width, height))
	ggImg := gg.NewContextForImage(img)
	ggImg.DrawImageAnchored(mprisImg, width/2, 35, 0.5, 0.5)
	return ggImg.Image()
}

func getAverageColour(img image.Image) string {
	imgSize := img.Bounds().Size()

	var redSum float64
	var greenSum float64
	var blueSum float64

	for x := 0; x < imgSize.X; x++ {
		for y := 0; y < imgSize.Y; y++ {
			pixel := img.At(x, y)
			col := color.RGBAModel.Convert(pixel).(color.RGBA)

			redSum += float64(col.R)
			greenSum += float64(col.G)
			blueSum += float64(col.B)
		}
	}

	imgArea := float64(imgSize.X * imgSize.Y)

	redAverage := math.Round(redSum / imgArea)
	greenAverage := math.Round(greenSum / imgArea)
	blueAverage := math.Round(blueSum / imgArea)

	return RGBToHex(int(redAverage), int(greenAverage), int(blueAverage))
}

func RGBToHex(r, g, b int) string {
	r = clamp(r, 0, 255)
	g = clamp(g, 0, 255)
	b = clamp(b, 0, 255)
	rHex := fmt.Sprintf("%02X", r)
	gHex := fmt.Sprintf("%02X", g)
	bHex := fmt.Sprintf("%02X", b)

	hex := "#" + rHex + gHex + bHex

	return hex

}

func clamp(value, min, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func playerNeedsRefreshing(player *mpris.Player) bool {
	if player == nil {
		return true
	}
	if player.GetShortName() == "playerctld" {
		return false
	}
	status, err := player.GetPlaybackStatus()
	if err != nil {
		return true
	}
	return status != mpris.PlaybackPlaying
}

func RegisterPlayerCtl() api.Module {

	return api.Module{
		Name: "Playerctl",
		NewForeground: func() api.ForegroundHandler {
			client, err := dbus.SessionBus()
			if err != nil {
				panic(err)
			}
			return &PlayerCtlHandler{Client: client}
		},
		ForegroundFields: []api.Field{
			{Title: "Icon", Name: "icon", Type: api.File, FileTypes: []string{".png", ".jpg", ".jpeg"}},
			{Title: "Accent Colour", Name: "colour", Type: api.Colour},
		},
		NewInput: func() api.InputHandler {
			client, err := dbus.SessionBus()
			if err != nil {
				panic(err)
			}
			return &PlayerCtlHandler{Running: true, Lock: semaphore.NewWeighted(1), Client: client}
		},
		InputFields: []api.Field{
			{Title: "Operation", Name: "operation", Type: api.Select, ListItems: slices.Collect(maps.Keys(operationsMap))},
		},
		LinkedFields: []api.Field{
			{Title: "Player Name", Name: "player_name", Type: api.Text},
			{Title: "Type", Name: "type", Type: api.Select, ListItems: []string{string(Playback), string(Volume)}},
		},
	}
}
