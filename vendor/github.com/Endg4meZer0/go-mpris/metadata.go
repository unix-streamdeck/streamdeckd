package mpris

import (
	"errors"
	"time"

	"github.com/godbus/dbus/v5"
)

// Some metadata fields are parsed using different types
// because of general inconsistency between different players.

const timeFormat = "2006-01-02T15:04-07:00"

// Returns a unique identity for the track within the context of an MPRIS object (e.g. tracklist).
func (md Metadata) TrackID() (dbus.ObjectPath, error) {
	variant := md["mpris:trackid"].Value()
	if variant == "" {
		return "", nil
	}

	switch v := variant.(type) {
	case dbus.ObjectPath:
		return v, nil
	case *dbus.ObjectPath:
		return *v, nil
	case string:
		return dbus.ObjectPath(v), nil
	default:
		return "", errors.New("could not parse mpris:trackid")
	}
}

// Returns the duration of the track in microseconds.
// Why int64 and not uint64: https://www.freedesktop.org/wiki/Specifications/mpris-spec/metadata/#mpris:length
func (md Metadata) Length() (int64, error) {
	variant := md["mpris:length"].Value()
	if variant == nil {
		return 0, nil
	}

	switch v := variant.(type) {
	case uint64:
		return int64(v), nil
	case int:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return 0, errors.New("could not parse mpris:length")
	}
}

// Returns the location of an image representing the track or album.
func (md Metadata) ArtURL() (string, error) {
	va := md["mpris:artUrl"].Value()
	if va == nil {
		return "", nil
	}

	v, ok := va.(string)
	if !ok {
		return "", errors.New("could not parse mpris:artUrl")
	}

	return v, nil
}

// Returns the album name.
func (md Metadata) Album() (string, error) {
	variant := md["xesam:album"].Value()
	if variant == nil {
		return "", nil
	}

	v, ok := variant.(string)
	if !ok {
		return "", errors.New("could not parse xesam:album")
	}

	return v, nil
}

// Returns the album artist(s)
func (md Metadata) AlbumArtist() ([]string, error) {
	variant := md["xesam:albumArtist"].Value()
	if variant == nil {
		return nil, nil
	}

	v, ok := variant.([]string)
	if !ok {
		return nil, errors.New("could not parse xesam:albumArtist")
	}

	return v, nil
}

// Returns the track artist(s).
func (md Metadata) Artist() ([]string, error) {
	variant := md["xesam:artist"].Value()
	if variant == nil {
		return nil, nil
	}

	v, ok := variant.([]string)
	if !ok {
		return nil, errors.New("could not parse xesam:artist")
	}

	return v, nil
}

// Returns the track lyrics.
func (md Metadata) AsText() (string, error) {
	variant := md["xesam:asText"].Value()
	if variant == nil {
		return "", nil
	}

	v, ok := variant.(string)
	if !ok {
		return "", errors.New("could not parse xesam:asText")
	}

	return v, nil
}

// Returns the speed of the music in beats per minute.
func (md Metadata) AudioBPM() (int, error) {
	variant := md["xesam:audioBPM"].Value()
	if variant == nil {
		return 0, nil
	}

	v, ok := variant.(int)
	if !ok {
		return 0, errors.New("could not parse xesam:audioBPM")
	}

	return v, nil
}

// Returns an automatically-generated rating, based on things such as how often it has been played.
// This should be in the range 0.0 to 1.0.
func (md Metadata) AutoRating() (float64, error) {
	variant := md["xesam:autoRating"].Value()
	if variant == nil {
		return 0, nil
	}

	v, ok := variant.(float64)
	if !ok {
		return 0, errors.New("could not parse xesam:autoRating")
	}
	return v, nil
}

// Comment returns a (list of) freeform comment(s).
func (md Metadata) Comment() ([]string, error) {
	variant := md["xesam:comment"].Value()
	if variant == nil {
		return nil, nil
	}

	v, ok := variant.([]string)
	if !ok {
		return nil, errors.New("could not parse xesam:comment")
	}
	return v, nil
}

// Returns the composer(s) of the track.
func (md Metadata) Composer() ([]string, error) {
	variant := md["xesam:composer"].Value()
	if variant == nil {
		return nil, nil
	}

	v, ok := variant.([]string)
	if !ok {
		return nil, errors.New("could not parse xesam:composer")
	}
	return v, nil
}

// Returns when the track was created.
func (md Metadata) ContentCreated() (time.Time, error) {
	variant := md["xesam:contentCreated"].Value()
	if variant == nil {
		return time.Time{}, nil
	}

	vs, ok := variant.(string)
	if !ok {
		return time.Time{}, errors.New("could not parse xesam:contentCreated as a string")
	}

	t, err := time.Parse(timeFormat, vs)
	if err != nil {
		return time.Time{}, errors.New("could not parse xesam:contentCreated as a time object")
	}

	return t, nil
}

// Returns the disc number on the album that this track is from.
func (md Metadata) DiscNumber() (int, error) {
	variant := md["xesam:discNumber"].Value()
	if variant == nil {
		return 0, nil
	}

	v, ok := variant.(int)
	if !ok {
		return 0, errors.New("could not parse xesam:discNumber")
	}
	return v, nil
}

// Returns when the track was first played.
func (md Metadata) FirstUsed() (time.Time, error) {
	variant := md["xesam:firstUsed"].Value()
	if variant == nil {
		return time.Time{}, nil
	}

	vs, ok := variant.(string)
	if !ok {
		return time.Time{}, errors.New("could not parse xesam:firstUsed as a string")
	}

	t, err := time.Parse(timeFormat, vs)
	if err != nil {
		return time.Time{}, errors.New("could not parse xesam:firstUsed as a time object")
	}

	return t, nil
}

// Returns the genre(s) of the track.
func (md Metadata) Genre() ([]string, error) {
	variant := md["xesam:genre"].Value()
	if variant == nil {
		return nil, nil
	}

	v, ok := variant.([]string)
	if !ok {
		return nil, errors.New("could not parse xesam:genre")
	}
	return v, nil
}

// Returns when the track was last played.
func (md Metadata) LastUsed() (time.Time, error) {
	variant := md["xesam:lastUsed"].Value()
	if variant == nil {
		return time.Time{}, nil
	}

	vString, ok := variant.(string)
	if !ok {
		return time.Time{}, errors.New("could not parse xesam:lastUsed as a string")
	}

	v, err := time.Parse(timeFormat, vString)
	if err != nil {
		return time.Time{}, errors.New("could not parse xesam:lastUsed as a time object")
	}

	return v, nil
}

// Returns the lyricist(s) of the track.
func (md Metadata) Lyricist() ([]string, error) {
	variant := md["xesam:lyricist"].Value()
	if variant == nil {
		return nil, nil
	}

	v, ok := variant.([]string)
	if !ok {
		return nil, errors.New("could not parse xesam:lyricist")
	}
	return v, nil
}

// Returns the track title.
func (md Metadata) Title() (string, error) {
	variant := md["xesam:title"].Value()
	if variant == nil {
		return "", nil
	}

	v, ok := variant.(string)
	if !ok {
		return "", errors.New("could not parse xesam:title")
	}

	return v, nil
}

// TrackNumber returns the track number on the album disc.
func (md Metadata) TrackNumber() (int, error) {
	variant := md["xesam:trackNumber"].Value()
	if variant == nil {
		return 0, nil
	}

	v, ok := variant.(int)
	if !ok {
		return 0, errors.New("could not parse xesam:trackNumber")
	}
	return v, nil
}

// Returns the location of the media file.
func (md Metadata) URL() (string, error) {
	variant := md["xesam:url"].Value()
	if variant == nil {
		return "", nil
	}

	v, ok := variant.(string)
	if !ok {
		return "", errors.New("could not parse xesam:url")
	}

	return v, nil
}

// Returns the number of times the track has been played.
func (md Metadata) UseCount() (int, error) {
	variant := md["xesam:useCount"].Value()
	if variant == nil {
		return 0, nil
	}

	v, ok := variant.(int)
	if !ok {
		return 0, errors.New("could not parse xesam:useCount")
	}
	return v, nil
}

// UserRating returns a user-specified rating. This should be in the range 0.0 to 1.0.
func (md Metadata) UserRating() (float64, error) {
	variant := md["xesam:userRating"].Value()
	if variant == nil {
		return 0, nil
	}

	v, ok := variant.(float64)
	if !ok {
		return 0, errors.New("could not parse xesam:userRating")
	}
	return v, nil
}

// Find returns a generic representation of the requested value when present.
func (md Metadata) Find(key string) (dbus.Variant, bool) {
	variant, found := md[key]
	return variant, found
}
