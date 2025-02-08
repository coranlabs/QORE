package nasType_test

import (
	"testing"

	nas "github.com/coranlabs/CORAN_LIB_NAS"
	"github.com/stretchr/testify/assert"

	"github.com/coranlabs/CORAN_LIB_NAS/nasType"
)

type nasTypeDeregistrationAcceptMessageIdentityData struct {
	in  uint8
	out uint8
}

var nasTypeDeregistrationAcceptMessageIdentityTable = []nasTypeDeregistrationAcceptMessageIdentityData{
	{nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration, nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration},
}

func TestNasTypeNewDeregistrationAcceptMessageIdentity(t *testing.T) {
	a := nasType.NewDeregistrationAcceptMessageIdentity()
	assert.NotNil(t, a)
}

func TestNasTypeGetSetDeregistrationAcceptMessageIdentity(t *testing.T) {
	a := nasType.NewDeregistrationAcceptMessageIdentity()
	for _, table := range nasTypeDeregistrationAcceptMessageIdentityTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
