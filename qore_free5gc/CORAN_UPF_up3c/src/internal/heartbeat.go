// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	"net"
	"time"

	infoElement "github.com/coranlabs/CORAN_GO_PFCP/ie"
	mes "github.com/coranlabs/CORAN_GO_PFCP/message"
	"github.com/coranlabs/HEXA_UPF/src/logger"
)

func HandlePfcpHeartbeatResponse(conn *PfcpConn, msg mes.Message, addr string) error {
	hbresp := msg.(*mes.HeartbeatResponse)
	ts, err := hbresp.RecoveryTimeStamp.RecoveryTimeStamp()
	if err != nil {
		logger.AppLog.Debugf("Got Heartbeat Response with invalid TS: %s, from: %s", err, addr)
		return err
	} else {
		logger.AppLog.Debugf("Got Heartbeat Response with TS: %s, from: %s", ts, addr)
	}
	return err
}

func (association *NodeAssociation) ScheduleHeartbeat(conn *PfcpConn) {
	association.HeartbeatsActive = true
	logger.Pfcplog.Infof("Heartbeat started for SMF with nodeid %s", association.Addr)
	for {
		sequence := association.NewSequenceID()
		SendHeartbeatRequest(conn, sequence, association.Addr)
		time.Sleep(5 * time.Second)
	}
}

func (association *NodeAssociation) HandleHeartbeat(sequence uint32) {
	association.Lock()
	defer association.Unlock()

	if association.HeartbeatChannel != nil {
		association.HeartbeatChannel <- sequence
	}
}

func SendHeartbeatRequest(conn *PfcpConn, sequenceID uint32, associationAddr string) {
	hbreq := mes.NewHeartbeatRequest(sequenceID, infoElement.NewRecoveryTimeStamp(conn.RecoveryTimestamp), nil)
	logger.Pfcplog.Debugf("Sent Heartbeat Request to: %s", associationAddr)
	udpAddr, err := net.ResolveUDPAddr("udp", associationAddr+":8805")
	if err == nil {
		if err := conn.SendRespose(hbreq, udpAddr); err != nil {
			logger.Pfcplog.Infof("Failed to send Heartbeat Request: %s\n", err.Error())
		}
	} else {
		logger.Pfcplog.Infof("Failed to send Heartbeat Request: %s\n", err.Error())
	}
}
