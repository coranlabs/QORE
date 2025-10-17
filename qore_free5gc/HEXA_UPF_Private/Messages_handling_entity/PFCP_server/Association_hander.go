package PFCP_server

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	config "github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/config"

	"github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/CORAN_GO_PFCP/message"
	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/logger"
)

type PfcpFunc func(conn *Pfcp_Link, msg message.Message, addr string) (message.Message, error)

type PfcpHandlerMap map[uint8]PfcpFunc

func setBit(n uint8, pos uint) uint8 {
	n |= (1 << pos)
	return n
}

// https://www.etsi.org/deliver/etsi_ts/129200_129299/129244/16.04.00_60/ts_129244v160400p.pdf page 95
func (conn *Pfcp_Link) PFCP_Association_Setup_Procedure(msg message.Message, addr string) (message.Message, error) {
	asreq := msg.(*message.AssociationSetupRequest)
	logger.PfcpAssocLog.Infof("Got Association Setup Request from: %s", addr)
	logger.PFCPSessLog.Infof("PFCP Message: %v.", msg)
	if asreq.NodeID == nil {
		logger.PfcpAssocLog.Warnf("Got Association Setup Request without NodeID from: %s", addr)
		// Reject with cause

		asres := message.NewAssociationSetupResponse(asreq.SequenceNumber,
			ie.NewCause(ie.CauseMandatoryIEMissing),
		)
		return asres, nil
	}
	printAssociationSetupRequest(asreq)
	// Get NodeID
	remoteNodeID, err := asreq.NodeID.NodeID()
	if err != nil {
		logger.PfcpAssocLog.Warnf("Got Association Setup Request with invalid NodeID from: %s", addr)
		asres := message.NewAssociationSetupResponse(asreq.SequenceNumber,
			ie.NewCause(ie.CauseMandatoryIEMissing),
		)
		return asres, nil
	}
	// Check if the PFCP Association Setup Request contains a Node ID for which a PFCP association was already established
	if _, ok := conn.NodeAssociations[remoteNodeID]; ok {
		logger.PfcpAssocLog.Warnf("Association Setup Request with NodeID: %s from: %s already exists", remoteNodeID, addr)
		// retain the PFCP sessions that were established with the existing PFCP association and that are requested to be retained, if the PFCP Session Retention Information IE was received in the request; otherwise, delete the PFCP sessions that were established with the existing PFCP association;
		logger.PfcpAssocLog.Warn("Session retention is not yet implemented")
	}

	// If the PFCP Association Setup Request contains a Node ID for which a PFCP association was already established
	// proceed with establishing the new PFCP association (regardless of the Recovery AssociationStart received in the request), overwriting the existing association;
	// if the request is accepted:
	// shall store the Node ID of the CP function as the identifier of the PFCP association;
	// Create RemoteNode from AssociationSetupRequest
	remoteNode := NewNodeAssociation(remoteNodeID, addr)
	// Add or replace RemoteNode to NodeAssociationMap
	conn.NodeAssociations[addr] = remoteNode
	logger.PfcpAssocLog.Infof("Association Saved for NodeID: %s", remoteNodeID)
	featuresOctets := []uint8{0, 0, 0}
	if config.Conf.FeatureFTUP {
		featuresOctets[0] = setBit(featuresOctets[0], 4)
	}
	if config.Conf.FeatureUEIP {
		featuresOctets[2] = setBit(featuresOctets[2], 2)
	}
	upFunctionFeaturesIE := ie.NewUPFunctionFeatures(featuresOctets[:]...)

	// shall send a PFCP Association Setup Response including:
	asres := message.NewAssociationSetupResponse(asreq.SequenceNumber,
		ie.NewCause(ie.CauseRequestAccepted), // a successful cause
		newIeNodeID(conn.NodeId),             // its Node ID;
		ie.NewRecoveryTimeStamp(conn.RecoveryTimestamp),
		upFunctionFeaturesIE,
	)

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Printf("Error resolving address: %v\n", err)
		return nil, err
	}
	fmt.Printf("UDP Address: %v\n", udpAddr)

	// Using as net.Addr
	var netaddr net.Addr = udpAddr
	conn.sendRspTo(asres, netaddr)
	// Send AssociationSetupResponse
	return asres, nil
}

func newIeNodeID(nodeID string) *ie.IE {
	ip := net.ParseIP(nodeID)
	if ip != nil {
		if ip.To4() != nil {
			return ie.NewNodeID(nodeID, "", "")
		}
		return ie.NewNodeID("", nodeID, "")
	}
	return ie.NewNodeID("", "", nodeID)
}

type NodeAssociation struct {
	ID               string
	Addr             string
	NextSessionID    uint64
	NextSequenceID   uint32
	Sessions         map[uint64]*Session
	HeartbeatChannel chan uint32
	FailedHeartbeats uint32
	HeartbeatsActive bool
	sync.Mutex
	// AssociationStart time.Time // Held until propper failure detection is implemented
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

func (association *NodeAssociation) NewLocalSEID() uint64 {
	association.NextSessionID += 1
	return association.NextSessionID
}

func (association *NodeAssociation) NewSequenceID() uint32 {
	association.NextSequenceID += 1
	return association.NextSequenceID
}

func (association *NodeAssociation) ScheduleHeartbeat(conn *Pfcp_Link) {
	association.HeartbeatsActive = true
	ctx := context.Background()

	for {
		sequence := association.NewSequenceID()
		SendHeartbeatRequest(conn, sequence, association.Addr)

		select {
		case <-time.After(time.Duration(config.Conf.HeartbeatTimeout) * time.Second):
			if !association.HandleHeartbeatTimeout() {
				logger.PfcpAssocLog.Warnf("the number of unanswered heartbeats has reached the limit, association deleted: %s", association.Addr)
				close(association.HeartbeatChannel)
				conn.DeleteAssociation(association.Addr)
				return
			}
		case seq := <-association.HeartbeatChannel:
			if sequence == seq {
				association.ResetFailedHeartbeats()
				<-time.After(time.Duration(config.Conf.HeartbeatInterval) * time.Second)
			}
		case <-ctx.Done():
			logger.InitLog.Infof("HeartbeatScheduler context done | association address: %s", association.Addr)
			return
		}
	}
}

func (association *NodeAssociation) ResetFailedHeartbeats() {
	association.Lock()
	association.FailedHeartbeats = 0
	association.Unlock()
}

func (association *NodeAssociation) HandleHeartbeatTimeout() bool {
	association.Lock()
	defer association.Unlock()

	association.FailedHeartbeats++
	return association.FailedHeartbeats < config.Conf.HeartbeatRetries
}

func (association *NodeAssociation) HandleHeartbeat(sequence uint32) {
	association.Lock()
	defer association.Unlock()

	if association.HeartbeatChannel != nil {
		association.HeartbeatChannel <- sequence
	}
}
