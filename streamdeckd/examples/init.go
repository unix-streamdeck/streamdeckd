package examples

import (
	"github.com/unix-streamdeck/streamdeckd/streamdeckd"
)

func RegisterBaseModules() {
	streamdeckd.RegisterModule(RegisterGif())
	streamdeckd.RegisterModule(RegisterTime())
	streamdeckd.RegisterModule(RegisterCounter())
	streamdeckd.RegisterModule(RegisterSpotify())
	streamdeckd.RegisterModule(RegisterToggle())
	streamdeckd.RegisterModule(RegisterPlayerCtl())
	streamdeckd.RegisterModule(RegisterVolume())
}
