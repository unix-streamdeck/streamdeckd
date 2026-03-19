package mpris

import (
	"github.com/godbus/dbus/v5"
)

// A representation of a MPRIS player.
type Player struct {
	// A pointer to DBus connection that is used to acquire info or control the player
	conn *dbus.Conn
	// DBus object of the player
	obj *dbus.Object
	// Name of the player
	name string
}

// The status of the playback. May be "Playing", "Paused" or "Stopped".
type PlaybackStatus string

// The status of the player loop. May be "None", "Track" or "Playlist".
type LoopStatus string

// The type of signal received by channel. May be "PropertiesChanged", "NameOwnerChanged" and "Seeked"
type SignalType string

// A representation of MPRIS metadata.
// See also: https://www.freedesktop.org/wiki/Specifications/mpris-spec/metadata
type Metadata map[string]dbus.Variant
