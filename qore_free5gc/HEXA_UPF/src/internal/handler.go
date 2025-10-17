// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	"log"
	"net"
	"time"

	"github.com/coranlabs/HEXA_UPF/src/logger"
	infoElement "github.com/wmnsk/go-pfcp/ie"
	mes "github.com/wmnsk/go-pfcp/message"
)

func Handle(conn *PfcpConn, buf []byte, addr *net.UDPAddr) error {
	go CheckAssociation(conn)
	msg, err := mes.Parse(buf)
	stringIpAddr := addr.IP.String()
	if err != nil {

		return err
	}

	for {
		switch msg.MessageType() {
		case mes.MsgTypeAssociationSetupRequest:
			Msg, err := HandlePfcpAssociationSetupRequest(conn, msg, stringIpAddr)
			if err != nil {

				return err
			}
			return conn.SendMessage(Msg, addr)
		default:
			log.Printf("Got unexpected mes %s: %s, from: %s", msg.MessageTypeName(), msg, addr)
			return nil
		}
	}
}

func HandlePfcpAssociationSetupRequest(conn *PfcpConn, msg mes.Message, addr string) (mes.Message, error) {
	req := msg.(*mes.AssociationSetupRequest)
	remoteNodeID, err := req.NodeID.NodeID()
	logger.AppLog.Infof("Association Setup Request from: %s with NodeID: %s", addr, remoteNodeID)

	logger.AppLog.Traceln("Nodeip: ", remoteNodeID)
	if err != nil {
		logger.AppLog.Infof("Got Association Setup Request with invalid NodeID from: %s", addr)
		asres := mes.NewAssociationSetupResponse(req.SequenceNumber,
			infoElement.NewCause(infoElement.CauseMandatoryIEMissing),
		)
		return asres, nil
	}

	if conn.NodeAssociations == nil {
		conn.NodeAssociations = make(map[string]*NodeAssociation)
	}

	logger.AppLog.Trace(conn.NodeAssociations)

	remoteNode := NewNodeAssociation(remoteNodeID, addr)
	logger.AppLog.Trace("remotenode value", remoteNode)

	logger.AppLog.Trace(conn.nodeId)
	conn.NodeAssociations[addr] = remoteNode
	featuresOctets := []uint8{0, 0, 0}
	upFunctionFeaturesIE := infoElement.NewUPFunctionFeatures(featuresOctets[:]...)
	res := mes.NewAssociationSetupResponse(req.SequenceNumber,
		infoElement.NewCause(infoElement.CauseRequestAccepted),
		infoElement.NewRecoveryTimeStamp(conn.RecoveryTimestamp),
		upFunctionFeaturesIE,
	)
	logger.AppLog.Traceln("response: ", res)
	logger.AppLog.Infof("Association Accepted")
	return res, nil
}

func NewNodeAssociation(remoteNodeID string, addr string) *NodeAssociation {
	return &NodeAssociation{
		ID:               remoteNodeID,
		Addr:             addr,
		NextSessionID:    1,
		NextSequenceID:   1,
		Sessions:         make(map[uint64]*Session),
		HeartbeatChannel: make(chan uint32),
	}
}

func (connection *PfcpConn) SendMessage(msg mes.Message, addr *net.UDPAddr) error {
	responseBytes := make([]byte, msg.MarshalLen())
	if err := msg.MarshalTo(responseBytes); err != nil {
		logger.AppLog.Infof(err.Error())
		return err
	}
	if _, err := connection.Send(responseBytes, addr); err != nil {
		logger.AppLog.Infof(err.Error())
		return err
	}
	return nil
}

func (connection *PfcpConn) Send(b []byte, addr *net.UDPAddr) (int, error) {
	return connection.udpConn.WriteTo(b, addr)
}

func CheckAssociation(conn *PfcpConn) {
	go func() {
		for {
			conn.RefreshAssociations()
			// 5 is hard coded the value should be provided through config
			time.Sleep(time.Duration(5) * time.Second)
		}
	}()
}

func (connection *PfcpConn) RefreshAssociations() {
	for _, assoc := range connection.NodeAssociations {
		if !assoc.HeartbeatsActive {
			go assoc.ScheduleHeartbeat(connection)
		}
	}
}

func (connection *PfcpConn) GetAssociation(assocAddr string) *NodeAssociation {
	if assoc, ok := connection.NodeAssociations[assocAddr]; ok {
		return assoc
	}
	return nil
}

func newIeNodeID(nodeID string) *infoElement.IE {
	ip := net.ParseIP(nodeID)
	logger.AppLog.Trace("node ip from newIeNodeID", ip)
	if ip != nil {
		if ip.To4() != nil {
			return infoElement.NewNodeID(nodeID, "", "")
		}
		return infoElement.NewNodeID("", nodeID, "")
	}
	return infoElement.NewNodeID("", "", nodeID)
}