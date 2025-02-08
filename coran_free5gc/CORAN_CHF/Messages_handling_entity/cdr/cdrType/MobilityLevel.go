package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	MobilityLevelPresentStationary         asn.Enumerated = 0
	MobilityLevelPresentNomadic            asn.Enumerated = 1
	MobilityLevelPresentRestrictedMobility asn.Enumerated = 2
	MobilityLevelPresentFullyMobility      asn.Enumerated = 3
)

type MobilityLevel struct {
	Value asn.Enumerated
}
