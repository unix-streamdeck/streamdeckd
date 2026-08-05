package streamdeckd

import (
	"fmt"

	"github.com/unix-streamdeck/api/v3"
)

type IHandlerPruner interface {
	OnPageChange()
	OnAppSwitch()
	StopAllHandlers()
}

type HandlerPruner struct {
	vdev IVirtualDev
}

func (hp *HandlerPruner) OnPageChange() {
	hp.vdev.PageManager().AttachListener(func(newPage, previousPage int) {
		if newPage != previousPage {
			hp.stopPageHandlers(previousPage)
		}
	})
}

func (hp *HandlerPruner) stopPageHandlers(pageNo int) {
	page := hp.vdev.Config().Pages[pageNo]

	if page.TouchPanelBackgroundHandler != nil {
		go hp.stopHandler(page.TouchPanelBackgroundHandler, page.TouchPanelBackground, fmt.Sprintf("page %d background", pageNo))
	}

	if page.KeyGridBackgroundHandler != nil {
		go hp.stopHandler(page.KeyGridBackgroundHandler, page.KeyGridBackground, fmt.Sprintf("page %d background", pageNo))
	}

	for i, key := range page.Keys {

		if key.KeyBackgroundHandler != nil {
			go hp.stopHandler(key.KeyBackgroundHandler, key.KeyBackground, fmt.Sprintf("page %d, key %d background", pageNo, i))
		}

		for j, layer := range key.Layers {
			go hp.stopHandler(layer.HandlerStruct, layer.Handler, fmt.Sprintf("page %d, key %d, global layer %d", pageNo, i, j))
		}

		for _, keyConfig := range key.Application {
			if keyConfig.KeyBackgroundHandler != nil {
				go hp.stopHandler(keyConfig.KeyBackgroundHandler, keyConfig.KeyBackground, fmt.Sprintf("page %d, key %d background", pageNo, i))
			}

			for j, layer := range keyConfig.Layers {
				if layer.Global {
					continue
				}
				go hp.stopHandler(layer.HandlerStruct, layer.Handler, fmt.Sprintf("page %d, key %d, layer %d", pageNo, i, j))
			}
		}
	}

	for i, knob := range page.Knobs {

		if knob.TouchPanelBackgroundHandler != nil {
			go hp.stopHandler(knob.TouchPanelBackgroundHandler, knob.TouchPanelBackground, fmt.Sprintf("page %d, knob %d background", pageNo, i))
		}

		for j, layer := range knob.Layers {
			go hp.stopHandler(layer.HandlerStruct, layer.Handler, fmt.Sprintf("page %d, knob %d, global layer %d", pageNo, i, j))
		}

		for _, knobConfig := range knob.Application {

			if knobConfig.TouchPanelBackgroundHandler != nil {
				go hp.stopHandler(knobConfig.TouchPanelBackgroundHandler, knobConfig.TouchPanelBackground, fmt.Sprintf("page %d, knob %d background", pageNo, i))
			}

			for j, layer := range knobConfig.Layers {
				if layer.Global {
					continue
				}
				go hp.stopHandler(layer.HandlerStruct, layer.Handler, fmt.Sprintf("page %d, knob %d, layer %d", pageNo, i, j))
			}
		}
	}
}

func (hp *HandlerPruner) OnAppSwitch() {
	applicationManager.AttachListener(func(activeApp string) {
		page := hp.vdev.Config().Pages[hp.vdev.PageManager().GetPage()]

		for i, key := range page.Keys {
			_, hasNewActiveApp := key.Application[activeApp]

			var newActiveApp string

			if hasNewActiveApp {
				newActiveApp = activeApp
			} else {
				newActiveApp = ""
			}

			for appName, keyConfig := range key.Application {
				if appName == newActiveApp {
					continue
				}

				if keyConfig.KeyBackgroundHandler != nil {
					go hp.stopHandler(keyConfig.KeyBackgroundHandler, keyConfig.KeyBackground, fmt.Sprintf("key %d app `%s`", i, appName))
				}

				for j, layer := range keyConfig.Layers {
					if layer.Global {
						continue
					}
					go hp.stopHandler(layer.HandlerStruct, layer.Handler, fmt.Sprintf("key %d app %s layer %d", i, appName, j))
				}
			}
		}

		for i, knob := range page.Knobs {

			_, hasNewActiveApp := knob.Application[activeApp]

			var newActiveApp string

			if hasNewActiveApp {
				newActiveApp = activeApp
			} else {
				newActiveApp = ""
			}

			for appName, knobConfig := range knob.Application {
				if appName == newActiveApp {
					continue
				}

				if knobConfig.TouchPanelBackgroundHandler != nil {
					go hp.stopHandler(knobConfig.TouchPanelBackgroundHandler, knobConfig.TouchPanelBackground, fmt.Sprintf("knob %d app `%s`", i, appName))
				}

				for j, layer := range knobConfig.Layers {
					if layer.Global {
						continue
					}
					go hp.stopHandler(layer.HandlerStruct, layer.Handler, fmt.Sprintf("knob %d app %s layer %d", i, appName, j))
				}
			}
		}
	})
}

func (hp *HandlerPruner) StopAllHandlers() {
	if hp.vdev.Config().KeyGridBackgroundHandler != nil {
		hp.vdev.Config().KeyGridBackgroundHandler.Stop()
	}
	if hp.vdev.Config().TouchPanelBackgroundHandler != nil {
		hp.vdev.Config().TouchPanelBackgroundHandler.Stop()
	}
	for page := range hp.vdev.Config().Pages {
		hp.stopPageHandlers(page)
	}
}

func (hp *HandlerPruner) stopHandler(handler api.Handler, name, context string) {
	if handler == nil {
		return
	}
	visHandler, ok := handler.(api.VisualHandler)
	if !ok {
		return
	}
	if !visHandler.IsRunning() {
		return
	}
	hp.vdev.Logger().Printf("Stopping handler: %s on %s\n", name, context)
	visHandler.Stop()
	hp.vdev.Logger().Printf("Stopped handler: %s on %s\n", name, context)
}
