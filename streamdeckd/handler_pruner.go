package streamdeckd

import (
	"fmt"

	"github.com/unix-streamdeck/api/v2"
)

type HandlerPruner struct {
	vdev *VirtualDev
}

func (hp *HandlerPruner) OnPageChange() {
	hp.vdev.pageManager.AttachListener(func(_, previousPage int) {
		hp.StopPageHandlers(previousPage)
	})
}

func (hp *HandlerPruner) StopPageHandlers(pageNo int) {
	page := hp.vdev.Config.Pages[pageNo]

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

		for _, keyConfig := range key.Application {
			if keyConfig.KeyBackgroundHandler != nil {
				go hp.stopHandler(keyConfig.KeyBackgroundHandler, keyConfig.KeyBackground, fmt.Sprintf("page %d, key %d background", pageNo, i))
			}

			if keyConfig.IconHandlerStruct != nil {
				go hp.stopHandler(keyConfig.IconHandlerStruct, keyConfig.IconHandler, fmt.Sprintf("page %d, key %d", pageNo, i))
			}
		}
	}

	for i, knob := range page.Knobs {

		if knob.TouchPanelBackgroundHandler != nil {
			go hp.stopHandler(knob.TouchPanelBackgroundHandler, knob.TouchPanelBackground, fmt.Sprintf("page %d, knob %d background", pageNo, i))
		}

		for _, knobConfig := range knob.Application {

			if knobConfig.TouchPanelBackgroundHandler != nil {
				go hp.stopHandler(knobConfig.TouchPanelBackgroundHandler, knobConfig.TouchPanelBackground, fmt.Sprintf("page %d, knob %d background", pageNo, i))
			}

			if knobConfig.LcdHandlerStruct != nil {
				go hp.stopHandler(knobConfig.LcdHandlerStruct, knobConfig.LcdHandler, fmt.Sprintf("page %d, knob %d", pageNo, i))
			}
		}
	}
}

func (hp *HandlerPruner) OnAppSwitch() {
	applicationManager.AttachListener(func(activeApp string) {
		page := hp.vdev.Config.Pages[hp.vdev.pageManager.page]

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

				if keyConfig.IconHandlerStruct != nil {
					go hp.stopHandler(keyConfig.IconHandlerStruct, keyConfig.IconHandler, fmt.Sprintf("key %d app `%s`", i, appName))
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

				if knobConfig.LcdHandlerStruct != nil {
					go hp.stopHandler(knobConfig.LcdHandlerStruct, knobConfig.LcdHandler, fmt.Sprintf("knob %d app `%s`", i, appName))
				}
			}
		}
	})
}

func (hp *HandlerPruner) StopAllHandlers() {
	for page := range hp.vdev.Config.Pages {
		hp.StopPageHandlers(page)
	}
}

func (hp *HandlerPruner) stopHandler(handler api.VisualHandler, name, context string) {
	if !handler.IsRunning() {
		return
	}
	hp.vdev.logger.Printf("Stopping handler: %s on %s\n", name, context)
	handler.Stop()
	hp.vdev.logger.Printf("Stopped handler: %s on %s\n", name, context)
}
