{ lib
, buildGoModule
, pkg-config
, libudev-zero
}:

buildGoModule rec {
  pname = "streamdeckd";
  version = "1.0.0";

  src = ./.;

  vendorHash = null;

  nativeBuildInputs = [ pkg-config ];

  buildInputs = [ libudev-zero ];

  meta = with lib; {
    description = "Elgato Streamdeck Driver for Linux";
    homepage = "https://github.com/unix-streamdeck/streamdeckd";
    license = licenses.bsd3;
    platforms = platforms.linux;
    maintainers = [ ];
  };
}
