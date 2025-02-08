package PFCP_server

import (
	"net"

	"github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/CORAN_GO_PFCP/message"
	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/logger"
)

func (conn *Pfcp_Link) Heartbeat_Procedure(msg message.Message, addr string) (message.Message, error) {
	hbreq := msg.(*message.HeartbeatRequest)
	if association := conn.GetAssociation(addr); association != nil {
		association.ResetFailedHeartbeats()
	}
	logger.Heartbeat_Procedure.Debugf("Got Heartbeat Request from: %s", addr)
	ts, err := hbreq.RecoveryTimeStamp.RecoveryTimeStamp()
	if err != nil {
		logger.Heartbeat_Procedure.Warnf("Got Heartbeat Request with invalid TS: %s, from: %s", err, addr)
		return nil, err
	} else {
		logger.Heartbeat_Procedure.Debugf("Got Heartbeat Request with TS: %s, from: %s connection: %v", ts, addr, conn)
	}

	hbres := message.NewHeartbeatResponse(hbreq.SequenceNumber, ie.NewRecoveryTimeStamp(conn.RecoveryTimestamp))
	logger.Heartbeat_Procedure.Debugf("Sent Heartbeat Response to: %s connection: %v", addr, conn)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		logger.Heartbeat_Procedure.Errorf("Error resolving address: %v\n", err)
		return nil, err
	}
	//fmt.Printf("UDP Address: %v\n", udpAddr)
	//fmt.Printf("connection: %v\n", conn)

	// Using as net.Addr
	var netaddr net.Addr = udpAddr
	conn.sendRspTo(hbres, netaddr)
	return hbres, nil
}

func (connection *Pfcp_Link) GetAssociation(assocAddr string) *NodeAssociation {
	if assoc, ok := connection.NodeAssociations[assocAddr]; ok {
		return assoc
	}
	return nil
}

func (conn *Pfcp_Link) Heartbeat_Response(msg message.Message, addr string) (message.Message, error) {
	hbresp := msg.(*message.HeartbeatResponse)
	ts, err := hbresp.RecoveryTimeStamp.RecoveryTimeStamp()
	if err != nil {
		logger.Heartbeat_Procedure.Warnf("Got Heartbeat Response with invalid TS: %s, from: %s", err, addr)
		return nil, err
	} else {
		logger.Heartbeat_Procedure.Debugf("Got Heartbeat Response with TS: %s, from: %s", ts, addr)
	}

	if association := conn.GetAssociation(addr); association != nil {
		association.HandleHeartbeat(msg.Sequence())
	}
	return nil, err
}

func SendHeartbeatRequest(conn *Pfcp_Link, sequenceID uint32, associationAddr string) {
	hbreq := message.NewHeartbeatRequest(sequenceID, ie.NewRecoveryTimeStamp(conn.RecoveryTimestamp), nil)
	logger.Heartbeat_Procedure.Debugf("Sent Heartbeat Request to: %s", associationAddr)
	udpAddr, err := net.ResolveUDPAddr("udp", associationAddr+":8805")
	if err == nil {
		if err := conn.SendMessage(hbreq, udpAddr); err != nil {
			logger.Heartbeat_Procedure.Infof("Failed to send Heartbeat Request: %s\n", err.Error())
		}
	} else {
		logger.Heartbeat_Procedure.Infof("Failed to send Heartbeat Request: %s\n", err.Error())
	}
}
