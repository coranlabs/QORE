package nasMessage_test

import (
	"bytes"
	"reflect"
	"testing"

	nas "github.com/coranlabs/CORAN_LIB_NAS"
	"github.com/stretchr/testify/assert"

	"github.com/coranlabs/CORAN_LIB_NAS/logger"
	"github.com/coranlabs/CORAN_LIB_NAS/nasMessage"
	"github.com/coranlabs/CORAN_LIB_NAS/nasType"
)

type nasMessagePDUSessionAuthenticationResultData struct {
	inExtendedProtocolDiscriminator                 uint8
	inPDUSessionID                                  uint8
	inPTI                                           uint8
	inPDUSESSIONAUTHENTICATIONRESULTMessageIdentity uint8
	inEAPMessage                                    nasType.EAPMessage
	inExtendedProtocolConfigurationOptions          nasType.ExtendedProtocolConfigurationOptions
}

var nasMessagePDUSessionAuthenticationResultTable = []nasMessagePDUSessionAuthenticationResultData{
	{
		inExtendedProtocolDiscriminator: nas.MsgTypePDUSessionAuthenticationResult,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSESSIONAUTHENTICATIONRESULTMessageIdentity: 0x01,
		inEAPMessage: nasType.EAPMessage{
			Iei:    nasMessage.PDUSessionAuthenticationResultEAPMessageType,
			Len:    4,
			Buffer: []uint8{0x01, 0x01, 0x01, 0x01},
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionAuthenticationResultExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionAuthenticationResult(t *testing.T) {
	a := nasMessage.NewPDUSessionAuthenticationResult(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewPDUSessionAuthenticationResultMessage(t *testing.T) {
	for i, table := range nasMessagePDUSessionAuthenticationResultTable {
		t.Logf("Test Cnt:%d", i)
		a := nasMessage.NewPDUSessionAuthenticationResult(0)
		b := nasMessage.NewPDUSessionAuthenticationResult(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.PDUSessionID.SetPDUSessionID(table.inPDUSessionID)
		a.PTI.SetPTI(table.inPTI)
		a.PDUSESSIONAUTHENTICATIONRESULTMessageIdentity.SetMessageType(table.inPDUSESSIONAUTHENTICATIONRESULTMessageIdentity)

		a.EAPMessage = nasType.NewEAPMessage(nasMessage.PDUSessionAuthenticationResultEAPMessageType)
		a.EAPMessage = &table.inEAPMessage

		a.ExtendedProtocolConfigurationOptions = nasType.NewExtendedProtocolConfigurationOptions(nasMessage.PDUSessionAuthenticationResultExtendedProtocolConfigurationOptionsType)
		a.ExtendedProtocolConfigurationOptions = &table.inExtendedProtocolConfigurationOptions

		buff := new(bytes.Buffer)
		a.EncodePDUSessionAuthenticationResult(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		logger.NasMsgLog.Debugln(data)
		b.DecodePDUSessionAuthenticationResult(&data)
		logger.NasMsgLog.Debugln("Decode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}

	}
}
