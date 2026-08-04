//go:build darwin

package streamdeckd

import (
	"errors"
	"log"

	"github.com/bendahl/uinput"
)

var kb uinput.Keyboard

func UpdateApplication() {
	log.Println("Application based configuration is not currently supported on macOS, due to the difficulty in getting the current active application, sorry :(")
}

func EnableVirtualKeyboard() {
	log.Println("Virtual keyboard is not currently supported on macOS, sorry :(")
}

func ConnectScreensaver() (*ScreensaverConnection, error) {
	return nil, errors.New("Screensaver detection is not currently supported on macOS, sorry :(")
}

func (c *ScreensaverConnection) RegisterScreensaverActiveListener() {

}
