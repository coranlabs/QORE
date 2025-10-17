package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	UnauthorizedLCSClientDiagnosticPresentNoAdditionalInformation                        asn.Enumerated = 0
	UnauthorizedLCSClientDiagnosticPresentClientNotInMSPrivacyExceptionList              asn.Enumerated = 1
	UnauthorizedLCSClientDiagnosticPresentCallToClientNotSetup                           asn.Enumerated = 2
	UnauthorizedLCSClientDiagnosticPresentPrivacyOverrideNotApplicable                   asn.Enumerated = 3
	UnauthorizedLCSClientDiagnosticPresentDisallowedByLocalRegulatoryRequirements        asn.Enumerated = 4
	UnauthorizedLCSClientDiagnosticPresentUnauthorizedPrivacyClass                       asn.Enumerated = 5
	UnauthorizedLCSClientDiagnosticPresentUnauthorizedCallSessionUnrelatedExternalClient asn.Enumerated = 6
	UnauthorizedLCSClientDiagnosticPresentUnauthorizedCallSessionRelatedExternalClient   asn.Enumerated = 7
)

type UnauthorizedLCSClientDiagnostic struct {
	Value asn.Enumerated
}
