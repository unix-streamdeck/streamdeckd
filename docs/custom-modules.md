# Custom Modules

## Overview

Custom modules are Go plugins (`.so` files) that expose handlers for:
- **Icon Handlers:** Generate dynamic button icons
- **Key Handlers:** Implement custom button press actions
- **Lcd Handlers:** Control the StreamDeck Plus' touch screen LCD Display
- **Touch or Knob Handlers:** Handle touch and knob events on the StreamDeck Plus

Modules are loaded at startup and can be configured through the standard JSON config file.

## Getting Started

### Prerequisites

- Go 1.25+
- streamdeckd API package: `github.com/unix-streamdeck/api/v2`

### Basic Module Structure

Every module must implement the `GetModule()` function:

> N.B - For the latest reference for all the types referenced below, check the api module:
> 
> [Module](https://github.com/unix-streamdeck/api/blob/master/Module.go)
> 
> [Handlers](https://github.com/unix-streamdeck/api/blob/master/handler.go)

```go
package main

import "github.com/unix-streamdeck/api/v2"

func GetModule() api.Module {
    return api.Module{
        Name:       "MyModule",
        NewIcon:    func() api.IconHandler { return &MyIconHandler{} },
        NewKey:     func() api.KeyHandler { return &MyKeyHandler{} },
        IconFields: []api.Field{ /* ... */ },
        KeyFields:  []api.Field{ /* ... */ },
    }
}
```

## See Also

- [Configuration Guide](configuration.md)
- [D-Bus API](dbus-api.md)
- [Example Plugins Repository](https://github.com/unix-streamdeck/example-plugins)
- [Go Plugin Package Documentation](https://pkg.go.dev/plugin)
