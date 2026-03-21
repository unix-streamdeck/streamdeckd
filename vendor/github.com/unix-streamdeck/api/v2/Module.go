package api

type Module struct {
	Name string `json:"name,omitempty"`

	NewForeground func() ForegroundHandler `json:"-"`
	NewInput      func() InputHandler      `json:"-"`
	NewBackground func() BackgroundHandler `json:"-"`

	ForegroundFields []Field `json:"icon_fields,omitempty"`
	InputFields      []Field `json:"key_fields,omitempty"`
	BackgroundFields []Field `json:"lcd_fields,omitempty"`
	LinkedFields     []Field `json:"linked_fields,omitempty"`

	IsForeground     bool `json:"is_foreground,omitempty"`
	IsInput          bool `json:"is_input,omitempty"`
	IsBackground     bool `json:"is_background,omitempty"`
	IsLinkedHandlers bool `json:"is_linked_handlers,omitempty"`
	Linked           bool `json:"linked,omitempty"`
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
