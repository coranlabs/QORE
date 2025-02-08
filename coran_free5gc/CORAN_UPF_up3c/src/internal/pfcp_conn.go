// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	"net"
	"time"

	mes "github.com/coranlabs/CORAN_GO_PFCP/message"
	"github.com/coranlabs/HEXA_UPF/src/logger"
)

type PfcpFunc func(conn *PfcpConn, msg mes.Message, addr string) (mes.Message, error)

type PfcpHandlerMap map[uint8]PfcpFunc

func CreatePfcpConnection(addr string, nodeid string, n3Ip string) (*PfcpConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		logger.AppLog.Infof("Can't resolve UDP address: %s", err.Error())
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.AppLog.Infof("Can't listen UDP address: %s", err.Error())
		return nil, err
	}

	n3Addr := net.ParseIP(n3Ip)
	if n3Addr == nil {
		logger.AppLog.Infof("failed to parse N3 IP address ID: %s", n3Ip)
		return nil, err
	}
	logger.AppLog.Infof("Starting the HEXA_UPF V0.1")
	logger.AppLog.Infof("Initializing PFCP connection on : %v with Node ID: %v", udpAddr, nodeid)

	return &PfcpConn{
		udpConn:           udpConn,
		nodeId:            nodeid,
		nodeAddrV4:        udpAddr.IP,
		n3Address:         n3Addr,
		RecoveryTimestamp: time.Now(),
	}, nil
}

func (connection *PfcpConn) Run() {
	logger.AppLog.Infof("Starting the Server")
	buf := make([]byte, 1500)
	logger.AppLog.Infof("Server Started")
	for {
		n, addr, err := connection.Receive(buf)
		if err != nil {
			logger.AppLog.Fatalf("Error reading from UDP socket: %s", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		logger.AppLog.Debugf("Received %d bytes from %s", n, addr)
		connection.Handle(buf[:n], addr)
	}
}

func (connection *PfcpConn) Receive(b []byte) (n int, addr *net.UDPAddr, err error) {
	return connection.udpConn.ReadFromUDP(b)
}

func (connection *PfcpConn) Handle(b []byte, addr *net.UDPAddr) {
	err := Handle(connection, b, addr)
	if err != nil {
		logger.AppLog.Infof("Error handling PFCP message: %s", err.Error())
	}
}

func (connection *PfcpConn) Close() {
	connection.udpConn.Close()
}
