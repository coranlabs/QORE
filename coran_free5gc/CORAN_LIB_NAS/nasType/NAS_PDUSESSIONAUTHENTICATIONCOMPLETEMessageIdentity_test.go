package nasType_test

import (
	"testing"

	nas "github.com/coranlabs/CORAN_LIB_NAS"
	"github.com/stretchr/testify/assert"

	"github.com/coranlabs/CORAN_LIB_NAS/nasType"
)

func TestNasTypeNewPDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentity(t *testing.T) {
	a := nasType.NewPDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentity()
	assert.NotNil(t, a)
}

type nasTypePDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentityMessageType struct {
	in  uint8
	out uint8
}

var nasTypePDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentityMessageTypeTable = []nasTypePDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentityMessageType{
	{nas.MsgTypePDUSessionAuthenticationComplete, nas.MsgTypePDUSessionAuthenticationComplete},
}

func TestNasTypeGetSetPDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentityMessageType(t *testing.T) {
	a := nasType.NewPDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentity()
	for _, table := range nasTypePDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentityMessageTypeTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
