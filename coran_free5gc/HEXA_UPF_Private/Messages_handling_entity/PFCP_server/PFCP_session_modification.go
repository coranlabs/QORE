package PFCP_server

import (
	"fmt"
	"net"
	"github.com/coranlabs/CORAN_GO_PFCP/message"

	"github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/logger"
	ebpf_datapath "github.com/coranlabs/CORAN_UPF_eBPF/eBPF_Datapath_entity"
	
)

func (conn *Pfcp_Link) PFCP_Session_Modification_Procedure(msg message.Message, addr string) (message.Message, error) {
	req := msg.(*message.SessionModificationRequest)
	logger.PFCPSessLog.Infof("Got Session Modification Request from: %s", addr)
	logger.PFCPSessLog.Infof("PFCP Message: %v.", msg)
	logger.PFCPSessLog.Infof("Finding association for %s", addr)
	association, ok := conn.NodeAssociations[addr]
	if !ok {
		logger.PFCPSessLog.Infof("Rejecting Session Modification Request from: %s (no association)", addr)
		modResp := message.NewSessionModificationResponse(0, 0, req.SEID(), req.Sequence(), 0, newIeNodeID(conn.NodeId), ie.NewCause(ie.CauseNoEstablishedPFCPAssociation))
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			fmt.Printf("Error resolving address: %v\n", err)
			return nil, err
		}
		fmt.Printf("UDP Address: %v\n", udpAddr)

		// Using as net.Addr
		var netaddr net.Addr = udpAddr
		conn.sendRspTo(modResp, netaddr)
		return modResp, nil
	}

	logger.PFCPSessLog.Infof("Finding session %d", req.SEID())
	session, ok := association.Sessions[req.SEID()]
	if !ok {
		logger.PFCPSessLog.Infof("Rejecting Session Modification Request from: %s (unknown SEID)", addr)
		modResp := message.NewSessionModificationResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.NodeId), ie.NewCause(ie.CauseSessionContextNotFound))
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			fmt.Printf("Error resolving address: %v\n", err)
			return nil, err
		}
		fmt.Printf("UDP Address: %v\n", udpAddr)

		// Using as net.Addr
		var netaddr net.Addr = udpAddr
		conn.sendRspTo(modResp, netaddr)
		return modResp, nil
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
		EBPFMapManager := conn.EBPFMapManager

		for _, far := range req.CreateFAR {
			farInfo, err := composeFarInfo(far, conn.N3Address, ebpf_datapath.FarInfo{})
			if err != nil {
				logger.PFCPSessLog.Infof("Error extracting FAR info: %s", err.Error())
				continue
			}

			farid, _ := far.FARID()
			logger.PFCPSessLog.Infof("Saving FAR info to session: %d, %+v", farid, farInfo)
			if internalId, err := EBPFMapManager.NewFar(farid, farInfo); err == nil {
				session.NewFar(farid, internalId, farInfo)
			} else {
				logger.PFCPSessLog.Infof("Can't put FAR: %s", err.Error())
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
			sFarInfo.FarInfo, err = composeFarInfo(far, conn.N3Address, sFarInfo.FarInfo)
			if err != nil {
				logger.PFCPSessLog.Infof("Error extracting FAR info: %s", err.Error())
				continue
			}
			logger.PFCPSessLog.Infof("Updating FAR info: %d, %+v", farid, sFarInfo)
			session.UpdateFar(farid, sFarInfo.FarInfo)
			gnbip = ConvertUint32ToReversedIP(sFarInfo.FarInfo.RemoteIP)
			logger.PFCPSessLog.Infof("Gnbip IP: %s", gnbip)
			if err := EBPFMapManager.UpdateFar(farid, sFarInfo.FarInfo); err != nil {
				logger.PFCPSessLog.Infof("Can't update FAR: %s", err.Error())
			}
		}

		for _, far := range req.RemoveFAR {
			farid, _ := far.FARID()
			logger.PFCPSessLog.Infof("Removing FAR: %d", farid)
			sFarInfo := session.RemoveFar(farid)
			if err := EBPFMapManager.DeleteFar(sFarInfo.GlobalId); err != nil {
				logger.PFCPSessLog.Infof("Can't remove FAR: %s", err.Error())
			}
		}

		for _, pdr := range req.CreatePDR {
			// PDR should be created last, because we need to reference FARs and QERs global id
			pdrId, err := pdr.PDRID()
			if err != nil {
				logger.PFCPSessLog.Infof("PDR ID missing")
				continue
			}
			logger.PFCPSessLog.Infof("creating PDR in modification to session: %d", pdrId)
			spdrInfo := SPDRInfo{PdrID: uint32(pdrId)}

			if err := pdrContext.extractPDR(pdr, &spdrInfo); err == nil {
				session.PutPDR(spdrInfo.PdrID, spdrInfo)
				//applyPDR(spdrInfo, EBPFMapManager)
				logger.PFCPSessLog.Infof("Updating gnbip %v", gnbip)
				applyPDR_for_update(gnbip, spdrInfo, EBPFMapManager)
				createdPDRs = append(createdPDRs, spdrInfo)
			} else {
				logger.PFCPSessLog.Infof("Error extracting PDR info: %s", err.Error())
			}
		}

		for _, pdr := range req.UpdatePDR {
			pdrId, err := pdr.PDRID()
			if err != nil {
				return fmt.Errorf("PDR ID missing")
			}
			logger.PFCPSessLog.Infof("Updating PDR ID: %d", pdrId)
			spdrInfo := session.GetPDR(pdrId)
			if err := pdrContext.extractPDR(pdr, &spdrInfo); err == nil {
				session.PutPDR(uint32(pdrId), spdrInfo)
				//applyPDR(spdrInfo, EBPFMapManager)
				logger.PFCPSessLog.Infof("Updating gnbip %v", gnbip)
				applyPDR_for_update(gnbip, spdrInfo, EBPFMapManager)
			} else {
				logger.PFCPSessLog.Infof("Error extracting PDR info: %v", err)
				logger.PFCPSessLog.Debugf("Error extracting PDR info: %s", err.Error())
			}
		}

		for _, pdr := range req.RemovePDR {
			pdrId, _ := pdr.PDRID()
			if _, ok := session.PDRs[uint32(pdrId)]; ok {
				logger.PFCPSessLog.Infof("Removing uplink PDR: %d", pdrId)
				sPDRInfo := session.RemovePDR(uint32(pdrId))

				if err := pdrContext.deletePDR(sPDRInfo, EBPFMapManager); err != nil {
					logger.PFCPSessLog.Infof("Failed to remove uplink PDR: %v", err)
				}
			}
		}

		return nil
	}()
	if err != nil {
		logger.PFCPSessLog.Infof("Rejecting Session Modification Request from: %s (failed to apply rules)", err)
		modResp := message.NewSessionModificationResponse(0, 0, session.RemoteSEID, req.Sequence(), 0, newIeNodeID(conn.NodeId), ie.NewCause(ie.CauseRuleCreationModificationFailure))
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			fmt.Printf("Error resolving address: %v\n", err)
			return nil, err
		}
		fmt.Printf("UDP Address: %v\n", udpAddr)

		// Using as net.Addr
		var netaddr net.Addr = udpAddr
		conn.sendRspTo(modResp, netaddr)
		return modResp, nil
	}

	association.Sessions[req.SEID()] = session

	additionalIEs := []*ie.IE{
		ie.NewCause(ie.CauseRequestAccepted),
		newIeNodeID(conn.NodeId),
	}

	pdrIEs := processCreatedPDRs(createdPDRs, conn.N3Address)
	additionalIEs = append(additionalIEs, pdrIEs...)

	// Send SessionEstablishmentResponse
	modResp := message.NewSessionModificationResponse(0, 0, session.RemoteSEID, req.Sequence(), 0, additionalIEs...)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Printf("Error resolving address: %v\n", err)
		return nil, err
	}
	fmt.Printf("UDP Address: %v\n", udpAddr)

	// Using as net.Addr
	var netaddr net.Addr = udpAddr
	conn.sendRspTo(modResp, netaddr)
	return modResp, nil
}
