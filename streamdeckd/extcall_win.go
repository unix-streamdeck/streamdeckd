//go:build windows
// +build windows

package streamdeckd

func UpdateApplication() {

}

func InitDBUS() {

}

func ConnectScreensaver() (*ScreensaverConnection, error) {
    return &ScreensaverConnection{}, nil
}

func (c *ScreensaverConnection) RegisterScreensaverActiveListener() {

}