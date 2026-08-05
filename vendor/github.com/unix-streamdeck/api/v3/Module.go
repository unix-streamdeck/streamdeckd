package api

type Module struct {
	Name string `json:"name,omitempty"`

	NewCombined   func() Handler           `json:"-"`
	NewBackground func() BackgroundHandler `json:"-"`

	CombinedFields   []Field `json:"combined_fields,omitempty"`
	BackgroundFields []Field `json:"lcd_fields,omitempty"`

	IsForeground bool `json:"is_foreground,omitempty"`
	IsInput      bool `json:"is_input,omitempty"`
	IsBackground bool `json:"is_background,omitempty"`
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
