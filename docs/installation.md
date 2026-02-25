# Installation Guide

Complete instructions for installing and setting up streamdeckd on Linux.

## NixOS Installation

streamdeckd includes native NixOS support with modules for both system-level and user-level configuration.

### Using Flakes (Recommended)

Add streamdeckd to your flake inputs:

```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    streamdeckd = {
      url = "github:unix-streamdeck/streamdeckd";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { nixpkgs, streamdeckd, ... }: {
    nixosConfigurations.yourhostname = nixpkgs.lib.nixosSystem {
      modules = [
        streamdeckd.nixosModules.default
        {
          hardware.streamdeck.enable = true;
          users.users.yourname.extraGroups = [ "plugdev" ];
        }
      ];
    };

    homeConfigurations.yourname = home-manager.lib.homeManagerConfiguration {
      modules = [
        streamdeckd.homeManagerModules.default
        {
          programs.streamdeckd = {
            enable = true;
            package = streamdeckd.packages.${system}.default;
            settings = {
              decks = [
                {
                  serial = "YOUR_SERIAL";
                  pages = [ ];
                }
              ];
            };
          };
        }
      ];
    };
  };
}
```

### What Gets Configured

**NixOS Module** (`hardware.streamdeck.enable`):
- Loads the `uinput` kernel module
- Creates the `plugdev` group
- Installs udev rules for all Stream Deck models
- Sets correct permissions for device access

**Home Manager Module** (`programs.streamdeckd.enable`):
- Creates a systemd user service that starts with your session
- Generates the JSON configuration file at `~/.config/.streamdeck-config.json`
- Provides type-safe configuration options
- Automatically restarts on failure

### Applying Configuration

```bash
# Rebuild your system
sudo nixos-rebuild switch

# Rebuild your home configuration
home-manager switch

# Check service status
systemctl --user status streamdeckd
```

## Prerequisites

### Optional Software

- **streamdeckui** - GUI configuration tool (highly recommended)
- **Go 1.25+** - Only needed if building from source

## udev Rules Setup

udev rules allow non-root users to access Stream Deck devices.

### Method 1: Automatic (Recommended)

Download and install the udev rules file:

```bash
sudo curl -o /etc/udev/rules.d/50-elgato.rules \
  https://raw.githubusercontent.com/unix-streamdeck/streamdeckd/master/50-elgato.rules
sudo udevadm control --reload-rules
```

### Method 2: Manual

Create `/etc/udev/rules.d/50-elgato.rules` with the following content:

```
SUBSYSTEM=="input", GROUP="input", MODE="0666"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0060", MODE:="666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0063", MODE:="666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="006c", MODE:="666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="006d", MODE:="666", GROUP="plugdev"
```

These rules cover:
- `0060` - Stream Deck Original (15 keys)
- `0063` - Stream Deck Mini (6 keys)
- `006c` - Stream Deck XL (32 keys)
- `006d` - Stream Deck MK.2 (15 keys)

Reload udev rules:
```bash
sudo udevadm control --reload-rules
```

### Verify udev Rules

Unplug and replug your Stream Deck, then verify permissions:

```bash
ls -l /dev/input/by-id/*Stream*
```

You should see `rw-rw-rw-` permissions.

## Building from Source

### Clone the Repository

```bash
git clone https://github.com/unix-streamdeck/streamdeckd.git
cd streamdeckd
```

### Build

```bash
go build
```

This creates the `streamdeckd` binary in the current directory.

### Install (Optional)

Copy to your local bin directory:

```bash
mkdir -p ~/.local/bin
cp streamdeckd ~/.local/bin/
```

Or install system-wide:

```bash
sudo cp streamdeckd /usr/local/bin/
```

## Running streamdeckd

### Manual Start

```bash
./streamdeckd
```

Or with a custom config file:

```bash
./streamdeckd -config /path/to/config.json
```

### Automatic Start (systemd)

Create `~/.config/systemd/user/streamdeckd.service`:

```ini
[Unit]
Description=Stream Deck Daemon
After=graphical-session.target

[Service]
Type=simple
ExecStart=%h/.local/bin/streamdeckd
Restart=on-failure
RestartSec=5

[Install]
WantedBy=default.target
```

Enable and start:

```bash
systemctl --user daemon-reload
systemctl --user enable streamdeckd.service
systemctl --user start streamdeckd.service
```

Check status:

```bash
systemctl --user status streamdeckd.service
```

View logs:

```bash
journalctl --user -u streamdeckd.service -f
```

## Initial Configuration

On first run, streamdeckd looks for a config file at:

```
$XDG_CONFIG_HOME/.streamdeck-config.json
```

If `XDG_CONFIG_HOME` is not set, it defaults to `~/.config/.streamdeck-config.json`.

### Option 1: Use streamdeckui (Recommended)

Install and run [streamdeckui](https://github.com/unix-streamdeck/streamdeckui) to configure your Stream Deck with a graphical interface.

### Option 2: Manual Configuration

Create a basic config file:

```bash
cat > ~/.config/.streamdeck-config.json << 'EOF'
{
  "modules": [],
  "decks": []
}
EOF
```

See the [Configuration Guide](configuration.md) for detailed configuration options.

## Troubleshooting

### Stream Deck Not Detected

1. **Check USB connection:**
   ```bash
   lsusb | grep 0fd9
   ```
   Should show Elgato device.

2. **Verify udev rules:**
   ```bash
   udevadm info -a -n /dev/hidraw0 | grep -i elgato
   ```

3. **Check permissions:**
   ```bash
   ls -l /dev/hidraw*
   ```

4. **Reload udev and reconnect device:**
   ```bash
   sudo udevadm control --reload-rules
   sudo udevadm trigger
   ```

### streamdeckd Won't Start

1. **Check for duplicate instances:**
   ```bash
   ps aux | grep streamdeckd
   ```
   streamdeckd prevents multiple instances. Kill any existing processes.

2. **Check D-Bus connection:**
   ```bash
   dbus-send --session --print-reply \
     --dest=org.freedesktop.DBus \
     /org/freedesktop/DBus \
     org.freedesktop.DBus.ListNames
   ```

3. **Run with verbose output:**
   ```bash
   ./streamdeckd 2>&1 | tee streamdeckd.log
   ```


### Wayland Issues

If using Wayland, some features may require additional setup:

- **KDE Plasma:** Should work with recent versions
- **Hyprland:** Supported
- **GNOME:** May have limitations with window detection

Consider using X11 session for full compatibility.

## Next Steps

- Read the [Configuration Guide](configuration.md) to set up your buttons
- Install [streamdeckui-wails](https://github.com/unix-streamdeck/streamdeckui-wails) for easier configuration
- Explore the [D-Bus API](dbus-api.md) for programmatic control
- Create [Custom Modules](custom-modules.md) to extend functionality
