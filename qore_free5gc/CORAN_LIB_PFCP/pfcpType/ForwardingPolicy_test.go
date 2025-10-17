package pfcpType

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalForwardingPolicy(t *testing.T) {
	testData := ForwardingPolicy{
		ForwardingPolicyIdentifierLength: 13,
		ForwardingPolicyIdentifier:       []byte("coranlabs.local"),
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{13, 102, 114, 101, 101, 53, 103, 99, 46, 108, 111, 99, 97, 108}, buf)
}

func TestUnmarshalForwardingPolicy(t *testing.T) {
	buf := []byte{13, 102, 114, 101, 101, 53, 103, 99, 46, 108, 111, 99, 97, 108}
	var testData ForwardingPolicy
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := ForwardingPolicy{
		ForwardingPolicyIdentifierLength: 13,
		ForwardingPolicyIdentifier:       []byte("coranlabs.local"),
	}
	assert.Equal(t, expectData, testData)
}
