package basic

import (
	"github.com/unix-streamdeck/api/v3"
)

func GetModuleGroup() api.ModuleGroup {
	return api.ModuleGroup{
		Name:        "builtin",
		DisplayName: "Builtin Basics",
		Description: "The basic handlers that ship with streamdeckd, these used to be the builtin actions on the keys",
		Modules: []api.Module{
			RegisterImage(),
			RegisterText(),
			RegisterCommand(),
			RegisterUrl(),
			RegisterKeybind(),
			RegisterKeyHold(),
			RegisterSwitchPage(),
			RegisterBrightness(),
		},
	}

}
