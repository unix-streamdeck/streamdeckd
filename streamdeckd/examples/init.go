package examples

import (
	"github.com/unix-streamdeck/streamdeckd/streamdeckd"
	"github.com/unix-streamdeck/streamdeckd/streamdeckd/examples/key"
)

func RegisterBaseModules() {
	streamdeckd.RegisterModule(key.RegisterGif())
	streamdeckd.RegisterModule(key.RegisterTime())
	streamdeckd.RegisterModule(key.RegisterCounter())
	streamdeckd.RegisterModule(key.RegisterSpotify())
}
