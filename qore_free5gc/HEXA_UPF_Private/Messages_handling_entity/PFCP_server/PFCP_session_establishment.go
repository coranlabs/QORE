package PFCP_server

import (
	"fmt"
	"net"

	ebpf_datapath "github.com/coranlabs/CORAN_UPF_eBPF/eBPF_Datapath_entity"

	"github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/CORAN_GO_PFCP/message"
	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/logger"
	"github.com/rs/zerolog/log"
)

var errMandatoryIeMissing = fmt.Errorf("mandatory IE missing")
var errNoEstablishedAssociation = fmt.Errorf("no established association")

func (conn *Pfcp_Link) PFCP_Session_Establishment_Procedure(msg message.Message, addr string) (message.Message, error) {
	req := msg.(*message.SessionEstablishmentRequest)
	logger.PFCPSessLog.Infof("Got Session Establishment Request from: %s.", addr)
	logger.PFCPSessLog.Infof("PFCP Message: %v.", msg)
	remoteSEID, err := validateRequest(req.NodeID, req.CPFSEID)
	logger.PFCPSessLog.Infof("Remote SEID: %v", conn.N3Address)
	if err != nil {
		logger.PFCPSessLog.Infof("Rejecting Session Establishment Request from: %s (missing NodeID or F-SEID)", addr)
		estResp := message.NewSessionEstablishmentResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.NodeId), convertErrorToIeCause(err))
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			fmt.Printf("Error resolving address: %v\n", err)
			return nil, err
		}
		fmt.Printf("UDP Address: %v\n", udpAddr)

		// Using as net.Addr
		var netaddr net.Addr = udpAddr
		conn.sendRspTo(estResp, netaddr)
		return estResp, nil
	}

	association, ok := conn.NodeAssociations[addr]
	if !ok {
		logger.PFCPSessLog.Infof("Rejecting Session Establishment Request from: %s (no association)", addr)
		estResp := message.NewSessionEstablishmentResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.NodeId), ie.NewCause(ie.CauseNoEstablishedPFCPAssociation))
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			fmt.Printf("Error resolving address: %v\n", err)
			return nil, err
		}
		fmt.Printf("UDP Address: %v\n", udpAddr)

		// Using as net.Addr
		var netaddr net.Addr = udpAddr
		conn.sendRspTo(estResp, netaddr)
		return estResp, nil
	}

	localSEID := association.NewLocalSEID()

	session := NewSession(localSEID, remoteSEID.SEID)

	printSessionEstablishmentRequest(req)
	// #TODO: Implement rollback on error
	createdPDRs := []SPDRInfo{}
	pdrContext := NewPDRCreationContext(session, conn.ResourceManager)

	err = func() error {
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

		for _, pdr := range req.CreatePDR {
			// PDR should be created last, because we need to reference FARs and QERs global id
			pdrId, err := pdr.PDRID()
			if err != nil {
				continue
			}

			logger.PFCPSessLog.Infof("Saving PDR info to session: %d", pdrId)
			spdrInfo := SPDRInfo{PdrID: uint32(pdrId)}

			if err := pdrContext.extractPDR(pdr, &spdrInfo); err == nil {
				logger.PFCPSessLog.Infof("PDR info extracted: %+v", spdrInfo)
				session.PutPDR(spdrInfo.PdrID, spdrInfo)
				applyPDR(spdrInfo, EBPFMapManager)
				createdPDRs = append(createdPDRs, spdrInfo)
			} else {
				log.Error().Msgf("error extracting PDR info: %s", err.Error())
			}
		}
		return nil
	}()

	if err != nil {
		logger.PFCPSessLog.Infof("Rejecting Session Establishment Request from: %s (error in applying IEs)", err)
		estResp := message.NewSessionEstablishmentResponse(0, 0, remoteSEID.SEID, req.Sequence(), 0, newIeNodeID(conn.NodeId), ie.NewCause(ie.CauseRuleCreationModificationFailure))
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			fmt.Printf("Error resolving address: %v\n", err)
			return nil, err
		}
		fmt.Printf("UDP Address: %v\n", udpAddr)

		// Using as net.Addr
		var netaddr net.Addr = udpAddr
		conn.sendRspTo(estResp, netaddr)
		return estResp, nil
	}

	// Reassigning is the best I can think of for now
	association.Sessions[localSEID] = session
	conn.NodeAssociations[addr] = association

	additionalIEs := []*ie.IE{
		newIeNodeID(conn.NodeId),
		ie.NewCause(ie.CauseRequestAccepted),
		ie.NewFSEID(localSEID, cloneIP(conn.NodeAddrV4), nil),
	}

	pdrIEs := processCreatedPDRs(createdPDRs, cloneIP(conn.N3Address))
	additionalIEs = append(additionalIEs, pdrIEs...)

	// Send SessionEstablishmentResponse
	estResp := message.NewSessionEstablishmentResponse(0, 0, remoteSEID.SEID, req.Sequence(), 0, additionalIEs...)
	logger.PFCPSessLog.Infof("Session Establishment Request from %s accepted.", addr)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Printf("Error resolving address: %v\n", err)
		return nil, err
	}
	fmt.Printf("UDP Address: %v\n", udpAddr)

	// Using as net.Addr
	var netaddr net.Addr = udpAddr
	conn.sendRspTo(estResp, netaddr)
	return estResp, nil
}
