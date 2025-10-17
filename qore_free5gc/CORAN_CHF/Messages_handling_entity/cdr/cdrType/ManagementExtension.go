package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type ManagementExtension struct { /* Sequence Type */
	Identifier   asn.ObjectIdentifier
	Significance *bool                          `ber:"tagNum:1,optional,default:FALSE"`
	Information  ManagementExtensionInformation `ber:"tagNum:2"`
}

const (
	ManagementExtensionPresentNothing int = iota /* No components present */
)

type ManagementExtensionInformation struct {
	Present int /* Open Type */
}
