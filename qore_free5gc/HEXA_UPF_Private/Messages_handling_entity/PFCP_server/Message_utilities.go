package PFCP_server

import (
	"encoding/binary"
	"fmt"
	"net"
	"slices"

	"github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/logger"
	ebpf_datapath "github.com/coranlabs/CORAN_UPF_eBPF/eBPF_Datapath_entity"
)

const flagPresentIPv4 = 2

func ConvertUint32ToReversedIP(n uint32) string {
	ip := make(net.IP, 4)
	ip[3] = byte(n >> 24)
	ip[2] = byte(n >> 16)
	ip[1] = byte(n >> 8)
	ip[0] = byte(n)
	return ip.String()
}

func convertErrorToIeCause(err error) *ie.IE {
	switch err {
	case errMandatoryIeMissing:
		return ie.NewCause(ie.CauseMandatoryIEMissing)
	case errNoEstablishedAssociation:
		return ie.NewCause(ie.CauseNoEstablishedPFCPAssociation)
	default:
		logger.PFCPSessLog.Infof("Unknown error: %s", err.Error())
		return ie.NewCause(ie.CauseRequestRejected)
	}
}

func validateRequest(nodeId *ie.IE, cpfseid *ie.IE) (fseid *ie.FSEIDFields, err error) {
	if nodeId == nil || cpfseid == nil {
		return nil, errMandatoryIeMissing
	}

	_, err = nodeId.NodeID()
	if err != nil {
		return nil, errMandatoryIeMissing
	}

	fseid, err = cpfseid.FSEID()
	if err != nil {
		return nil, errMandatoryIeMissing
	}

	return fseid, nil
}

func findIEindex(ieArr []*ie.IE, ieType uint16) int {
	arrIndex := slices.IndexFunc(ieArr, func(ie *ie.IE) bool {
		return ie.Type == ieType
	})
	return arrIndex
}

func causeToString(cause uint8) string {
	switch cause {
	case ie.CauseRequestAccepted:
		return "RequestAccepted"
	case ie.CauseRequestRejected:
		return "RequestRejected"
	case ie.CauseSessionContextNotFound:
		return "SessionContextNotFound"
	case ie.CauseMandatoryIEMissing:
		return "MandatoryIEMissing"
	case ie.CauseConditionalIEMissing:
		return "ConditionalIEMissing"
	case ie.CauseInvalidLength:
		return "InvalidLength"
	case ie.CauseMandatoryIEIncorrect:
		return "MandatoryIEIncorrect"
	case ie.CauseInvalidForwardingPolicy:
		return "InvalidForwardingPolicy"
	case ie.CauseInvalidFTEIDAllocationOption:
		return "InvalidFTEIDAllocationOption"
	case ie.CauseNoEstablishedPFCPAssociation:
		return "NoEstablishedPFCPAssociation"
	case ie.CauseRuleCreationModificationFailure:
		return "RuleCreationModificationFailure"
	case ie.CausePFCPEntityInCongestion:
		return "PFCPEntityInCongestion"
	case ie.CauseNoResourcesAvailable:
		return "NoResourcesAvailable"
	case ie.CauseServiceNotSupported:
		return "ServiceNotSupported"
	case ie.CauseSystemFailure:
		return "SystemFailure"
	case ie.CauseRedirectionRequested:
		return "RedirectionRequested"
	default:
		return "UnknownCause"
	}
}

func cloneIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}

func composeFarInfo(far *ie.IE, localIp net.IP, farInfo ebpf_datapath.FarInfo) (ebpf_datapath.FarInfo, error) {
	logger.PFCPSessLog.Infof("localIp: %v", localIp)
	farInfo.LocalIP = binary.LittleEndian.Uint32(localIp)
	if applyAction, err := far.ApplyAction(); err == nil {
		farInfo.Action = applyAction[0]
	}
	var forward []*ie.IE
	var err error
	if far.Type == ie.CreateFAR {
		forward, err = far.ForwardingParameters()
	} else if far.Type == ie.UpdateFAR {
		forward, err = far.UpdateForwardingParameters()
	} else {
		return ebpf_datapath.FarInfo{}, fmt.Errorf("unsupported IE type")
	}
	if err == nil {
		outerHeaderCreationIndex := findIEindex(forward, 84) // IE Type Outer Header Creation
		if outerHeaderCreationIndex == -1 {
			logger.PFCPSessLog.Warnf("WARN: No OuterHeaderCreation")
		} else {
			outerHeaderCreation, _ := forward[outerHeaderCreationIndex].OuterHeaderCreation()
			if outerHeaderCreation != nil {
				farInfo.OuterHeaderCreation = uint8(outerHeaderCreation.OuterHeaderCreationDescription >> 8)
				farInfo.Teid = outerHeaderCreation.TEID
				if outerHeaderCreation.HasIPv4() {
					farInfo.RemoteIP = binary.LittleEndian.Uint32(outerHeaderCreation.IPv4Address)
				}
				if outerHeaderCreation.HasIPv6() {
					logger.InitLog.Infof("WARN: IPv6 not supported yet, ignoring")
					return ebpf_datapath.FarInfo{}, fmt.Errorf("IPv6 not supported yet")
				}
			}
		}
	}
	transportLevelMarking, err := GetTransportLevelMarking(far)
	if err == nil {
		farInfo.TransportLevelMarking = transportLevelMarking
	}
	return farInfo, nil
}

func updateQer(qerInfo *ebpf_datapath.QerInfo, qer *ie.IE) {

	gateStatusDL, err := qer.GateStatusDL()
	if err == nil {
		qerInfo.GateStatusDL = gateStatusDL
	}
	gateStatusUL, err := qer.GateStatusUL()
	if err == nil {
		qerInfo.GateStatusUL = gateStatusUL
	}
	maxBitrateDL, err := qer.MBRDL()
	if err == nil {
		qerInfo.MaxBitrateDL = uint32(maxBitrateDL) * 1000
	}
	maxBitrateUL, err := qer.MBRUL()
	if err == nil {
		qerInfo.MaxBitrateUL = uint32(maxBitrateUL) * 1000
	}
	qfi, err := qer.QFI()
	if err == nil {
		qerInfo.Qfi = qfi
	}
	qerInfo.StartUL = 0
	qerInfo.StartDL = 0
}

func GetTransportLevelMarking(far *ie.IE) (uint16, error) {
	for _, informationalElement := range far.ChildIEs {
		if informationalElement.Type == ie.TransportLevelMarking {
			return informationalElement.TransportLevelMarking()
		}
	}
	return 0, fmt.Errorf("no TransportLevelMarking found")
}

func applyPDR(spdrInfo SPDRInfo, EBPFMapManager ebpf_datapath.EBPFMapInterface) {

	if spdrInfo.Ipv4 != nil {
		if err := EBPFMapManager.PutPdrDownlink(spdrInfo.Ipv4, spdrInfo.PdrInfo); err != nil {
			logger.InitLog.Infof("Can't apply IPv4 PDR: %s", err.Error())
		}
	} else if spdrInfo.Ipv6 != nil {

		logger.InitLog.Infof("Can't apply IPv6 PDR: ")

	} else {
		if err := EBPFMapManager.PutPdrUplink(spdrInfo.Teid, spdrInfo.PdrInfo); err != nil {
			logger.InitLog.Infof("Can't apply GTP PDR: %s", err.Error())
		}
	}
}
func applyPDR_for_update(gnbip string, spdrInfo SPDRInfo, EBPFMapManager ebpf_datapath.EBPFMapInterface) {

	if err := EBPFMapManager.Update_m_arpUplink(); err != nil {
		logger.InitLog.Infof("Can't apply IPv4 PDR: %s", err.Error())
	}
	if err := EBPFMapManager.Update_m_arpDownlink(gnbip); err != nil {
		logger.InitLog.Infof("Can't apply IPv4 PDR: %s", err.Error())
	}
	if spdrInfo.Ipv4 != nil {
		if err := EBPFMapManager.PutPdrDownlink(spdrInfo.Ipv4, spdrInfo.PdrInfo); err != nil {
			logger.InitLog.Infof("Can't apply IPv4 PDR: %s", err.Error())
		}
	} else if spdrInfo.Ipv6 != nil {
		logger.InitLog.Infof("Can't apply IPv6 PDR")

		// if err := EBPFMapManager.PutDownlinkPdrIp6(spdrInfo.Ipv6, spdrInfo.PdrInfo); err != nil {
		// 	logger.InitLog.Infof("Can't apply IPv6 PDR: %s", err.Error())
		// }
	} else {
		if err := EBPFMapManager.PutPdrUplink(spdrInfo.Teid, spdrInfo.PdrInfo); err != nil {
			logger.InitLog.Infof("Can't apply GTP PDR: %s", err.Error())
		}
	}
}

func processCreatedPDRs(createdPDRs []SPDRInfo, n3Address net.IP) []*ie.IE {
	var additionalIEs []*ie.IE
	for _, pdr := range createdPDRs {
		if pdr.Allocated {
			if pdr.Ipv4 != nil {
				additionalIEs = append(additionalIEs, ie.NewCreatedPDR(ie.NewPDRID(uint16(pdr.PdrID)), ie.NewUEIPAddress(flagPresentIPv4, pdr.Ipv4.String(), "", 0, 0)))
			} else if pdr.Ipv6 != nil {

			} else {
				additionalIEs = append(additionalIEs, ie.NewCreatedPDR(ie.NewPDRID(uint16(pdr.PdrID)), ie.NewFTEID(0x01, pdr.Teid, cloneIP(n3Address), nil, 0)))
			}
		}
	}
	return additionalIEs
}
