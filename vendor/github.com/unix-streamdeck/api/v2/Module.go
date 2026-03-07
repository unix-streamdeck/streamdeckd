package api

type Module struct {
	Name              string                    `json:"name,omitempty"`
	NewIcon           func() IconHandler        `json:"-"`
	NewKey            func() KeyHandler         `json:"-"`
	NewLcd            func() LcdHandler         `json:"-"`
	NewKnobOrTouch    func() KnobOrTouchHandler `json:"-"`
	IconFields        []Field                   `json:"icon_fields,omitempty"`
	KeyFields         []Field                   `json:"key_fields,omitempty"`
	LcdFields         []Field                   `json:"lcd_fields,omitempty"`
	KnobOrTouchFields []Field                   `json:"knob_or_touch_fields,omitempty"`
	LinkedFields      []Field                   `json:"linked_fields,omitempty"`
	IsIcon            bool                      `json:"is_icon,omitempty"`
	IsKey             bool                      `json:"is_key,omitempty"`
	IsLcd             bool                      `json:"is_lcd,omitempty"`
	IsKnob            bool                      `json:"is_knob,omitempty"`
	IsLinkedHandlers  bool                      `json:"is_linked_handlers,omitempty"`
	Linked            bool                      `json:"linked,omitempty"`
}

type FieldType string

const (
	File          FieldType = "File"
	Text          FieldType = "Text"
	Number        FieldType = "Number"
	TextAlignment FieldType = "TextAlignment"
	FontFace      FieldType = "FontFace"
	Select        FieldType = "Select"
	Colour        FieldType = "Colour"
)

type Field struct {
	Title     string    `json:"title,omitempty"`
	Name      string    `json:"name,omitempty"`
	Type      FieldType `json:"type,omitempty"`
	FileTypes []string  `json:"file_types,omitempty"`
	ListItems []string  `json:"list_items,omitempty"`
}
