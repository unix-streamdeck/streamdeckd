# go-x11-client

**Description**:
go-x11-client是一个用go语言实现与X11协议进行绑定的项目。

## 依赖
请查看“debian/control”文件中提供的“Depends”。

### 编译依赖
请查看“debian/control”文件中提供的“Build-Depends”。

## 安装

### Deepin

go-x11-client需要预安装以下包
```
$ sudo apt-get build-dep go-x11-client
```

构建
```
$ GOPATH=/usr/share/gocode make
```

安装
如果你有独立的测试构建环境（比如一个 docker 容器），你可以直接安装它。

```
$ sudo make install
```

生成包文件并安装 go-x11-client
```
$ debuild -uc -us ...
$ sudo dpkg -i ../golang-github-linuxdeepin-go-x11-client.deb
```

## 用法
go get github.com/linuxdeepin/go-x11-client

## 获得帮助

如果您遇到任何其他问题，您可能还会发现这些渠道很有用：

* [Gitter](https://gitter.im/orgs/linuxdeepin/rooms)
* [IRC channel](https://webchat.freenode.net/?channels=deepin)
* [Forum](https://bbs.deepin.org)
* [WiKi](https://wiki.deepin.org/)

## 贡献指南

我们鼓励您报告问题并做出更改

* [Contribution guide for developers](https://github.com/linuxdeepin/developer-center/wiki/Contribution-Guidelines-for-Developers-en). (English)
* [开发者代码贡献指南](https://github.com/linuxdeepin/developer-center/wiki/Contribution-Guidelines-for-Developers) (中文)

## License

go-x11-client项目在 [GPL-3.0-or-later](LICENSE)下发布。


