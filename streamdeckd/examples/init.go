package examples

import (
	"github.com/unix-streamdeck/streamdeckd/streamdeckd"
)

func RegisterBaseModules() {
	streamdeckd.RegisterModule(RegisterGif())
	streamdeckd.RegisterModule(RegisterTime())
	streamdeckd.RegisterModule(RegisterCounter())
	streamdeckd.RegisterModule(RegisterToggle())
	streamdeckd.RegisterModule(RegisterPlayerCtl())
	streamdeckd.RegisterModule(RegisterVolume())
	streamdeckd.RegisterModule(RegisterImage())
	streamdeckd.RegisterModule(RegisterText())
	streamdeckd.RegisterModule(RegisterCommand())
	streamdeckd.RegisterModule(RegisterUrl())
	streamdeckd.RegisterModule(RegisterKeybind())
	streamdeckd.RegisterModule(RegisterKeyHold())
	streamdeckd.RegisterModule(RegisterSwitchPage())
	streamdeckd.RegisterModule(RegisterBrightness())
}
