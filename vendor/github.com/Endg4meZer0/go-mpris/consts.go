package mpris

// DBus consts
const (
	dbusObjectPath = "/org/mpris/MediaPlayer2"

	dbusInterface           = "org.freedesktop.DBus"
	dbusPropertiesInterface = dbusInterface + ".Properties"

	BaseInterface      = "org.mpris.MediaPlayer2"
	PlayerInterface    = "org.mpris.MediaPlayer2.Player"
	TrackListInterface = "org.mpris.MediaPlayer2.TrackList"
	PlaylistsInterface = "org.mpris.MediaPlayer2.Playlists"

	getPropertyMethod = dbusPropertiesInterface + ".Get"
	setPropertyMethod = dbusPropertiesInterface + ".Set"

	propertiesChangedSignal    = dbusPropertiesInterface + ".PropertiesChanged"
	nameOwnerChangedSignal     = dbusInterface + ".NameOwnerChanged"
	seekedSignal               = PlayerInterface + ".Seeked"
	trackListReplacedSignal    = TrackListInterface + ".TrackListReplaced"
	trackAddedSignal           = TrackListInterface + ".TrackAdded"
	trackRemovedSignal         = TrackListInterface + ".TrackRemoved"
	trackMetadataChangedSignal = TrackListInterface + ".TrackMetadataChanged"
)

// Playback statuses
const (
	PlaybackPlaying PlaybackStatus = "Playing"
	PlaybackPaused  PlaybackStatus = "Paused"
	PlaybackStopped PlaybackStatus = "Stopped"
)

// Loop statuses
const (
	LoopNone     LoopStatus = "None"
	LoopTrack    LoopStatus = "Track"
	LoopPlaylist LoopStatus = "Playlist"
)

// Signal types
const (
	SignalNotSupported         SignalType = ""
	SignalPropertiesChanged    SignalType = "PropertiesChanged"
	SignalNameOwnerChanged     SignalType = "NameOwnerChanged"
	SignalSeeked               SignalType = "Seeked"
	SignalTrackListReplaced    SignalType = "TrackListReplaced"
	SignalTrackAdded           SignalType = "TrackAdded"
	SignalTrackRemoved         SignalType = "TrackRemoved"
	SignalTrackMetadataChanged SignalType = "TrackMetadataChanged"
)
