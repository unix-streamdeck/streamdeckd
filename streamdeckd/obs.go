package streamdeckd

import (
	"errors"
	obsws "github.com/christopher-dG/go-obs-websocket"
	"github.com/unix-streamdeck/api/v2"
	"log"
	"strconv"
)

var paramlessObsCommands = map[string]func() obsws.Request{
	"StartStopRecording": func() obsws.Request {
		request := obsws.NewStartStopRecordingRequest()
		return &request
	},
	"StartRecording": func() obsws.Request {
		request := obsws.NewStartRecordingRequest()
		return &request
	},
	"StopRecording": func() obsws.Request {
		request := obsws.NewStopRecordingRequest()
		return &request
	},
	"PauseRecording": func() obsws.Request {
		request := obsws.NewPauseRecordingRequest()
		return &request
	},
	"ResumeRecording": func() obsws.Request {
		request := obsws.NewResumeRecordingRequest()
		return &request
	},
	"StartStopReplayBuffer": func() obsws.Request {
		request := obsws.NewStartStopReplayBufferRequest()
		return &request
	},
	"StartReplayBuffer": func() obsws.Request {
		request := obsws.NewStartReplayBufferRequest()
		return &request
	},
	"StopReplayBuffer": func() obsws.Request {
		request := obsws.NewStopReplayBufferRequest()
		return &request
	},
	"SaveReplayBuffer": func() obsws.Request {
		request := obsws.NewSaveReplayBufferRequest()
		return &request
	},
	"EnableStudioMode": func() obsws.Request {
		request := obsws.NewEnableStudioModeRequest()
		return &request
	},
	"DisableStudioMode": func() obsws.Request {
		request := obsws.NewDisableStudioModeRequest()
		return &request
	},
	"ToggleStudioMode": func() obsws.Request {
		request := obsws.NewToggleStudioModeRequest()
		return &request
	},
	"StartStopStreaming": func() obsws.Request {
		request := obsws.NewStartStopStreamingRequest()
		return &request
	},
	"StopStreaming": func() obsws.Request {
		request := obsws.NewStopStreamingRequest()
		return &request
	},
}
var obs obsws.Client

func tryConnectObs() error {
	if config.ObsConnectionInfo.Host != "" && config.ObsConnectionInfo.Port != 0 {
		log.Println("Found connection info")
		obs = obsws.Client{Host: config.ObsConnectionInfo.Host, Port: config.ObsConnectionInfo.Port}
		log.Println("Attempting connection")
		if err := obs.Connect(); err != nil {
			return err
		}
	} else {
		return errors.New("No Obs Connection Info Provided")
	}
	return nil
}

func runObsCommand(command string, params map[string]string) {
	if !obs.Connected() {
		err := tryConnectObs()
		if err != nil {
			log.Println(err)
		}
	}
	var req obsws.Request
	requestLambda, exists := paramlessObsCommands[command]

	if exists {
		req = requestLambda()
	} else if command == "SetVolume" {
		volString, exists := params["volume"]
		if !exists {
			log.Println("No volume parameter set")
			return
		}
		vol, err := strconv.Atoi(volString)
		if err != nil {
			log.Println(err)
			return
		}
		src, exists := params["source"]
		if !exists {
			log.Println("No source parameter set")
			return
		}
		request := obsws.NewSetVolumeRequest(src, float64(vol))
		req = &request
	} else if command == "SetMute" {
		src, exists := params["source"]
		if !exists {
			log.Println("No source parameter set")
			return
		}
		mute, exists := params["mute"]
		if !exists {
			log.Println("No mute parameter set")
			return
		}
		request := obsws.NewSetMuteRequest(src, mute == "true")
		req = &request
	} else if command == "ToggleMute" {
		src, exists := params["source"]
		if !exists {
			log.Println("No source parameter set")
			return
		}
		getRequest := obsws.NewGetMuteRequest(src)
		if err := getRequest.Send(obs); err != nil {
			log.Println(err)
		} else {
			resp, err := getRequest.Receive()
			if err != nil {
				log.Println(err)
			} else {
				request := obsws.NewSetMuteRequest(src, !resp.Muted)
				req = &request
			}
		}
	}
	if err := req.Send(obs); err != nil {
		log.Println(err)
	}
}

func getObsHandlerFields() ([]api.Module, error) {
	var modules []api.Module
	for key := range paramlessObsCommands {
		modules = append(modules, api.Module{
			Name:   key,
			IsIcon: false,
			IsKey:  true,
		})
	}

	var sources []string

	if !obs.Connected() {
		err := tryConnectObs()
		if err != nil {
			return nil, err
		}
	}

	request := obsws.NewGetSpecialSourcesRequest()
	if err := request.Send(obs); err != nil {
		return nil, err
	} else {
		resp, err := request.Receive()
		if err != nil {
			return nil, err
		} else {
			for _, source := range []string{resp.Desktop1, resp.Desktop2, resp.Mic1, resp.Mic2, resp.Mic3} {
				if source != "" {
					sources = append(sources, source)
				}
			}
		}
	}

	modules = append(modules, api.Module{
		Name:   "SetVolume",
		IsIcon: false,
		IsKey:  true,
		KeyFields: []api.Field{
			{Title: "Volume", Name: "volume", Type: "Number"},
			{Title: "Audio Source", Name: "source", Type: "List", ListItems: sources},
		},
	})

	modules = append(modules, api.Module{
		Name:   "SetMute",
		IsIcon: false,
		IsKey:  true,
		KeyFields: []api.Field{
			{Title: "Mute", Name: "mute", Type: "Checkbox"},
			{Title: "Audio Source", Name: "source", Type: "List", ListItems: sources},
		},
	})

	modules = append(modules, api.Module{
		Name:   "ToggleMute",
		IsIcon: false,
		IsKey:  true,
		KeyFields: []api.Field{
			{Title: "Mute", Name: "mute", Type: "Checkbox"},
			{Title: "Audio Source", Name: "source", Type: "List", ListItems: sources},
		},
	})

	return modules, nil
}
