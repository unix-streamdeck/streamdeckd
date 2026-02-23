package api

type Module struct {
	Name              string
	NewIcon           func() IconHandler
	NewKey            func() KeyHandler
	NewLcd            func() LcdHandler
	NewKnobOrTouch    func() KnobOrTouchHandler
	IconFields        []Field `json:"icon_fields,omitempty"`
	KeyFields         []Field `json:"key_fields,omitempty"`
	LcdFields         []Field `json:"lcd_fields,omitempty"`
	KnobOrTouchFields []Field `json:"knob_or_touch_fields,omitempty"`
	LinkedFields      []Field `json:"linked_fields,omitempty"`
	IsIcon            bool    `json:"is_icon,omitempty"`
	IsKey             bool    `json:"is_key,omitempty"`
	IsLcd             bool    `json:"is_lcd,omitempty"`
	IsKnob            bool    `json:"is_knob,omitempty"`
	IsLinkedHandlers  bool    `json:"is_linked_handlers,omitempty"`
	Linked            bool    `json:"linked,omitempty"`
}

type Field struct {
	Title     string   `json:"title,omitempty"`
	Name      string   `json:"name,omitempty"`
	Type      string   `json:"type,omitempty"`
	FileTypes []string `json:"file_types,omitempty"`
	ListItems []string `json:"list_items,omitempty"`
}
