package pfcpType

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalApplicationID(t *testing.T) {
	testData := ApplicationID{
		ApplicationIdentifier: []byte("coranlabs.local"),
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{102, 114, 101, 101, 53, 103, 99, 46, 108, 111, 99, 97, 108}, buf)
}

func TestUnmarshalApplicationID(t *testing.T) {
	buf := []byte{102, 114, 101, 101, 53, 103, 99, 46, 108, 111, 99, 97, 108}
	var testData ApplicationID
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := ApplicationID{
		ApplicationIdentifier: []byte("coranlabs.local"),
	}
	assert.Equal(t, expectData, testData)
}
