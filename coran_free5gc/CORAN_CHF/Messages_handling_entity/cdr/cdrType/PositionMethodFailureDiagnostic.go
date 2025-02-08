package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	PositionMethodFailureDiagnosticPresentCongestion                               asn.Enumerated = 0
	PositionMethodFailureDiagnosticPresentInsufficientResources                    asn.Enumerated = 1
	PositionMethodFailureDiagnosticPresentInsufficientMeasurementData              asn.Enumerated = 2
	PositionMethodFailureDiagnosticPresentInconsistentMeasurementData              asn.Enumerated = 3
	PositionMethodFailureDiagnosticPresentLocationProcedureNotCompleted            asn.Enumerated = 4
	PositionMethodFailureDiagnosticPresentLocationProcedureNotSupportedByTargetMS  asn.Enumerated = 5
	PositionMethodFailureDiagnosticPresentQoSNotAttainable                         asn.Enumerated = 6
	PositionMethodFailureDiagnosticPresentPositionMethodNotAvailableInNetwork      asn.Enumerated = 7
	PositionMethodFailureDiagnosticPresentPositionMethodNotAvailableInLocationArea asn.Enumerated = 8
)

type PositionMethodFailureDiagnostic struct {
	Value asn.Enumerated
}
