package pulseaudio

import (
	errors2 "errors"
	"io"
    "math"
)

//Source contains information about a source in pulseaudio, e.g. a microphone
type Source struct {
    Index              uint32
    Name               string
    Description        string
    SampleSpec         sampleSpec
    ChannelMap         channelMap
    ModuleIndex        uint32
    Cvolume            cvolume
    Muted              bool
    MonitorSourceIndex uint32
    MonitorSourceName  string
    Latency            uint64
    Driver             string
    Flags              uint32
    PropList           map[string]string
    RequestedLatency   uint64
    BaseVolume         uint32
    SinkState          uint32
    NVolumeSteps       uint32
    CardIndex          uint32
    Ports              []sinkPort
    ActivePortName     string
    Formats            []formatInfo
    Client             *Client
}

//ReadFrom deserialized a PA source packet
func (s *Source) ReadFrom(r io.Reader) (int64, error) {
    var portCount uint32
    err := bread(r,
        uint32Tag, &s.Index,
        stringTag, &s.Name,
        stringTag, &s.Description,
        &s.SampleSpec,
        &s.ChannelMap,
        uint32Tag, &s.ModuleIndex,
        &s.Cvolume,
        &s.Muted,
        uint32Tag, &s.MonitorSourceIndex,
        stringTag, &s.MonitorSourceName,
        usecTag, &s.Latency,
        stringTag, &s.Driver,
        uint32Tag, &s.Flags,
        &s.PropList,
        usecTag, &s.RequestedLatency,
        volumeTag, &s.BaseVolume,
        uint32Tag, &s.SinkState,
        uint32Tag, &s.NVolumeSteps,
        uint32Tag, &s.CardIndex,
        uint32Tag, &portCount)
    if err != nil {
        return 0, err
    }
    s.Ports = make([]sinkPort, portCount)
    for i := uint32(0); i < portCount; i++ {
        err = bread(r, &s.Ports[i])
        if err != nil {
            return 0, err
        }
    }
    if portCount == 0 {
        err = bread(r, stringNullTag)
        if err != nil {
            return 0, err
        }
    } else {
        err = bread(r, stringTag, &s.ActivePortName)
        if err != nil {
            return 0, err
        }
    }

    var formatCount uint8
    err = bread(r,
        uint8Tag, &formatCount)
    if err != nil {
        return 0, err
    }
    s.Formats = make([]formatInfo, formatCount)
    for i := uint8(0); i < formatCount; i++ {
        err = bread(r, &s.Formats[i])
        if err != nil {
            return 0, err
        }
    }
    return 0, nil
}

func (s Source) SetVolume(volume float32) error {
    _, err := s.Client.request(commandSetSourceVolume, uint32Tag, uint32(0xffffffff), stringTag, []byte(s.Name), byte(0), cvolume{uint32(volume * 0xffff)})
    return err
}

func (s Source) SetMute(b bool) error {
    muteCmd := '0'
    if b {
        muteCmd = '1'
    }
    _, err := s.Client.request(commandSetSourceMute, uint32Tag, uint32(0xffffffff), stringTag, []byte(s.Name), byte(0), uint8(muteCmd))
    return err
}

func (s Source) ToggleMute() error {
    return s.SetMute(!s.Muted)
}

func (s Source) IsMute() bool {
    return s.Muted
}

func (s Source) GetVolume() float32 {
    return float32(math.Round(float64(float32(s.Cvolume[0])/0xffff) * 100)) / 100
}

// Sources queries pulseaudio for a list of all it's sources and returns an array of them
func (c *Client) Sources() ([]Source, error) {
    b, err := c.request(commandGetSourceInfoList)
    if err != nil {
        return nil, err
    }
    var sources []Source
    for b.Len() > 0 {
        var source Source
        err = bread(b, &source)
        if err != nil {
            return nil, err
        }
		source.Client = c
        sources = append(sources, source)
    }
    return sources, nil
}

func (c *Client) GetDefaultSource() (Source, error) {
	s, err := c.ServerInfo()
	if err != nil {
		return Source{}, err
	}
	sources, err := c.Sources()
	if err != nil {
		return Source{}, err
	}
	for _, source := range sources{
		if source.Name == s.DefaultSource {
			return source, nil
		}
	}
	return Source{}, errors2.New("Could not get default sink")
}