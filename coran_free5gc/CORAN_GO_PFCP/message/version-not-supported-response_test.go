// Copyright 2019-2024 go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"testing"

	"github.com/coranlabs/CORAN_GO_PFCP/message"

	"github.com/coranlabs/CORAN_GO_PFCP/internal/testutil"
)

func TestVersionNotSupportedResponse(t *testing.T) {
	cases := []testutil.TestCase{
		{
			Description: "Normal",
			Structured:  message.NewVersionNotSupportedResponse(seq),
			Serialized: []byte{
				0x20, 0x0b, 0x00, 0x04, 0x11, 0x22, 0x33, 0x00,
			},
		},
	}

	testutil.Run(t, cases, func(b []byte) (testutil.Serializable, error) {
		v, err := message.ParseVersionNotSupportedResponse(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}
