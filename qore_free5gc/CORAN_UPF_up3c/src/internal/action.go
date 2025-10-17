// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

const (
	APPLY_ACT_DROP = 1 << iota
	APPLY_ACT_FORW
	APPLY_ACT_BUFF
)

type ApplyAction struct {
	Flags uint16
}

func (a *ApplyAction) Unmarshal(b []byte) error {
	var v []byte
	if len(b) < 1 {
		return errors.Errorf("ApplyAction Unmarshal: less than 1 bytes")
	} else if len(b) < 2 {
		v = make([]byte, len(b)+1)
		copy(v, b)
	}
	a.Flags = binary.LittleEndian.Uint16(v)
	return nil
}
