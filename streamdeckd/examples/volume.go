package examples

import (
	"context"
	"errors"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/the-jonsey/pulseaudio"
	"github.com/unix-streamdeck/api/v2"
	"golang.org/x/sync/semaphore"
)

type VolumeHandler struct {
	Running    bool
	Quit       chan bool
	Lock       *semaphore.Weighted
	MuteBuff   image.Image
	UnmuteBuff image.Image
	DevType    string
	InputName  string
	Props      map[string]string
	Mute       bool
	Volume     int
	FirstLoop  bool
	client     *pulseaudio.Client
}

func (v *VolumeHandler) Start(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1, callback func(image image.Image)) {

	if v.Quit == nil {
		v.Quit = make(chan bool)
	}
	if v.Lock == nil {
		v.Lock = semaphore.NewWeighted(1)
	}
	if v.MuteBuff == nil {
		v.MuteBuff = v.GetImage("mute_icon", fields, handlerType, info)
	}
	if v.UnmuteBuff == nil {
		v.UnmuteBuff = v.GetImage("unmute_icon", fields, handlerType, info)
	}
	devType, ok := fields["device_type"]
	if !ok {
		devType, ok = fields["device_type"]
		if !ok {
			log.Println("Device type missing")
			return
		}
	}
	if devType == "sink_input" || devType == "source_output" {
		var inputName string
		props := make(map[string]string)
		inputName, ok := fields["input_name"].(string)
		if !ok {
			propss, ok := fields["props"]
			if !ok {
				log.Println("No Input Name or Props")
				return
			}
			for key, value := range propss.(map[string]interface{}) {
				if value == nil {
					continue
				}
				props[key] = value.(string)
			}
		}
		v.InputName = inputName
		v.Props = props
	}
	v.DevType = devType.(string)
	v.Running = true
	v.Run(handlerType, info, callback)
}
func (v *VolumeHandler) IsRunning() bool {
	return v.Running
}

func (v *VolumeHandler) SetRunning(running bool) {
	v.Running = running
}

func (v *VolumeHandler) Stop() {
	v.Running = false
	v.Quit <- true
	v.FirstLoop = true
	v.Mute = false
	v.Volume = 0
}

func (v *VolumeHandler) GetImage(index string, fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1) image.Image {
	path, ok := fields[index]
	w, h := info.GetDimensions(handlerType)
	if !ok {
		log.Println("image missing: " + index)
		return image.NewNRGBA(image.Rect(0, 0, w, h))
	}
	f, err := os.Open(path.(string))
	defer f.Close()
	if err != nil {
		log.Println(err)
		return image.NewNRGBA(image.Rect(0, 0, w, h))
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Println(err)
		return image.NewNRGBA(image.Rect(0, 0, w, h))
	}
	return api.ResizeImageWH(img, w, h)
}

func (v *VolumeHandler) Run(handlerType api.HandlerType, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	ctx := context.Background()
	err := v.Lock.Acquire(ctx, 1)
	if err != nil {
		return
	}
	defer v.Lock.Release(1)
	var subscriptionMask pulseaudio.SubscriptionMask
	if v.DevType == "sink" {
		subscriptionMask = pulseaudio.SubscriptionMaskSink
	} else if v.DevType == "source" {
		subscriptionMask = pulseaudio.SubscriptionMaskSource
	} else if v.DevType == "sink_input" {
		subscriptionMask = pulseaudio.SubscriptionMaskSinkInput
	} else if v.DevType == "source_output" {
		subscriptionMask = pulseaudio.SubscriptionMaskSourceOutput
	}
	err = v.client.Subscribe(subscriptionMask)
	if err != nil {
		return
	}
	err = update(v, handlerType, info, callback)
	if err != nil {
		log.Println(err)
	}
	for {
		select {
		case <-v.Quit:
			log.Println("Volume Quit")
			return
		case <-v.client.Events:
			err := update(v, handlerType, info, callback)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func update(v *VolumeHandler, handlerType api.HandlerType, info api.StreamDeckInfoV1, callback func(image image.Image)) error {
	var device pulseaudio.Device
	var err error
	w, h := info.GetDimensions(handlerType)
	if v.DevType == "sink" {
		device, err = v.client.GetDefaultSink()
	} else if v.DevType == "source" {
		device, err = v.client.GetDefaultSource()
	} else if v.DevType == "sink_input" {
		var devices []pulseaudio.SinkInput
		if v.InputName == "" {
			devices, err = v.client.GetSinkInputsByProps(v.Props)
		} else {
			devices, err = v.client.GetSinkInputsByName(v.InputName)
		}
		if len(devices) < 1 {
			err = errors.New("No Device Found")
		} else {
			device = devices[0]
		}
	} else if v.DevType == "source_output" {
		var devices []pulseaudio.SourceOutput
		if v.InputName == "" {
			devices, err = v.client.GetSourceOutputsByProps(v.Props)
		} else {
			devices, err = v.client.GetSourceOutputsByName(v.InputName)
		}
		if len(devices) < 1 {
			err = errors.New("No Device Found")
		} else {
			device = devices[0]
		}
	}
	if err != nil {
		img := image.NewNRGBA(image.Rect(0, 0, w, h))
		var text string
		if v.DevType == "sink" || v.DevType == "source" {
			text = "Could not find default " + v.DevType
		} else if v.DevType == "sink_input" {
			text = "Could not find sink input"
		} else if v.DevType == "source_output" {
			text = "Could not find source output"
		}
		imgParsed, err2 := api.DrawText(img, text, api.DrawTextOptions{
			FontSize:          24,
			VerticalAlignment: api.Center,
		})
		if err2 != nil {
			log.Println(err2)
		} else {
			callback(imgParsed)
		}
		return err
	} else {
		var text string
		mute := device.IsMute()
		if mute == true && mute == v.Mute && !v.FirstLoop {
			return nil
		}
		v.Mute = mute
		var img image.Image
		var vol int

		vol = int(math.Round(float64(device.GetVolume()) * 100))
		if vol == v.Volume && !v.FirstLoop {
			return nil
		}
		v.Volume = vol

		if mute {
			text = "Muted"
			img = v.MuteBuff
		} else {
			text = strconv.Itoa(vol) + "%"
			img = v.UnmuteBuff
		}

		if img == nil {
			image.NewNRGBA(image.Rect(0, 0, w, h))
		}

		imgParsed, err := api.DrawProgressBarWithAccent(img, text, 5, float64(h-25), 20, float64(w-10), float64(vol), "#cc3333")

		if err != nil {
			log.Println(err)
		} else {
			callback(imgParsed)
		}
		return nil
	}
}

func (v *VolumeHandler) Input(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfoV1, event api.InputEvent) {
	dev, ok := fields["device_type"]
	if !ok {
		dev, ok = fields["device_type"]
		if !ok {
			log.Println("Device type missing")
			return
		}
	}
	var err error
	if dev == "sink_input" || dev == "source_output" {
		var inputName string
		props := make(map[string]string)
		inputName, ok := fields["input_name"].(string)
		if !ok {
			propss, ok := fields["props"]
			if !ok {
				log.Println("No Input Name or Props")
				return
			}
			for key, value := range propss.(map[string]interface{}) {
				if value == nil {
					continue
				}
				props[key] = value.(string)
			}
		}
		if dev == "sink_input" {
			var devices []pulseaudio.SinkInput
			if inputName == "" {
				devices, err = v.client.GetSinkInputsByProps(props)
			} else {
				devices, err = v.client.GetSinkInputsByName(inputName)
			}
			if err != nil {
				log.Println(err)
				return
			}
			for _, device := range devices {
				updateDevice(device, event)
			}
		} else {
			var devices []pulseaudio.SourceOutput
			if inputName == "" {
				devices, err = v.client.GetSourceOutputsByProps(props)
			} else {
				devices, err = v.client.GetSourceOutputsByName(inputName)
			}
			if err != nil {
				log.Println(err)
				return
			}
			for _, device := range devices {
				updateDevice(device, event)
			}
		}

	} else {
		var device pulseaudio.Device
		if dev == "sink" {
			device, err = v.client.GetDefaultSink()
		} else if dev == "source" {
			device, err = v.client.GetDefaultSource()
		}
		if device == nil {
			err = errors.New("No device found")
		}
		if err != nil {
			log.Println(err)
			return
		}
		updateDevice(device, event)
	}
}

func HexColor(hex string) color.RGBA {
	values, _ := strconv.ParseUint(string(hex[1:]), 16, 32)
	return color.RGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}

func updateDevice(device pulseaudio.Device, event api.InputEvent) {
	muted := device.IsMute()
	if event.EventType == api.KNOB_CCW {
		if !muted && device.GetVolume() > 0 {
			device.SetVolume(device.GetVolume() - (float32(event.RotateNotches) * 0.01))
		}
	} else if event.EventType == api.KNOB_CW {
		if !muted && device.GetVolume() < 1 {
			device.SetVolume(device.GetVolume() + (float32(event.RotateNotches) * 0.01))
		}
	} else if event.EventType == api.KNOB_PRESS || event.EventType == api.SCREEN_SHORT_TAP || event.EventType == api.KEY_PRESS {
		device.ToggleMute()
	}
}

func RegisterVolume() api.Module {
	return api.Module{
		NewForeground: func() api.ForegroundHandler {
			client, err := pulseaudio.NewClient()
			if err != nil {
				panic(err)
			}
			return &VolumeHandler{Running: true, Lock: semaphore.NewWeighted(1), FirstLoop: true, client: client}
		},
		ForegroundFields: []api.Field{
			{Title: "Unmuted Icon", Name: "unmute_icon", Type: api.File, FileTypes: []string{".png", ".jpg", ".jpeg"}},
			{Title: "Muted Icon", Name: "mute_icon", Type: api.File},
			{Title: "Device Type", Name: "device_type", Type: api.Text},
			{Title: "Input Name", Name: "input_name", Type: api.Text},
			{Title: "Props", Name: "props", Type: api.Text},
		},
		NewInput: func() api.InputHandler {
			client, err := pulseaudio.NewClient()
			if err != nil {
				panic(err)
			}
			return &VolumeHandler{client: client}
		},
		InputFields: []api.Field{
			{Title: "Device Type", Name: "device_type", Type: api.Text},
			{Title: "Input Name", Name: "input_name", Type: api.Text},
			{Title: "Props", Name: "props", Type: api.Text},
		},
		Name: "Volume",
	}
}
