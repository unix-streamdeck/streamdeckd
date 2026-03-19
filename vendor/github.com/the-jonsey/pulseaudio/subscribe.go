package pulseaudio

import (
	goErrors "errors"
	"fmt"
	"strings"
)

type SubscriptionMask uint32

const (
	SubscriptionMaskSink         SubscriptionMask = 0x0001
	SubscriptionMaskSource       SubscriptionMask = 0x0002
	SubscriptionMaskSinkInput    SubscriptionMask = 0x0004
	SubscriptionMaskSourceOutput SubscriptionMask = 0x0008
	SubscriptionMaskModule       SubscriptionMask = 0x0010
	SubscriptionMaskClient       SubscriptionMask = 0x0020
	SubscriptionMaskSampleCache  SubscriptionMask = 0x0040
	SubscriptionMaskServer       SubscriptionMask = 0x0080
	SubscriptionMaskCard         SubscriptionMask = 0x0200
	SubscriptionMaskAll          SubscriptionMask = 0x02ff
)

func StrToSubscriptionMask(s string) (SubscriptionMask, error) {
	var mask SubscriptionMask
	parts := strings.Split(s, ",")
	for _, part := range parts {
		part = strings.TrimSpace(strings.ToLower(part))
		switch part {
		case "sink":
			mask |= SubscriptionMaskSink
		case "source":
			mask |= SubscriptionMaskSource
		case "sinkinput":
			mask |= SubscriptionMaskSinkInput
		case "sourceoutput":
			mask |= SubscriptionMaskSourceOutput
		case "module":
			mask |= SubscriptionMaskModule
		case "client":
			mask |= SubscriptionMaskClient
		case "samplecache":
			mask |= SubscriptionMaskSampleCache
		case "server":
			mask |= SubscriptionMaskServer
		case "card":
			mask |= SubscriptionMaskCard
		case "all":
			mask |= SubscriptionMaskAll
		default:
			return 0, goErrors.New("unknown SubscriptionMask: " + part)
		}
	}
	return mask, nil
}

// ////////////////////////////////////////////////////////
type SubscriptionEventFacility uint32

const (
	FacilitySink         SubscriptionEventFacility = 0
	FacilitySource       SubscriptionEventFacility = 1
	FacilitySinkInput    SubscriptionEventFacility = 2
	FacilitySourceOutput SubscriptionEventFacility = 3
	FacilityModule       SubscriptionEventFacility = 4
	FacilityClient       SubscriptionEventFacility = 5
	FacilitySampleCache  SubscriptionEventFacility = 6
	FacilityServer       SubscriptionEventFacility = 7
	FacilityAutoload     SubscriptionEventFacility = 8
	FacilityCard         SubscriptionEventFacility = 9
)

func (f SubscriptionEventFacility) String() string {
	switch f {
	case FacilitySink:
		return "Sink"
	case FacilitySource:
		return "Source"
	case FacilitySinkInput:
		return "SinkInput"
	case FacilitySourceOutput:
		return "SourceOutput"
	case FacilityModule:
		return "Module"
	case FacilityClient:
		return "Client"
	case FacilitySampleCache:
		return "SampleCache"
	case FacilityServer:
		return "Server"
	case FacilityAutoload:
		return "Autoload"
	case FacilityCard:
		return "Card"
	default:
		return fmt.Sprintf("UnknownFacility(%d)", f)
	}
}

// ////////////////////////////////////////////////////////
type SubscriptionEventType uint32

const (
	EventTypeNew     SubscriptionEventType = 0x00
	EventTypeChanged SubscriptionEventType = 0x10
	EventTypeRemoved SubscriptionEventType = 0x20
)

func (t SubscriptionEventType) String() string {
	switch t {
	case EventTypeNew:
		return "New"
	case EventTypeChanged:
		return "Changed"
	case EventTypeRemoved:
		return "Removed"
	default:
		return fmt.Sprintf("UnknownType(%d)", t)
	}
}

type SubscriptionEvent struct {
	EventFacility string
	EventType     string
	Index         *uint32
}

// ////////////////////////////////////////////////////////

// Deprecated. Use Subscribe() and read client.Events
func (c *Client) Updates() (updates <-chan SubscriptionEvent, err error) {
	_, err = c.request(commandSubscribe, uint32Tag, uint32(SubscriptionMaskAll))
	if err != nil {
		return nil, err
	}
	return c.Events, nil
}

// All events will be sent to client.Updates channel
func (c *Client) Subscribe(mask SubscriptionMask) (err error) {
	_, err = c.request(commandSubscribe, uint32Tag, uint32(mask))
	return err
}
