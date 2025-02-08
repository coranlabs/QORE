// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	"net"
	"time"

	"github.com/coranlabs/HEXA_UPF/src/logger"
	infoElement "github.com/wmnsk/go-pfcp/ie"
	mes "github.com/wmnsk/go-pfcp/message"
)

func (association *NodeAssociation) NewSequenceID() uint32 {
	association.NextSequenceID += 1
	return association.NextSequenceID
}

func (association *NodeAssociation) ScheduleHeartbeat(conn *PfcpConn) {
	association.HeartbeatsActive = true
	logger.AppLog.Infof("Heartbeat started for SMF with nodeid %s", association.Addr)
	for {
		sequence := association.NewSequenceID()
		SendHeartbeatRequest(conn, sequence, association.Addr)
		time.Sleep(5 * time.Second)
	}
}
func SendHeartbeatRequest(conn *PfcpConn, sequenceID uint32, associationAddr string) {
	req := mes.NewHeartbeatRequest(sequenceID, infoElement.NewRecoveryTimeStamp(conn.RecoveryTimestamp), nil)
	logger.AppLog.Debugf("Sent Heartbeat Request to: %s", associationAddr)
	udpAddr, err := net.ResolveUDPAddr("udp", associationAddr+":8805")
	if err == nil {
		if err := conn.SendMessage(req, udpAddr); err != nil {
			logger.AppLog.Infof("Failed to send Heartbeat Request: %s\n", err.Error())
		}
	} else {
		logger.AppLog.Infof("Failed to send Heartbeat Request: %s\n", err.Error())
	}
}

func (association *NodeAssociation) HandleHeartbeatTimeout() bool {
	association.Lock()
	defer association.Unlock()

	association.FailedHeartbeats++
	return association.FailedHeartbeats < 5 // value should be config provided
}

func HandlePfcpHeartbeatResponse(conn *PfcpConn, msg mes.Message, addr string) error {
	return nil
}

func (association *NodeAssociation) HandleHeartbeat(sequence uint32) {}
