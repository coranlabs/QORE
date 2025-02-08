package util_test

import (
	"testing"

	"github.com/coranlabs/CORAN_SMF/Application_entity/util"
)

func TestBitRateToKbpsWithValidBpsBitRateShouldReturnValidKbpsBitRate(t *testing.T) {
	var bitrate string = "1000 bps"
	var correctBitRateKbps uint64 = 1

	bitrateKbps, err := util.BitRateTokbps(bitrate)

	t.Log("Check: err should be nil since act should work correctly.")
	if err != nil {
		t.Errorf("Error: err should be nil but it returns %s", err)
	}
	t.Log("Check: convert should act correctly.")
	if bitrateKbps != correctBitRateKbps {
		t.Errorf("Error: bitrate convert failed. Expect: %d. Actually: %d", correctBitRateKbps, bitrateKbps)
	}
	t.Log("Passed.")
}

func TestBitRateToKbpsWithValidKbpsBitRateShouldReturnValidKbpsBitRate(t *testing.T) {
	var bitrate string = "1000 Kbps"
	var correctBitRateKbps uint64 = 1000

	bitrateKbps, err := util.BitRateTokbps(bitrate)

	t.Log("Check: err should be nil since act should work correctly.")
	if err != nil {
		t.Errorf("Error: err should be nil but it returns %s", err)
	}
	t.Log("Check: convert should act correctly.")
	if bitrateKbps != correctBitRateKbps {
		t.Errorf("Error: bitrate convert failed. Expect: %d. Actually: %d", correctBitRateKbps, bitrateKbps)
	}
	t.Log("Passed.")
}

func TestBitRateToKbpsWithValidMbpsBitRateShouldReturnValidKbpsBitRate(t *testing.T) {
	var bitrate string = "1000 Mbps"
	var correctBitRateKbps uint64 = 1000000

	bitrateKbps, err := util.BitRateTokbps(bitrate)

	t.Log("Check: err should be nil since act should work correctly.")
	if err != nil {
		t.Errorf("Error: err should be nil but it returns %s", err)
	}
	t.Log("Check: convert should act correctly.")
	if bitrateKbps != correctBitRateKbps {
		t.Errorf("Error: bitrate convert failed. Expect: %d. Actually: %d", correctBitRateKbps, bitrateKbps)
	}
	t.Log("Passed.")
}

func TestBitRateToKbpsWithValidGbpsBitRateShouldReturnValidKbpsBitRate(t *testing.T) {
	var bitrate string = "1000 Gbps"
	var correctBitRateKbps uint64 = 1000000000

	bitrateKbps, err := util.BitRateTokbps(bitrate)

	t.Log("Check: err should be nil since act should work correctly.")
	if err != nil {
		t.Errorf("Error: err should be nil but it returns %s", err)
	}
	t.Log("Check: convert should act correctly.")
	if bitrateKbps != correctBitRateKbps {
		t.Errorf("Error: bitrate convert failed. Expect: %d. Actually: %d", correctBitRateKbps, bitrateKbps)
	}
	t.Log("Passed.")
}

func TestBitRateToKbpsWithValidTbpsBitRateShouldReturnValidKbpsBitRate(t *testing.T) {
	var bitrate string = "1000 Tbps"
	var correctBitRateKbps uint64 = 1000000000000

	bitrateKbps, err := util.BitRateTokbps(bitrate)

	t.Log("Check: err should be nil since act should work correctly.")
	if err != nil {
		t.Errorf("Error: err should be nil but it returns %s", err)
	}
	t.Log("Check: convert should act correctly.")
	if bitrateKbps != correctBitRateKbps {
		t.Errorf("Error: bitrate convert failed. Expect: %d. Actually: %d", correctBitRateKbps, bitrateKbps)
	}
	t.Log("Passed.")
}

func TestBitRateToKbpsWithInvalidBitRateShouldReturnError(t *testing.T) {
	var bitrate string = "1000" // The unit is absent. It should raise error for `BitRateToKbps`.

	_, err := util.BitRateTokbps(bitrate)

	t.Log("Check: err should not be nil.")
	if err == nil {
		t.Error("Error: err should not be nil.")
	}
	t.Log("Passed.")
}
