# Stream Deck API for Unix

This is a Go library for the streamdeckd daemon, not a standalone application. It provides interfaces for connecting to the daemon, accessing configuration objects, and creating custom handlers or GUI config editors for Elgato Stream Deck devices on Unix systems.

The library enables handling Stream Deck inputs (buttons, knobs, touch), managing icons and images, and communicating with the Stream Deck daemon via DBus.

## Features

- Connect to Stream Deck devices through the streamdeckd daemon
- Handle button presses, knob rotations, and touch inputs
- Create and manipulate images for Stream Deck displays
- Manage device configurations and pages
- Draw text on Stream Deck buttons with customizable fonts and alignments
- Resize images to fit Stream Deck displays
- OBS integration support

## Installation

```bash
go get github.com/unix-streamdeck/api
```

## Usage

### Connecting to the Stream Deck daemon

```go
import "github.com/unix-streamdeck/api"

// Connect to the Stream Deck daemon
conn, err := api.Connect()
if err != nil {
    // Handle error
}
defer conn.Close()

// Get information about connected Stream Deck devices
devices, err := conn.GetInfo()
if err != nil {
    // Handle error
}
```

### Working with images

```go
import (
    "github.com/unix-streamdeck/api"
    "image"
    _ "image/png" // Import for PNG support
    "os"
)

// Load an image
file, _ := os.Open("icon.png")
img, _, _ := image.Decode(file)

// Resize image to fit a Stream Deck key
resizedImg := api.ResizeImage(img, 72) // 72x72 pixels

// Add text to an image
imgWithText, _ := api.DrawText(resizedImg, "Hello", 0, "CENTER")
```

### Implementing handlers

```go
// Implement a key handler
type MyKeyHandler struct{}

func (h *MyKeyHandler) Key(key api.KeyConfigV3, info api.StreamDeckInfoV1) {
    // Handle key press
}

// Implement an icon handler
type MyIconHandler struct {
    running bool
}

func (h *MyIconHandler) Start(key api.KeyConfigV3, info api.StreamDeckInfoV1, callback func(image image.Image)) {
    h.running = true
    // Generate and update icon
    // Call callback with new images when needed
}

func (h *MyIconHandler) IsRunning() bool {
    return h.running
}

func (h *MyIconHandler) SetRunning(running bool) {
    h.running = running
}

func (h *MyIconHandler) Stop() {
    h.running = false
    // Clean up resources
}
```

## API Documentation

The API provides several interfaces for handling Stream Deck interactions:

- `Handler`: Base interface for all handlers
- `IconHandler`: For handling dynamic icons/images
- `KeyHandler`: For handling key press events
- `LcdHandler`: For handling LCD displays
- `KnobOrTouchHandler`: For handling knob rotations and touch events

Key components:

- `Connection`: Manages DBus communication with the Stream Deck daemon
- Image utilities: Functions for drawing text and resizing images
- Configuration management: Functions for getting and setting device configurations

### Configuration Objects

The library exposes configuration objects that can be used to interact with the streamdeckd daemon:

```go
// Get the current configuration
config, err := conn.GetConfig()
if err != nil {
    // Handle error
}

// Modify configuration
// ...

// Set the updated configuration
err = conn.SetConfig(config)
if err != nil {
    // Handle error
}

// Reload configuration from disk
err = conn.ReloadConfig()

// Commit configuration changes to disk
err = conn.CommitConfig()
```

### Custom GUI Config Editors

The library provides the `Module` and `Field` types that can be used to create custom GUI configuration editors:

```go
// Get available modules
modules, err := conn.GetModules()
if err != nil {
    // Handle error
}

// Get OBS-specific fields
obsFields, err := conn.GetObsFields()
if err != nil {
    // Handle error
}
```

## Help Wanted!

If you want to help with the development of streamdeckd and its related repos, either by submitting code, finding/fixing bugs, or just replying to issues, please join this discord server: https://discord.gg/nyhuVEJWMQ

## License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.
