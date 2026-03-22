package streamdeckd

type IPageManager interface {
	SetPage(page int)
	AttachListener(channel func(newPage, previousPage int))
	GetPage() int
}

type PageManager struct {
	vdev      IVirtualDev
	page      int
	listeners []func(newPage, previousPage int)
}

func (pm *PageManager) SetPage(page int) {

	if page != pm.page || page == 0 {

		oldPage := pm.page

		pm.page = page

		pm.vdev.SdInfo().Page = page
		EmitPage(pm.vdev, page)

		for _, listener := range pm.listeners {
			go listener(pm.page, oldPage)
		}
	}
}

func (pm *PageManager) AttachListener(channel func(newPage, previousPage int)) {
	pm.listeners = append(pm.listeners, channel)
}

func (pm *PageManager) GetPage() int {
	return pm.page
}
