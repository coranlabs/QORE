package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	PresenceReportingAreaStatusPresentInsideArea  asn.Enumerated = 0
	PresenceReportingAreaStatusPresentOutsideArea asn.Enumerated = 1
	PresenceReportingAreaStatusPresentInactive    asn.Enumerated = 2
	PresenceReportingAreaStatusPresentUnknown     asn.Enumerated = 3
)

type PresenceReportingAreaStatus struct {
	Value asn.Enumerated
}
