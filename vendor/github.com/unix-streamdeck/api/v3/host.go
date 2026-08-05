package api

type IHost interface {
	SetPage(page int)
	SetBrightness(brightness uint8) error
	RunCommand(command string)
	ExecuteKeybind(keybind string) error
	KeyDown(keycode int) error
	KeyUp(keycode int) error
}
