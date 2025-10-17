package ngap_handler

import (
	"fmt"

	//"github.com/coranlabs/CORAN_AMF/context"

	//"github.com/coranlabs/CORAN_LIB_NGAP/ngapConvert"
	ngap "github.com/coranlabs/CORAN_LIB_NGAP"
	"github.com/coranlabs/CORAN_LIB_NGAP/ngapType"
)

func BuildNGSetupResponse() ([]byte, error) {
	//amfSelf := context.GetSelf()
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeNGSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject
	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentNGSetupResponse
	successfulOutcome.Value.NGSetupResponse = new(ngapType.NGSetupResponse)

	// nGSetupResponse := successfulOutcome.Value.NGSetupResponse
	// nGSetupResponseIEs := &nGSetupResponse.ProtocolIEs

	// // AMFName
	// ie := ngapType.NGSetupResponseIEs{}
	// ie.Id.Value = ngapType.ProtocolIEIDAMFName
	// ie.Criticality.Value = ngapType.CriticalityPresentReject
	// ie.Value.Present = ngapType.NGSetupResponseIEsPresentAMFName
	// ie.Value.AMFName = new(ngapType.AMFName)

	// aMFName := ie.Value.AMFName
	// aMFName.Value = amfSelf.Name

	// nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	// //ServedGUAMIList
	// ie = ngapType.NGSetupResponseIEs{}
	// ie.Id.Value = ngapType.ProtocolIEIDServedGUAMIList
	// ie.Criticality.Value = ngapType.CriticalityPresentReject
	// ie.Value.Present = ngapType.NGSetupResponseIEsPresentServedGUAMIList
	// ie.Value.ServedGUAMIList = new(ngapType.ServedGUAMIList)

	// servedGUAMIList := ie.Value.ServedGUAMIList
	// for _, guami := range amfSelf.ServedGuamiList {
	// 	servedGUAMIItem := ngapType.ServedGUAMIItem{}
	// 	servedGUAMIItem.GUAMI.PLMNIdentity = ngapConvert.PlmnIdToNgap(*guami.PlmnId)
	// 	regionId, setId, prtId := ngapConvert.AmfIdToNgap(guami.AmfId)
	// 	servedGUAMIItem.GUAMI.AMFRegionID.Value = regionId
	// 	servedGUAMIItem.GUAMI.AMFSetID.Value = setId
	// 	servedGUAMIItem.GUAMI.AMFPointer.Value = prtId
	// 	servedGUAMIList.List = append(servedGUAMIList.List, servedGUAMIItem)
	// }

	// nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	// //relativeAMFCapacity
	// ie = ngapType.NGSetupResponseIEs{}
	// ie.Id.Value = ngapType.ProtocolIEIDRelativeAMFCapacity
	// ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	// ie.Value.Present = ngapType.NGSetupResponseIEsPresentRelativeAMFCapacity
	// ie.Value.RelativeAMFCapacity = new(ngapType.RelativeAMFCapacity)
	// relativeAMFCapacity := ie.Value.RelativeAMFCapacity
	// relativeAMFCapacity.Value = amfSelf.RelativeCapacity

	// nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	// //ServedGUAMIList
	// ie = ngapType.NGSetupResponseIEs{}
	// ie.Id.Value = ngapType.ProtocolIEIDPLMNSupportList
	// ie.Criticality.Value = ngapType.CriticalityPresentReject
	// ie.Value.Present = ngapType.NGSetupResponseIEsPresentPLMNSupportList
	// ie.Value.PLMNSupportList = new(ngapType.PLMNSupportList)

	// pLMNSupportList := ie.Value.PLMNSupportList
	// for _, plmnItem := range amfSelf.PlmnSupportList {
	// 	pLMNSupportItem := ngapType.PLMNSupportItem{}
	// 	pLMNSupportItem.PLMNIdentity = ngapConvert.PlmnIdToNgap(*plmnItem.PlmnId)
	// 	for _, snssai := range plmnItem.SNssaiList {
	// 		sliceSupportItem := ngapType.SliceSupportItem{}
	// 		sliceSupportItem.SNSSAI = ngapConvert.SNssaiToNgap(snssai)
	// 		pLMNSupportItem.SliceSupportList.List = append(pLMNSupportItem.SliceSupportList.List, sliceSupportItem)
	// 	}
	// 	pLMNSupportList.List = append(pLMNSupportList.List, pLMNSupportItem)
	// }

	//nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)
	fmt.Printf("\nPDU Content: %+v\n", pdu)
	fmt.Printf("SuccessfulOutcome: %+v\n", pdu.SuccessfulOutcome)
	fmt.Printf("ProcedureCode: %v\n", pdu.SuccessfulOutcome.ProcedureCode.Value)
	fmt.Printf("Criticality: %v\n", pdu.SuccessfulOutcome.Criticality.Value)
	fmt.Printf("Value.Present: %v\n", pdu.SuccessfulOutcome.Value.Present)

	return ngap.Encoder(pdu)
}
