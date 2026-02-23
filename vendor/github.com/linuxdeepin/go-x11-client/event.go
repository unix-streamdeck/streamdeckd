// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package x

import (
	"fmt"
)

type GenericEvent []byte

func (ge GenericEvent) GetEventCode() uint8 {
	return ge[0] &^ 0x80
}

func (ge GenericEvent) Real() bool {
	return ge[0]&0x80 == 0
}

func (ge GenericEvent) String() string {
	return fmt.Sprintf("GenericEvent{ EventCode: %d, Real: %v }", ge.GetEventCode(), ge.Real())
}
