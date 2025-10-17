package nasType_test

import (
	"testing"

	nas "github.com/coranlabs/CORAN_LIB_NAS"
	"github.com/stretchr/testify/assert"

	"github.com/coranlabs/CORAN_LIB_NAS/nasType"
)

func TestNasTypeNewNotificationMessageIdentity(t *testing.T) {
	a := nasType.NewNotificationMessageIdentity()
	assert.NotNil(t, a)
}

type nasTypeNotificationMessageIdentityMessageType struct {
	in  uint8
	out uint8
}

var nasTypeNotificationMessageIdentityMessageTypeTable = []nasTypeNotificationMessageIdentityMessageType{
	{nas.MsgTypeNotification, nas.MsgTypeNotification},
}

func TestNasTypeGetSetNotificationMessageIdentityMessageType(t *testing.T) {
	a := nasType.NewNotificationMessageIdentity()
	for _, table := range nasTypeNotificationMessageIdentityMessageTypeTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
