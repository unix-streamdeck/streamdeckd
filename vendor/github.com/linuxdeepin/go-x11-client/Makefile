PREFIX = /usr
GOSITE_DIR = ${PREFIX}/share/gocode
GOPKG_PERFIX = github.com/linuxdeepin/go-x11-client
SRC_DIR=${DESTDIR}${GOSITE_DIR}/src/${GOPKG_PERFIX}

all: build

build:
	echo ignore build

print_gopath:
	GOPATH="${CURDIR}/${GOPATH_DIR}:${GOPATH}"

install:
	mkdir -p ${SRC_DIR}
	cp *.go ${SRC_DIR}
	cp -r ext ${SRC_DIR}
	cp -r util ${SRC_DIR}
