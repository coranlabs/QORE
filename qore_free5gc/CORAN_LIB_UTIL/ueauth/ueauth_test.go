package ueauth

import (
	"bytes"
	"encoding/hex"
	"testing"
)

type TestKDF struct {
	NetworkName string
	SQNxorAK    string
	CK          string
	IK          string
	FC          string
	DerivedKey  string
}

func TestGetKDFValue(t *testing.T) {
	// Only the network name is different, which should yet yield different derived results
	// RFC 5448 Test Vector 1
	TestKDFTable := []TestKDF{
		// SUCCESS case
		{
			NetworkName: "WLAN",
			SQNxorAK:    "bb52e91c747a",
			CK:          "5349fbe098649f948f5d2e973a81c00f",
			IK:          "9744871ad32bf9bbd1dd5ce54e3e2e5a",
			FC:          FC_FOR_CK_PRIME_IK_PRIME_DERIVATION,
			DerivedKey:  "0093962d0dd84aa5684b045c9edffa04ccfc230ca74fcc96c0a5d61164f5a76c",
		},
		// FAILURE case
		{
			NetworkName: "WLANNNNNNNNNNNNNNN",
			SQNxorAK:    "bb52e91c747a",
			CK:          "5349fbe098649f948f5d2e973a81c00f",
			IK:          "9744871ad32bf9bbd1dd5ce54e3e2e5a",
			FC:          FC_FOR_CK_PRIME_IK_PRIME_DERIVATION,
			DerivedKey:  "0093962d0dd84aa5684b045c9edffa04ccfc230ca74fcc96c0a5d61164f5a76c",
		},
	}

	for i, tc := range TestKDFTable {
		P0 := []byte(tc.NetworkName)

		P1, err := hex.DecodeString(tc.SQNxorAK)
		if err != nil {
			t.Errorf("TestGetKDFValue[%d] error: %+v\n", i, err)
		}

		ckik, err := hex.DecodeString(tc.CK + tc.IK)
		if err != nil {
			t.Errorf("TestGetKDFValue[%d] error: %+v\n", i, err)
		}

		val, err := GetKDFValue(ckik, tc.FC, P0, KDFLen(P0), P1, KDFLen(P1))
		if err != nil {
			t.Errorf("TestGetKDFValue[%d] error: %+v\n", i, err)
		}
		// fmt.Printf("val = %x\n", val)

		dk, err := hex.DecodeString(tc.DerivedKey)
		if err != nil {
			t.Errorf("TestGetKDFValue[%d] error: %+v\n", i, err)
		}
		if (i == 0 && !bytes.Equal(val, dk)) || (i == 1 && bytes.Equal(val, dk)) {
			t.Errorf("TestGetKDFValue[%d] failed\n", i)
		}
	}
}
