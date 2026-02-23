package key

import (
	"errors"

	"image"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/unix-streamdeck/api/v2"
)

type SpotifyIconHandler struct {
	Running bool
	oldUrl  string
	Quit    chan bool
}

func (s *SpotifyIconHandler) Start(key api.KeyConfigV3, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	s.Running = true
	if s.Quit == nil {
		s.Quit = make(chan bool)
	}
	c, err := Connect()
	if err != nil {
		if err.Error() != "The name org.mpris.MediaPlayer2.spotify was not provided by any .service files" {
			log.Println(err)
		}
		return
	}
	go s.run(c, callback)
}

func (s *SpotifyIconHandler) IsRunning() bool {
	return s.Running
}

func (s *SpotifyIconHandler) SetRunning(running bool) {
	s.Running = running
}

func (s *SpotifyIconHandler) Stop() {
	s.Running = false
	s.Quit <- true
	s.oldUrl = ""
}

func (s *SpotifyIconHandler) run(c *Connection, callback func(image image.Image)) {
	defer c.Close()
	for {
		select {
		case <-s.Quit:
			return
		default:
			url, err := c.GetAlbumArtUrl()
			if err != nil {
				if err.Error() != "The name org.mpris.MediaPlayer2.spotify was not provided by any .service files" {
					log.Println(err)
				}
				time.Sleep(time.Second)
				continue
			}
			if url == s.oldUrl {
				time.Sleep(time.Second)
				continue
			}
			img, err := getImage(url)
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second)
				continue
			}
			callback(img)
			s.oldUrl = url
			time.Sleep(time.Second)
		}
	}
}

func RegisterSpotify() api.Module {
	return api.Module{NewIcon: func() api.IconHandler {
		return &SpotifyIconHandler{Running: true}
	}, Name: "Spotify"}
}

// region DBus
func getImage(url string) (image.Image, error) {
	response, err := http.Get(url)
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

type Connection struct {
	busobj dbus.BusObject
	conn   *dbus.Conn
}

func Connect() (*Connection, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}
	return &Connection{
		conn:   conn,
		busobj: conn.Object("org.mpris.MediaPlayer2.spotify", "/org/mpris/MediaPlayer2"),
	}, nil
}

func (c *Connection) GetAlbumArtUrl() (string, error) {
	variant, err := c.busobj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
	if err != nil {
		return "", err
	}
	metadataMap := variant.Value().(map[string]dbus.Variant)
	var url string
	for key, val := range metadataMap {
		if key == "mpris:artUrl" {
			url = val.String()
		}
	}
	if url == "" {
		return "", errors.New("Couldn't get URL from DBus")
	}
	url = strings.ReplaceAll(url, "\"", "")
	url = strings.ReplaceAll(url, "https://open.spotify.com/image/", "https://i.scdn.co/image/")
	return url, nil
}

func (c *Connection) Close() {
	c.conn.Close()
}

// endregion
