// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	"net"
	"time"

	infoElement "github.com/coranlabs/CORAN_GO_PFCP/ie"
	mes "github.com/coranlabs/CORAN_GO_PFCP/message"
	"github.com/coranlabs/HEXA_UPF/internal/logger"

)

func Handle(conn *PfcpConn, buf []byte, addr *net.UDPAddr) error {
	go CheckAssociation(conn)

	logger.Pfcplog.Debugf("Handling PFCP mes from %s", addr)
	stringIpAddr := addr.IP.String()
	msg, err := mes.Parse(buf)

	logger.Pfcplog.Trace("mes: ", msg)
	if err != nil {

		return err
	}
	var s *Session
	for {
		switch msg.MessageType() {
		case mes.MsgTypeAssociationSetupRequest:
			logger.Pfcplog.Trace(stringIpAddr)
			Msg, err := HandlePfcpAssociationSetupRequest(conn, msg, stringIpAddr)
			if err != nil {

				return err
			}

			// logger.MainLog.Tracef("outgoing mes: %s", Msg)
			logger.Pfcplog.Tracef("outgoing mes: %s", Msg)
			return conn.SendRespose(Msg, addr)

		case mes.MsgTypeHeartbeatRequest:
			stringIpAddr := addr.IP.String()
			logger.Pfcplog.Infof("its a heartbeat request")
			logger.Pfcplog.Traceln("ip addr :", stringIpAddr)
			return nil
		case mes.MsgTypeSessionEstablishmentRequest:
			Msg, err := s.HandleSessionEstablishmentRequest(conn, msg, stringIpAddr)
			if err != nil {

				return err
			}

			logger.Pfcplog.Tracef("outgoing mes: %s", Msg)

			return conn.SendRespose(Msg, addr)
		case mes.MsgTypeSessionModificationRequest:
			Msg, err := s.HandleSessionModificationRequest(conn, msg, stringIpAddr)
			if err != nil {

				return err
			}

			logger.Pfcplog.Tracef("outgoing mes: %s", Msg)
			return conn.SendRespose(Msg, addr)
		case mes.MsgTypeSessionDeletionRequest:
			Msg, err := s.HandleSessionDeletionRequest(conn, msg, stringIpAddr)
			if err != nil {

				return err
			}

			logger.Pfcplog.Tracef("outgoing mes: %s", Msg)
			return conn.SendRespose(Msg, addr)

		case mes.MsgTypeHeartbeatResponse:
			err := HandlePfcpHeartbeatResponse(conn, msg, stringIpAddr)
			if err != nil {
				logger.Pfcplog.Infof("Error handling PFCP mes: %s", err.Error())
				return err
			}

			return nil

		default:
			logger.Pfcplog.Tracef("Got unexpected mes %s: %s, from: %s", msg.MessageTypeName(), msg, addr)
			return nil
		}
	}
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

func (association *NodeAssociation) NewSequenceID() uint32 {
	association.NextSequenceID += 1
	return association.NextSequenceID
}

func (association *NodeAssociation) HandleHeartbeatTimeout() bool {
	association.Lock()
	defer association.Unlock()

	association.FailedHeartbeats++
	return association.FailedHeartbeats < 5 // value should be config provided
}

func HandlePfcpAssociationSetupRequest(conn *PfcpConn, msg mes.Message, addr string) (mes.Message, error) {
	asreq := msg.(*mes.AssociationSetupRequest)
	remoteNodeID, err := asreq.NodeID.NodeID()
	logger.Pfcplog.Infof("Handling association Setup Request from: %s with NodeID: %s", addr, remoteNodeID)

	logger.Pfcplog.Traceln("nodeip: ", remoteNodeID)
	if err != nil {
		logger.Pfcplog.Infof("Got Association Setup Request with invalid NodeID from: %s", addr)
		asres := mes.NewAssociationSetupResponse(asreq.SequenceNumber,
			infoElement.NewCause(infoElement.CauseMandatoryIEMissing),
		)
		return asres, nil
	}

	if conn.NodeAssociations == nil {
		conn.NodeAssociations = make(map[string]*NodeAssociation)
	}

	logger.Pfcplog.Trace(conn.NodeAssociations)

	remoteNode := NewNodeAssociation(remoteNodeID, addr)
	logger.Pfcplog.Trace("remotenode value", remoteNode)

	logger.Pfcplog.Trace(conn.nodeId)
	// Add or replace RemoteNode to NodeAssociationMap
	conn.NodeAssociations[addr] = remoteNode
	featuresOctets := []uint8{0, 0, 0}
	upFunctionFeaturesIE := infoElement.NewUPFunctionFeatures(featuresOctets[:]...)
	asres := mes.NewAssociationSetupResponse(asreq.SequenceNumber,
		infoElement.NewCause(infoElement.CauseRequestAccepted),
		infoElement.NewRecoveryTimeStamp(conn.RecoveryTimestamp),
		upFunctionFeaturesIE,
	)
	logger.Pfcplog.Traceln("response: ", asres)
	logger.Pfcplog.Tracef("out of HandlePfcpAssociationSetupRequest")
	logger.Pfcplog.Infof("Association Accepted")
	return asres, nil
}

func NewNodeAssociation(remoteNodeID string, addr string) *NodeAssociation {
	return &NodeAssociation{
		ID:               remoteNodeID,
		Addr:             addr,
		NextSessionID:    1,
		NextSequenceID:   1,
		Sessions:         make(map[uint64]*Session),
		HeartbeatChannel: make(chan uint32),
		// AssociationStart: time.Now(),
	}
}

func newIeNodeID(nodeID string) *infoElement.IE {
	ip := net.ParseIP(nodeID)
	logger.Pfcplog.Trace("node ip from newIeNodeID", ip)
	if ip != nil {
		if ip.To4() != nil {
			return infoElement.NewNodeID(nodeID, "", "")
		}
		return infoElement.NewNodeID("", nodeID, "")
	}
	return infoElement.NewNodeID("", "", nodeID)
}

func (connection *PfcpConn) Send(b []byte, addr *net.UDPAddr) (int, error) {
	return connection.udpConn.WriteTo(b, addr)
}

func (connection *PfcpConn) SendRespose(msg mes.Message, addr *net.UDPAddr) error {
	responseBytes := make([]byte, msg.MarshalLen())
	if err := msg.MarshalTo(responseBytes); err != nil {
		logger.Pfcplog.Infof(err.Error())
		return err
	}
	if _, err := connection.Send(responseBytes, addr); err != nil {
		logger.Pfcplog.Infof(err.Error())
		return err
	}
	return nil
}
