// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	"log"
	"net"
	"time"
)

func CreateConn(addr string, nodeid string, n3Ip string) (*PfcpConn, error) {
	udpAddr := "192.168.1.2"
	udpAddrIp, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp", udpAddrIp)
	if err != nil {
		return nil, err
	}
	return &PfcpConn{
		udpConn:           udpConn,
		nodeId:            nodeid,
		nodeAddrV4:        udpAddrIp.IP,
		RecoveryTimestamp: time.Now(),
	}, nil
}

func (connection *PfcpConn) Run() {
	log.Printf("Starting the Server")
	buf := make([]byte, 1500)
	log.Printf("Server Started")
	for {
		n, addr, err := connection.Receive(buf)
		if err != nil {
			log.Printf("Error reading from UDP socket: %s", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		log.Printf("Received %d bytes from %s", n, addr)
		connection.Handle(buf[:n], addr)
	}
}

func (connection *PfcpConn) Receive(b []byte) (n int, addr *net.UDPAddr, err error) {
	return connection.udpConn.ReadFromUDP(b)
}

func (connection *PfcpConn) Handle(b []byte, addr *net.UDPAddr) {
	err := Handle(connection, b, addr)
	if err != nil {
		log.Printf("Error handling PFCP message: %s", err.Error())
	}
}

func (connection *PfcpConn) Close() {}
