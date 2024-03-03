module github.com/unix-streamdeck/streamdeckd

go 1.19

require (
	github.com/bendahl/uinput v1.7.0
	github.com/christopher-dG/go-obs-websocket v0.0.0-20200720193653-c4fed10356a5
	github.com/godbus/dbus/v5 v5.0.4-0.20200513180336-df5ef3eb7cca
	github.com/linuxdeepin/go-x11-client v0.0.0-20230710064023-230ea415af17
	github.com/shirou/gopsutil/v3 v3.21.9
	github.com/unix-streamdeck/api v1.0.1
	github.com/unix-streamdeck/driver v0.0.0-20211119182210-fc6b90443bcd
	golang.org/x/sync v0.1.0

)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/bearsh/hid v1.4.2-0.20220627123055-35af594cb5a7 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/tklauser/numcpus v0.3.0 // indirect
	golang.org/x/image v0.0.0-20201208152932-35266b937fa6 // indirect
	golang.org/x/sys v0.0.0-20220624220833-87e55d714810 // indirect
)

replace github.com/unix-streamdeck/api v1.0.1 => ../api

replace github.com/unix-streamdeck/driver v0.0.0-20211119182210-fc6b90443bcd => ../driver
