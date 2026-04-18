package streamdeckd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unix-streamdeck/api/v2"
	"go.uber.org/mock/gomock"
)

func TestGetPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockVDev := NewMockIVirtualDev(ctrl)
	assertions := assert.New(t)

	pm := PageManager{
		vdev: mockVDev,
	}

	mockSDInfo := &api.StreamDeckInfoV1{}
	mockVDev.EXPECT().SdInfo().Return(mockSDInfo).Times(1)

	assertions.Equal(0, pm.GetPage())

	pm.SetPage(3)

	assertions.Equal(3, pm.GetPage())
	assertions.Equal(3, mockSDInfo.Page)
}

func TestPageListener(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockVDev := NewMockIVirtualDev(ctrl)
	assertions := assert.New(t)

	pm := PageManager{
		vdev: mockVDev,
	}

	mockSDInfo := &api.StreamDeckInfoV1{}
	mockVDev.EXPECT().SdInfo().Return(mockSDInfo).Times(1)

	spyChan := make(chan []int)
	pm.AttachListener(func(newPage, previousPage int) {
		spyChan <- []int{newPage, previousPage}
	})

	pm.SetPage(3)

	select {

	case args := <-spyChan:
		assertions.Equal(3, args[0])
		assertions.Equal(0, args[1])

	case <-time.After(5 * time.Second):
		assertions.Fail("callback function not called")
	}
}
