package nasType_test

import (
	"testing"

	nas "github.com/coranlabs/CORAN_LIB_NAS"
	"github.com/stretchr/testify/assert"

	"github.com/coranlabs/CORAN_LIB_NAS/nasType"
)

type nasTypeIdentityRequestMessageIdentityData struct {
	in  uint8
	out uint8
}

var nasTypeIdentityRequestMessageIdentityTable = []nasTypeIdentityRequestMessageIdentityData{
	{nas.MsgTypeIdentityRequest, nas.MsgTypeIdentityRequest},
}

func TestNasTypeNewIdentityRequestMessageIdentity(t *testing.T) {
	a := nasType.NewIdentityRequestMessageIdentity()
	assert.NotNil(t, a)
}

func TestNasTypeIdentityRequestMessageIdentity(t *testing.T) {
	a := nasType.NewIdentityRequestMessageIdentity()
	for _, table := range nasTypeIdentityRequestMessageIdentityTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
