package pfcpType

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalPFDContents(t *testing.T) {
	testData := PFDContents{
		FlowDescription: "permit in tcp from 60.60.0.1 29500-29600 to any",
		DomainName:      "example.com",
		URL:             "https://example.com/coranlabs",
	}

	expectedHex := "0700002f7065726d697420696e207463702066726f6d2036302e36302e302e312032393530302d323936303020746f20616e79001b68747470733a2f2f6578616d706c652e636f6d2f66726565356763000b6578616d706c652e636f6d"
	buf, err := testData.MarshalBinary()
	testHex := hex.EncodeToString(buf)
	assert.Nil(t, err)
	assert.Equal(t, expectedHex, testHex)
}

func TestUnmarshalPFDContents(t *testing.T) {
	testHex := "0700002f7065726d697420696e207463702066726f6d2036302e36302e302e312032393530302d323936303020746f20616e79001b68747470733a2f2f6578616d706c652e636f6d2f66726565356763000b6578616d706c652e636f6d"
	var testData PFDContents

	expectedObject := PFDContents{
		FlowDescription: "permit in tcp from 60.60.0.1 29500-29600 to any",
		DomainName:      "example.com",
		URL:             "https://example.com/coranlabs",
	}

	buf, _ := hex.DecodeString(testHex)
	err := testData.UnmarshalBinary(buf)
	assert.Nil(t, err)
	assert.Equal(t, expectedObject, testData)
}
