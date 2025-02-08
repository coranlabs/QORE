package PFCP_server

import (
	"fmt"
	"net"

	"github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/CORAN_GO_PFCP/message"
	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/logger"
)

func (conn *Pfcp_Link) PFCP_Session_Deletion_Procedure(msg message.Message, addr string) (message.Message, error) {
	req := msg.(*message.SessionDeletionRequest)
	logger.PFCPSessLog.Infof("Got Session Deletion Request from: %s", addr)
	logger.PFCPSessLog.Infof("PFCP Message: %v.", msg)
	association, ok := conn.NodeAssociations[addr]
	if !ok {
		logger.PFCPSessLog.Infof("Rejecting Session Deletion Request from: %s (no association)", addr)
		delrsp := message.NewSessionDeletionResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.NodeId), ie.NewCause(ie.CauseNoEstablishedPFCPAssociation))

		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			fmt.Printf("Error resolving address: %v\n", err)
			return nil, err
		}
		fmt.Printf("UDP Address: %v\n", udpAddr)

		// Using as net.Addr
		var netaddr net.Addr = udpAddr
		conn.sendRspTo(delrsp, netaddr)
		return delrsp, nil
	}
	printSessionDeleteRequest(req)

	session, ok := association.Sessions[req.SEID()]
	if !ok {
		logger.PFCPSessLog.Infof("Rejecting Session Deletion Request from: %s (unknown SEID)", addr)

		delrsp := message.NewSessionDeletionResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.NodeId), ie.NewCause(ie.CauseSessionContextNotFound))
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			fmt.Printf("Error resolving address: %v\n", err)
			return nil, err
		}
		fmt.Printf("UDP Address: %v\n", udpAddr)

		// Using as net.Addr
		var netaddr net.Addr = udpAddr
		conn.sendRspTo(delrsp, netaddr)
		return delrsp, nil
	}
	EBPFMapManager := conn.EBPFMapManager
	pdrContext := NewPDRCreationContext(session, conn.ResourceManager)
	for _, pdrInfo := range session.PDRs {
		if err := pdrContext.deletePDR(pdrInfo, EBPFMapManager); err != nil {
			delrsp := message.NewSessionDeletionResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.NodeId), ie.NewCause(ie.CauseRuleCreationModificationFailure))
			udpAddr, err := net.ResolveUDPAddr("udp", addr)
			if err != nil {
				fmt.Printf("Error resolving address: %v\n", err)
				return nil, err
			}
			fmt.Printf("UDP Address: %v\n", udpAddr)

			// Using as net.Addr
			var netaddr net.Addr = udpAddr
			conn.sendRspTo(delrsp, netaddr)
			return delrsp, nil
		}
	}
	for id := range session.FARs {
		if err := EBPFMapManager.DeleteFar(id); err != nil {
			delrsp := message.NewSessionDeletionResponse(0, 0, 0, req.Sequence(), 0, newIeNodeID(conn.NodeId), ie.NewCause(ie.CauseRuleCreationModificationFailure))
			udpAddr, err := net.ResolveUDPAddr("udp", addr)
			if err != nil {
				fmt.Printf("Error resolving address: %v\n", err)
				return nil, err
			}
			fmt.Printf("UDP Address: %v\n", udpAddr)

			// Using as net.Addr
			var netaddr net.Addr = udpAddr
			conn.sendRspTo(delrsp, netaddr)
			return delrsp, nil
		}
	}

	logger.PFCPSessLog.Infof("Deleting session: %d", req.SEID())
	delete(association.Sessions, req.SEID())

	conn.ReleaseResources(req.SEID())

	delrsp := message.NewSessionDeletionResponse(0, 0, session.RemoteSEID, req.Sequence(), 0, newIeNodeID(conn.NodeId), ie.NewCause(ie.CauseRequestAccepted))
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Printf("Error resolving address: %v\n", err)
		return nil, err
	}
	fmt.Printf("UDP Address: %v\n", udpAddr)

	// Using as net.Addr
	var netaddr net.Addr = udpAddr
	conn.sendRspTo(delrsp, netaddr)
	return delrsp, nil
}
