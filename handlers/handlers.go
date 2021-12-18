package handlers

import (
	"github.com/unix-streamdeck/api"
	"log"
	"plugin"
)

type Module struct {
	Name	string
	NewIcon func() api.IconHandler
	NewKey func() api.KeyHandler
	IconFields []api.Field
	KeyFields []api.Field
}



var modules []Module


func AvailableModules() []Module {
	return modules
}

func RegisterModule(m Module) {
	for _, module := range modules {
		if module.Name == m.Name {
			log.Println("Module already loaded: " + m.Name)
			return
		}
	}
	log.Println("Loaded module " + m.Name)
	modules = append(modules, m)
}

func LoadModule(path string) {
	plug, err := plugin.Open(path)
	if err != nil {
		//log.Println("Failed to load module: " + path)
		log.Println(err)
		return
	}
	mods, err := plug.Lookup("GetModules")
	if err == nil {
		modsMethod, ok := mods.(func() []Module)
		if !ok {
			log.Println("Failed to get list of modules: ", path)
			return
		}
		modules := modsMethod()
		for idx := range modules {
			RegisterModule(modules[idx])
		}
	} else {
		mod, err := plug.Lookup("GetModule")
		if err != nil {
			log.Println(err)
			return
		}
		var modMethod func() Module
		modMethod, ok := mod.(func() Module)
		if !ok {
			log.Println("Failed to load module: " + path)
			return
		}
		RegisterModule(modMethod())
	}
}
