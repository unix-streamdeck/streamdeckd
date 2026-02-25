# streamdeckd

A Linux daemon for Elgato Stream Deck devices. Control your Stream Deck with custom buttons, macros, and per-application configurations.

## What is it?

streamdeckd runs in the background and manages your Elgato Stream Deck devices on Linux. It lets you:

- **Create custom button layouts** with icons and actions
- **Switch between pages** for unlimited buttons
- **Set per-application profiles** - buttons change based on what app you're using
- **Control multiple Stream Decks** at once
- **Extend with plugins** for advanced functionality

## Quick Start


1. **Set up device permissions:**
   ```bash
   # Copy udev rules
   sudo curl -o /etc/udev/rules.d/50-elgato.rules \
     https://raw.githubusercontent.com/unix-streamdeck/streamdeckd/master/50-elgato.rules

   # Reload rules
   sudo udevadm control --reload-rules
   ```

2. **Build and run:**
   ```bash
   go build
   ./streamdeckd
   ```

3. **Configure your buttons** by editing `$XDG_CONFIG_HOME/.streamdeck-config.json`

## Button Actions

Each button can perform different actions:

- **Command** - Run any shell command
- **Keybind** - Simulate keyboard shortcuts
- **URL** - Open websites
- **Page Switch** - Navigate between button layouts
- **Brightness** - Adjust Stream Deck brightness
- **Key Hold** - Simulate holding down a key
- **Custom Plugin Actions** - Any action you could want

## Example Configuration

```json
{
  "decks": [{
    "serial": "AB12C3D45678",
    "pages": [[
      {
        "application": {
          "": {
            "command": "notify-send 'Hello!'",
            "icon": "~/icons/hello.png"
          }
        }
      }
    ]]
  }]
}
```

## GUI Configuration Tool

For a graphical interface, use **[streamdeckui-wails](https://github.com/unix-streamdeck/streamdeckui-wails)**.

## Documentation

- **[Installation Guide](docs/installation.md)** - Detailed setup instructions
- **[Configuration Guide](docs/configuration.md)** - Complete configuration reference
- **[D-Bus API](docs/dbus-api.md)** - Programmatic control via D-Bus
- **[Custom Modules](docs/custom-modules.md)** - Creating plugins and extensions

## Community

- **Report issues:** [GitHub Issues](https://github.com/unix-streamdeck/streamdeckd/issues)
- **Example plugins:** [example-plugins repository](https://github.com/unix-streamdeck/example-plugins)

## License

See [LICENSE](LICENSE) for details.
