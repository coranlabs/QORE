package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	MAPDUSessionIndicatorPresentMAPDURequest               asn.Enumerated = 0
	MAPDUSessionIndicatorPresentMAPDUNetworkUpgradeAllowed asn.Enumerated = 1
)

type MAPDUSessionIndicator struct {
	Value asn.Enumerated
}
