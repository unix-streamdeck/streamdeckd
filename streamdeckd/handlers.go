package streamdeckd

import (
	"github.com/unix-streamdeck/api/v2"
	"log"
	"plugin"
)

type Module struct {
	Name              string
	NewIcon           func() api.IconHandler
	NewKey            func() api.KeyHandler
	NewLcd            func() api.LcdHandler
	NewKnobOrTouch    func() api.KnobOrTouchHandler
	IconFields        []api.Field
	KeyFields         []api.Field
	LcdFields         []api.Field
	KnobOrTouchFields []api.Field
	Linked            bool
	LinkedFields      []api.Field
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
