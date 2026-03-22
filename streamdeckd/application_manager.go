package streamdeckd

import "log"

type IApplicationManager interface {
	SetApplication(application string)
	AttachListener(listener func(application string))
	GetApplication() string
}

type ApplicationManager struct {
	listeners         []func(application string)
	activeApplication string
}

func (am *ApplicationManager) SetApplication(application string) {

	if am.activeApplication != application {
		log.Println("Application updated to: " + application)
		am.activeApplication = application
		for _, listener := range am.listeners {
			go listener(application)
		}
	}
}

func (am *ApplicationManager) AttachListener(listener func(application string)) {
	am.listeners = append(am.listeners, listener)
}

func (am *ApplicationManager) GetApplication() string {
	return am.activeApplication
}
