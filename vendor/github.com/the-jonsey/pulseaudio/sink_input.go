package pulseaudio

import (
    "io"
    "math"
    "strings"
)

// SinkInput contains information about a sink in pulseaudio
type SinkInput struct {
    Index          uint32
    Name           string
    OwnerModule    uint32
    ClientIndex    uint32
    Sink           uint32
    SampleSpec     sampleSpec
    ChannelMap     channelMap
    Cvolume        cvolume
    BufferUsec     uint64
    SinkUsec       uint64
    ResampleMethod string
    Driver         string
    Muted          bool
    PropList       map[string]string
    Corked         bool
    HasVolume      bool
    VolumeWritable bool
    Format         formatInfo
    Client         *Client
}

// ReadFrom deserializes a sink packet from pulseaudio
func (s *SinkInput) ReadFrom(r io.Reader) (int64, error) {
    err := bread(r,
        uint32Tag, &s.Index,
        stringTag, &s.Name,
        uint32Tag, &s.OwnerModule,
        uint32Tag, &s.ClientIndex,
        uint32Tag, &s.Sink,
        &s.SampleSpec,
        &s.ChannelMap,
        &s.Cvolume,
        usecTag, &s.BufferUsec,
        usecTag, &s.SinkUsec,
        stringTag, &s.ResampleMethod,
        stringTag, &s.Driver,
        &s.Muted,
        &s.PropList,
        &s.Corked,
        &s.HasVolume,
        &s.VolumeWritable)
    if err != nil {
        return 0, err
    }
    err = bread(r, &s.Format)
    return 0, nil
}

func (s SinkInput) SetVolume(volume float32) error {
    _, err := s.Client.request(commandSetSinkInputVolume, uint32Tag, s.Index, cvolume{uint32(volume * 0xffff)})
    return err
}

func (s SinkInput) SetMute(b bool) error {
    muteCmd := '0'
    if b {
        muteCmd = '1'
    }
    _, err := s.Client.request(commandSetSinkInputMute, uint32Tag, s.Index, uint8(muteCmd))
    return err
}

func (s SinkInput) ToggleMute() error {
    return s.SetMute(!s.Muted)
}

func (s SinkInput) IsMute() bool {
    return s.Muted
}

func (s SinkInput) GetVolume() float32 {
    return float32(math.Round(float64(float32(s.Cvolume[0])/0xffff)*100)) / 100
}

// Sinks queries PulseAudio for a list of sinks and returns an array
func (c *Client) SinkInputs() ([]SinkInput, error) {
    b, err := c.request(commandGetSinkInputInfoList)
    if err != nil {
        return nil, err
    }
    var sinkInputs []SinkInput
    for b.Len() > 0 {
        var sinkInput SinkInput
        err = bread(b, &sinkInput)
        if err != nil {
            return nil, err
        }
        sinkInput.Client = c
        sinkInputs = append(sinkInputs, sinkInput)
    }
    return sinkInputs, nil
}

func (c *Client) GetSinkInputsByName(name string) ([]SinkInput, error) {
    sinkInputs, err := c.SinkInputs()
    if err != nil {
        return []SinkInput{}, err
    }
    var inputs []SinkInput
    for _, sinkInput := range sinkInputs {
        if strings.ToLower(sinkInput.Name) == strings.ToLower(name) {
            inputs = append(inputs, sinkInput)
        }
    }
    return inputs, nil
}

func (c *Client) GetSinkInputsByProps(props map[string]string) ([]SinkInput, error) {
    sinkInputs, err := c.SinkInputs()
    if err != nil {
        return []SinkInput{}, err
    }
    var inputs []SinkInput
    for _, sinkInput := range sinkInputs {
        for key, val := range props {
            inpVal, ok := sinkInput.PropList[key]
            if ok && strings.ToLower(inpVal) == strings.ToLower(val) {
                inputs = append(inputs, sinkInput)
            }
        }
        //if strings.ToLower(sinkInput.Name) == strings.ToLower(name) {
        //    inputs = append(inputs, sinkInput)
        //}
    }
    return inputs, nil
}
