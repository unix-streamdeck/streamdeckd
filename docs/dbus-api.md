# D-Bus API Reference

## Connection Details

- **Service Name:** `com.unixstreamdeck.streamdeckd`
- **Object Path:** `/com/unixstreamdeck/streamdeckd`
- **Interface:** `com.unixstreamdeck.streamdeckd`
- **Bus Type:** Session bus

## Methods

### GetConfig

Retrieve the current running configuration.

**Parameters:** None

**Returns:** JSON string containing the complete configuration

**Example:**
```bash
dbus-send --print-reply --session \
  --dest=com.unixstreamdeck.streamdeckd \
  /com/unixstreamdeck/streamdeckd \
  com.unixstreamdeck.streamdeckd.GetConfig
```

**Response:**
```json
{
  "modules": ["/path/to/module.so"],
  "decks": [
    {
      "serial": "AB12C3D45678",
      "pages": [ /* ... */ ]
    }
  ]
}
```

---

### SetConfig

Update the configuration in memory without saving to disk.

**Parameters:**
- `config` (string): JSON string containing the new configuration

**Returns:** Error message if invalid, empty on success

**Example:**
```bash
dbus-send --session --dest=com.unixstreamdeck.streamdeckd \
  /com/unixstreamdeck/streamdeckd \
  com.unixstreamdeck.streamdeckd.SetConfig \
  string:'{"modules":[],"decks":[]}'
```

**Notes:**
- Changes take effect immediately
- Configuration is **not** saved to disk (use `CommitConfig` to save)
- Invalid JSON will return an error

---

### ReloadConfig

Reload the configuration from disk, discarding any unsaved changes.

**Parameters:** None

**Returns:** None

**Example:**
```bash
dbus-send --session --dest=com.unixstreamdeck.streamdeckd \
  /com/unixstreamdeck/streamdeckd \
  com.unixstreamdeck.streamdeckd.ReloadConfig
```

---

### CommitConfig

Save the current in-memory configuration to disk.

**Parameters:** None

**Returns:** Error message if save fails, empty on success

**Example:**
```bash
dbus-send --session --dest=com.unixstreamdeck.streamdeckd \
  /com/unixstreamdeck/streamdeckd \
  com.unixstreamdeck.streamdeckd.CommitConfig
```

---

### GetDeckInfo

Get information about all connected Stream Deck devices.

**Parameters:** None

**Returns:** JSON array of deck information objects

**Example:**
```bash
dbus-send --print-reply --session \
  --dest=com.unixstreamdeck.streamdeckd \
  /com/unixstreamdeck/streamdeckd \
  com.unixstreamdeck.streamdeckd.GetDeckInfo
```

**Response:**
```json
[
  {
    "icon_size": 72,
    "rows": 3,
    "cols": 5,
    "page": 0,
    "serial": "AB12C3D45678"
  },
  {
    "icon_size": 80,
    "rows": 4,
    "cols": 8,
    "page": 1,
    "serial": "XYZ987654321"
  }
]
```

**Fields:**
- `icon_size`: Button icon dimensions in pixels
- `rows`: Number of button rows
- `cols`: Number of button columns
- `page`: Currently active page number
- `serial`: Device serial number

---

### SetPage

Switch the Stream Deck to a specific page.

**Parameters:**
- `serial` (string): Device serial number
- `page` (int): Page number to switch to (zero-indexed)

**Returns:** Error message if invalid, empty on success

**Example:**
```bash
dbus-send --session --dest=com.unixstreamdeck.streamdeckd \
  /com/unixstreamdeck/streamdeckd \
  com.unixstreamdeck.streamdeckd.SetPage \
  string:'AB12C3D45678' int32:2
```

**Notes:**
- Page numbers start at 0
- Returns error if page doesn't exist
- Emits `Page` signal on success

---

### GetModules

Get information about loaded custom modules.

**Parameters:** None

**Returns:** Array of module information including available config fields

**Example:**
```bash
dbus-send --print-reply --session \
  --dest=com.unixstreamdeck.streamdeckd \
  /com/unixstreamdeck/streamdeckd \
  com.unixstreamdeck.streamdeckd.GetModules
```

**Response:**
```json
[
  {
    "name": "CustomModule",
    "icon_fields": [
      {
        "title": "Background Color",
        "name": "bg_color",
        "type": "Text"
      }
    ],
    "key_fields": [
      {
        "title": "API Key",
        "name": "api_key",
        "type": "Text"
      }
    ]
  }
]
```

---

### PressButton

Simulate a button press on a Stream Deck device.

**Parameters:**
- `serial` (string): Device serial number
- `key_index` (int): Button index to press

**Returns:** Error message if invalid, empty on success

**Example:**
```bash
# Press button 5 on device AB12C3D45678
dbus-send --session --dest=com.unixstreamdeck.streamdeckd \
  /com/unixstreamdeck/streamdeckd \
  com.unixstreamdeck.streamdeckd.PressButton \
  string:'AB12C3D45678' int32:5
```

**Notes:**
- Button indices start at 0
- Indices go left-to-right, top-to-bottom
- Triggers the same action as physically pressing the button

---

### GetHandlerExample & GetKnobHandlerExample

Simulate a handler config, and get an example response of what image that handler would generate with that config

**Parameters:**
- `serial` (string): Device serial number
- `key/knobConfig` (string): JSON string containing a key config

**Returns** Error message if invalid, b64 encoded png string if valid

---

## Signals

### Page

Emitted when the active page changes on any Stream Deck device.

**Parameters:**
- `serial` (string): Device serial number
- `page` (int): New page number

**Example - Monitor Page Changes:**
```bash
dbus-monitor "type='signal',\
interface='com.unixstreamdeck.streamdeckd',\
member='Page'"
```

**Use Case:** Update external UI or trigger actions when pages change



## See Also

- [Configuration Guide](configuration.md)
- [Installation Guide](installation.md)
- [Custom Modules](custom-modules.md)
