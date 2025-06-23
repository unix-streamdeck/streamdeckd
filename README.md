# Streamdeckd

### Installation

- create the file `/etc/udev/rules.d/50-elgato.rules` with the following config

```  
SUBSYSTEM=="input", GROUP="input", MODE="0666"  
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0060", MODE:="666", GROUP="plugdev"  
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0063", MODE:="666", GROUP="plugdev"  
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="006c", MODE:="666", GROUP="plugdev"  
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="006d", MODE:="666", GROUP="plugdev"  
```  

- run `sudo udevadm control --reload-rules` to reload the udev rules

Then xdotool will be required to simulate keypresses, to install this run:

#### Arch

`sudo pacman -S xdotool`

#### Debian based

`sudo apt install xdotool`

### Configuration

#### Manual configuration

## Warning:

If you are updating from v1.0.0, the config file is now being set in the location as below, instead of where it used to be, in the home dir, either consider moving the config file to that dir, or running streamdeckd with the `-config` flag, which allows you to point to a config file in a custom location

---


The configuration file streamdeckd uses is a JSON file found at `$XDG_CONFIG_HOME/.streamdeck-config.json`

An example config would be something like:

```json
{
  "modules": [
    "/home/user/module.so"
  ],
  "decks": [
    {
      "serial": "AB12C3D45678",
      "pages": [
        [
          {
            "switch_page": 1,
            "icon": "~/icon.png"
          }
        ]
      ]
    }
  ]
}
```

At the top is the list of custom modules, these are go plugins in the .so format, following that is the list of deck
objects, each represents a different streamdeck device, and contains its serial, and its list of pages

The outer array in a deck is the list of pages, the inner array is the list of button on that page, with the buttons
going in a right to left order.

The actions you can have on a button are:

- `command`: runs a native shell command, something like `notify-send "Hello World"`
- `keybind`: simulates the indicated keybind via xdtotool
- `url`: opens a url in your default browser via xdg
- `brightness`: set the brightness of the streamdeck as a percentage
- `switch_page`: change the active page on the streamdeck

### D-Bus

There is a D-Bus interface built into the daemon, the service name and interface for D-Bus
are `com.unixstreamdeck.streamdeckd` and `com/unixstreamdeck/streamdeckd` respectively, and is made up of the following
methods/signals

#### Methods

- GetConfig - returns the current running config
- SetConfig - sets the config, without saving to disk, takes in Stringified json, returns an error if anything breaks
- ReloadConfig - reloads the config from disk
- GetDeckInfo - Returns information about all the active streamdecks in the format of

```json
[
  {
    "icon_size": 72,
    "rows": 3,
    "cols": 5,
    "page": 0,
    "serial": "AB12C3D45678"
  }
]
```

- SetPage - Set the page on the streamdeck to the number passed to it, returns an error if anything breaks
- CommitConfig - Commits the currently active config to disk, returns an error if anything breaks
- GetModules - Get the list of loaded modules, and the config fields those modules use
- PressButton - Simulates a button press on the streamdeck device, consumes a device serial, and a key index


#### Signals

- Page - sends the number of the page switched to on the StreamDeck

### Custom Modules

To create custom modules, I suggest looking at the gif, counter, and time modules in the example handlers package in streamdeckd, they should be in a file with the GetModule method as shown below

#### Loading Modules into streamdeckd

Modules require a method on them in the main package called "GetModule" that returns an instance of [handler.Module](https://github.com/unix-streamdeck/streamdeckd/blob/575e672c26f275d35a016be6406ceb8480ccfff5/handlers/handlers.go#L9) e.g

```go
package main

type CustomIconHandler struct {
	
}
...

type CustomKeyHandler struct {
	
}
...

func GetModule() handlers.Module {
	return handlers.Module{
		Name: "CustomModule", // the name that will be used in the icon_handler/key_handler field in the config, and that will be shown in the handler dropdown in streamdeckui
		NewIcon: func() api.IconHandler { return &CustomerIconHandler{}}, // Method to create a new instance of the Icon handler, if left empty, streamdeckui will not include it in the icon handler fields
		NewKey: func() api.KeyHandler { return &CustomerKeyHandler{}}, // Method to create a new instance of the Key Handler, if left empty, streamdeckui will not include it in the key handler fields
		IconFields: []api.Field{ // list of fields to be shown in streamdeckui when the icon handler is selected
			{
				Title: "Icon", // name of field to show in UI
				Name: "icon", // name of field that will be included in the iconHandlerFields map
				Type: "File" // type of input to show on streamdeckui, options are Text, File, TextAlignment, and Number
				FileTypes: []string{".png", ".jpg"} // Allowed file types if a File input type is used
			}
		},
		KeyFields: []api.Field{}, // Same as IconFields
	}
}

```
