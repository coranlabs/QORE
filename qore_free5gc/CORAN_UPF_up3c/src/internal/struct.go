// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	"net"
	"sync"
	"time"
)

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
}

type Session struct {
	LocalSEID  uint64
	RemoteSEID uint64
	pdrs       *pdr
	fars       *far
	qers       *qer
}

type pdr struct {
	LocalSEID          uint64
	FARID              uint32
	QERID              uint32
	PDRID              uint32
	Precedence         uint32
	PDI                pdi
	outerHeaderRemoval uint8
	FTEID              fteid
}

type fteid struct {
	TEID        uint32
	IPv4Address net.IP
}

type pdi struct {
	SourceInterface uint8
	NetworkInstance string
	UeIpAddress     net.IP
}

type far struct {
	farID                uint32
	LocalSEID            uint64
	// sendEndMarker        bool
	// tunnelType           uint8
	// tunnelPort           uint16
	forwardingparameters ForwardingParameters
	applyAction          uint16
	NetworkInstance      string
}
type ForwardingParameters struct {
	DestinationInterface uint8
	OuterHeaderCreation  OuterHeaderCreation
}

type OuterHeaderCreation struct {
	OuterHeaderCreationDescription uint16
	TEID                           uint32
	IPv4                           net.IP
	Port                           uint16
}

type qer struct {
	// qerID    uint32
	// qosLevel uint8
	// qfi      uint8
	// ulStatus uint8
	// dlStatus uint8
	// ulMbr    uint64
	// dlMbr    uint64
	// ulGbr    uint64
	// dlGbr    uint64
	// fseID    uint64
	// fseidIP  uint32
}

// just define qer like ```type qer struct {}``` 

type SdfFilter struct {
	Protocol     uint8
	SrcAddress   IpWMask
	SrcPortRange PortRange
	DstAddress   IpWMask
	DstPortRange PortRange
}

type IpWMask struct {
	Type uint8
	Ip   net.IP
	Mask net.IPMask
}

type PortRange struct {
	LowerBound uint16
	UpperBound uint16
}

// type nodeid struct {
// 	local  string
// 	remote string
// }

type PfcpConn struct {
	udpConn *net.UDPConn
	// pfcpHandlerMap   PfcpHandlerMap
	NodeAssociations  map[string]*NodeAssociation
	nodeId            string
	nodeAddrV4        net.IP
	n3Address         net.IP
	RecoveryTimestamp time.Time
}
