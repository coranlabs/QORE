package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	CauseProtocolPresentTransferSyntaxError                          aper.Enumerated = 0
	CauseProtocolPresentAbstractSyntaxErrorReject                    aper.Enumerated = 1
	CauseProtocolPresentAbstractSyntaxErrorIgnoreAndNotify           aper.Enumerated = 2
	CauseProtocolPresentMessageNotCompatibleWithReceiverState        aper.Enumerated = 3
	CauseProtocolPresentSemanticError                                aper.Enumerated = 4
	CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage aper.Enumerated = 5
	CauseProtocolPresentUnspecified                                  aper.Enumerated = 6
)

type CauseProtocol struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:6"`
}
