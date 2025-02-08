// Copyright 2019-2024 go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package testutil is an internal package to be used for unit tests. Don't use this.
package testutil

import (
	"testing"

	"github.com/coranlabs/CORAN_GO_PFCP/message"
	"github.com/google/go-cmp/cmp"
)

// Serializable is just for testing Message. Don't use this.
type Serializable interface {
	Marshal() ([]byte, error)
	MarshalLen() int
}

// TestCase is just for testing Message. Don't use this.
type TestCase struct {
	Description string
	Structured  Serializable
	Serialized  []byte
}

// ParseFunc is just for testing Message. Don't use this.
type ParseFunc func([]byte) (Serializable, error)

// TestBearerInfo is just for testing Message. Don't use this.
var TestBearerInfo = struct {
	SEID uint64
	Seq  uint32
}{0x1122334455667788, 0x00000001}

// Run is just for testing Message. Don't use this.
func Run(t *testing.T, cases []TestCase, decode ParseFunc) {
	t.Helper()

	for _, c := range cases {
		t.Run(c.Description, func(t *testing.T) {
			t.Run("Parse", func(t *testing.T) {
				v, err := decode(c.Serialized)
				if err != nil {
					t.Fatal(err)
				}

				if got, want := v, c.Structured; !cmp.Equal(got, want) {
					t.Fail()
				}
			})

			t.Run("Marshal", func(t *testing.T) {
				b, err := c.Structured.Marshal()
				if err != nil {
					t.Fatal(err)
				}

				if got, want := b, c.Serialized; !cmp.Equal(got, want) {
					t.Fail()
				}
			})

			t.Run("Len", func(t *testing.T) {
				if got, want := c.Structured.MarshalLen(), len(c.Serialized); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
			})

			t.Run("Interface", func(t *testing.T) {
				// Ignore *Header and Generic in this tests.
				if _, ok := c.Structured.(*message.Header); ok {
					return
				}

				if _, ok := c.Structured.(*message.Generic); ok {
					return
				}

				decoded, err := message.Parse(c.Serialized)
				if err != nil {
					t.Fatal(err)
				}

				if got, want := decoded.Version(), c.Structured.(message.Message).Version(); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
				if got, want := decoded.MessageType(), c.Structured.(message.Message).MessageType(); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
				if got, want := decoded.MessageTypeName(), c.Structured.(message.Message).MessageTypeName(); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
				if got, want := decoded.SEID(), c.Structured.(message.Message).SEID(); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
				if got, want := decoded.Sequence(), c.Structured.(message.Message).Sequence(); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
			})
		})
	}
}
