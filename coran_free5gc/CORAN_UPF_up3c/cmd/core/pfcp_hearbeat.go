package core

import (
	"net"

	"github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/CORAN_GO_PFCP/message"
	"github.com/coranlabs/HEXA_UPF/internal/logger"
)

func HandlePfcpHeartbeatRequest(conn *PfcpConnection, msg message.Message, addr string) (message.Message, error) {
	hbreq := msg.(*message.HeartbeatRequest)
	if association := conn.GetAssociation(addr); association != nil {
		association.ResetFailedHeartbeats()
	}
	ts, err := hbreq.RecoveryTimeStamp.RecoveryTimeStamp()
	if err != nil {
		logger.Pfcplog.Warnf("Got Heartbeat Request with invalid TS: %s, from: %s", err, addr)
		return nil, err
	} else {
		logger.MainLog.Debugf("Got Heartbeat Request with TS: %s, from: %s", ts, addr)
	}

	hbres := message.NewHeartbeatResponse(hbreq.SequenceNumber, ie.NewRecoveryTimeStamp(conn.RecoveryTimestamp))
	logger.MainLog.Debugf("Sent Heartbeat Response to: %s", addr)
	return hbres, nil
}

func HandlePfcpHeartbeatResponse(conn *PfcpConnection, msg message.Message, addr string) (message.Message, error) {
	hbresp := msg.(*message.HeartbeatResponse)
	ts, err := hbresp.RecoveryTimeStamp.RecoveryTimeStamp()
	if err != nil {
		logger.Pfcplog.Warnf("Got Heartbeat Response with invalid TS: %s, from: %s", err, addr)
		return nil, err
	} else {
		logger.MainLog.Debugf("Got Heartbeat Response with TS: %s, from: %s", ts, addr)
	}

	if association := conn.GetAssociation(addr); association != nil {
		association.HandleHeartbeat(msg.Sequence())
	}
	return nil, err
}

func SendHeartbeatRequest(conn *PfcpConnection, sequenceID uint32, associationAddr string) {
	hbreq := message.NewHeartbeatRequest(sequenceID, ie.NewRecoveryTimeStamp(conn.RecoveryTimestamp), nil)
	logger.MainLog.Debugf("Sent Heartbeat Request to: %s", associationAddr)
	udpAddr, err := net.ResolveUDPAddr("udp", associationAddr+":8805")
	if err == nil {
		if err := conn.SendMessage(hbreq, udpAddr); err != nil {
			logger.MainLog.Infof("Failed to send Heartbeat Request: %s\n", err.Error())
		}
	} else {
		logger.MainLog.Infof("Failed to send Heartbeat Request: %s\n", err.Error())
	}
}
