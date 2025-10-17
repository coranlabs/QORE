package tlv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	testCases := []struct {
		name      string
		instance  interface{}
		expectHex string
	}{
		{
			name: "int64",
			instance: struct {
				Value int64 `tlv:"30"`
			}{Value: 32324},
			expectHex: "001E00080000000000007E44",
		},
		{
			name: "int32",
			instance: struct {
				Value int32 `tlv:"31"`
			}{Value: 32324},
			expectHex: "001F000400007E44",
		},
		{
			name: "int16",
			instance: struct {
				Value int16 `tlv:"32"`
			}{Value: 32324},
			expectHex: "002000027E44",
		},
		{
			name: "int8",
			instance: struct {
				Value int8 `tlv:"33"`
			}{Value: 68},
			expectHex: "0021000144",
		},
		{
			name: "slice of struct",
			instance: struct {
				List []struct {
					Name     []byte `tlv:"20"`
					Sequence uint16 `tlv:"40"`
				} `tlv:"80"`
			}{
				List: []struct {
					Name     []byte `tlv:"20"`
					Sequence uint16 `tlv:"40"`
				}{
					{Name: []byte("Hello"), Sequence: 1},
					{Name: []byte("World"), Sequence: 2},
					{Name: []byte("coranlabs"), Sequence: 3},
				},
			},
			expectHex: "0050000F0014000548656C6C6F002800020001" +
				"0050000F00140005576F726C64002800020002" +
				"005000110014000766726565356763002800020003",
		},
		{
			name: "slice of binary",
			instance: struct {
				List []BinaryMarshalTest `tlv:"123"`
			}{
				List: []BinaryMarshalTest{
					{
						Value: 1100,
					},
					{
						Value: 1200,
					},
					{
						Value: 3244,
					},
				},
			},
			expectHex: "007B000431313030007B000431323030007B000433323434",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// testInstance := reflect.New(reflect.TypeOf(tc.instance)).Interface()
			buf, err := Marshal(tc.instance)
			require.NoError(t, err)
			require.Equal(t, tc.expectHex, fmt.Sprintf("%X", buf))
		})
	}
}
