# go-x11-client

**Description**:
go-x11-client it is a X11 protocol go language binding.

## Dependencies
You can also check the "Depends" provided in the debian/control file.

### Build dependencies
You can also check the "Build-Depends" provided in the debian/control file.

## Installation

### Deepin

Install prerequisites
```
$ sudo apt-get build-dep go-x11-client
```

Build
```
$ GOPATH=/usr/share/gocode make
```

Install
If you have isolated testing build environment (say a docker container), you can install it directly.

```
$ sudo make install
```

generate package files and install go-x11-client with it.
```
$ debuild -uc -us ...
$ sudo dpkg -i ../golang-github-linuxdeepin-go-x11-client.deb
```

## Usage
go get github.com/linuxdeepin/go-x11-client

## Getting help

Any usage issues can ask for help via

* [Gitter](https://gitter.im/orgs/linuxdeepin/rooms)
* [IRC channel](https://webchat.freenode.net/?channels=deepin)
* [Forum](https://bbs.deepin.org)
* [WiKi](https://wiki.deepin.org/)

## Getting involved

We encourage you to report issues and contribute changes.

* [Contribution guide for developers](https://github.com/linuxdeepin/developer-center/wiki/Contribution-Guidelines-for-Developers-en). (English)
* [开发者代码贡献指南](https://github.com/linuxdeepin/developer-center/wiki/Contribution-Guidelines-for-Developers) (中文)

## License

go-x11-client is licensed under [GPL-3.0-or-later](LICENSE).


