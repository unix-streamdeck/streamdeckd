package streamdeckd

import (
	"log"
	"plugin"

	"github.com/unix-streamdeck/api/v2"
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
	m.IsKey = m.NewKey != nil
	m.IsIcon = m.NewIcon != nil
	m.IsKnob = m.NewKnobOrTouch != nil
	m.IsLcd = m.NewLcd != nil
	m.IsLinkedHandlers = m.LinkedFields != nil && len(m.LinkedFields) != 0
	modules = append(modules, m)
}

func LoadModule(path string) {
	plug, err := plugin.Open(path)
	if err != nil {
		//log.Println("Failed to load module: " + path)
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
		dev.UnmountHandlers()
	}
}

func UnmountKeyHandler(keyConfig *api.KeyConfigV3) {
	keyConfig.IconHandlerStruct.Stop()
	log.Printf("Stopped %s\n", keyConfig.IconHandler)
}
func UnmountKnobHandler(keyConfig *api.KnobConfigV3) {
	keyConfig.LcdHandlerStruct.Stop()
	log.Printf("Stopped %s\n", keyConfig.LcdHandler)
}
