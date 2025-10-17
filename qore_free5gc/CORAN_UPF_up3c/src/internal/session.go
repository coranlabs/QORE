// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	"encoding/binary"
	"fmt"
	"net"

	infoElement "github.com/coranlabs/CORAN_GO_PFCP/ie"
	mes "github.com/coranlabs/CORAN_GO_PFCP/message"
	"github.com/coranlabs/HEXA_UPF/internal/logger"

)

func (s *Session) HandleSessionEstablishmentRequest(conn *PfcpConn, msg mes.Message, addr string) (mes.Message, error) {
	// TODO: error response
	logger.Pfcplog.Trace("handleSessionEstablishmentRequest")

	logger.Pfcplog.Infof("Got Session Establishment Request from: %s.", addr)

	association := conn.NodeAssociations[addr]

	req, ok := msg.(*mes.SessionEstablishmentRequest)
	if !ok {
		return nil, fmt.Errorf("error at line 25 session.go")

	}

	errUnmarshalReply := func(err error, offendingIE *infoElement.IE) (mes.Message, error) {
		// Build response message
		pfdres := mes.NewSessionEstablishmentResponse(0,
			0,
			0,
			req.SequenceNumber,
			0,
			infoElement.NewCause(infoElement.CauseRequestRejected),
			offendingIE,
		)

		return pfdres, fmt.Errorf("line 40 in session.go %s", err)
	}

	nodeID, err := req.NodeID.NodeID()
	if err != nil {
		return errUnmarshalReply(err, req.NodeID)
	}
	logger.Pfcplog.Trace("node id in session establishment: ", nodeID)
	fseid, err := req.CPFSEID.FSEID()
	if err != nil {
		return errUnmarshalReply(err, req.CPFSEID)
	}

	remoteSEID := fseid.SEID
	fseidIP := ip2int(fseid.IPv4Address)

	logger.Pfcplog.Traceln(fseidIP)

	var localSEID uint64
	localSEID = 0xc000057f40

	session := NewSession(localSEID, remoteSEID)
	session.LocalSEID = localSEID

	logger.Pfcplog.Debug("yha tak aya")
	for _, i := range req.CreatePDR {
		err = session.CreatePDR(localSEID, i)
		logger.Pfcplog.Infof("PDR created")
		if err != nil {
			logger.Pfcplog.Fatalf("Est CreatePDR error: %+v", err)
		}
	}

	for _, i := range req.CreateFAR {
		err = session.CreateFAR(localSEID, i)
		logger.Pfcplog.Infof("FAR created")
		if err != nil {
			logger.Pfcplog.Fatalf("Est CreateFAR error: %+v", err)
		}
	}

	for _, i := range req.CreateQER {
		err = session.CreateQER(localSEID, i)
		logger.Pfcplog.Infof("QER not created skipping")
		if err != nil {
			logger.Pfcplog.Fatalf("Est CreateFAR error: %+v", err)
		}
	}

	association.Sessions[localSEID] = session
	conn.NodeAssociations[addr] = association
	var v6 net.IP

	rsp := mes.NewSessionEstablishmentResponse(
		0,          // mp
		0,          // fo
		remoteSEID, // seid
		req.Header.SequenceNumber,
		0, // pri
		newIeNodeID(conn.nodeAddrV4.String()),
		infoElement.NewCause(infoElement.CauseRequestAccepted),
		infoElement.NewFSEID(localSEID, conn.nodeAddrV4, v6),
		// infoElement.NewCreatedPDR(req.CreatePDR...),
	)
	logger.Pfcplog.Infof("Session Establishment Request from %s accepted.", addr)
	return rsp, nil
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}

	return binary.BigEndian.Uint32(ip)
}

func NewSession(localSEID uint64, remoteSEID uint64) *Session {
	s := &Session{
		LocalSEID:  localSEID,
		RemoteSEID: remoteSEID,
		pdrs:       &pdr{},
		fars:       &far{},
		qers:       &qer{},
	}
	return s
}

func GetTransportLevelMarking(far *infoElement.IE) (uint16, error) {
	for _, informationalElement := range far.ChildIEs {
		if informationalElement.Type == infoElement.TransportLevelMarking {
			return informationalElement.TransportLevelMarking()
		}
	}
	return 0, fmt.Errorf("no TransportLevelMarking found")
}

func (s *Session) HandleSessionModificationRequest(conn *PfcpConn, msg mes.Message, addr string) (mes.Message, error) {
	logger.Pfcplog.Trace("handleSessionModificationRequest")
	// req := msg.(*mes.SessionEstablishmentRequest)
	logger.Pfcplog.Infof("Got Session modification Request from: %s.", addr)
	req, ok := msg.(*mes.SessionModificationRequest)
	if !ok {
		return nil, fmt.Errorf("error at line 356 session.go")
	}

	// nodeID, err := req.NodeID.NodeID()
	// println("node id in session establishment: ", nodeID)

	var remoteSEID uint64

	remoteSEID = 1

	localSEID := req.SEID

	var LocalSEID uint64
	LocalSEID = 0xc000057f40

	session := NewSession(LocalSEID, remoteSEID)
	session.LocalSEID = LocalSEID

	if localSEID() == LocalSEID {
		logger.Pfcplog.Trace("local seid matched")
	} else {
		logger.Pfcplog.Trace("local seid not matched")
	}

	for _, i := range req.CreatePDR {
		err := session.CreatePDR(LocalSEID, i)
		logger.Pfcplog.Infof("pdr created")
		if err != nil {
			logger.Pfcplog.Fatalf("Est CreatePDR error: %+v", err)
		}
	}

	for _, i := range req.CreateFAR {
		err := session.CreateFAR(LocalSEID, i)
		logger.Pfcplog.Infof("far created")
		if err != nil {
			logger.Pfcplog.Fatalf("Est CreateFAR error: %+v", err)
		}
	}

	for _, i := range req.CreateQER {
		err := session.CreateQER(LocalSEID, i)
		logger.Pfcplog.Infof("qer not created skipping")
		if err != nil {
			logger.Pfcplog.Fatalf("Est CreateFAR error: %+v", err)
		}
	}

	for _, i := range req.UpdatePDR {
		err := session.UpdatePDR(LocalSEID, i)
		logger.Pfcplog.Infof("PDR updated")
		if err != nil {
			logger.Pfcplog.Fatalf("Est CreatePDR error: %+v", err)
		}
	}

	for _, i := range req.UpdateFAR {
		err := session.UpdateFAR(LocalSEID, i)
		logger.Pfcplog.Infof("FAR updated")
		if err != nil {
			logger.Pfcplog.Fatalf("Est CreatePDR error: %+v", err)
		}
	}
	logger.Pfcplog.Infof("Session Modification Request from %s accepted", addr)
	resp := mes.NewSessionModificationResponse(0, /* MO?? <-- what's this */
		0,                  /* FO <-- what's this? */
		remoteSEID,         /* seid */
		req.SequenceNumber, /* seq # */
		0,                  /* priority */
		infoElement.NewCause(infoElement.CauseRequestAccepted), /* accept it blindly for the time being */
	)
	return resp, nil
}

func (s *Session) HandleSessionDeletionRequest(conn *PfcpConn, msg mes.Message, addr string) (mes.Message, error) {
	logger.Pfcplog.Tracef("handleSessionDeletionRequest")
	// req := msg.(*mes.SessionEstablishmentRequest)
	logger.Pfcplog.Tracef("Got Session Deletion Request from: %s.", addr)

	req, ok := msg.(*mes.SessionDeletionRequest)
	if !ok {
		return nil, fmt.Errorf("error at line 455 session.go")
	}

	var remoteSEID uint64
	remoteSEID = 2

	localSEID := req.SEID

	var LocalSEID uint64
	LocalSEID = 0xc000057f40

	if localSEID() == LocalSEID {
		logger.Pfcplog.Traceln("match hogyi yayyyy!!!")
	} else {
		logger.Pfcplog.Traceln("match nhi hui")
	}

	resp := mes.NewSessionDeletionResponse(0,
		0,
		remoteSEID,
		req.SequenceNumber,
		0,
		infoElement.NewCause(infoElement.CauseRequestAccepted),
	)

	return resp, nil
}
