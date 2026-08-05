package streamdeckd

import (
	"log"
	"plugin"

	"github.com/unix-streamdeck/api/v3"
)

var modules []api.Module

func AvailableModules() []api.Module {
	return modules
}

func RegisterModule(m api.Module) {
	for _, module := range modules {
		if module.Name == m.Name {
			log.Println("Module already loaded: " + m.Name)
			return
		}
	}
	log.Println("Loaded module " + m.Name)
	if m.NewCombined != nil {
		handler := m.NewCombined()
		_, m.IsForeground = handler.(api.ForegroundHandler)
		_, m.IsInput = handler.(api.InputHandler)
	}
	m.IsBackground = m.NewBackground != nil
	modules = append(modules, m)
}

func LoadModule(path string) {
	plug, err := plugin.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	mod, err := plug.Lookup("GetModule")
	if err != nil {
		log.Println(err)
		return
	}
	var modMethod func() api.Module
	modMethod, ok := mod.(func() api.Module)
	if !ok {
		log.Println("Failed to load module: " + path)
		return
	}
	RegisterModule(modMethod())
}

func UnmountHandlers() {
	for s := range Devs {
		dev := Devs[s]
		dev.HandlerPruner().StopAllHandlers()
	}
}
