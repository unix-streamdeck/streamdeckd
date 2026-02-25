{ config, lib, pkgs, ... }:

let
  cfg = config.hardware.streamdeck;
in
{
  options.hardware.streamdeck.enable =
    lib.mkEnableOption "Streamdeckd hardware support";

  config = lib.mkIf cfg.enable {

    boot.kernelModules = lib.mkAfter [ "uinput" ];

    users.groups.plugdev = {};

    services.udev.packages = [
      (pkgs.writeTextFile {
        name = "streamdeckd-udev-rules";
        destination = "/lib/udev/rules.d/99-streamdeckd.rules";
        text = ''
          KERNEL=="uinput", GROUP="plugdev", MODE:="0660"
          SUBSYSTEM=="input", GROUP="plugdev", MODE="0666"
          SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0060", MODE:="666", GROUP="plugdev"
          SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0063", MODE:="666", GROUP="plugdev"
          SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="006c", MODE:="666", GROUP="plugdev"
          SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="006d", MODE:="666", GROUP="plugdev"
          SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0080", MODE:="666", GROUP="plugdev"
          SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0084", MODE:="666", GROUP="plugdev"
          SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0086", MODE:="666", GROUP="plugdev"
          SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="008f", MODE:="666", GROUP="plugdev"
          SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0090", MODE:="666", GROUP="plugdev"
          SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="009a", MODE:="666", GROUP="plugdev"
          SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="00a5", MODE:="666", GROUP="plugdev"
        '';
      })
    ];
  };
}
