package nasType_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coranlabs/CORAN_LIB_NAS/nasType"
)

func TestNasTypeNewExtendedProtocolDiscriminatort(t *testing.T) {
	a := nasType.NewExtendedProtocolDiscriminator()
	assert.NotNil(t, a)
}

type nasTypeExtendedProtocolDiscriminatorData struct {
	in  uint8
	out uint8
}

var nasTypeExtendedProtocolDiscriminatorTable = []nasTypeExtendedProtocolDiscriminatorData{
	{2, 2},
}

func TestNasTypeGetSetExtendedProtocolDiscriminator(t *testing.T) {
	a := nasType.NewExtendedProtocolDiscriminator()
	for _, table := range nasTypeExtendedProtocolDiscriminatorTable {
		a.SetExtendedProtocolDiscriminator(table.in)
		assert.Equal(t, table.out, a.GetExtendedProtocolDiscriminator())
	}
}
