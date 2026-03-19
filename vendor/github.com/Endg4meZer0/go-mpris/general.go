package mpris

import (
	"strings"

	"github.com/godbus/dbus/v5"
)

func getProperty(obj *dbus.Object, iface string, prop string) (dbus.Variant, error) {
	result := dbus.Variant{}
	err := obj.Call(getPropertyMethod, 0, iface, prop).Store(&result)
	if err != nil {
		return dbus.Variant{}, err
	}
	return result, nil
}

func setProperty(obj *dbus.Object, iface string, prop string, val interface{}) error {
	call := obj.Call(setPropertyMethod, 0, iface, prop, dbus.MakeVariant(val))
	return call.Err
}

// Connects to the player with the specified name using the specified DBus connection.
func New(conn *dbus.Conn, name string) *Player {
	obj := conn.Object(name, dbusObjectPath).(*dbus.Object)

	return &Player{conn, obj, name}
}

// Lists the available players in alphabetical order.
func List(conn *dbus.Conn) ([]string, error) {
	var names []string
	err := conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)
	if err != nil {
		return nil, err
	}

	var mprisNames []string
	for _, name := range names {
		if strings.HasPrefix(name, BaseInterface) {
			mprisNames = append(mprisNames, name)
		}
	}
	return mprisNames, nil
}

func RegisterNameOwnerChanged(conn *dbus.Conn, ch chan<- *dbus.Signal) (err error) {
	// Add NameOwnerChanged handler
	err = conn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.DBus"),
		dbus.WithMatchMember("NameOwnerChanged"),
	)
	if err != nil {
		return
	}

	conn.Signal(ch)
	return nil
}

func UnregisterNameOwnerChanged(conn *dbus.Conn, ch chan<- *dbus.Signal) (err error) {
	// Add NameOwnerChanged handler
	err = conn.RemoveMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.DBus"),
		dbus.WithMatchMember("NameOwnerChanged"),
	)
	if err != nil {
		return
	}

	conn.RemoveSignal(ch)
	return nil
}

// Registers a new signal receiver channel that will be able to get signals as specified in SignalType definition.
func (i *Player) RegisterSignalReceiver(ch chan<- *dbus.Signal) (err error) {
	// Add PropertiesChanged handler
	err = i.conn.AddMatchSignal(
		dbus.WithMatchSender(i.name),
		dbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		dbus.WithMatchMember("PropertiesChanged"),
	)
	if err != nil {
		return
	}

	// Add Seeked handler
	// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Signal:Seeked
	err = i.conn.AddMatchSignal(
		dbus.WithMatchSender(i.name),
		dbus.WithMatchInterface(PlayerInterface),
		dbus.WithMatchMember("Seeked"),
	)
	if err != nil {
		return
	}

	// Add TrackListReplaced handler
	// See also: https://specifications.freedesktop.org/mpris-spec/latest/Track_List_Interface.html#Signal:TrackListReplaced
	err = i.conn.AddMatchSignal(
		dbus.WithMatchInterface(TrackListInterface),
		dbus.WithMatchMember("TrackListReplaced"),
		dbus.WithMatchSender(i.name),
	)
	if err != nil {
		return
	}

	// Add TrackAdded handler
	// See also: https://specifications.freedesktop.org/mpris-spec/latest/Track_List_Interface.html#Signal:TrackAdded
	err = i.conn.AddMatchSignal(
		dbus.WithMatchInterface(TrackListInterface),
		dbus.WithMatchMember("TrackAdded"),
		dbus.WithMatchSender(i.name),
	)
	if err != nil {
		return
	}

	// Add TrackRemoved handler.
	// See also: https://specifications.freedesktop.org/mpris-spec/latest/Track_List_Interface.html#Signal:TrackRemoved
	err = i.conn.AddMatchSignal(
		dbus.WithMatchInterface(TrackListInterface),
		dbus.WithMatchMember("TrackRemoved"),
		dbus.WithMatchSender(i.name),
	)
	if err != nil {
		return
	}

	// Add TrackMetadataChanged handler.
	// See also: https://specifications.freedesktop.org/mpris-spec/latest/Track_List_Interface.html#Signal:TrackMetadataChanged
	err = i.conn.AddMatchSignal(
		dbus.WithMatchInterface(TrackListInterface),
		dbus.WithMatchMember("TrackMetadataChanged"),
		dbus.WithMatchSender(i.name),
	)
	if err != nil {
		return
	}

	i.conn.Signal(ch)
	return
}

func (i *Player) UnregisterSignalReceiver(ch chan *dbus.Signal) (err error) {
	// Remove PropertiesChanged handler
	err = i.conn.RemoveMatchSignal(
		dbus.WithMatchSender(i.name),
		dbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		dbus.WithMatchMember("PropertiesChanged"),
	)
	if err != nil {
		return
	}

	// Remove Seeked handler
	// See also: https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Signal:Seeked
	err = i.conn.RemoveMatchSignal(
		dbus.WithMatchSender(i.name),
		dbus.WithMatchInterface(PlayerInterface),
		dbus.WithMatchMember("Seeked"),
	)
	if err != nil {
		return
	}

	// Remove TrackListReplaced handler
	// See also: https://specifications.freedesktop.org/mpris-spec/latest/Track_List_Interface.html#Signal:TrackListReplaced
	err = i.conn.RemoveMatchSignal(
		dbus.WithMatchInterface(TrackListInterface),
		dbus.WithMatchMember("TrackListReplaced"),
		dbus.WithMatchSender(i.name),
	)
	if err != nil {
		return
	}

	// Remove TrackAdded handler
	// See also: https://specifications.freedesktop.org/mpris-spec/latest/Track_List_Interface.html#Signal:TrackAdded
	err = i.conn.RemoveMatchSignal(
		dbus.WithMatchInterface(TrackListInterface),
		dbus.WithMatchMember("TrackAdded"),
		dbus.WithMatchSender(i.name),
	)
	if err != nil {
		return
	}

	// Remove TrackRemoved handler.
	// See also: https://specifications.freedesktop.org/mpris-spec/latest/Track_List_Interface.html#Signal:TrackRemoved
	err = i.conn.RemoveMatchSignal(
		dbus.WithMatchInterface(TrackListInterface),
		dbus.WithMatchMember("TrackRemoved"),
		dbus.WithMatchSender(i.name),
	)
	if err != nil {
		return
	}

	// Remove TrackMetadataChanged handler.
	// See also: https://specifications.freedesktop.org/mpris-spec/latest/Track_List_Interface.html#Signal:TrackMetadataChanged
	err = i.conn.RemoveMatchSignal(
		dbus.WithMatchInterface(TrackListInterface),
		dbus.WithMatchMember("TrackMetadataChanged"),
		dbus.WithMatchSender(i.name),
	)
	if err != nil {
		return
	}

	i.conn.RemoveSignal(ch)
	return nil
}

// Gets the player full name (including base interface name).
func (i *Player) GetName() string {
	return i.name
}

// Gets the player short name (without the base interface name).
func (i *Player) GetShortName() string {
	return strings.ReplaceAll(i.name, BaseInterface+".", "")
}

// Gets the supported signal type from *dbus.Signal
func GetSignalType(signal *dbus.Signal) SignalType {
	switch signal.Name {
	case propertiesChangedSignal:
		return SignalPropertiesChanged
	case nameOwnerChangedSignal:
		return SignalNameOwnerChanged
	case seekedSignal:
		return SignalSeeked
	case trackListReplacedSignal:
		return SignalTrackListReplaced
	case trackAddedSignal:
		return SignalTrackAdded
	case trackRemovedSignal:
		return SignalTrackRemoved
	case trackMetadataChangedSignal:
		return SignalTrackMetadataChanged
	default:
		return SignalNotSupported
	}
}
