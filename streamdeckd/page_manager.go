package streamdeckd

type PageManager struct {
	vdev      *VirtualDev
	page      int
	listeners []func(newPage, previousPage int)
}

func (pm *PageManager) SetPage(page int) {

	if page != pm.page || page == 0 {

		oldPage := pm.page

		pm.page = page

		pm.vdev.sdInfo.Page = page
		EmitPage(pm.vdev, page)

		for _, listener := range pm.listeners {
			go listener(pm.page, oldPage)
		}
	}
}

func (pm *PageManager) AttachListener(channel func(newPage, previousPage int)) {
	pm.listeners = append(pm.listeners, channel)
}
