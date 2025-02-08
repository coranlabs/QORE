package nasType_test

import (
	"testing"

	nas "github.com/coranlabs/CORAN_LIB_NAS"
	"github.com/stretchr/testify/assert"

	"github.com/coranlabs/CORAN_LIB_NAS/nasType"
)

func TestNasTypeNewPDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity(t *testing.T) {
	a := nasType.NewPDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity()
	assert.NotNil(t, a)
}

type nasTypePDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentityMessageType struct {
	in  uint8
	out uint8
}

var nasTypePDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentityMessageTypeTable = []nasTypePDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentityMessageType{
	{nas.MsgTypePDUSessionModificationCommandReject, nas.MsgTypePDUSessionModificationCommandReject},
}

func TestNasTypeGetSetPDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentityMessageType(t *testing.T) {
	a := nasType.NewPDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity()
	for _, table := range nasTypePDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentityMessageTypeTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
