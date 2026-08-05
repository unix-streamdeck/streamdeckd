{ config, lib, pkgs, ... }:

let
  cfg = config.programs.streamdeckd;

  inherit (lib) mkOption types;

  strOption = description:
    mkOption {
      type = types.nullOr types.str;
      default = null;
      inherit description;
    };

  attrsOption = description:
    mkOption {
      type = types.nullOr types.attrs;
      default = null;
      inherit description;
    };

  inputEventIds = {
    knob_ccw = 0;
    knob_cw = 1;
    knob_press = 2;
    screen_short_tap = 3;
    screen_long_tap = 4;
    screen_swipe = 5;
    key_press = 6;
    key_release = 7;
  };

  layerModule = types.submodule {
    options = {
      handler = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = ''
          Name of the handler backing this layer, for example "Image", "Text",
          "Command" or the name of a module. Leave unset only when {option}`global`
          is true.
        '';
      };

      fields = attrsOption "Configuration fields for the handler. Both sides of the layer see the same fields.";

      visual = mkOption {
        type = types.bool;
        default = false;
        description = "Whether this layer draws. Requires a handler that renders.";
      };

      input = mkOption {
        type = types.bool;
        default = false;
        description = ''
          Whether this layer handles input. Requires a handler that takes input.
          A layer with both {option}`visual` and {option}`input` set uses one
          handler instance for both, which is how a handler shares state between
          what it draws and what it does.
        '';
      };

      input_events = mkOption {
        type = types.nullOr (types.listOf (types.enum (builtins.attrNames inputEventIds)));
        default = null;
        example = [ "knob_press" ];
        description = ''
          Restricts which events reach this layer. Unset means every event, which
          for a key means both press and release. This is how a knob's press,
          clockwise and anticlockwise actions are told apart.
        '';
      };

      global = mkOption {
        type = types.bool;
        default = false;
        description = ''
          Marks the position in this stack where the key or knob's application
          independent layers are spliced in. A layer with this set carries no
          handler of its own. Without it, those layers are placed at the bottom.
        '';
      };

      label = strOption "Optional name for this layer, used by the editor only.";
    };
  };

  layersOption = description:
    mkOption {
      type = types.listOf layerModule;
      default = [ ];
      inherit description;
    };

  keyBackgroundOptions = {
    background = strOption "Background image for the key, or the name of a background handler.";
    key_background_handler_fields = attrsOption "Key background handler configuration fields.";
  };

  panelBackgroundOptions = {
    touch_panel_background = strOption "Background image for the touch panel, or the name of a background handler.";
    touch_panel_background_handler_fields = attrsOption "Touch panel background handler configuration fields.";
  };

  gridBackgroundOptions = panelBackgroundOptions // {
    key_grid_background = strOption "Background image spanning the key grid, or the name of a background handler.";
    key_grid_background_handler_fields = attrsOption "Key grid background handler configuration fields.";
  };

  applicationOption = profileModule: kind:
    mkOption {
      type = types.nullOr (types.attrsOf profileModule);
      default = null;
      description = ''
        Per application ${kind} configuration, keyed by application name. The
        empty string is the default profile, used when no other profile matches
        the focused application. A matching profile replaces the default rather
        than adding to it.
      '';
    };

  keyConfigModule = types.submodule {
    options = keyBackgroundOptions // {
      layers = layersOption "Foreground layers for this profile, painted bottom first.";
    };
  };

  keyModule = types.submodule {
    options = keyBackgroundOptions // {
      layers = layersOption ''
        Foreground layers that apply whichever application is focused. They are
        spliced into the active profile's stack at its global marker, or placed
        at the bottom when it has none.
      '';
      application = applicationOption keyConfigModule "key";
    };
  };

  knobConfigModule = types.submodule {
    options = panelBackgroundOptions // {
      layers = layersOption "Foreground layers for this profile, painted bottom first.";
    };
  };

  knobModule = types.submodule {
    options = panelBackgroundOptions // {
      layers = layersOption ''
        Foreground layers that apply whichever application is focused. They are
        spliced into the active profile's stack at its global marker, or placed
        at the bottom when it has none.
      '';
      application = applicationOption knobConfigModule "knob";
    };
  };

  pageModule = types.submodule {
    options = gridBackgroundOptions // {
      keys = mkOption {
        type = types.listOf keyModule;
        default = [ ];
        description = "Keys on this page, in device order.";
      };
      knobs = mkOption {
        type = types.listOf knobModule;
        default = [ ];
        description = "Knobs and LCD segments on this page, in device order.";
      };
    };
  };

  deckModule = types.submodule {
    options = gridBackgroundOptions // {
      serial = mkOption {
        type = types.str;
        description = "Serial number of the Stream Deck.";
      };
      pages = mkOption {
        type = types.listOf pageModule;
        default = [ ];
        description = "Pages for this deck.";
      };
    };
  };

  obsModule = types.submodule {
    options = {
      host = strOption "OBS WebSocket host.";
      port = mkOption {
        type = types.nullOr types.int;
        default = null;
        description = "OBS WebSocket port.";
      };
    };
  };


  renderLayer = layer:
    lib.filterAttrs (_: value: value != null) (
      layer // lib.optionalAttrs (layer.input_events != null) {
        input_events = map (name: inputEventIds.${name}) layer.input_events;
      }
    );

  renderLayers = map renderLayer;

  renderProfiles = profiles:
    if profiles == null then
      null
    else
      lib.mapAttrs (_: profile: profile // { layers = renderLayers profile.layers; }) profiles;

  renderElement = element:
    element // {
      layers = renderLayers element.layers;
      application = renderProfiles element.application;
    };

  renderPage = page:
    page // {
      keys = map renderElement page.keys;
      knobs = map renderElement page.knobs;
    };

  renderDeck = deck: deck // { pages = map renderPage deck.pages; };

  settingsJson = builtins.toJSON (cfg.settings // {
    decks = map renderDeck cfg.settings.decks;
  });

in {
  options.programs.streamdeckd = {
    enable = lib.mkEnableOption "Streamdeckd user daemon";

    package = mkOption {
      type = types.package;
      default = pkgs.streamdeckd;
      description = "The streamdeckd package to use.";
    };

    settings = mkOption {
      type = types.submodule {
        options = {
          modules = mkOption {
            type = types.nullOr (types.listOf types.str);
            default = null;
            description = "Paths of module plugins to load.";
          };
          decks = mkOption {
            type = types.listOf deckModule;
            default = [ ];
            description = "Stream Deck configurations.";
          };
          obs_connection_info = mkOption {
            type = types.nullOr obsModule;
            default = null;
            description = "OBS WebSocket connection information.";
          };
        };
      };
      default = { };
      description = "Streamdeckd JSON configuration.";
    };
  };

  config = lib.mkIf cfg.enable {

    xdg.configFile."streamdeck-config.json".text = settingsJson;

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

      Install = { WantedBy = [ "default.target" ]; };
    };
  };
}
