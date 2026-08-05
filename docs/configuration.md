# Configuration Guide

Complete reference for configuring streamdeckd.

## Configuration File Location

Default: `$XDG_CONFIG_HOME/.streamdeck-config.json` (usually `~/.config/.streamdeck-config.json`)

Custom location: `./streamdeckd -config /path/to/config.json`

## Configuration Structure

```json
{
  "modules": [
    "/path/to/module1.so",
    "/path/to/module2.so"
  ],
  "decks": [
    {
      "serial": "AB12C3D45678",
      "pages": [
        [ /* Page 0 buttons */ ],
        [ /* Page 1 buttons */ ],
        [ /* Page 2 buttons */ ]
      ]
    }
  ]
}
```

### Top-Level Fields

| Field     | Type             | Description                               |
|-----------|------------------|-------------------------------------------|
| `modules` | Array of strings | Paths to custom plugin `.so` files        |
| `decks`   | Array of objects | Configuration for each Stream Deck device |

## Deck Configuration

Each deck object represents one physical Stream Deck device.

```json
{
  "serial": "AB12C3D45678",
  "pages": [ /* ... */ ]
}
```

### Getting Your Device Serial

Use the D-Bus `GetDeckInfo` method:

```bash
dbus-send --print-reply --dest=com.unixstreamdeck.streamdeckd \
  /com/unixstreamdeck/streamdeckd \
  com.unixstreamdeck.streamdeckd.GetDeckInfo
```

Or use streamdeckui, which displays the serial automatically.

## Pages and Buttons

### Page Structure

Pages are nested arrays:
- **Outer array:** List of pages
- **Inner array:** List of buttons on each page

```json
"pages": [
  [
    { /* Button 0 */ },
    { /* Button 1 */ },
    { /* Button 2 */ }
  ],
  [
    { /* Button 0 on page 2 */ },
    { /* Button 1 on page 2 */ }
  ]
]
```

### Button Order

Buttons are indexed **left-to-right, top-to-bottom**:

```
Stream Deck Original (15 keys):
┌─────┬─────┬─────┬─────┬─────┐
│  0  │  1  │  2  │  3  │  4  │
├─────┼─────┼─────┼─────┼─────┤
│  5  │  6  │  7  │  8  │  9  │
├─────┼─────┼─────┼─────┼─────┤
│ 10  │ 11  │ 12  │ 13  │ 14  │
└─────┴─────┴─────┴─────┴─────┘

Stream Deck Mini (6 keys):
┌─────┬─────┬─────┐
│  0  │  1  │  2  │
├─────┼─────┼─────┤
│  3  │  4  │  5  │
└─────┴─────┴─────┘
```

## Button Configuration

Each button has an `application` object that maps application names to actions, nested inside is a list of "layers", those layers map to handlers that can update images on the device, handle inputs from the device, or for more complex handlers, do both

### Basic Button

```json
{
  "application": {
    "": {
      "layers": [
        ...
      ]
    }
  }
}
```

The `""` (empty string) key is the default configuration for all applications.

### Per-Application Buttons

Different actions based on the active application:

```json
{
  "application": {
    "": {
      "layers": [
        {
          "handler": "Command",
          "fields": {
            "command": "echo 'default'"
          }
        },
        {
          "handler": "Image",
          "fields": {
            "path": "/path/to/default.png"
          }
        }
      ]
    },
    "firefox": {
      "layers": [
        {
          "handler": "Keybind",
          "fields": {
            "keybind": "ctrl+t"
          }
        },
        {
          "handler": "Image",
          "fields": {
            "path": "/path/to/firefox.png"
          }
        }
      ]
    },
    "spotify": {
      "layers": [
        {
          "handler": "Keybind",
          "fields": {
            "keybind": "XF86AudioPlay"
          }
        },
        {
          "handler": "Image",
          "fields": {
            "path": "/path/to/spotify.png"
          }
        }
      ]
    }
  }
}
```

When Firefox is active, the button sends `Ctrl+T`. When Spotify is active, it toggles play/pause. Otherwise, it runs the default command.

### Application Class Detection

Applications are detected via their classes, as these tend to stay relatively consistent and unique. Currently only Hyprland, KDE, and X11 are supported for the application class detection, but pull requests are welcome.


**Tip:** Use streamdeckui to see detected application classes in real-time.

## Actions

### Command

Execute shell commands.

```json
{
  "handler": "Command",
  "fields": {
    "command": "notify-send 'Button Pressed'"
  }
}
```
### Keybind

Simulate keyboard input using xdotool syntax.

```json
{
  "handler": "Keybind",
  "fields": {
    "keybind": "ctrl+t"
  }
}
```

**Special Keys:**
- `ctrl`, `shift`, `alt`, `super`
- `Tab`, `Return`, `Escape`, `space`
- `Up`, `Down`, `Left`, `Right`
- `F1` through `F12`
- `XF86Audio*` (media keys)

### URL

Open URLs in the default browser.

```json
{
  "handler": "Url",
  "fields": {
    "url": "https://github.com"
  }
}
```

### Switch Page

Navigate to a different button page.

```json
{
  "handler": "SwitchPage",
  "fields": {
    "page": 1
  },
}
```

Pages are zero-indexed (first page is 0).

### Brightness

Adjust Stream Deck display brightness (0-100).

```json
{
  "handler": "Brightness",
  "fields": {
    "brightness": 50
  },
}
```

### Icon

Set the button icon image.

```json
{
  "handler": "Image",
  "fields": {
    "path": "/path/to/image.png"
  },
}
```

## Complete Examples

### Media Control Page

```json
{
  "decks": [{
    "serial": "AB12C3D45678",
    "pages": [[
      {
        "application": {
          "": {
            "layers": [
              {
                "handler": "Image",
                "fields": {
                  "path": ~/icons/play-pause.png"
                },
              },
              {
                "handler": "Keybind",
                "fields": {
                  "keybind": "XF86AudioPlay"
                }
              }
            ]
          }
        }
      },
      {
        "application": {
          "": {
            "layers": [
              {
                "handler": "Image",
                "fields": {
                  "path": "~/icons/previous.png"
                },
              },
              {
                "handler": "Keybind",
                "fields": {
                  "keybind": "XF86AudioPrev"
                }
              }
            ]
          }
        }
      },
      {
        "application": {
          "": {
            "layers": [
              {
                "handler": "Image",
                "fields": {
                  "path": "~/icons/next.png"
                },
              },
              {
                "handler": "Keybind",
                "fields": {
                  "keybind": "XF86AudioNext"
                }
              }
            ]
          }
        }
      },
      {
        "application": {
          "": {
            "layers": [
              {
                "handler": "Image",
                "fields": {
                  "path": "~/icons/mute.png"
                },
              },
              {
                "handler": "Keybind",
                "fields": {
                  "keybind": "XF86AudioMute"
                }
              }
            ]
          }
        }
      }
    ]]
  }]
}
```

### Application-Specific Controls

```json
{
  "decks": [{
    "serial": "AB12C3D45678",
    "pages": [[
      {
        "application": {
          "": {
            "command": "notify-send 'No app active'",
            "icon": "~/icons/default.png",
            "layers": [
              {
                "handler": "Image",
                "fields": {
                  "path": "~/icons/default.png"
                },
              },
              {
                "handler": "Command",
                "fields": {
                  "command": "notify-send 'No app active'"
                }
              }
            ]
          },
          "firefox": {
            "layers": [
              {
                "handler": "Image",
                "fields": {
                  "path": "~/icons/new-tab.png"
                },
              },
              {
                "handler": "Keybind",
                "fields": {
                  "keybind": "ctrl+t"
                }
              }
            ]
          },
          "chrome": {
            "layers": [
              {
                "handler": "Image",
                "fields": {
                  "path": "~/icons/new-tab.png"
                },
              },
              {
                "handler": "Keybind",
                "fields": {
                  "keybind": "ctrl+t"
                }
              }
            ]
          }
        }
      },
      {
        "application": {
          "": {
            "command": "notify-send 'No app active'",
            "icon": "~/icons/default.png",
            "layers": [
              {
                "handler": "Image",
                "fields": {
                  "path": "~/icons/default.png"
                },
              },
              {
                "handler": "Command",
                "fields": {
                  "command": "notify-send 'No app active'"
                }
              }
            ]
          },
          "code": {
            "layers": [
              {
                "handler": "Image",
                "fields": {
                  "path": "~/icons/terminal.png"
                },
              },
              {
                "handler": "Keybind",
                "fields": {
                  "keybind": "ctrl+grave"
                }
              }
            ]
          }
        }
      }
    ]]
  }]
}
```

## Dynamic Configuration

### Reload Configuration

Reload from disk without restarting:

```bash
dbus-send --dest=com.unixstreamdeck.streamdeckd \
  /com/unixstreamdeck/streamdeckd \
  com.unixstreamdeck.streamdeckd.ReloadConfig
```

### Set Configuration Programmatically

See the [D-Bus API documentation](dbus-api.md) for `SetConfig` and `CommitConfig` methods.

## Best Practices

### Icon Guidelines

- **Size:** 72x72px for most Stream Decks
- **Format:** PNG with transparency recommended
- **Design:** High contrast, simple icons work best
- **Resources:**
  - [Elgato's Streamdeck Key Creator](https://www.elgato.com/uk/en/s/keycreator)


## Advanced Configuration

### Custom Modules

For advanced functionality, see [Custom Modules](custom-modules.md).

### Multiple Devices

Configure multiple Stream Decks with different serials:

```json
{
  "decks": [
    {
      "serial": "ABC123",
      "pages": [ /* ... */ ]
    },
    {
      "serial": "DEF456",
      "pages": [ /* ... */ ]
    }
  ]
}
```

### Empty Buttons

Leave a button empty by omitting it or using an empty application object:

```json
{
  "application": {}
}
```

## See Also

- [Installation Guide](installation.md)
- [D-Bus API](dbus-api.md)
- [Custom Modules](custom-modules.md)
