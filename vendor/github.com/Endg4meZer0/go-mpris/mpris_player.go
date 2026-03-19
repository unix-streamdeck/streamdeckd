package mpris

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
    __  _________________  ______  ____  _____
   /  |/  / ____/_  __/ / / / __ \/ __ \/ ___/
  / /|_/ / __/   / / / /_/ / / / / / / /\__ \
 / /  / / /___  / / / __  / /_/ / /_/ /___/ /
/_/  /_/_____/ /_/ /_/ /_/\____/_____//____/
*/

// Skips to the next track in the tracklist.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Method:Next
func (i *Player) Next() error {
	return i.obj.Call(PlayerInterface+".Next", 0).Err
}

// Skips to the previous track in the tracklist.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Method:Previous
func (i *Player) Previous() error {
	return i.obj.Call(PlayerInterface+".Previous", 0).Err
}

// Pauses playback.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Method:Pause
func (i *Player) Pause() error {
	return i.obj.Call(PlayerInterface+".Pause", 0).Err
}

// Resumes playback if paused and pauses playback if playing.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Method:PlayPause
func (i *Player) PlayPause() error {
	return i.obj.Call(PlayerInterface+".PlayPause", 0).Err
}

// Stops playback.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Method:Stop
func (i *Player) Stop() error {
	return i.obj.Call(PlayerInterface+".Stop", 0).Err
}

// Starts or resumes playback.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Method:Play
func (i *Player) Play() error {
	return i.obj.Call(PlayerInterface+".Play", 0).Err
}

// Seeks in the current track position by the specified offset. The offset should be in microseconds.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Method:Seek
func (i *Player) SeekBy(offset int64) error {
	return i.obj.Call(PlayerInterface+".Seek", 0, offset).Err
}

// Sets the specified track's position in microseconds (if it's playing).
// Perhaps you would like to use SetPosition instead.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Method:SetPosition
func (i *Player) SetTrackPosition(trackId dbus.ObjectPath, position int64) error {
	return i.obj.Call(PlayerInterface+".SetPosition", 0, trackId, position).Err
}

// Opens the Uri for a playback, if supported.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Method:OpenUri
func (i *Player) OpenUri(uri string) error {
	return i.obj.Call(PlayerInterface+".OpenUri", 0, uri).Err
}

/*
    ____  ____  ____  ____  __________  _____________________
   / __ \/ __ \/ __ \/ __ \/ ____/ __ \/_  __/  _/ ____/ ___/
  / /_/ / /_/ / / / / /_/ / __/ / /_/ / / /  / // __/  \__ \
 / ____/ _, _/ /_/ / ____/ /___/ _, _/ / / _/ // /___ ___/ /
/_/   /_/ |_|\____/_/   /_____/_/ |_| /_/ /___/_____//____/
*/

// Returns the playback status.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:PlaybackStatus
func (i *Player) GetPlaybackStatus() (PlaybackStatus, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "PlaybackStatus")
	if err != nil {
		return "", err
	}
	if variant.Value() == nil {
		return "", errors.New("variant value is nil")
	}
	value, ok := variant.Value().(string)
	if !ok {
		return "", errors.New("variant type is not string")
	}

	return PlaybackStatus(value), nil
}

// Returns the loop status.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:LoopStatus
func (i *Player) GetLoopStatus() (LoopStatus, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "LoopStatus")
	if err != nil {
		return LoopStatus(""), err
	}
	if variant.Value() == nil {
		return "", errors.New("variant value is nil")
	}
	value, ok := variant.Value().(string)
	if !ok {
		return "", errors.New("variant type is not string")
	}
	return LoopStatus(value), nil
}

// Sets the loop status to the specified value.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:LoopStatus
func (i *Player) SetLoopStatus(loopStatus LoopStatus) error {
	return setProperty(i.obj, PlayerInterface, "LoopStatus", loopStatus)
}

// Returns the current playback rate.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:Rate
func (i *Player) GetRate() (float64, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "Rate")
	if err != nil {
		return 0.0, err
	}
	if variant.Value() == nil {
		return 0.0, errors.New("variant value is nil")
	}
	value, ok := variant.Value().(float64)
	if !ok {
		return 0.0, errors.New("variant type is not float64")
	}
	return value, nil
}

// Returns if shuffle is enabled or not.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:Shuffle
func (i *Player) GetShuffle() (bool, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "Shuffle")
	if err != nil {
		return false, err
	}
	if variant.Value() == nil {
		return false, errors.New("variant value is nil")
	}
	value, ok := variant.Value().(bool)
	if !ok {
		return false, errors.New("variant type is not bool")
	}
	return value, nil
}

// Sets shuffle on/off.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:Shuffle
func (i *Player) SetShuffle(value bool) error {
	return setProperty(i.obj, PlayerInterface, "Shuffle", value)
}

// Returns the metadata.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:Metadata
func (i *Player) GetMetadata() (Metadata, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "Metadata")
	if err != nil {
		return nil, err
	}
	if variant.Value() == nil {
		return nil, errors.New("variant value is nil")
	}
	value, ok := variant.Value().(map[string]dbus.Variant)
	if !ok {
		return nil, errors.New("variant type is not map[string]dbus.Variant")
	}
	return value, nil
}

// Returns the volume.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:Volume
func (i *Player) GetVolume() (float64, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "Volume")
	if err != nil {
		return 0.0, err
	}
	if variant.Value() == nil {
		return 0.0, errors.New("variant value is nil")
	}
	value, ok := variant.Value().(float64)
	if !ok {
		return 0.0, errors.New("variant type is not float64")
	}
	return value, nil
}

// Sets the volume to the specified value.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:Volume
func (i *Player) SetVolume(value float64) error {
	return setProperty(i.obj, PlayerInterface, "Volume", value)
}

// Returns the currently playing track's position in microseconds. If there isn't any, returns 0.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:Position
func (i *Player) GetPosition() (int64, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "Position")
	if err != nil {
		return 0.0, err
	}
	if variant.Value() == nil {
		return 0.0, errors.New("variant value is nil")
	}
	value, ok := variant.Value().(int64)
	if !ok {
		return 0.0, errors.New("variant type is not int64")
	}
	return value, nil
}

// Sets the currently playing track's position in microseconds (if there is any).
// Not to confuse with SetTrackPosition.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:Position
func (i *Player) SetPosition(position int64) error {
	metadata, err := i.GetMetadata()
	if err != nil {
		return err
	}
	if metadata == nil {
		return errors.New("metadata is nil")
	}
	trackId, err := metadata.TrackID()
	if err != nil {
		return err
	}
	i.SetTrackPosition(trackId, position)
	return nil
}

// Returns the minimum value that Rate can take.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:MinimumRate
func (i *Player) GetMinimumRate() (float64, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "MinimumRate")
	if err != nil {
		return 1.0, err
	}
	if variant.Value() == nil {
		return 1.0, errors.New("variant value is nil")
	}
	value, ok := variant.Value().(float64)
	if !ok {
		return 0.0, errors.New("variant type is not float64")
	}
	return value, nil
}

// Returns the maximum value that Rate can take.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:MaximumRate
func (i *Player) GetMaximumRate() (float64, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "MaximumRate")
	if err != nil {
		return 1.0, err
	}
	if variant.Value() == nil {
		return 1.0, errors.New("variant value is nil")
	}
	value, ok := variant.Value().(float64)
	if !ok {
		return 0.0, errors.New("variant type is not float64")
	}
	return value, nil
}

// Returns if the player can switch to the next track using the Next call.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:CanGoNext
func (i *Player) CanGoNext() (bool, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "CanGoNext")
	if err != nil {
		return false, err
	}
	if variant.Value() == nil {
		return false, errors.New("variant value is nil")
	}
	value, ok := variant.Value().(bool)
	if !ok {
		return false, errors.New("variant type is not bool")
	}
	return value, nil
}

// Returns if the player can switch to the previous track using the Previous call.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:CanGoPrevious
func (i *Player) CanGoPrevious() (bool, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "CanGoPrevious")
	if err != nil {
		return false, err
	}
	if variant.Value() == nil {
		return false, errors.New("variant value is nil")
	}
	value, ok := variant.Value().(bool)
	if !ok {
		return false, errors.New("variant type is not bool")
	}
	return value, nil
}

// Returns if the player can be started by Play or PlayPause calls.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:CanPlay
func (i *Player) CanPlay() (bool, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "CanPlay")
	if err != nil {
		return false, err
	}
	if variant.Value() == nil {
		return false, errors.New("variant value is nil")
	}
	value, ok := variant.Value().(bool)
	if !ok {
		return false, errors.New("variant type is not bool")
	}
	return value, nil
}

// Returns if the player can be paused by Pause or PlayPause calls.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:CanPause
func (i *Player) CanPause() (bool, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "CanPause")
	if err != nil {
		return false, err
	}
	if variant.Value() == nil {
		return false, errors.New("variant value is nil")
	}
	value, ok := variant.Value().(bool)
	if !ok {
		return false, errors.New("variant type is not bool")
	}
	return value, nil
}

// Returns if the position can be controlled by Seek and SetPosition calls.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:CanSeek
func (i *Player) CanSeek() (bool, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "CanSeek")
	if err != nil {
		return false, err
	}
	if variant.Value() == nil {
		return false, errors.New("variant value is nil")
	}
	value, ok := variant.Value().(bool)
	if !ok {
		return false, errors.New("variant type is not bool")
	}
	return value, nil
}

// Returns if the player can be controlled by calls.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Property:CanControl
func (i *Player) CanControl() (bool, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "CanControl")
	if err != nil {
		return false, err
	}
	if variant.Value() == nil {
		return false, errors.New("variant value is nil")
	}
	value, ok := variant.Value().(bool)
	if !ok {
		return false, errors.New("variant type is not bool")
	}
	return value, nil
}
