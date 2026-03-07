# StreamDeck API for Unix

This is a Go library for the streamdeckd daemon, not a standalone application. It provides interfaces for connecting to the daemon, accessing configuration objects, and creating custom handlers or GUI config editors for Elgato StreamDeck devices on Unix systems.

The library exposes an API for creating custom plugins to handle Stream Deck inputs (buttons, knobs, touch), managing icons and images, and also for communicating with the streamdeckd daemon via DBus.

## Installation

```bash
go get github.com/unix-streamdeck/api
```

## Usage

### Connecting to the Stream Deck daemon

```go
import (
    "fmt"
    "github.com/unix-streamdeck/api"
)

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

// Listen for page changes
err = conn.RegisterPageListener(func(serial string, page int32) {
    fmt.Printf("Device %s changed to page %d\n", serial, page)
})
```

### Working with images

```go
import (
    "github.com/unix-streamdeck/api"
    "image"
    _ "image/png"
    "os"
)

// Load an image
file, _ := os.Open("icon.png")
img, _, _ := image.Decode(file)

// Resize image to fit a Stream Deck key, ideally pass the `IconSize` field from `StreamDeckInfoV1` rather than magic num,ber
resizedImg := api.ResizeImage(img, 72) // 72x72 pixels

// Resize image with specific width and height (e.g., for LCD display) (`LcdWidth` and `LcdHeight` in `StreamDeckInfoV1`)
resizedLcdImg := api.ResizeImageWH(img, 800, 100)

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
    // Clean up resources and stop calling callback
}

// Implement an LCD handler (for Stream Deck Plus)
type MyLcdHandler struct {
    running bool
}

func (h *MyLcdHandler) Start(knob api.KnobConfigV3, info api.StreamDeckInfoV1, callback func(image image.Image)) {
    h.running = true
    // Generate and update LCD image
}

func (h *MyLcdHandler) IsRunning() bool {
    return h.running
}

func (h *MyLcdHandler) SetRunning(running bool) {
    h.running = running
}

func (h *MyLcdHandler) Stop() {
    h.running = false
}

// Implement a knob or touch handler (for Stream Deck Plus)
type MyKnobHandler struct{}

func (h *MyKnobHandler) Input(knob api.KnobConfigV3, info api.StreamDeckInfoV1, event api.InputEvent) {
    switch event.EventType {
    case api.KNOB_CW:
        // Handle clockwise rotation
    case api.KNOB_CCW:
        // Handle counter-clockwise rotation
    case api.KNOB_PRESS:
        // Handle knob press
    }
}
```

## API Documentation

The API provides several interfaces for handling Stream Deck interactions:

- `Handler`: Base interface for all handlers
- `IconHandler`: For handling dynamic icons/images
- `KeyHandler`: For handling key press events
- `LcdHandler`: For handling LCD displays (Stream Deck Plus)
- `KnobOrTouchHandler`: For handling knob rotations and touch events (Stream Deck Plus)

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


// Commit configuration changes to disk
err = conn.CommitConfig()

// Or Reload configuration from disk if the changes weren't correct
err = conn.ReloadConfig()

```


## License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.
