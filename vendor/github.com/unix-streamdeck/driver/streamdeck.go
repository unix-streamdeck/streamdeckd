package streamdeck

import (
    "bytes"
    "context"
    "encoding/binary"
    "errors"
    "fmt"
    "image"
    "image/color"
    "image/jpeg"
    "math"
    "sync"
    "time"

    "github.com/bearsh/hid"
    "golang.org/x/image/draw"
    "golang.org/x/sync/semaphore"
)

const (
    // 30 fps fade animation.
    fadeDelay = time.Second / 30
)

//nolint:revive
// Stream Deck Vendor & Product IDs.
const (
    VID_ELGATO              = 0x0fd9
    PID_STREAMDECK          = 0x0060
    PID_STREAMDECK_V2       = 0x006d
    PID_STREAMDECK_MK2      = 0x0080
    PID_STREAMDECK_MINI     = 0x0063
    PID_STREAMDECK_MINI_MK2 = 0x0090
    PID_STREAMDECK_XL       = 0x006c
    PID_STREAMDECK_PEDAL    = 0x0086
    PID_STREAMDECK_PLUS     = 0x0084
)

type InputEventType uint8

const (
    KNOB_CCW InputEventType = iota
    KNOB_CW
    KNOB_PRESS
    SCREEN_SHORT_TAP
    SCREEN_LONG_TAP
    SCREEN_SWIPE
    KEY_PRESS
    KEY_RELEASE
)

type InputEvent struct {
    EventType     InputEventType
    Index         uint8
    RotateNotches uint8
    ScreenX       uint16
    ScreenY       uint16
    ScreenEndX    uint16
    ScreenEndY    uint16
}

//nolint:revive
// Firmware command IDs.
var (
    c_REV1_FIRMWARE   = []byte{0x04}
    c_REV1_RESET      = []byte{0x0b, 0x63}
    c_REV1_BRIGHTNESS = []byte{0x05, 0x55, 0xaa, 0xd1, 0x01}

    c_REV2_FIRMWARE   = []byte{0x05}
    c_REV2_RESET      = []byte{0x03, 0x02}
    c_REV2_BRIGHTNESS = []byte{0x03, 0x08}
)

type IHidDevice interface {
    Open() (*hid.Device, error)
    Close() error
    GetFeatureReport(payload []byte) (int, error)
    SendFeatureReport(payload []byte) (int, error)
    Write(payload []byte) (int, error)
    Read(payload []byte) (int, error)
    GetManufacturer() (string, error)
    GetProduct() (string, error)
    GetSerial() (string, error)
}

type HidDevice struct {
    Device   *hid.Device
    info     hid.DeviceInfo
    WriteSem *semaphore.Weighted
}

func (h *HidDevice) Open() (*hid.Device, error) {
    dev, err := h.info.Open()
    if err == nil {
        h.Device = dev
    }
    return dev, err
}

func (h *HidDevice) Close() error {
    if h.Device == nil {
        return errors.New("Device not available")
    }
    err := h.Device.Close()
    return err
}

func (h *HidDevice) GetFeatureReport(payload []byte) (int, error) {
    ctx := context.Background()
    err := h.WriteSem.Acquire(ctx, 1)
    if err != nil {
        return 0, err
    }
    defer h.WriteSem.Release(1)
    if h.Device == nil {
        return 0, errors.New("Device not available")
    }
    report, err := h.Device.GetFeatureReport(payload)
    return report, err
}

func (h *HidDevice) SendFeatureReport(payload []byte) (int, error) {
    ctx := context.Background()
    err := h.WriteSem.Acquire(ctx, 1)
    if err != nil {
        return 0, err
    }
    defer h.WriteSem.Release(1)
    if h.Device == nil {
        return 0, errors.New("Device not available")
    }
    report, err := h.Device.SendFeatureReport(payload)
    return report, err
}

func (h *HidDevice) Write(payload []byte) (int, error) {
    ctx := context.Background()
    err := h.WriteSem.Acquire(ctx, 1)
    if err != nil {
        return 0, err
    }
    defer h.WriteSem.Release(1)
    if h.Device == nil {
        return 0, errors.New("Device not available")
    }

    write, err := h.Device.Write(payload)
    return write, err
}

func (h *HidDevice) Read(payload []byte) (int, error) {
    if h.Device == nil {
        return 0, errors.New("Device not available")
    }
    write, err := h.Device.Read(payload)
    return write, err
}

func (h *HidDevice) GetManufacturer() (string, error) {
    if h.Device == nil {
        return "", errors.New("Device not available")
    }
    return h.Device.Manufacturer, nil
}

func (h *HidDevice) GetProduct() (string, error) {
    if h.Device == nil {
        return "", errors.New("Device not available")
    }
    return h.Device.Product, nil
}

func (h *HidDevice) GetSerial() (string, error) {
    if h.Device == nil {
        return "", errors.New("Device not available")
    }
    return h.Device.Serial, nil
}

// Device represents a single Stream Deck device.
type Device struct {
    ID     string
    Serial string

    Columns    uint8
    Rows       uint8
    LcdColumns uint8
    Keys       uint8
    Knobs      uint8
    Pixels     uint
    LcdWidth   uint
    LcdHeight  uint
    DPI        uint
    Padding    uint

    featureReportSize   int
    firmwareOffset      int
    KeyStateOffset      int
    TranslateKeyIndex   func(index, columns uint8) uint8
    imagePageSize       int
    imagePageHeaderSize int
    flipImage           func(image.Image) image.Image
    toImageFormat       func(image.Image) ([]byte, error)
    imagePageHeader     func(pageIndex int, keyIndex uint8, payloadLength int, lastPage bool) []byte

    getFirmwareCommand   []byte
    resetCommand         []byte
    setBrightnessCommand []byte

    KeyState []byte

    Device IHidDevice

    lastActionTime time.Time
    asleep         bool
    sleepCancel    context.CancelFunc
    sleepMutex     *sync.RWMutex
    fadeDuration   time.Duration

    brightness         uint8
    preSleepBrightness uint8

    HasScreen bool
    HasLCD    bool
    HasKnobs  bool
    LCDBuffer image.Image
    WriteSem  *semaphore.Weighted

    InputHandler func(device *Device, cback func(event InputEvent))
}

// Key holds the current status of a key on the device.
type Key struct {
    Index   uint8
    Pressed bool
}

// Devices returns all attached Stream Decks.
func Devices() ([]Device, error) {
    var dd []Device

    devs := hid.Enumerate(VID_ELGATO, 0)
    for _, d := range devs {
        dd = append(dd, GetDevInfo(d))
    }

    return dd, nil
}

func GetDevInfo(d hid.DeviceInfo) Device {
    var dev Device

    switch {
    case d.VendorID == VID_ELGATO && d.ProductID == PID_STREAMDECK:
        dev = Device{
            ID:                   d.Path,
            Serial:               d.Serial,
            Columns:              5,
            Rows:                 3,
            Keys:                 15,
            Pixels:               72,
            DPI:                  124,
            Padding:              16,
            featureReportSize:    17,
            firmwareOffset:       5,
            KeyStateOffset:       1,
            TranslateKeyIndex:    translateRightToLeft,
            imagePageSize:        7819,
            imagePageHeaderSize:  16,
            imagePageHeader:      rev1ImagePageHeader,
            flipImage:            flipHorizontally,
            toImageFormat:        toBMP,
            getFirmwareCommand:   c_REV1_FIRMWARE,
            resetCommand:         c_REV1_RESET,
            setBrightnessCommand: c_REV1_BRIGHTNESS,
            HasScreen:            true,
            InputHandler:         basicHandleInput,
            WriteSem:             semaphore.NewWeighted(1),
        }
    case d.VendorID == VID_ELGATO && (d.ProductID == PID_STREAMDECK_MINI || d.ProductID == PID_STREAMDECK_MINI_MK2):
        dev = Device{
            ID:                   d.Path,
            Serial:               d.Serial,
            Columns:              3,
            Rows:                 2,
            Keys:                 6,
            Pixels:               80,
            DPI:                  138,
            Padding:              16,
            featureReportSize:    17,
            firmwareOffset:       5,
            KeyStateOffset:       1,
            TranslateKeyIndex:    identity,
            imagePageSize:        1024,
            imagePageHeaderSize:  16,
            imagePageHeader:      miniImagePageHeader,
            flipImage:            rotateCounterclockwise,
            toImageFormat:        toBMP,
            getFirmwareCommand:   c_REV1_FIRMWARE,
            resetCommand:         c_REV1_RESET,
            setBrightnessCommand: c_REV1_BRIGHTNESS,
            HasScreen:            true,
            InputHandler:         basicHandleInput,
            WriteSem:             semaphore.NewWeighted(1),
        }
    case d.VendorID == VID_ELGATO && (d.ProductID == PID_STREAMDECK_V2 || d.ProductID == PID_STREAMDECK_MK2):
        dev = Device{
            ID:                   d.Path,
            Serial:               d.Serial,
            Columns:              5,
            Rows:                 3,
            Keys:                 15,
            Pixels:               72,
            DPI:                  124,
            Padding:              16,
            featureReportSize:    32,
            firmwareOffset:       6,
            KeyStateOffset:       4,
            TranslateKeyIndex:    identity,
            imagePageSize:        1024,
            imagePageHeaderSize:  8,
            imagePageHeader:      rev2ImagePageHeader,
            flipImage:            flipHorizontallyAndVertically,
            toImageFormat:        toJPEG,
            getFirmwareCommand:   c_REV2_FIRMWARE,
            resetCommand:         c_REV2_RESET,
            setBrightnessCommand: c_REV2_BRIGHTNESS,
            HasScreen:            true,
            InputHandler:         basicHandleInput,
            WriteSem:             semaphore.NewWeighted(1),
        }
    case d.VendorID == VID_ELGATO && d.ProductID == PID_STREAMDECK_XL:
        dev = Device{
            ID:                   d.Path,
            Serial:               d.Serial,
            Columns:              8,
            Rows:                 4,
            Keys:                 32,
            Pixels:               96,
            DPI:                  166,
            Padding:              16,
            featureReportSize:    32,
            firmwareOffset:       6,
            KeyStateOffset:       4,
            TranslateKeyIndex:    identity,
            imagePageSize:        1024,
            imagePageHeaderSize:  8,
            imagePageHeader:      rev2ImagePageHeader,
            flipImage:            flipHorizontallyAndVertically,
            toImageFormat:        toJPEG,
            getFirmwareCommand:   c_REV2_FIRMWARE,
            resetCommand:         c_REV2_RESET,
            setBrightnessCommand: c_REV2_BRIGHTNESS,
            HasScreen:            true,
            InputHandler:         basicHandleInput,
            WriteSem:             semaphore.NewWeighted(1),
        }
    case d.VendorID == VID_ELGATO && d.ProductID == PID_STREAMDECK_PEDAL:
        dev = Device{
            ID:                 d.Path,
            Serial:             d.Serial,
            Columns:            3,
            Rows:               1,
            Keys:               3,
            Pixels:             72,
            DPI:                166,
            featureReportSize:  32,
            firmwareOffset:     6,
            KeyStateOffset:     4,
            TranslateKeyIndex:  identity,
            getFirmwareCommand: c_REV2_FIRMWARE,
            resetCommand:       c_REV2_RESET,
            HasScreen:          false,
            InputHandler:       basicHandleInput,
            WriteSem:           semaphore.NewWeighted(1),
        }
    case d.VendorID == VID_ELGATO && d.ProductID == PID_STREAMDECK_PLUS:
        dev = Device{
            ID:                  d.Path,
            Serial:              d.Serial,
            Columns:             4,
            Rows:                2,
            LcdColumns:          4,
            Keys:                8,
            Knobs:               4,
            Pixels:              120,
            LcdWidth:            200,
            LcdHeight:           100,
            DPI:                 166,
            Padding:             16,
            featureReportSize:   32,
            firmwareOffset:      6,
            KeyStateOffset:      4,
            TranslateKeyIndex:   identity,
            imagePageSize:       1024,
            imagePageHeaderSize: 8,
            imagePageHeader:     rev2ImagePageHeader,
            flipImage: func(i image.Image) image.Image {
                return i
            },
            toImageFormat:        toJPEG,
            getFirmwareCommand:   c_REV2_FIRMWARE,
            resetCommand:         c_REV2_RESET,
            setBrightnessCommand: c_REV2_BRIGHTNESS,
            HasScreen:            true,
            HasLCD:               true,
            HasKnobs:             true,
            InputHandler:         PlusInputHandler,
            WriteSem:             semaphore.NewWeighted(1),
        }

    }

    if dev.ID != "" {
        dev.KeyState = make([]byte, dev.Columns*dev.Rows)
        dev.Device = &HidDevice{info: d, WriteSem: semaphore.NewWeighted(1)}
    }
    return dev
}

// Open the device for input/output. This must be called before trying to
// communicate with the device.
func (d *Device) Open() error {
    var err error
    _, err = d.Device.Open()
    d.lastActionTime = time.Now()
    d.sleepMutex = &sync.RWMutex{}
    return err
}

// Close the connection with the device.
func (d *Device) Close() error {
    d.cancelSleepTimer()
    return d.Device.Close()
}

func (d *Device) HandleInput(cback func(event InputEvent)) {
    d.InputHandler(d, cback)
}

// FirmwareVersion returns the firmware version of the device.
func (d Device) FirmwareVersion() (string, error) {
    result, err := d.getFeatureReport(d.getFirmwareCommand)
    if err != nil {
        return "", err
    }
    return string(result[d.firmwareOffset:]), nil
}

// Resets the Stream Deck, clears all button images and shows the standby image.
func (d Device) Reset() error {
    if !d.HasScreen {
        return nil
    }
    return d.sendFeatureReport(d.resetCommand)
}

// Clears the Stream Deck, setting a black image on all buttons.
func (d Device) Clear() error {
    if !d.HasScreen {
        return nil
    }
    img := image.NewRGBA(image.Rect(0, 0, int(d.Pixels), int(d.Pixels)))
    draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{0, 0, 0, 255}), image.Point{}, draw.Src)
    for i := uint8(0); i <= d.Columns*d.Rows; i++ {
        err := d.SetImage(i, img)
        if err != nil {
            fmt.Println(err)
            return err
        }
    }

    return nil
}

func (d *Device) cancelSleepTimer() {
    if d.sleepCancel == nil {
        return
    }

    d.sleepCancel()
    d.sleepCancel = nil
}

// SetSleepFadeDuration sets the duration of the fading animation when the
// device is put to sleep or wakes up.
func (d *Device) SetSleepFadeDuration(t time.Duration) {
    if !d.HasScreen {
        return
    }
    d.fadeDuration = t
}

// Fade fades the brightness in or out.
func (d *Device) Fade(start uint8, end uint8, duration time.Duration) error {
    if !d.HasScreen {
        return nil
    }
    step := (float64(end) - float64(start)) / float64(duration/fadeDelay)
    if step == math.Inf(1) || step == math.Inf(-1) {
        return nil
    }

    for current := float64(start); ; current += step {
        if !((start < end && int8(current) < int8(end)) ||
            (start > end && int8(current) > int8(end))) {
            break
        }
        if err := d.SetBrightness(uint8(current)); err != nil {
            return err
        }

        time.Sleep(fadeDelay)
    }
    return nil
}

// SetBrightness sets the background lighting brightness from 0 to 100 percent.
func (d *Device) SetBrightness(percent uint8) error {
    if !d.HasScreen {
        return nil
    }
    if percent > 100 {
        percent = 100
    }

    d.brightness = percent
    if d.asleep && percent > 0 {
        // if the device is asleep, remember the brightness, but don't set it
        d.sleepMutex.Lock()
        d.preSleepBrightness = percent
        d.sleepMutex.Unlock()
        return nil
    }

    report := make([]byte, len(d.setBrightnessCommand)+1)
    copy(report, d.setBrightnessCommand)
    report[len(report)-1] = percent

    return d.sendFeatureReport(report)
}

// SetImage sets the image of a button on the Stream Deck. The provided image
// needs to be in the correct resolution for the device. The index starts with
// 0 being the top-left button.
func (d Device) SetImage(index uint8, img image.Image) error {
    ctx := context.Background()
    err := d.WriteSem.Acquire(ctx, 1)
    if err != nil {
        return err
    }
    defer d.WriteSem.Release(1)
    if !d.HasScreen {
        return nil
    }
    if img.Bounds().Dy() != int(d.Pixels) ||
        img.Bounds().Dx() != int(d.Pixels) {
        return fmt.Errorf("supplied image has wrong dimensions, expected %[1]dx%[1]d pixels", d.Pixels)
    }

    imageBytes, err := d.toImageFormat(d.flipImage(img))
    if err != nil {
        return fmt.Errorf("cannot convert image data: %v", err)
    }
    imageData := imageData{
        image:    imageBytes,
        pageSize: d.imagePageSize - d.imagePageHeaderSize,
    }

    data := make([]byte, d.imagePageSize)

    var page int
    var lastPage bool
    for !lastPage {
        var payload []byte
        payload, lastPage = imageData.Page(page)
        header := d.imagePageHeader(page, d.TranslateKeyIndex(index, d.Columns), len(payload), lastPage)

        copy(data, header)
        copy(data[len(header):], payload)

        _, err := d.Device.Write(data)
        if err != nil {
            return fmt.Errorf("cannot write image page %d of %d (%d image bytes) %d bytes: %v",
                page, imageData.PageCount(), imageData.Length(), len(data), err)
        }

        page++
    }
    return nil
}

func (d Device) SetLcdImage(index int, img image.Image) error {
    ctx := context.Background()
    err := d.WriteSem.Acquire(ctx, 1)
    if err != nil {
        return err
    }
    defer d.WriteSem.Release(1)
    if !d.HasLCD {
        return nil
    }
    if img.Bounds().Dx() != int(d.LcdWidth) || img.Bounds().Dy() != int(d.LcdHeight) {
        return fmt.Errorf("supplied image has wrong dimensions, expected %[1]dx%[1]d pixels", d.Pixels)
    }
    imageBytes, err := toJPEG(d.flipImage(img))
    if err != nil {
        return fmt.Errorf("cannot convert image data: %v", err)
    }
    imageData := imageData{
        image:    imageBytes,
        pageSize: d.imagePageSize - 16,
    }

    var page int
    var lastPage bool

    for !lastPage {
        data := make([]byte, d.imagePageSize)
        var payload []byte
        payload, lastPage = imageData.Page(page)
        var header []byte
        header = append(header, 0x02)
        header = append(header, 0x0c)
        header = binary.LittleEndian.AppendUint16(header, uint16(index*int(d.LcdWidth)))
        header = binary.LittleEndian.AppendUint16(header, 0)
        header = binary.LittleEndian.AppendUint16(header, uint16(d.LcdWidth))
        header = binary.LittleEndian.AppendUint16(header, uint16(d.LcdHeight))
        if lastPage {
            header = append(header, 0x01)
        } else {
            header = append(header, 0x00)
        }
        header = binary.LittleEndian.AppendUint16(header, uint16(page))
        header = binary.LittleEndian.AppendUint16(header, uint16(math.Min(float64(len(payload)), float64(d.imagePageSize-16))))
        header = append(header, 0x00)

        copy(data, header)
        copy(data[len(header):], payload)

        _, err := d.Device.Write(data)
        if err != nil {
            return fmt.Errorf("cannot write image page %d of %d (%d image bytes) %d bytes: %v",
                page, imageData.PageCount(), imageData.Length(), len(data), err)
        }

        page++
    }
    return nil
}

// getFeatureReport from the device without worries about the correct payload
// size.
func (d Device) getFeatureReport(payload []byte) ([]byte, error) {
    b := make([]byte, d.featureReportSize)
    copy(b, payload)
    _, err := d.Device.GetFeatureReport(b)
    if err != nil {
        return nil, err
    }
    return b, nil
}

// sendFeatureReport to the device without worries about the correct payload
// size.
func (d Device) sendFeatureReport(payload []byte) error {
    b := make([]byte, d.featureReportSize)
    copy(b, payload)
    _, err := d.Device.SendFeatureReport(b)
    return err
}

// translateRightToLeft translates the given key index from right-to-left to
// left-to-right, based on the given number of columns.
func translateRightToLeft(index, columns uint8) uint8 {
    keyCol := index % columns
    return (index - keyCol) + (columns - 1) - keyCol
}

// identity returns the given key index as it is.
func identity(index, _ uint8) uint8 {
    return index
}

// toRGBA converts an image.Image to an image.RGBA.
func toRGBA(img image.Image) *image.RGBA {
    switch img := img.(type) {
    case *image.RGBA:
        return img
    }
    out := image.NewRGBA(img.Bounds())
    draw.Copy(out, image.Pt(0, 0), img, img.Bounds(), draw.Src, nil)
    return out
}

// flipHorizontally returns the given image horizontally flipped.
func flipHorizontally(img image.Image) image.Image {
    flipped := image.NewRGBA(img.Bounds())
    draw.Copy(flipped, image.Point{}, img, img.Bounds(), draw.Src, nil)
    for y := 0; y < flipped.Bounds().Dy(); y++ {
        for x := 0; x < flipped.Bounds().Dx()/2; x++ {
            xx := flipped.Bounds().Max.X - x - 1
            c := flipped.RGBAAt(x, y)
            flipped.SetRGBA(x, y, flipped.RGBAAt(xx, y))
            flipped.SetRGBA(xx, y, c)
        }
    }
    return flipped
}

// flipHorizontallyAndVertically returns the given image horizontally and
// vertically flipped.
func flipHorizontallyAndVertically(img image.Image) image.Image {
    flipped := image.NewRGBA(img.Bounds())
    draw.Copy(flipped, image.Point{}, img, img.Bounds(), draw.Src, nil)
    for y := 0; y < flipped.Bounds().Dy()/2; y++ {
        yy := flipped.Bounds().Max.Y - y - 1
        for x := 0; x < flipped.Bounds().Dx(); x++ {
            xx := flipped.Bounds().Max.X - x - 1
            c := flipped.RGBAAt(x, y)
            flipped.SetRGBA(x, y, flipped.RGBAAt(xx, yy))
            flipped.SetRGBA(xx, yy, c)
        }
    }
    return flipped
}

// rotateCounterclockwise returns the given image rotated counterclockwise.
func rotateCounterclockwise(img image.Image) image.Image {
    flipped := image.NewRGBA(img.Bounds())
    draw.Copy(flipped, image.Point{}, img, img.Bounds(), draw.Src, nil)
    for y := 0; y < flipped.Bounds().Dy(); y++ {
        for x := y + 1; x < flipped.Bounds().Dx(); x++ {
            c := flipped.RGBAAt(x, y)
            flipped.SetRGBA(x, y, flipped.RGBAAt(y, x))
            flipped.SetRGBA(y, x, c)
        }
    }
    for y := 0; y < flipped.Bounds().Dy()/2; y++ {
        yy := flipped.Bounds().Max.Y - y - 1
        for x := 0; x < flipped.Bounds().Dx(); x++ {
            c := flipped.RGBAAt(x, y)
            flipped.SetRGBA(x, y, flipped.RGBAAt(x, yy))
            flipped.SetRGBA(x, yy, c)
        }
    }
    return flipped
}

// toBMP returns the raw bytes of the given image in BMP format.
func toBMP(img image.Image) ([]byte, error) {
    rgba := toRGBA(img)

    // this is a BMP file header followed by a BPM bitmap info header
    // find more information here: https://en.wikipedia.org/wiki/BMP_file_format
    header := []byte{
        0x42, 0x4d, 0xf6, 0x3c, 0x00, 0x00, 0x00, 0x00,
        0x00, 0x00, 0x36, 0x00, 0x00, 0x00, 0x28, 0x00,
        0x00, 0x00, 0x48, 0x00, 0x00, 0x00, 0x48, 0x00,
        0x00, 0x00, 0x01, 0x00, 0x18, 0x00, 0x00, 0x00,
        0x00, 0x00, 0xc0, 0x3c, 0x00, 0x00, 0xc4, 0x0e,
        0x00, 0x00, 0xc4, 0x0e, 0x00, 0x00, 0x00, 0x00,
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    }

    buffer := make([]byte, len(header)+rgba.Bounds().Dx()*rgba.Bounds().Dy()*3)
    copy(buffer, header)

    i := len(header)
    for y := 0; y < rgba.Bounds().Dy(); y++ {
        for x := 0; x < rgba.Bounds().Dx(); x++ {
            c := rgba.RGBAAt(x, y)
            buffer[i] = c.B
            buffer[i+1] = c.G
            buffer[i+2] = c.R
            i += 3
        }
    }
    return buffer, nil
}

// toJPEG returns the raw bytes of the given image in JPEG format.
func toJPEG(img image.Image) ([]byte, error) {
    buffer := bytes.NewBuffer([]byte{})
    opts := jpeg.Options{
        Quality: 100,
    }
    err := jpeg.Encode(buffer, img, &opts)
    if err != nil {
        return nil, err
    }
    return buffer.Bytes(), err
}

// rev1ImagePageHeader returns the image page header sequence used by the
// Stream Deck v1.
func rev1ImagePageHeader(pageIndex int, keyIndex uint8, payloadLength int, lastPage bool) []byte {
    var lastPageByte byte
    if lastPage {
        lastPageByte = 1
    }
    return []byte{
        0x02, 0x01,
        byte(pageIndex + 1), 0x00,
        lastPageByte,
        keyIndex + 1,
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    }
}

// miniImagePageHeader returns the image page header sequence used by the
// Stream Deck Mini.
func miniImagePageHeader(pageIndex int, keyIndex uint8, payloadLength int, lastPage bool) []byte {
    var lastPageByte byte
    if lastPage {
        lastPageByte = 1
    }
    return []byte{
        0x02, 0x01,
        byte(pageIndex), 0x00,
        lastPageByte,
        keyIndex + 1,
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    }
}

// rev2ImagePageHeader returns the image page header sequence used by Stream
// Deck XL and Stream Deck v2.
func rev2ImagePageHeader(pageIndex int, keyIndex uint8, payloadLength int, lastPage bool) []byte {
    var lastPageByte byte
    if lastPage {
        lastPageByte = 1
    }
    return []byte{
        0x02, 0x07, keyIndex, lastPageByte,
        byte(payloadLength), byte(payloadLength >> 8),
        byte(pageIndex), byte(pageIndex >> 8),
    }
}

func basicHandleInput(d *Device, cback func(event InputEvent)) {
    for {
        keyBuff := make([]byte, d.KeyStateOffset+len(d.KeyState))

        if _, err := d.Device.Read(keyBuff); err != nil {
            panic(err)
        }

        for i := d.KeyStateOffset; i < len(keyBuff); i++ {
            keyIndex := uint8(i - d.KeyStateOffset)
            if keyBuff[i] != d.KeyState[keyIndex] {
                if keyBuff[i] == 0x01 {
                    d.KeyState[keyIndex] = 0x01
                    cback(InputEvent{
                        EventType: KEY_PRESS,
                        Index:     d.TranslateKeyIndex(keyIndex, d.Columns),
                    })
                } else {
                    d.KeyState[keyIndex] = 0x00
                    cback(InputEvent{
                        EventType: KEY_RELEASE,
                        Index:     d.TranslateKeyIndex(keyIndex, d.Columns),
                    })
                }
            }
        }
    }
}

func PlusInputHandler(d *Device, cback func(event InputEvent)) {
    for {
        keyBuff := make([]byte, d.KeyStateOffset+len(d.KeyState)+3)

        copy(d.KeyState, keyBuff[d.KeyStateOffset:])

        _, err := d.Device.Read(keyBuff)

        if err != nil {
            panic(err)
        }
        if keyBuff[1] == 0x00 {
            // keys
            for i := d.KeyStateOffset; i < len(keyBuff); i++ {
                keyIndex := uint8(i - d.KeyStateOffset)
                if keyBuff[i] == 0x01 {
                    cback(InputEvent{
                        EventType: KEY_PRESS,
                        Index:     d.TranslateKeyIndex(keyIndex, d.Columns),
                    })
                }
            }
        } else if keyBuff[1] == 0x02 {
            x := binary.LittleEndian.Uint16(keyBuff[6:])
            y := binary.LittleEndian.Uint16(keyBuff[8:])
            buttonIndex, err := virtButtonPressed(x)
            if err != nil {
                fmt.Println(err)
                continue
            }
            if keyBuff[4] == 0x01 {
                cback(InputEvent{
                    EventType: SCREEN_SHORT_TAP,
                    Index:     uint8(buttonIndex),
                })
            } else if keyBuff[4] == 0x02 {
                cback(InputEvent{
                    EventType: SCREEN_LONG_TAP,
                    Index:     uint8(buttonIndex),
                })
            } else if keyBuff[4] == 0x03 {
                x2 := binary.LittleEndian.Uint16(keyBuff[10:])
                y2 := binary.LittleEndian.Uint16(keyBuff[12:])
                cback(InputEvent{
                    EventType:  SCREEN_SWIPE,
                    ScreenX:    x,
                    ScreenY:    y,
                    ScreenEndX: x2,
                    ScreenEndY: y2,
                })
            }
        } else if keyBuff[1] == 0x03 {
            isPress := keyBuff[4] == 0x00
            for i := 5; i < len(keyBuff); i++ {
                knobIndex := uint8(i - 5)
                if isPress && keyBuff[i] == 0x01 {
                    cback(InputEvent{
                        EventType: KNOB_PRESS,
                        Index:     knobIndex,
                    })
                } else if keyBuff[i] != 0x00 {
                    angleChange := int8(binary.LittleEndian.Uint16(keyBuff[i:]))
                    eventType := KNOB_CW
                    if angleChange < 0 {
                        eventType = KNOB_CCW
                    }
                    cback(InputEvent{
                        EventType:     eventType,
                        RotateNotches: uint8(math.Abs(float64(angleChange))),
                        Index:         knobIndex,
                    })
                }
            }
        }
    }
}

func virtButtonPressed(x uint16) (int, error) {
    if x < 0 || x > 800 {
        return 0, errors.New("Invalid x position")
    }
    return int(math.Floor(float64(x) / 200)), nil
}

// imageData allows to access raw image data in a byte array through pages of a
// given size.
type imageData struct {
    image    []byte
    pageSize int
}

// Page returns the page with the given index and an indication if this is the
// last page.
func (d imageData) Page(pageIndex int) ([]byte, bool) {
    offset := pageIndex * d.pageSize
    if offset >= len(d.image) {
        return []byte{}, true
    }

    length := d.pageLength(pageIndex)
    if offset+length > len(d.image) {
        length = len(d.image) - offset
    }

    return d.image[offset : offset+length], pageIndex == d.PageCount()-1
}

func (d imageData) pageLength(pageIndex int) int {
    remaining := len(d.image) - (pageIndex * d.pageSize)
    if remaining > d.pageSize {
        return d.pageSize
    }
    if remaining > 0 {
        return remaining
    }
    return 0
}

// PageCount returns the total number of pages.
func (d imageData) PageCount() int {
    count := len(d.image) / d.pageSize
    if len(d.image)%d.pageSize != 0 {
        return count + 1
    }
    return count
}

// Length of the raw image data in bytes.
func (d imageData) Length() int {
    return len(d.image)
}
