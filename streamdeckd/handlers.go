package streamdeckd

import (
	"errors"
	"fmt"
	"log"
	"plugin"
	"regexp"
	"strings"

	"github.com/unix-streamdeck/api/v3"
)

var modules api.AvailableModules

func AvailableModules() api.AvailableModules {
	return modules
}

const moduleNameRegex = "^[a-zA-Z0-9\\-]+$"

func RegisterModule(m api.Module) {
	m, err := validateModule(m)
	if err != nil {
		log.Println("[ERROR]", err)
		return
	}
	for _, module := range modules.IndividualModules {
		if module.Name == m.Name {
			log.Println("Module already loaded: " + m.Name)
			return
		}
	}
	modules.IndividualModules = append(modules.IndividualModules, m)
}

func validateModule(m api.Module) (api.Module, error) {
	match, _ := regexp.MatchString(moduleNameRegex, m.Name)
	if !match {
		return m, errors.New(fmt.Sprintf("error loading module: `%s`, name does not match format", m.Name))
	}
	if m.NewCombined != nil {
		handler := m.NewCombined()
		_, m.IsForeground = handler.(api.ForegroundHandler)
		_, m.IsInput = handler.(api.InputHandler)
	}

	m.IsBackground = m.NewBackground != nil

	if m.DisplayName == "" {
		m.DisplayName = m.Name
	}
	return m, nil
}

func RegisterModuleGroup(mg api.ModuleGroup) {
	match, _ := regexp.MatchString(moduleNameRegex, mg.Name)
	if !match {
		log.Printf("[ERROR] Error loading module: `%s`, name does not match format", mg.Name)
	}
	for _, module := range modules.GroupModules {
		if module.Name == mg.Name {
			log.Println("ModuleGroup already loaded: " + mg.Name)
			return
		}
	}
	for i, m := range mg.Modules {
		m, err := validateModule(m)
		if err != nil {
			log.Println("[ERROR]", err)
			return
		}
		mg.Modules[i] = m
	}

	if mg.DisplayName == "" {
		mg.DisplayName = mg.Name
	}

	var moduleNames []string

	for _, module := range mg.Modules {
		moduleNames = append(moduleNames, mg.Name + "." + module.Name)
	}

	log.Println("Loaded module group " + mg.Name + " with child modules: " + strings.Join(moduleNames, ", "))
	modules.GroupModules = append(modules.GroupModules, mg)
}

func LoadModule(path string) {
	plug, err := plugin.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	mod, err := plug.Lookup("GetModule")
	if err == nil {
		var modMethod func() api.Module
		modMethod, ok := mod.(func() api.Module)
		if !ok {
			log.Println("Failed to load module: " + path)
			return
		}
		RegisterModule(modMethod())
		return
	}

	modGroup, err := plug.Lookup("GetModuleGroup")
	if err == nil {
		var modGroupMethod func() api.ModuleGroup
		modGroupMethod, ok := modGroup.(func() api.ModuleGroup)
		if !ok {
			log.Println("Failed to load modulegroup: " + path)
			return
		}
		RegisterModuleGroup(modGroupMethod())
		return
	}
	log.Println("[ERROR] Failed to load module/modulegroup: " + path + ", no GetModule or GetModuleGroup method found")
}

func UnmountHandlers() {
	for s := range Devs {
		dev := Devs[s]
		dev.HandlerPruner().StopAllHandlers()
	}
}

type HandlerContext uint

const (
	FOREGROUND HandlerContext = 1
	INPUT      HandlerContext = 2
	BACKGROUND HandlerContext = 3
)

func moduleHasContext(module api.Module, context HandlerContext) bool {
	switch context {
	case FOREGROUND:
		return module.IsForeground
	case INPUT:
		return module.IsInput
	case BACKGROUND:
		return module.IsBackground
	default:
		return false
	}
}

func GetModule(moduleName string, context HandlerContext) (api.Module, error) {
	if strings.Contains(moduleName, ".") {
		nameSplit := strings.Split(moduleName, ".")
		if len(nameSplit) != 2 {
			return api.Module{}, errors.New("malformed module name")
		}
		for _, group := range modules.GroupModules {
			if group.Name == nameSplit[0] {
				for _, module := range group.Modules {
					if module.Name == nameSplit[1] {
						if !moduleHasContext(module, context) {
							return module, errors.New("module does not support context")
						}
						return module, nil
					}
				}
			}
		}
	} else {
		for _, module := range modules.IndividualModules {
			if module.Name == moduleName {
				if !moduleHasContext(module, context) {
					return module, errors.New("module does not support context")
				}
				return module, nil
			}
		}
	}
	return api.Module{}, errors.New("could not find module: " + moduleName)
}
