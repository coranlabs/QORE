package nasType_test

import (
	"testing"

	nas "github.com/coranlabs/CORAN_LIB_NAS"
	"github.com/stretchr/testify/assert"

	"github.com/coranlabs/CORAN_LIB_NAS/nasType"
)

func TestNasTypeNewPDUSESSIONESTABLISHMENTREJECTMessageIdentity(t *testing.T) {
	a := nasType.NewPDUSESSIONESTABLISHMENTREJECTMessageIdentity()
	assert.NotNil(t, a)
}

type nasTypePDUSESSIONESTABLISHMENTREJECTMessageIdentityMessageType struct {
	in  uint8
	out uint8
}

var nasTypePDUSESSIONESTABLISHMENTREJECTMessageIdentityMessageTypeTable = []nasTypePDUSESSIONESTABLISHMENTREJECTMessageIdentityMessageType{
	{nas.MsgTypePDUSessionEstablishmentReject, nas.MsgTypePDUSessionEstablishmentReject},
}

func TestNasTypeGetSetPDUSESSIONESTABLISHMENTREJECTMessageIdentityMessageType(t *testing.T) {
	a := nasType.NewPDUSESSIONESTABLISHMENTREJECTMessageIdentity()
	for _, table := range nasTypePDUSESSIONESTABLISHMENTREJECTMessageIdentityMessageTypeTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
