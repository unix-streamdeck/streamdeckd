# pulseaudio [![GoDoc](https://godoc.org/github.com/the-jonsey/pulseaudio?status.svg)](https://godoc.org/github.com/the-jonsey/pulseaudio)
Package pulseaudio is a pure-Go (no libpulse) implementation of the PulseAudio native protocol.

Download:
```shell
go get github.com/the-jonsey/pulseaudio
```

* * *
Package pulseaudio is a pure-Go (no libpulse) implementation of the PulseAudio native protocol.

This library is a fork of https://github.com/mafik/pulseaudio
The original library deliberately tries to hide pulseaudio internals and doesn't expose them.

Rather than exposing the PulseAudio protocol directly this library attempts to hide
the PulseAudio complexity behind a Go interface.
Some of the things which are deliberately not exposed in the API are:

→ backwards compatibility for old PulseAudio servers

→ transport mechanism used for the connection (Unix sockets / memfd / shm)

→ encoding used in the pulseaudio-native protocol

→ wors with pipewire as long as pipewire-pulse is installed and running

## Working features
Querying and setting the volume.

Querying and setting mute.

Listing audio sinks/sources/outputs/inputs.

Changing the default audio output.

Notifications on config updates.

Filtering config update notifications by event type.