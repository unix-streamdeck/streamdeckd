package examples

import (
	"github.com/unix-streamdeck/api/v3"
)

func GetModuleGroup() api.ModuleGroup {
	return api.ModuleGroup{
		Name:        "examples",
		DisplayName: "Builtin Example Modules",
		Description: "Basic example handlers that ship with streamdeckd",
		Modules: []api.Module{
			RegisterGif(),
			RegisterTime(),
			RegisterCounter(),
			RegisterToggle(),
			RegisterPlayerCtl(),
			RegisterVolume(),
		},
	}

}
