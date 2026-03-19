package mpris

import "errors"

/*
    __  _________________  ______  ____  _____
   /  |/  / ____/_  __/ / / / __ \/ __ \/ ___/
  / /|_/ / __/   / / / /_/ / / / / / / /\__ \
 / /  / / /___  / / / __  / /_/ / /_/ /___/ /
/_/  /_/_____/ /_/ /_/ /_/\____/_____//____/
*/

// Raises player priority.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Method:Raise
func (i *Player) Raise() error {
	return i.obj.Call(BaseInterface+".Raise", 0).Err
}

// Closes the player.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Method:Quit
func (i *Player) Quit() error {
	return i.obj.Call(BaseInterface+".Quit", 0).Err
}

/*
    ____  ____  ____  ____  __________  _____________________
   / __ \/ __ \/ __ \/ __ \/ ____/ __ \/_  __/  _/ ____/ ___/
  / /_/ / /_/ / / / / /_/ / __/ / /_/ / / /  / // __/  \__ \
 / ____/ _, _/ /_/ / ____/ /___/ _, _/ / / _/ // /___ ___/ /
/_/   /_/ |_|\____/_/   /_____/_/ |_| /_/ /___/_____//____/
*/

// Returns the CanQuit property of the player.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:CanQuit
func (i *Player) CanQuit() (bool, error) {
	variant, err := getProperty(i.obj, BaseInterface, "CanQuit")
	if err != nil {
		return false, err
	}

	if variant.Value() == nil {
		return false, errors.New("variant value is nil")
	}

	return variant.Value().(bool), nil
}

// Added in MPRIS v2.2. Returns the Fullscreen property of the player.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:Fullscreen
func (i *Player) GetFullscreen() (bool, error) {
	variant, err := getProperty(i.obj, BaseInterface, "Fullscreen")
	if err != nil {
		return false, err
	}

	if variant.Value() == nil {
		return false, errors.New("variant value is nil")
	}

	return variant.Value().(bool), nil
}

// Added in MPRIS v2.2. Sets the Fullscreen property of the player.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:Fullscreen
func (i *Player) SetFullscreen(value bool) error {
	err := setProperty(i.obj, BaseInterface, "Fullscreen", value)
	if err != nil {
		return err
	}

	return nil
}

// Added in MPRIS v2.2. Returns the CanSetFullscreen property of the player.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:CanSetFullscreen
func (i *Player) CanSetFullscreen() (bool, error) {
	variant, err := getProperty(i.obj, BaseInterface, "CanSetFullscreen")
	if err != nil {
		return false, err
	}

	if variant.Value() == nil {
		return false, errors.New("variant value is nil")
	}

	return variant.Value().(bool), nil
}

// Returns the CanRaise property of the player.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:CanRaise
func (i *Player) CanRaise() (bool, error) {
	variant, err := getProperty(i.obj, BaseInterface, "CanRaise")
	if err != nil {
		return false, err
	}

	if variant.Value() == nil {
		return false, errors.New("variant value is nil")
	}

	return variant.Value().(bool), nil
}

// Returns the HasTrackList property of the player.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:HasTrackList
func (i *Player) HasTrackList() (bool, error) {
	variant, err := getProperty(i.obj, BaseInterface, "HasTrackList")
	if err != nil {
		return false, err
	}

	if variant.Value() == nil {
		return false, errors.New("variant value is nil")
	}

	return variant.Value().(bool), nil
}

// Returns the Identity property, which is the friendly name of the player.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:Identity
func (i *Player) GetIdentity() (string, error) {
	variant, err := getProperty(i.obj, BaseInterface, "Identity")

	if err != nil {
		return "", err
	}

	if variant.Value() == nil {
		return "", errors.New("variant value is nil")
	}

	return variant.Value().(string), err
}

// Returns the DesktopEntry property of the player.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:DesktopEntry
func (i *Player) GetDesktopEntry() (string, error) {
	variant, err := getProperty(i.obj, BaseInterface, "DesktopEntry")

	if err != nil {
		return "", err
	}

	if variant.Value() == nil {
		return "", errors.New("variant value is nil")
	}

	return variant.Value().(string), err
}

// Returns the SupportedUriSchemes property of the player.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:SupportedUriSchemes
func (i *Player) GetSupportedUriSchemes() ([]string, error) {
	variant, err := getProperty(i.obj, BaseInterface, "SupportedUriSchemes")

	if err != nil {
		return nil, err
	}

	if variant.Value() == nil {
		return nil, errors.New("variant value is nil")
	}

	return variant.Value().([]string), err
}

// Returns the SupportedMimeTypes property of the player.
// See also: https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:SupportedMimeTypes
func (i *Player) GetSupportedMimeTypes() ([]string, error) {
	variant, err := getProperty(i.obj, BaseInterface, "SupportedMimeTypes")

	if err != nil {
		return nil, err
	}

	if variant.Value() == nil {
		return nil, errors.New("variant value is nil")
	}

	return variant.Value().([]string), err
}
