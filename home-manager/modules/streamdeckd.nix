{ config, lib, pkgs, ... }:

let
  cfg = config.programs.streamdeckd;

in {
  options.programs.streamdeckd = {
    enable = lib.mkEnableOption "Streamdeckd user daemon";

    package = lib.mkOption {
      type = lib.types.package;
      default = pkgs.streamdeckd;
    };

    settings = lib.mkOption {
      type = lib.types.submodule {
        options = {
          modules = lib.mkOption {
            type = lib.types.nullOr (lib.types.listOf lib.types.str);
            default = null;
            description = "List of modules to load";
          };
          decks = lib.mkOption {
            type = lib.types.listOf (lib.types.submodule {
              options = {
                serial = lib.mkOption {
                  type = lib.types.str;
                  description = "Serial number of the Stream Deck";
                };
                pages = lib.mkOption {
                  type = lib.types.listOf (lib.types.submodule {
                    options = {
                      keys = lib.mkOption {
                        type = lib.types.listOf (lib.types.submodule {
                          options = {
                            application = lib.mkOption {
                              type = lib.types.nullOr (lib.types.attrsOf (lib.types.submodule {
                                options = {
                                  icon = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "Path to icon file";
                                  };
                                  switch_page = lib.mkOption {
                                    type = lib.types.nullOr lib.types.int;
                                    default = null;
                                    description = "Page number to switch to";
                                  };
                                  text = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "Text to display";
                                  };
                                  text_size = lib.mkOption {
                                    type = lib.types.nullOr lib.types.int;
                                    default = null;
                                    description = "Text size";
                                  };
                                  text_alignment = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "Text alignment (left, center, right)";
                                  };
                                  keybind = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "Keybind to trigger";
                                  };
                                  command = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "Command to execute";
                                  };
                                  brightness = lib.mkOption {
                                    type = lib.types.nullOr lib.types.int;
                                    default = null;
                                    description = "Brightness level";
                                  };
                                  url = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "URL to open";
                                  };
                                  key_hold = lib.mkOption {
                                    type = lib.types.nullOr lib.types.int;
                                    default = null;
                                    description = "Key hold duration in milliseconds";
                                  };
                                  obs_command = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "OBS command to execute";
                                  };
                                  obs_command_params = lib.mkOption {
                                    type = lib.types.nullOr (lib.types.attrsOf lib.types.str);
                                    default = null;
                                    description = "OBS command parameters";
                                  };
                                  icon_handler = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "Icon handler name";
                                  };
                                  key_handler = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "Key handler name";
                                  };
                                  icon_handler_fields = lib.mkOption {
                                    type = lib.types.nullOr lib.types.attrs;
                                    default = null;
                                    description = "Icon handler configuration fields";
                                  };
                                  key_handler_fields = lib.mkOption {
                                    type = lib.types.nullOr lib.types.attrs;
                                    default = null;
                                    description = "Key handler configuration fields";
                                  };
                                  shared_handler_fields = lib.mkOption {
                                    type = lib.types.nullOr lib.types.attrs;
                                    default = null;
                                    description = "Shared handler configuration fields";
                                  };
                                };
                              }));
                              default = null;
                              description = "Application-specific key configurations";
                            };
                          };
                        });
                        default = [];
                        description = "List of keys on this page";
                      };
                      knobs = lib.mkOption {
                        type = lib.types.listOf (lib.types.submodule {
                          options = {
                            application = lib.mkOption {
                              type = lib.types.nullOr (lib.types.attrsOf (lib.types.submodule {
                                options = {
                                  icon = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "Path to icon file";
                                  };
                                  text = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "Text to display";
                                  };
                                  text_size = lib.mkOption {
                                    type = lib.types.nullOr lib.types.int;
                                    default = null;
                                    description = "Text size";
                                  };
                                  text_alignment = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "Text alignment (left, center, right)";
                                  };
                                  lcd_handler = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "LCD handler name";
                                  };
                                  knob_or_touch_handler = lib.mkOption {
                                    type = lib.types.nullOr lib.types.str;
                                    default = null;
                                    description = "Knob or touch handler name";
                                  };
                                  lcd_handler_fields = lib.mkOption {
                                    type = lib.types.nullOr lib.types.attrs;
                                    default = null;
                                    description = "LCD handler configuration fields";
                                  };
                                  knob_or_touch_handler_fields = lib.mkOption {
                                    type = lib.types.nullOr lib.types.attrs;
                                    default = null;
                                    description = "Knob or touch handler configuration fields";
                                  };
                                  shared_handler_fields = lib.mkOption {
                                    type = lib.types.nullOr lib.types.attrs;
                                    default = null;
                                    description = "Shared handler configuration fields";
                                  };
                                  knob_press_action = lib.mkOption {
                                    type = lib.types.nullOr (lib.types.submodule {
                                      options = {
                                        switch_page = lib.mkOption {
                                          type = lib.types.nullOr lib.types.int;
                                          default = null;
                                          description = "Page number to switch to";
                                        };
                                        keybind = lib.mkOption {
                                          type = lib.types.nullOr lib.types.str;
                                          default = null;
                                          description = "Keybind to trigger";
                                        };
                                        command = lib.mkOption {
                                          type = lib.types.nullOr lib.types.str;
                                          default = null;
                                          description = "Command to execute";
                                        };
                                        brightness = lib.mkOption {
                                          type = lib.types.nullOr lib.types.int;
                                          default = null;
                                          description = "Brightness level";
                                        };
                                        url = lib.mkOption {
                                          type = lib.types.nullOr lib.types.str;
                                          default = null;
                                          description = "URL to open";
                                        };
                                        obs_command = lib.mkOption {
                                          type = lib.types.nullOr lib.types.str;
                                          default = null;
                                          description = "OBS command to execute";
                                        };
                                        obs_command_params = lib.mkOption {
                                          type = lib.types.nullOr (lib.types.attrsOf lib.types.str);
                                          default = null;
                                          description = "OBS command parameters";
                                        };
                                      };
                                    });
                                    default = null;
                                    description = "Action on knob press";
                                  };
                                  knob_turn_up_action = lib.mkOption {
                                    type = lib.types.nullOr (lib.types.submodule {
                                      options = {
                                        switch_page = lib.mkOption {
                                          type = lib.types.nullOr lib.types.int;
                                          default = null;
                                          description = "Page number to switch to";
                                        };
                                        keybind = lib.mkOption {
                                          type = lib.types.nullOr lib.types.str;
                                          default = null;
                                          description = "Keybind to trigger";
                                        };
                                        command = lib.mkOption {
                                          type = lib.types.nullOr lib.types.str;
                                          default = null;
                                          description = "Command to execute";
                                        };
                                        brightness = lib.mkOption {
                                          type = lib.types.nullOr lib.types.int;
                                          default = null;
                                          description = "Brightness level";
                                        };
                                        url = lib.mkOption {
                                          type = lib.types.nullOr lib.types.str;
                                          default = null;
                                          description = "URL to open";
                                        };
                                        obs_command = lib.mkOption {
                                          type = lib.types.nullOr lib.types.str;
                                          default = null;
                                          description = "OBS command to execute";
                                        };
                                        obs_command_params = lib.mkOption {
                                          type = lib.types.nullOr (lib.types.attrsOf lib.types.str);
                                          default = null;
                                          description = "OBS command parameters";
                                        };
                                      };
                                    });
                                    default = null;
                                    description = "Action on knob turn up";
                                  };
                                  knob_turn_down_action = lib.mkOption {
                                    type = lib.types.nullOr (lib.types.submodule {
                                    options = {
                                      switch_page = lib.mkOption {
                                        type = lib.types.nullOr lib.types.int;
                                        default = null;
                                        description = "Page number to switch to";
                                      };
                                      keybind = lib.mkOption {
                                        type = lib.types.nullOr lib.types.str;
                                        default = null;
                                        description = "Keybind to trigger";
                                      };
                                      command = lib.mkOption {
                                        type = lib.types.nullOr lib.types.str;
                                        default = null;
                                        description = "Command to execute";
                                      };
                                      brightness = lib.mkOption {
                                        type = lib.types.nullOr lib.types.int;
                                        default = null;
                                        description = "Brightness level";
                                      };
                                      url = lib.mkOption {
                                        type = lib.types.nullOr lib.types.str;
                                        default = null;
                                        description = "URL to open";
                                      };
                                      obs_command = lib.mkOption {
                                        type = lib.types.nullOr lib.types.str;
                                        default = null;
                                        description = "OBS command to execute";
                                      };
                                      obs_command_params = lib.mkOption {
                                        type = lib.types.nullOr (lib.types.attrsOf lib.types.str);
                                        default = null;
                                        description = "OBS command parameters";
                                      };
                                    };
                                  });
                                    default = null;
                                    description = "Action on knob turn down";
                                  };
                                };
                              }));
                              default = null;
                              description = "Application-specific knob configurations";
                            };
                          };
                        });
                        default = [];
                        description = "List of knobs on this page";
                      };
                    };
                  });
                  default = [];
                  description = "List of pages for this deck";
                };
              };
            });
            default = [];
            description = "List of Stream Deck configurations";
          };
          obs_connection_info = lib.mkOption {
            type = lib.types.nullOr (lib.types.submodule {
              options = {
                host = lib.mkOption {
                  type = lib.types.nullOr lib.types.str;
                  default = null;
                  description = "OBS WebSocket host";
                };
                port = lib.mkOption {
                  type = lib.types.nullOr lib.types.int;
                  default = null;
                  description = "OBS WebSocket port";
                };
              };
            });
            default = null;
            description = "OBS WebSocket connection information";
          };
        };
      };
      default = {};
      description = "Streamdeckd JSON configuration";
    };
  };

  config = lib.mkIf cfg.enable {

    xdg.configFile."streamdeck-config.json".text =
      builtins.toJSON cfg.settings;

    systemd.user.services.streamdeckd = {
      Unit = {
        Description = "Streamdeck user daemon";
        After = [ "graphical-session.target" ];
      };

      Service = {
        ExecStart =
          "${cfg.package}/bin/streamdeckd --config %h/.config/streamdeck-config.json";
        Restart = "always";
      };

      Install = {
        WantedBy = [ "default.target" ];
      };
    };
  };
}
