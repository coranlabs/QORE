package nasType_test

import (
	"testing"

	nas "github.com/coranlabs/CORAN_LIB_NAS"
	"github.com/stretchr/testify/assert"

	"github.com/coranlabs/CORAN_LIB_NAS/nasType"
)

func TestNasTypeNewPDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity(t *testing.T) {
	a := nasType.NewPDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity()
	assert.NotNil(t, a)
}

type nasTypePDUSESSIONAUTHENTICATIONCOMMANDMessageIdentityMessageType struct {
	in  uint8
	out uint8
}

var nasTypePDUSESSIONAUTHENTICATIONCOMMANDMessageIdentityMessageTypeTable = []nasTypePDUSESSIONAUTHENTICATIONCOMMANDMessageIdentityMessageType{
	{nas.MsgTypePDUSessionAuthenticationCommand, nas.MsgTypePDUSessionAuthenticationCommand},
}

func TestNasTypeGetSetPDUSESSIONAUTHENTICATIONCOMMANDMessageIdentityMessageType(t *testing.T) {
	a := nasType.NewPDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity()
	for _, table := range nasTypePDUSESSIONAUTHENTICATIONCOMMANDMessageIdentityMessageTypeTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
