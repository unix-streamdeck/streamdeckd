#!/bin/sh

set -x
set -e

goClient='python2 tools/go_client.py'
reqCodeGen=tools/gen/gen

$goClient -o xproto_auto.go /usr/share/xcb/xproto.xml
env GOPACKAGE=x $reqCodeGen xproto.go > xproto_auto_req.go

$goClient -o ext/record/auto.go /usr/share/xcb/record.xml
env GOPACKAGE=record $reqCodeGen -e ext/record/record.go > ext/record/record_req_auto.go

$goClient -o ext/test/auto.go /usr/share/xcb/xtest.xml
env GOPACKAGE=test $reqCodeGen -e ext/test/test.go > ext/test/test_req_auto.go

$goClient -o ext/damage/auto.go /usr/share/xcb/damage.xml
env GOPACKAGE=damage $reqCodeGen -e ext/damage/damage.go > ext/damage/damage_req_auto.go

$goClient -o ext/composite/auto.go /usr/share/xcb/composite.xml
env GOPACKAGE=composite $reqCodeGen -e ext/composite/composite.go > ext/composite/composite_req_auto.go

$goClient -o ext/screensaver/auto.go /usr/share/xcb/screensaver.xml
env GOPACKAGE=screensaver $reqCodeGen -e ext/screensaver/screensaver.go > ext/screensaver/screensaver_req_auto.go

$goClient -o ext/dpms/auto.go /usr/share/xcb/dpms.xml
env GOPACKAGE=dpms $reqCodeGen -e ext/dpms/dpms.go > ext/dpms/dpms_req_auto.go

$goClient -o ext/randr/auto.go /usr/share/xcb/randr.xml
env GOPACKAGE=randr $reqCodeGen -e -extra-exts render ext/randr/randr.go > ext/randr/randr_req_auto.go

$goClient -o ext/xfixes/auto.go /usr/share/xcb/xfixes.xml
env GOPACKAGE=xfixes $reqCodeGen -e ext/xfixes/xfixes.go > ext/xfixes/xfixes_req_auto.go

$goClient -o ext/input/auto.go /usr/share/xcb/xinput.xml
env GOPACKAGE=input $reqCodeGen -e ext/input/input.go > ext/input/input_req_auto.go
env GOPACKAGE=input $reqCodeGen -e ext/input/input1.go > ext/input/input1_req_auto.go

$goClient -p ge -o ext/ge/auto.go /usr/share/xcb/ge.xml
env GOPACKAGE=ge $reqCodeGen -e ext/ge/ge.go > ext/ge/ge_req_auto.go

$goClient -o ext/render/auto.go /usr/share/xcb/render.xml
env GOPACKAGE=render $reqCodeGen -e ext/render/render.go > ext/render/render_req_auto.go

$goClient -o ext/shm/auto.go /usr/share/xcb/shm.xml
env GOPACKAGE=shm $reqCodeGen -e ext/shm/shm.go > ext/shm/shm_req_auto.go

$goClient -o ext/bigrequests/auto.go /usr/share/xcb/bigreq.xml
env GOPACKAGE=bigrequests $reqCodeGen -e ext/bigrequests/bigreq.go > ext/bigrequests/bigreq_req_auto.go

$goClient -o ext/xkb/auto.go /usr/share/xcb/xkb.xml
env GOPACKAGE=xkb $reqCodeGen -e ext/xkb/xkb.go > ext/xkb/xkb_req_auto.go
