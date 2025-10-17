package core

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/coranlabs/CORAN_UPF_eBPF/ebpf"

	"github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/CORAN_GO_PFCP/message"
	"github.com/coranlabs/CORAN_UPF_eBPF/logger"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

var errMandatoryIeMissing = fmt.Errorf("mandatory IE missing")
var errNoEstablishedAssociation = fmt.Errorf("no established association")

func HandlePfcpSessionEstablishmentRequest(conn *PfcpConnection, msg message.Message, addr string) (message.Message, error) {
	req := msg.(*message.SessionEstablishmentRequest)
	logger.Pfcplog.Infof("Got Session Establishment Request from: %s.", addr)
	remoteSEID, err := validateRequest(req.NodeID, req.CPFSEID)
	if err != nil {
		logger.Pfcplog.Infof("Rejecting Session Establishment Request from: %s (missing NodeID or F-SEID)", addr)
		PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseMandatoryIEMissing)).Inc()
		return message.NewSessionEstablishmentResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.nodeId), convertErrorToIeCause(err)), nil
	}

	association, ok := conn.NodeAssociations[addr]
	if !ok {
		logger.Pfcplog.Infof("Rejecting Session Establishment Request from: %s (no association)", addr)
		PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseNoEstablishedPFCPAssociation)).Inc()
		return message.NewSessionEstablishmentResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.nodeId), ie.NewCause(ie.CauseNoEstablishedPFCPAssociation)), nil
	}

	localSEID := association.NewLocalSEID()

	session := NewSession(localSEID, remoteSEID.SEID)

	printSessionEstablishmentRequest(req)
	// #TODO: Implement rollback on error
	createdPDRs := []SPDRInfo{}
	pdrContext := NewPDRCreationContext(session, conn.ResourceManager)

	err = func() error {
		mapOperations := conn.mapOperations
		for _, far := range req.CreateFAR {
			farInfo, err := composeFarInfo(far, conn.n3Address.To4(), ebpf.FarInfo{})
			if err != nil {
				logger.Pfcplog.Infof("Error extracting FAR info: %s", err.Error())
				continue
			}

			farid, _ := far.FARID()
			logger.Pfcplog.Infof("Saving FAR info to session: %d, %+v", farid, farInfo)
			if internalId, err := mapOperations.NewFar(farid, farInfo); err == nil {
				session.NewFar(farid, internalId, farInfo)
			} else {
				logger.Pfcplog.Infof("Can't put FAR: %s", err.Error())
				return err
			}
		}

		for _, qer := range req.CreateQER {
			qerInfo := ebpf.QerInfo{}
			qerId, err := qer.QERID()
			if err != nil {
				return fmt.Errorf("QER ID missing")
			}
			updateQer(&qerInfo, qer)
			logger.Pfcplog.Infof("Saving QER info to session: %d, %+v", qerId, qerInfo)
			if internalId, err := mapOperations.NewQer(qerInfo); err == nil {
				session.NewQer(qerId, internalId, qerInfo)
			} else {
				logger.Pfcplog.Infof("Can't put QER: %s", err.Error())
				return err
			}
		}

		for _, pdr := range req.CreatePDR {
			// PDR should be created last, because we need to reference FARs and QERs global id
			pdrId, err := pdr.PDRID()
			if err != nil {
				continue
			}

			logger.Pfcplog.Infof("Saving PDR info to session: %d", pdrId)
			spdrInfo := SPDRInfo{PdrID: uint32(pdrId)}

			if err := pdrContext.extractPDR(pdr, &spdrInfo); err == nil {
				logger.Pfcplog.Infof("PDR info extracted: %+v", spdrInfo)
				session.PutPDR(spdrInfo.PdrID, spdrInfo)
				applyPDR(spdrInfo, mapOperations)
				createdPDRs = append(createdPDRs, spdrInfo)
			} else {
				log.Error().Msgf("error extracting PDR info: %s", err.Error())
			}
		}
		return nil
	}()

	if err != nil {
		logger.Pfcplog.Infof("Rejecting Session Establishment Request from: %s (error in applying IEs)", err)
		PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseRuleCreationModificationFailure)).Inc()
		return message.NewSessionEstablishmentResponse(0, 0, remoteSEID.SEID, req.Sequence(), 0, newIeNodeID(conn.nodeId), ie.NewCause(ie.CauseRuleCreationModificationFailure)), nil
	}

	// Reassigning is the best I can think of for now
	association.Sessions[localSEID] = session
	conn.NodeAssociations[addr] = association

	additionalIEs := []*ie.IE{
		newIeNodeID(conn.nodeId),
		ie.NewCause(ie.CauseRequestAccepted),
		ie.NewFSEID(localSEID, cloneIP(conn.nodeAddrV4), nil),
	}

	pdrIEs := processCreatedPDRs(createdPDRs, cloneIP(conn.n3Address))
	additionalIEs = append(additionalIEs, pdrIEs...)

	// Send SessionEstablishmentResponse
	estResp := message.NewSessionEstablishmentResponse(0, 0, remoteSEID.SEID, req.Sequence(), 0, additionalIEs...)
	PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseRequestAccepted)).Inc()
	logger.Pfcplog.Infof("Session Establishment Request from %s accepted.", addr)
	return estResp, nil
}

func HandlePfcpSessionDeletionRequest(conn *PfcpConnection, msg message.Message, addr string) (message.Message, error) {
	req := msg.(*message.SessionDeletionRequest)
	logger.Pfcplog.Infof("Got Session Deletion Request from: %s", addr)
	association, ok := conn.NodeAssociations[addr]
	if !ok {
		logger.Pfcplog.Infof("Rejecting Session Deletion Request from: %s (no association)", addr)
		PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseNoEstablishedPFCPAssociation)).Inc()
		return message.NewSessionDeletionResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.nodeId), ie.NewCause(ie.CauseNoEstablishedPFCPAssociation)), nil
	}
	printSessionDeleteRequest(req)

	session, ok := association.Sessions[req.SEID()]
	if !ok {
		logger.Pfcplog.Infof("Rejecting Session Deletion Request from: %s (unknown SEID)", addr)
		PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseSessionContextNotFound)).Inc()
		return message.NewSessionDeletionResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.nodeId), ie.NewCause(ie.CauseSessionContextNotFound)), nil
	}
	mapOperations := conn.mapOperations
	pdrContext := NewPDRCreationContext(session, conn.ResourceManager)
	for _, pdrInfo := range session.PDRs {
		if err := pdrContext.deletePDR(pdrInfo, mapOperations); err != nil {
			PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseRuleCreationModificationFailure)).Inc()
			return message.NewSessionDeletionResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.nodeId), ie.NewCause(ie.CauseRuleCreationModificationFailure)), err
		}
	}
	for id := range session.FARs {
		if err := mapOperations.DeleteFar(id); err != nil {
			PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseRuleCreationModificationFailure)).Inc()
			return message.NewSessionDeletionResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.nodeId), ie.NewCause(ie.CauseRuleCreationModificationFailure)), err
		}
	}
	for id := range session.QERs {
		if err := mapOperations.DeleteQer(id); err != nil {
			PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseRuleCreationModificationFailure)).Inc()
			return message.NewSessionDeletionResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.nodeId), ie.NewCause(ie.CauseRuleCreationModificationFailure)), err
		}
	}
	logger.Pfcplog.Infof("Deleting session: %d", req.SEID())
	delete(association.Sessions, req.SEID())

	conn.ReleaseResources(req.SEID())

	PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseRequestAccepted)).Inc()
	return message.NewSessionDeletionResponse(0, 0, session.RemoteSEID, req.Sequence(), 0, newIeNodeID(conn.nodeId), ie.NewCause(ie.CauseRequestAccepted)), nil
}

func ConvertUint32ToReversedIP(n uint32) string {
	ip := make(net.IP, 4)
	ip[3] = byte(n >> 24)
	ip[2] = byte(n >> 16)
	ip[1] = byte(n >> 8)
	ip[0] = byte(n)
	return ip.String()
}
func HandlePfcpSessionModificationRequest(conn *PfcpConnection, msg message.Message, addr string) (message.Message, error) {
	req := msg.(*message.SessionModificationRequest)
	logger.Pfcplog.Infof("Got Session Modification Request from: %s", addr)

	logger.Pfcplog.Infof("Finding association for %s", addr)
	association, ok := conn.NodeAssociations[addr]
	if !ok {
		logger.Pfcplog.Infof("Rejecting Session Modification Request from: %s (no association)", addr)
		PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseNoEstablishedPFCPAssociation)).Inc()
		return message.NewSessionModificationResponse(0, 0, req.SEID(), req.Sequence(), 0, newIeNodeID(conn.nodeId), ie.NewCause(ie.CauseNoEstablishedPFCPAssociation)), nil
	}

	logger.Pfcplog.Infof("Finding session %d", req.SEID())
	session, ok := association.Sessions[req.SEID()]
	if !ok {
		logger.Pfcplog.Infof("Rejecting Session Modification Request from: %s (unknown SEID)", addr)
		PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseSessionContextNotFound)).Inc()
		return message.NewSessionModificationResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.nodeId), ie.NewCause(ie.CauseSessionContextNotFound)), nil
	}

	// This IE shall be present if the CP function decides to change its F-SEID for the PFCP session. The UP function
	// shall use the new CP F-SEID for subsequent PFCP Session related messages for this PFCP Session
	if req.CPFSEID != nil {
		remoteSEID, err := req.CPFSEID.FSEID()
		if err == nil {
			session.RemoteSEID = remoteSEID.SEID

			association.Sessions[req.SEID()] = session // FIXME
			conn.NodeAssociations[addr] = association  // FIXME
		}
	}

	printSessionModificationRequest(req)

	// #TODO: Implement rollback on error
	createdPDRs := []SPDRInfo{}
	pdrContext := NewPDRCreationContext(session, conn.ResourceManager)

	err := func() error {
		mapOperations := conn.mapOperations

		for _, far := range req.CreateFAR {
			farInfo, err := composeFarInfo(far, conn.n3Address.To4(), ebpf.FarInfo{})
			if err != nil {
				logger.Pfcplog.Infof("Error extracting FAR info: %s", err.Error())
				continue
			}

			farid, _ := far.FARID()
			logger.Pfcplog.Infof("Saving FAR info to session: %d, %+v", farid, farInfo)
			if internalId, err := mapOperations.NewFar(farid, farInfo); err == nil {
				session.NewFar(farid, internalId, farInfo)
			} else {
				logger.Pfcplog.Infof("Can't put FAR: %s", err.Error())
				return err
			}
		}
		var gnbip string
		for _, far := range req.UpdateFAR {
			farid, err := far.FARID()
			if err != nil {
				return err
			}
			sFarInfo := session.GetFar(farid)
			sFarInfo.FarInfo, err = composeFarInfo(far, conn.n3Address.To4(), sFarInfo.FarInfo)
			if err != nil {
				logger.Pfcplog.Infof("Error extracting FAR info: %s", err.Error())
				continue
			}
			logger.Pfcplog.Infof("Updating FAR info: %d, %+v", farid, sFarInfo)
			session.UpdateFar(farid, sFarInfo.FarInfo)
			gnbip = ConvertUint32ToReversedIP(sFarInfo.FarInfo.RemoteIP)
			logger.Pfcplog.Infof("Gnbip IP: %s", gnbip)
			if err := mapOperations.UpdateFar(farid, sFarInfo.FarInfo); err != nil {
				logger.Pfcplog.Infof("Can't update FAR: %s", err.Error())
			}
		}

		for _, far := range req.RemoveFAR {
			farid, _ := far.FARID()
			logger.Pfcplog.Infof("Removing FAR: %d", farid)
			sFarInfo := session.RemoveFar(farid)
			if err := mapOperations.DeleteFar(sFarInfo.GlobalId); err != nil {
				logger.Pfcplog.Infof("Can't remove FAR: %s", err.Error())
			}
		}

		for _, qer := range req.CreateQER {
			qerInfo := ebpf.QerInfo{}
			qerId, err := qer.QERID()
			if err != nil {
				return fmt.Errorf("QER ID missing")
			}
			updateQer(&qerInfo, qer)
			logger.Pfcplog.Infof("Saving QER info to session: %d, %+v", qerId, qerInfo)
			if internalId, err := mapOperations.NewQer(qerInfo); err == nil {
				session.NewQer(qerId, internalId, qerInfo)
			} else {
				logger.Pfcplog.Infof("Can't put QER: %s", err.Error())
				return err
			}
		}

		for _, qer := range req.UpdateQER {
			qerId, err := qer.QERID() // Probably will be used as ebpf map key
			if err != nil {
				return fmt.Errorf("QER ID missing")
			}
			sQerInfo := session.GetQer(qerId)
			updateQer(&sQerInfo.QerInfo, qer)
			logger.Pfcplog.Infof("Updating QER ID: %d, QER Info: %+v", qerId, sQerInfo)
			session.UpdateQer(qerId, sQerInfo.QerInfo)
			if err := mapOperations.UpdateQer(sQerInfo.GlobalId, sQerInfo.QerInfo); err != nil {
				logger.Pfcplog.Infof("Can't update QER: %s", err.Error())
				return err
			}
		}

		for _, qer := range req.RemoveQER {
			qerId, err := qer.QERID()
			if err != nil {
				return fmt.Errorf("QER ID missing")
			}
			logger.Pfcplog.Infof("Removing QER ID: %d", qerId)
			sQerInfo := session.RemoveQer(qerId)
			if err := mapOperations.DeleteQer(sQerInfo.GlobalId); err != nil {
				logger.Pfcplog.Infof("Can't remove QER: %s", err.Error())
				return err
			}
		}

		for _, pdr := range req.CreatePDR {
			// PDR should be created last, because we need to reference FARs and QERs global id
			pdrId, err := pdr.PDRID()
			if err != nil {
				logger.Pfcplog.Infof("PDR ID missing")
				continue
			}
			logger.Pfcplog.Infof("creating PDR in modification to session: %d", pdrId)
			spdrInfo := SPDRInfo{PdrID: uint32(pdrId)}

			if err := pdrContext.extractPDR(pdr, &spdrInfo); err == nil {
				session.PutPDR(spdrInfo.PdrID, spdrInfo)
				//applyPDR(spdrInfo, mapOperations)
				logger.Pfcplog.Infof("Updating gnbip %v", gnbip)
				applyPDR_for_update(gnbip, spdrInfo, mapOperations)
				createdPDRs = append(createdPDRs, spdrInfo)
			} else {
				logger.Pfcplog.Infof("Error extracting PDR info: %s", err.Error())
			}
		}

		for _, pdr := range req.UpdatePDR {
			pdrId, err := pdr.PDRID()
			if err != nil {
				return fmt.Errorf("PDR ID missing")
			}
			logger.Pfcplog.Infof("Updating PDR ID: %d", pdrId)
			spdrInfo := session.GetPDR(pdrId)
			if err := pdrContext.extractPDR(pdr, &spdrInfo); err == nil {
				session.PutPDR(uint32(pdrId), spdrInfo)
				//applyPDR(spdrInfo, mapOperations)
				logger.Pfcplog.Infof("Updating gnbip %v", gnbip)
				applyPDR_for_update(gnbip, spdrInfo, mapOperations)
			} else {
				logger.Pfcplog.Infof("Error extracting PDR info: %v", err)
				logger.Pfcplog.Tracef("Error extracting PDR info: %s", err.Error())
			}
		}

		for _, pdr := range req.RemovePDR {
			pdrId, _ := pdr.PDRID()
			if _, ok := session.PDRs[uint32(pdrId)]; ok {
				logger.Pfcplog.Infof("Removing uplink PDR: %d", pdrId)
				sPDRInfo := session.RemovePDR(uint32(pdrId))

				if err := pdrContext.deletePDR(sPDRInfo, mapOperations); err != nil {
					logger.Pfcplog.Infof("Failed to remove uplink PDR: %v", err)
				}
			}
		}

		return nil
	}()
	if err != nil {
		logger.Pfcplog.Infof("Rejecting Session Modification Request from: %s (failed to apply rules)", err)
		PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseRuleCreationModificationFailure)).Inc()
		return message.NewSessionModificationResponse(0, 0, session.RemoteSEID, req.Sequence(), 0, newIeNodeID(conn.nodeId), ie.NewCause(ie.CauseRuleCreationModificationFailure)), nil
	}

	association.Sessions[req.SEID()] = session

	additionalIEs := []*ie.IE{
		ie.NewCause(ie.CauseRequestAccepted),
		newIeNodeID(conn.nodeId),
	}

	pdrIEs := processCreatedPDRs(createdPDRs, conn.n3Address)
	additionalIEs = append(additionalIEs, pdrIEs...)

	// Send SessionEstablishmentResponse
	modResp := message.NewSessionModificationResponse(0, 0, session.RemoteSEID, req.Sequence(), 0, additionalIEs...)
	PfcpMessageRxErrors.WithLabelValues(msg.MessageTypeName(), causeToString(ie.CauseRequestAccepted)).Inc()
	return modResp, nil
}

func convertErrorToIeCause(err error) *ie.IE {
	switch err {
	case errMandatoryIeMissing:
		return ie.NewCause(ie.CauseMandatoryIEMissing)
	case errNoEstablishedAssociation:
		return ie.NewCause(ie.CauseNoEstablishedPFCPAssociation)
	default:
		logger.Pfcplog.Infof("Unknown error: %s", err.Error())
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

func composeFarInfo(far *ie.IE, localIp net.IP, farInfo ebpf.FarInfo) (ebpf.FarInfo, error) {
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
		return ebpf.FarInfo{}, fmt.Errorf("unsupported IE type")
	}
	if err == nil {
		outerHeaderCreationIndex := findIEindex(forward, 84) // IE Type Outer Header Creation
		if outerHeaderCreationIndex == -1 {
			logger.Pfcplog.Warnf("WARN: No OuterHeaderCreation")
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
					return ebpf.FarInfo{}, fmt.Errorf("IPv6 not supported yet")
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

func updateQer(qerInfo *ebpf.QerInfo, qer *ie.IE) {

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
