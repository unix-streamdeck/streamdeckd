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

The configuration file streamdeckd uses is a JSON file found at `~/.streamdeck-config.json`

An example config would be something like:

```json
{
  "pages": [
    [
      {
        "switch_page": 1,
        "icon": "~/icon.png"
      }
    ]
  ]
}
```

The outer array is the list of pages, the inner array is the list of button on that page, with the buttons going in a right to left order.

The actions you can have on a button are:

- `command`: runs a native shell command, something like `notify-send "Hello World"`
- `keybind`: simulates the indicated keybind via xdtotool
- `url`: opens a url in your default browser via xdg
- `brightness`: set the brightness of the streamdeck as a percentage
- `switch_page`: change the active page on the streamdeck


### D-Bus

There is a D-Bus interface built into the daemon, the service name and interface for D-Bus are `com.unixstreamdeck.streamdeckd` and `com/unixstreamdeck/streamdeckd` respectively, and is made up of the following methods/signals

#### Methods

- GetConfig  - returns the current running config
- SetConfig  - sets the config, without saving to disk, takes in Stringified json, returns an error if anything breaks
- ReloadConfig  - reloads the config from disk
- GetDeckInfo  - Returns information about the active streamdeck in the format of 
```json
{
  "icon_size": 72,
  "rows": 3,
  "cols": 5,
  "page": 0
}
```
- SetPage - Set the page on the streamdeck to the number passed to it, returns an error if anything breaks
- CommitConfig  - Commits the currently active config to disk, returns an error if anything breaks

#### Signals

- Page - sends the number of the page switched to on the StreamDeck


### Custom Handlers



#### Button Press Handler


#### Icon Handler