package ngap_handler

// import (
// 	"encoding/hex"
// 	"fmt"
// 	"strconv"

// 	"github.com/coranlabs/CORAN_AMF/Messages_controller/context"
// 	"github.com/coranlabs/CORAN_AMF/Messages_handling_entity/gmm/common"
// 	gmm_message "github.com/coranlabs/CORAN_AMF/Messages_handling_entity/gmm/message"
// 	"github.com/coranlabs/CORAN_AMF/Messages_handling_entity/nas/nas_security"

// 	//"github.com/coranlabs/CORAN_AMF/Application_entity/config/factory"
// 	"github.com/coranlabs/CORAN_AMF/Application_entity/config/factory"
// 	"github.com/coranlabs/CORAN_LIB_APER"
// 	"github.com/coranlabs/CORAN_LIB_NAS"
// 	"github.com/coranlabs/CORAN_LIB_NAS/nasMessage"
// 	"github.com/coranlabs/CORAN_LIB_NGAP/ngapConvert"
// 	"github.com/coranlabs/CORAN_LIB_NGAP/ngapType"
// 	//"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
// )

// const (
// 	MaxNumOfTAI                       int   = 16
// 	MaxNumOfBroadcastPLMNs            int   = 12
// 	MaxNumOfPLMNs                     int   = 12
// 	MaxNumOfSlice                     int   = 1024
// 	MaxNumOfAllowedSnssais            int   = 8
// 	MaxValueOfAmfUeNgapId             int64 = 1099511627775
// 	MaxNumOfServedGuamiList           int   = 256
// 	MaxNumOfPDUSessions               int   = 256
// 	MaxNumOfDRBs                      int   = 32
// 	MaxNumOfAOI                       int   = 64
// 	MaxT3513RetryTimes                int   = 4
// 	MaxT3522RetryTimes                int   = 4
// 	MaxT3550RetryTimes                int   = 4
// 	MaxT3560RetryTimes                int   = 4
// 	MaxT3565RetryTimes                int   = 4
// 	MAxNumOfAlgorithm                 int   = 8
// 	DefaultT3502                      int   = 720  // 12 min
// 	DefaultT3512                      int   = 3240 // 54 min
// 	DefaultNon3gppDeregistrationTimer int   = 3240 // 54 min
// )

// func buildCriticalityDiagnostics(
// 	procedureCode *int64,
// 	triggeringMessage *aper.Enumerated,
// 	procedureCriticality *aper.Enumerated,
// 	iesCriticalityDiagnostics *ngapType.CriticalityDiagnosticsIEList) (
// 	criticalityDiagnostics ngapType.CriticalityDiagnostics,
// ) {
// 	if procedureCode != nil {
// 		criticalityDiagnostics.ProcedureCode = new(ngapType.ProcedureCode)
// 		criticalityDiagnostics.ProcedureCode.Value = *procedureCode
// 	}

// 	if triggeringMessage != nil {
// 		criticalityDiagnostics.TriggeringMessage = new(ngapType.TriggeringMessage)
// 		criticalityDiagnostics.TriggeringMessage.Value = *triggeringMessage
// 	}

// 	if procedureCriticality != nil {
// 		criticalityDiagnostics.ProcedureCriticality = new(ngapType.Criticality)
// 		criticalityDiagnostics.ProcedureCriticality.Value = *procedureCriticality
// 	}

// 	if iesCriticalityDiagnostics != nil {
// 		criticalityDiagnostics.IEsCriticalityDiagnostics = iesCriticalityDiagnostics
// 	}

// 	return criticalityDiagnostics
// }

// func buildCriticalityDiagnosticsIEItem(ieCriticality aper.Enumerated, ieID int64, typeOfErr aper.Enumerated) (
// 	item ngapType.CriticalityDiagnosticsIEItem,
// ) {
// 	item = ngapType.CriticalityDiagnosticsIEItem{
// 		IECriticality: ngapType.Criticality{
// 			Value: ieCriticality,
// 		},
// 		IEID: ngapType.ProtocolIEID{
// 			Value: ieID,
// 		},
// 		TypeOfError: ngapType.TypeOfError{
// 			Value: typeOfErr,
// 		},
// 	}

// 	return item
// }

// func Ngsetuprequest(ran *context.AmfRan, initmsg *ngapType.InitiatingMessage) {

// 	var globalRANNodeID *ngapType.GlobalRANNodeID
// 	var rANNodeName *ngapType.RANNodeName
// 	var supportedTAList *ngapType.SupportedTAList
// 	var defaultPagingDRX *ngapType.PagingDRX
// 	var uERetentionInformation *ngapType.UERetentionInformation

// 	var syntaxCause *ngapType.Cause
// 	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList
// 	abort := false
// 	nGSetupRequest := initmsg.Value.NGSetupRequest
// 	fmt.Printf("message %v", nGSetupRequest)

// 	for _, ie := range nGSetupRequest.ProtocolIEs.List {
// 		switch ie.Id.Value {
// 		case ngapType.ProtocolIEIDGlobalRANNodeID: // mandatory, reject
// 			if globalRANNodeID != nil {
// 				fmt.Printf("Duplicate IE GlobalRANNodeID")
// 				syntaxCause = &ngapType.Cause{
// 					Present: ngapType.CausePresentProtocol,
// 					Protocol: &ngapType.CauseProtocol{
// 						Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 					},
// 				}
// 				abort = true
// 				break
// 			}
// 			globalRANNodeID = ie.Value.GlobalRANNodeID
// 			fmt.Printf("Decode IE GlobalRANNodeID")
// 		case ngapType.ProtocolIEIDRANNodeName: // optional, ignore
// 			if rANNodeName != nil {
// 				fmt.Printf("Duplicate IE RANNodeName")
// 				syntaxCause = &ngapType.Cause{
// 					Present: ngapType.CausePresentProtocol,
// 					Protocol: &ngapType.CauseProtocol{
// 						Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 					},
// 				}
// 				abort = true
// 				break
// 			}
// 			rANNodeName = ie.Value.RANNodeName
// 			fmt.Printf("Decode IE RANNodeName")
// 		case ngapType.ProtocolIEIDSupportedTAList: // mandatory, reject
// 			if supportedTAList != nil {
// 				fmt.Printf("Duplicate IE SupportedTAList")
// 				syntaxCause = &ngapType.Cause{
// 					Present: ngapType.CausePresentProtocol,
// 					Protocol: &ngapType.CauseProtocol{
// 						Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 					},
// 				}
// 				abort = true
// 				break
// 			}
// 			supportedTAList = ie.Value.SupportedTAList
// 			fmt.Printf("Decode IE SupportedTAList")
// 		case ngapType.ProtocolIEIDDefaultPagingDRX: // mandatory, ignore
// 			if defaultPagingDRX != nil {
// 				fmt.Printf("Duplicate IEdefaultPagingDRX ")
// 				syntaxCause = &ngapType.Cause{
// 					Present: ngapType.CausePresentProtocol,
// 					Protocol: &ngapType.CauseProtocol{
// 						Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 					},
// 				}
// 				abort = true
// 				break
// 			}
// 			defaultPagingDRX = ie.Value.DefaultPagingDRX
// 			fmt.Printf("Decode IEdefaultPagingDRX ")
// 		case ngapType.ProtocolIEIDUERetentionInformation: // optional, ignore
// 			if uERetentionInformation != nil {
// 				fmt.Printf("Duplicate IE UERetentionInformation")
// 				syntaxCause = &ngapType.Cause{
// 					Present: ngapType.CausePresentProtocol,
// 					Protocol: &ngapType.CauseProtocol{
// 						Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 					},
// 				}
// 				abort = true
// 				break
// 			}
// 			uERetentionInformation = ie.Value.UERetentionInformation
// 			fmt.Printf("Decode IE UERetentionInformation")
// 		default:
// 			switch ie.Criticality.Value {
// 			case ngapType.CriticalityPresentReject:
// 				fmt.Printf("Not comprehended IE ID 0x%04x (criticality: reject)", ie.Id.Value)
// 			case ngapType.CriticalityPresentIgnore:
// 				fmt.Printf("Not comprehended IE ID 0x%04x (criticality: ignore)", ie.Id.Value)
// 			case ngapType.CriticalityPresentNotify:
// 				fmt.Printf("Not comprehended IE ID 0x%04x (criticality: notify)", ie.Id.Value)
// 			}
// 			if ie.Criticality.Value != ngapType.CriticalityPresentIgnore {
// 				item := buildCriticalityDiagnosticsIEItem(ie.Criticality.Value, ie.Id.Value, ngapType.TypeOfErrorPresentNotUnderstood)
// 				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
// 				if ie.Criticality.Value == ngapType.CriticalityPresentReject {
// 					abort = true
// 				}
// 			}
// 		}
// 	}
// 	if abort {
// 		return
// 	}

// 	if globalRANNodeID == nil {
// 		fmt.Printf("Missing IE GlobalRANNodeID")
// 		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDGlobalRANNodeID, ngapType.TypeOfErrorPresentMissing)
// 		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
// 		abort = true
// 	}
// 	if supportedTAList == nil {
// 		fmt.Printf("Missing IE SupportedTAList")
// 		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDSupportedTAList, ngapType.TypeOfErrorPresentMissing)
// 		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
// 		abort = true
// 	}

// 	if syntaxCause != nil || len(iesCriticalityDiagnostics.List) > 0 {
// 		fmt.Printf("Has IE error")
// 		procedureCode := ngapType.ProcedureCodeNGSetup
// 		triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
// 		procedureCriticality := ngapType.CriticalityPresentReject
// 		var pIesCriticalityDiagnostics *ngapType.CriticalityDiagnosticsIEList
// 		if len(iesCriticalityDiagnostics.List) > 0 {
// 			pIesCriticalityDiagnostics = &iesCriticalityDiagnostics
// 		}
// 		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality, pIesCriticalityDiagnostics)
// 		if syntaxCause == nil {
// 			syntaxCause = &ngapType.Cause{
// 				Present: ngapType.CausePresentProtocol,
// 				Protocol: &ngapType.CauseProtocol{
// 					Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 				},
// 			}
// 		}
// 		fmt.Printf("%v , %v ,%v , %v ", ran, *syntaxCause, nil, &criticalityDiagnostics)
// 	}

// 	// start of packet plmn mapping and responding
// 	var cause ngapType.Cause

// 	ran.SetRanId(globalRANNodeID)
// 	if rANNodeName != nil {
// 		ran.Name = rANNodeName.Value
// 	}
// 	if defaultPagingDRX != nil {
// 		fmt.Printf("PagingDRX[%d]", defaultPagingDRX.Value)
// 	}

// 	for i := 0; i < len(supportedTAList.List); i++ {
// 		supportedTAItem := supportedTAList.List[i]
// 		tac := hex.EncodeToString(supportedTAItem.TAC.Value)
// 		capOfSupportTai := cap(ran.SupportedTAList)
// 		for j := 0; j < len(supportedTAItem.BroadcastPLMNList.List); j++ {
// 			supportedTAI := context.NewSupportedTAI()
// 			supportedTAI.Tai.Tac = tac
// 			broadcastPLMNItem := supportedTAItem.BroadcastPLMNList.List[j]
// 			plmnId := ngapConvert.PlmnIdToModels(broadcastPLMNItem.PLMNIdentity)
// 			supportedTAI.Tai.PlmnId = &plmnId
// 			capOfSNssaiList := cap(supportedTAI.SNssaiList)
// 			for k := 0; k < len(broadcastPLMNItem.TAISliceSupportList.List); k++ {
// 				tAISliceSupportItem := broadcastPLMNItem.TAISliceSupportList.List[k]
// 				if len(supportedTAI.SNssaiList) < capOfSNssaiList {
// 					supportedTAI.SNssaiList = append(supportedTAI.SNssaiList, ngapConvert.SNssaiToModels(tAISliceSupportItem.SNSSAI))
// 				} else {
// 					break
// 				}
// 			}
// 			fmt.Printf("PLMN_ID[MCC:%s MNC:%s] TAC[%s]", plmnId.Mcc, plmnId.Mnc, tac)
// 			if len(ran.SupportedTAList) < capOfSupportTai {
// 				ran.SupportedTAList = append(ran.SupportedTAList, supportedTAI)
// 			} else {
// 				break
// 			}
// 		}
// 	}

// 	if len(ran.SupportedTAList) == 0 {
// 		fmt.Printf("NG-Setup failure: No supported TA exist in NG-Setup request")
// 		cause.Present = ngapType.CausePresentMisc
// 		cause.Misc = &ngapType.CauseMisc{
// 			Value: ngapType.CauseMiscPresentUnspecified,
// 		}
// 	} else {
// 		var found bool
// 		for i, tai := range ran.SupportedTAList {
// 			// if context.InTaiList(tai.Tai, context.GetSelf().SupportTaiLists) {
// 			// 	fmt.Printf("SERVED_TAI_INDEX[%d]", i)
// 			// 	found = true
// 			// 	break
// 			// }
// 			fmt.Printf("\n SERVED_TAI_INDEX[%d] tai:%v", i, tai)
// 			found = true
// 			break
// 		}
// 		if !found {
// 			fmt.Printf("NG-Setup failure: Cannot find Served TAI in AMF")
// 			cause.Present = ngapType.CausePresentMisc
// 			cause.Misc = &ngapType.CauseMisc{
// 				Value: ngapType.CauseMiscPresentUnknownPLMN,
// 			}
// 		}
// 	}

// 	if cause.Present == ngapType.CausePresentNothing {
// 		SendNGSetupResponse(ran)
// 	} else {
// 		return
// 	}

// }

// func SendNGSetupResponse(ran *context.AmfRan) {
// 	fmt.Printf("Send NG-Setup response")

// 	pkt, err := BuildNGSetupResponse()
// 	if err != nil {
// 		fmt.Printf("Build NGSetupResponse failed : %s", err.Error())
// 		return
// 	}
// 	SendToRan(ran, pkt)
// }

// func InitialUEMessage(ran *context.AmfRan, msg *ngapType.NGAPPDU, initmsg *ngapType.InitiatingMessage) {
// 	var rANUENGAPID *ngapType.RANUENGAPID
// 	var nASPDU *ngapType.NASPDU
// 	var userLocationInformation *ngapType.UserLocationInformation
// 	var rRCEstablishmentCause *ngapType.RRCEstablishmentCause
// 	var fiveGSTMSI *ngapType.FiveGSTMSI
// 	var aMFSetID *ngapType.AMFSetID
// 	var uEContextRequest *ngapType.UEContextRequest
// 	var allowedNSSAI *ngapType.AllowedNSSAI

// 	var syntaxCause *ngapType.Cause
// 	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList
// 	abort := false

// 	initialUEMessage := initmsg.Value.InitialUEMessage
// 	if initialUEMessage == nil {
// 		ran.Log.Error("InitialUEMessage is nil")
// 		return
// 	}

// 	ran.Log.Info("Handle InitialUEMessage")

// 	for _, ie := range initialUEMessage.ProtocolIEs.List {
// 		switch ie.Id.Value {
// 		case ngapType.ProtocolIEIDRANUENGAPID: // mandatory, reject
// 			if rANUENGAPID != nil {
// 				ran.Log.Error("Duplicate IE RAN-UE-NGAP-ID")
// 				syntaxCause = &ngapType.Cause{
// 					Present: ngapType.CausePresentProtocol,
// 					Protocol: &ngapType.CauseProtocol{
// 						Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 					},
// 				}
// 				abort = true
// 				break
// 			}
// 			rANUENGAPID = ie.Value.RANUENGAPID
// 			ran.Log.Trace("Decode IE RAN-UE-NGAP-ID")
// 		case ngapType.ProtocolIEIDNASPDU: // mandatory, reject
// 			if nASPDU != nil {
// 				ran.Log.Error("Duplicate IE NAS-PDU")
// 				syntaxCause = &ngapType.Cause{
// 					Present: ngapType.CausePresentProtocol,
// 					Protocol: &ngapType.CauseProtocol{
// 						Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 					},
// 				}
// 				abort = true
// 				break
// 			}
// 			nASPDU = ie.Value.NASPDU
// 			ran.Log.Trace("Decode IE NAS-PDU")
// 		case ngapType.ProtocolIEIDUserLocationInformation: // mandatory, reject
// 			if userLocationInformation != nil {
// 				ran.Log.Error("Duplicate IE UserLocationInformation")
// 				syntaxCause = &ngapType.Cause{
// 					Present: ngapType.CausePresentProtocol,
// 					Protocol: &ngapType.CauseProtocol{
// 						Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 					},
// 				}
// 				abort = true
// 				break
// 			}
// 			userLocationInformation = ie.Value.UserLocationInformation
// 			ran.Log.Trace("Decode IE UserLocationInformation")
// 		case ngapType.ProtocolIEIDRRCEstablishmentCause: // mandatory, ignore
// 			if rRCEstablishmentCause != nil {
// 				ran.Log.Error("Duplicate IE RRCEstablishmentCause")
// 				syntaxCause = &ngapType.Cause{
// 					Present: ngapType.CausePresentProtocol,
// 					Protocol: &ngapType.CauseProtocol{
// 						Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 					},
// 				}
// 				abort = true
// 				break
// 			}
// 			rRCEstablishmentCause = ie.Value.RRCEstablishmentCause
// 			ran.Log.Trace("Decode IE RRCEstablishmentCause")
// 		case ngapType.ProtocolIEIDFiveGSTMSI: // optional, reject
// 			if fiveGSTMSI != nil {
// 				ran.Log.Error("Duplicate IE FiveG-S-TMSI")
// 				syntaxCause = &ngapType.Cause{
// 					Present: ngapType.CausePresentProtocol,
// 					Protocol: &ngapType.CauseProtocol{
// 						Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 					},
// 				}
// 				abort = true
// 				break
// 			}
// 			fiveGSTMSI = ie.Value.FiveGSTMSI
// 			ran.Log.Trace("Decode IE FiveG-S-TMSI")
// 		case ngapType.ProtocolIEIDAMFSetID: // optional, ignore
// 			if aMFSetID != nil {
// 				ran.Log.Error("Duplicate IE AMFSetID")
// 				syntaxCause = &ngapType.Cause{
// 					Present: ngapType.CausePresentProtocol,
// 					Protocol: &ngapType.CauseProtocol{
// 						Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 					},
// 				}
// 				abort = true
// 				break
// 			}
// 			aMFSetID = ie.Value.AMFSetID
// 			ran.Log.Trace("Decode IE AMFSetID")
// 		case ngapType.ProtocolIEIDUEContextRequest: // optional, ignore
// 			if uEContextRequest != nil {
// 				ran.Log.Error("Duplicate IE UEContextRequest")
// 				syntaxCause = &ngapType.Cause{
// 					Present: ngapType.CausePresentProtocol,
// 					Protocol: &ngapType.CauseProtocol{
// 						Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 					},
// 				}
// 				abort = true
// 				break
// 			}
// 			uEContextRequest = ie.Value.UEContextRequest
// 			ran.Log.Trace("Decode IE UEContextRequest")
// 		case ngapType.ProtocolIEIDAllowedNSSAI: // optional, reject
// 			if allowedNSSAI != nil {
// 				ran.Log.Error("Duplicate IE AllowedNSSAI")
// 				syntaxCause = &ngapType.Cause{
// 					Present: ngapType.CausePresentProtocol,
// 					Protocol: &ngapType.CauseProtocol{
// 						Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
// 					},
// 				}
// 				abort = true
// 				break
// 			}
// 			allowedNSSAI = ie.Value.AllowedNSSAI
// 			ran.Log.Trace("Decode IE AllowedNSSAI")
// 		default:
// 			switch ie.Criticality.Value {
// 			case ngapType.CriticalityPresentReject:
// 				ran.Log.Errorf("Not comprehended IE ID 0x%04x (criticality: reject)", ie.Id.Value)
// 			case ngapType.CriticalityPresentIgnore:
// 				ran.Log.Infof("Not comprehended IE ID 0x%04x (criticality: ignore)", ie.Id.Value)
// 			case ngapType.CriticalityPresentNotify:
// 				ran.Log.Warnf("Not comprehended IE ID 0x%04x (criticality: notify)", ie.Id.Value)
// 			}
// 			if ie.Criticality.Value != ngapType.CriticalityPresentIgnore {
// 				item := buildCriticalityDiagnosticsIEItem(ie.Criticality.Value, ie.Id.Value, ngapType.TypeOfErrorPresentNotUnderstood)
// 				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
// 				if ie.Criticality.Value == ngapType.CriticalityPresentReject {
// 					abort = true
// 				}
// 			}
// 		}
// 	}

// 	if rANUENGAPID == nil {
// 		ran.Log.Error("Missing IE RAN-UE-NGAP-ID")
// 		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDRANUENGAPID, ngapType.TypeOfErrorPresentMissing)
// 		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
// 		abort = true
// 	}
// 	if nASPDU == nil {
// 		ran.Log.Error("Missing IE NAS-PDU")
// 		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDNASPDU, ngapType.TypeOfErrorPresentMissing)
// 		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
// 		abort = true
// 	}
// 	if userLocationInformation == nil {
// 		ran.Log.Error("Missing IE UserLocationInformation")
// 		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDUserLocationInformation, ngapType.TypeOfErrorPresentMissing)
// 		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
// 		abort = true
// 	}

// 	if syntaxCause != nil || len(iesCriticalityDiagnostics.List) > 0 {
// 		ran.Log.Trace("Has IE error")
// 		// procedureCode := ngapType.ProcedureCodeInitialUEMessage
// 		// triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
// 		// procedureCriticality := ngapType.CriticalityPresentIgnore
// 		// var pIesCriticalityDiagnostics *ngapType.CriticalityDiagnosticsIEList
// 		// if len(iesCriticalityDiagnostics.List) > 0 {
// 		// 	pIesCriticalityDiagnostics = &iesCriticalityDiagnostics
// 		// }
// 		// criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality, pIesCriticalityDiagnostics)
// 		// //message.SendErrorIndication(ran, nil, rANUENGAPID, syntaxCause, &criticalityDiagnostics)
// 	}

// 	if abort {
// 		return
// 	}

// 	if rANUENGAPID == nil {
// 		ran.Log.Error("Missing IE RAN-UE-NGAP-ID")
// 		return
// 	}
// 	if nASPDU == nil {
// 		ran.Log.Error("Missing IE NAS-PDU")
// 		return
// 	}
// 	if userLocationInformation == nil {
// 		ran.Log.Error("Missing IE UserLocationInformation")
// 		return
// 	}
// 	if rRCEstablishmentCause == nil {
// 		ran.Log.Warn("Missing IE RRCEstablishmentCause")
// 	}
// 	if aMFSetID != nil {
// 		ran.Log.Warn("IE AMFSetID is not implemented")
// 	}
// 	if allowedNSSAI != nil {
// 		ran.Log.Warn("IE AllowedNSSAI is not implemented")
// 	}

// 	// func handleInitialUEMessageMain(ran *context.AmfRan,
// 	//	message *ngapType.NGAPPDU,
// 	//	rANUENGAPID *ngapType.RANUENGAPID,
// 	//	nASPDU *ngapType.NASPDU,
// 	//	userLocationInformation *ngapType.UserLocationInformation,
// 	//	rRCEstablishmentCause *ngapType.RRCEstablishmentCause,
// 	//	fiveGSTMSI *ngapType.FiveGSTMSI,
// 	//	uEContextRequest *ngapType.UEContextRequest) {
// 	handleInitialUEMessageMain(ran, msg, rANUENGAPID, nASPDU, userLocationInformation, rRCEstablishmentCause /* may be nil */, fiveGSTMSI /* may be nil */, uEContextRequest /* may be nil */)

// }

// func handleInitialUEMessageMain(ran *context.AmfRan,
// 	message *ngapType.NGAPPDU,
// 	rANUENGAPID *ngapType.RANUENGAPID,
// 	nASPDU *ngapType.NASPDU,
// 	userLocationInformation *ngapType.UserLocationInformation,
// 	rRCEstablishmentCause *ngapType.RRCEstablishmentCause,
// 	fiveGSTMSI *ngapType.FiveGSTMSI,
// 	uEContextRequest *ngapType.UEContextRequest,
// ) {
// 	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
// 	if ranUe != nil {
// 		amfUe := ranUe.AmfUe
// 		if amfUe != nil {
// 			// The fact that an amfUe having N2 connection (ranUE) is receiving
// 			// an Initial UE Message indicates there is something wrong,
// 			// so the ranUe with wrong RAN-UE-NGAP-IP should be cleared and detached from the amfUe.
// 			common.StopAll5GSMMTimers(amfUe)
// 			amfUe.DetachRanUe(ran.AnType)
// 			ranUe.DetachAmfUe()
// 		}
// 		err := ranUe.Remove()
// 		if err != nil {
// 			ran.Log.Errorln(err.Error())
// 		}
// 	}

// 	var err error
// 	ranUe, err = ran.NewRanUe(rANUENGAPID.Value)
// 	if err != nil {
// 		ran.Log.Errorf("NewRanUe Error: %+v", err)
// 	}
// 	ran.Log.Debugf("New RanUe [RanUeNgapID: %d]", ranUe.RanUeNgapId)

// 	// Try to get identity from 5G-S-TMSI IE first; if not available, try to get identity from the plain NAS.
// 	var id, idType string
// 	var gmmMessage *nas.GmmMessage
// 	var nasMsgType, regReqType uint8
// 	// Get nasMsgType to send corresponding NAS reject to UE when amfUe is not found.
// 	nasMsg, err := nas_security.DecodePlainNasNoIntegrityCheck(nASPDU.Value)
// 	if err == nil && nasMsg.GmmMessage != nil {
// 		gmmMessage = nasMsg.GmmMessage
// 		nasMsgType = gmmMessage.GmmHeader.GetMessageType()
// 		if gmmMessage.RegistrationRequest != nil {
// 			regReqType = gmmMessage.RegistrationRequest.NgksiAndRegistrationType5GS.GetRegistrationType5GS()
// 		}
// 	}

// 	if fiveGSTMSI != nil {
// 		// <5G-S-TMSI> := <AMF Set ID><AMF Pointer><5G-TMSI>
// 		// GUAMI := <MCC><MNC><AMF Region ID><AMF Set ID><AMF Pointer>
// 		// 5G-GUTI := <GUAMI><5G-TMSI>
// 		amfSetPtrID := hex.EncodeToString([]byte{
// 			fiveGSTMSI.AMFSetID.Value.Bytes[0],
// 			(fiveGSTMSI.AMFSetID.Value.Bytes[1] & 0xc0) | (fiveGSTMSI.AMFPointer.Value.Bytes[0] >> 2),
// 		})
// 		tmsi := hex.EncodeToString(fiveGSTMSI.FiveGTMSI.Value)

// 		id = amfSetPtrID + tmsi
// 		idType = "5G-S-TMSI"
// 		ranUe.Log.Infof("Find 5G-S-TMSI [%q] in InitialUEMessage", id)
// 	} else if regReqType == nasMessage.RegistrationType5GSInitialRegistration {
// 		// NGAP 5G-S-TMSI IE might not be present in InitialUEMessage carrying Initial Registration.
// 		// Need to get 5GSMobileIdentity from Initial Registration.

// 		id, idType, err = amf_nas.GetNas5GSMobileIdentity(gmmMessage)
// 		ran.Log.Infof("5GSMobileIdentity [%q:%q, err: %v]", idType, id, err)
// 	} else {
// 		// Missing NGAP 5G-S-TMSI IE
// 		var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList
// 		ranUe.Log.Warnf("Missing 5G-S-TMSI IE in InitialUEMessage; send ErrorIndication")
// 		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
// 			ngapType.ProtocolIEIDFiveGSTMSI, ngapType.TypeOfErrorPresentMissing)
// 		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
// 		sendErrorMessage(ran, nil, rANUENGAPID, iesCriticalityDiagnostics)

// 		ngap_message.SendUEContextReleaseCommand(ranUe, context.UeContextN2NormalRelease,
// 			ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentUnspecified)
// 		return
// 	}

// 	// If id type is GUTI, since MAC can't be checked here (no amfUe context), the GUTI may not direct to the right amfUe.
// 	// In this case, create a new amfUe to handle the following registration procedure.
// 	var isInvalidGUTI bool = (idType == "5G-GUTI")
// 	amfUe, ok := findAmfUe(ran, id, idType)
// 	if ok && !isInvalidGUTI {
// 		// TODO: invoke Namf_Communication_UEContextTransfer if serving AMF has changed since
// 		// last Registration Request procedure
// 		// Described in TS 23.502 4.2.2.2.2 step 4 (without UDSF deployment)
// 		ranUe.Log.Infof("find AmfUe [%q:%q]", idType, id)
// 		ranUe.Log.Debugf("AmfUe Attach RanUe [RanUeNgapID: %d]", ranUe.RanUeNgapId)
// 		common.AttachRanUeToAmfUeAndReleaseOldIfAny(amfUe, ranUe)
// 	} else if regReqType != nasMessage.RegistrationType5GSInitialRegistration {
// 		if regReqType == nasMessage.RegistrationType5GSPeriodicRegistrationUpdating ||
// 			regReqType == nasMessage.RegistrationType5GSMobilityRegistrationUpdating {
// 			gmm_message.SendRegistrationReject(
// 				ranUe, nasMessage.Cause5GMMImplicitlyDeregistered, "")
// 			ranUe.Log.Warn("Send RegistrationReject [Cause5GMMImplicitlyDeregistered]")
// 		} else if nasMsgType == nas.MsgTypeServiceRequest {
// 			gmm_message.SendServiceReject(
// 				ranUe, nil, nasMessage.Cause5GMMImplicitlyDeregistered)
// 			ranUe.Log.Warn("Send ServiceReject [Cause5GMMImplicitlyDeregistered]")
// 		}

// 		ngap_message.SendUEContextReleaseCommand(ranUe, context.UeContextN2NormalRelease,
// 			ngapType.CausePresentNas, ngapType.CauseNasPresentNormalRelease)
// 		return
// 	}

// 	if userLocationInformation != nil {
// 		ranUe.UpdateLocation(userLocationInformation)
// 	}

// 	if rRCEstablishmentCause != nil {
// 		ranUe.Log.Tracef("[Initial UE Message] RRC Establishment Cause[%d]", rRCEstablishmentCause.Value)
// 		ranUe.RRCEstablishmentCause = strconv.Itoa(int(rRCEstablishmentCause.Value))
// 	}

// 	if uEContextRequest != nil {
// 		ran.Log.Debug("Trigger initial Context Setup procedure")
// 		ranUe.UeContextRequest = true
// 		// TODO: Trigger Initial Context Setup procedure
// 	} else {
// 		ranUe.UeContextRequest = factory.AmfConfig.Configuration.DefaultUECtxReq
// 	}

// 	// TS 23.502 4.2.2.2.3 step 6a Nnrf_NFDiscovery_Request (NF type, AMF Set)
// 	// if aMFSetID != nil {
// 	// TODO: This is a rerouted message
// 	// TS 38.413: AMF shall, if supported, use the IE as described in TS 23.502
// 	// }

// 	// ng-ran propagate allowedNssai in the rerouted initial ue message (TS 38.413 8.6.5)
// 	// TS 23.502 4.2.2.2.3 step 4a Nnssf_NSSelection_Get
// 	// if allowedNSSAI != nil {
// 	// TODO: AMF should use it as defined in TS 23.502
// 	// }

// 	pdu, err := libngap.Encoder(*message)
// 	if err != nil {
// 		ran.Log.Errorf("libngap Encoder Error: %+v", err)
// 	}
// 	ranUe.InitialUEMessage = pdu
// 	amf_nas.HandleNAS(ranUe, ngapType.ProcedureCodeInitialUEMessage, nASPDU.Value, true)
// }
