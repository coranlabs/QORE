// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package server

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/cilium/ebpf/link"
	"github.com/coranlabs/HEXA_UPF/ebpf"
	"github.com/coranlabs/HEXA_UPF/internal/logger"
	data "github.com/coranlabs/HEXA_UPF/src/config"
	"github.com/coranlabs/HEXA_UPF/src/internal"

	upfapp "github.com/coranlabs/HEXA_UPF/pkg/app"
	"github.com/coranlabs/HEXA_UPF/pkg/factory"

	"github.com/sirupsen/logrus"
)

func runCommand(command string, args []string) {
	// Execute the command
	cmd := exec.Command(command, args...)

	// Capture output or error
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.MainLog.Tracef("Failed to execute %s: %s", command, err)
	}

	fmt.Printf("Output: %s\n", string(output))
}
func Action(conf *factory.Config) error {
	logger.MainLog.Tracef("in action starting")
	cfg := conf
	logger.MainLog.Tracef("after readconfig %v", cfg)
	// if err != nil {

	// 	logger.MainLog.Tracef("error is %v ", err)
	// 	time.Sleep(2000 * time.Second)
	// 	return err
	// }
	// 1. Install iptables using apk
	// fmt.Println("Installing iptables...")
	// runCommand("apk", []string{"add", "iptables"})

	// // 2. Add MASQUERADE rule in NAT table for eth0
	// fmt.Println("Adding NAT MASQUERADE rule...")
	// runCommand("iptables", []string{"-t", "nat", "-A", "POSTROUTING", "-o", "eth0", "-j", "MASQUERADE"})

	// // 3. Allow forwarding of packets
	// fmt.Println("Allowing packet forwarding...")
	// runCommand("iptables", []string{"-I", "FORWARD", "1", "-j", "ACCEPT"})

	// fmt.Println("All commands executed successfully.")

	upf, err := upfapp.NewApp(cfg)
	if err != nil {
		logger.MainLog.Tracef("error in newapp creation: %v", err)
		time.Sleep(30000)
		return err
	}

	if err := upf.Run(); err != nil {
		logger.MainLog.Tracef("error in run : %v", err)
		time.Sleep(30000)
		return err
	}

	return nil
}

func Service() {
	logger.InitializeLogger(logrus.InfoLevel)

	// addr := data.Upfdata{}
	// addr.PfcpAddress = "10.100.200.14:8805"
	// addr.N3Address = "192.168.5.203"
	// addr.InterfaceName = []string{"eth0"}
	logger.MainLog.Tracef("error after this")
	data.Init()
	n3 := data.Conf.N3Address
	ip := data.Conf.PfcpAddress

	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)

	logger.MainLog.Tracef("didnt come here")
	pfcpConn, err := internal.CreatePfcpConnection(ip, ip, n3)
	if err != nil {
		logger.AppLog.Fatalf("Could not create PFCP connection: %s", err.Error())
	}
	go pfcpConn.Run()
	defer pfcpConn.Close()

	if err := ebpf.IncreaseResourceLimits(); err != nil {
		logger.MainLog.Errorf("Can't increase resource limits: %s", err.Error())
	}

	bpfObjects := ebpf.NewBpfObjects()
	if err := bpfObjects.Load(); err != nil {
		logger.MainLog.Errorf("Loading bpf objects failed: %s", err.Error())
	}

	if data.Conf.EbpfMapResize {
		if err := bpfObjects.ResizeAllMaps(data.Conf.QerMapSize, data.Conf.FarMapSize, data.Conf.PdrMapSize); err != nil {
			logger.MainLog.Errorf("Failed to set ebpf map sizes: %s", err)
		}
	}

	defer bpfObjects.Close()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for _, ifaceName := range data.Conf.InterfaceName {
		iface, err := net.InterfaceByName(ifaceName)
		if err != nil {
			logger.MainLog.Errorf("Lookup network iface %q: %s", ifaceName, err.Error())
		}

		// Attach the program.
		l, err := link.AttachXDP(link.XDPOptions{
			Program:   bpfObjects.UpfIpEntrypointFunc,
			Interface: iface.Index,
			Flags:     StringToXDPAttachMode(data.Conf.XDPAttachMode),
		})
		if err != nil {
			logger.MainLog.Errorf("Could not attach XDP program: %s", err.Error())
		}
		defer l.Close()

		logger.InitLog.Infof("Attached XDP program to iface %q (index %d)", iface.Name, iface.Index)
	}

	for {
		select {
		case <-ticker.C:

		case <-stopper:

			logger.AppLog.Infof("Received signal, exiting program..")

			return
		}
	}
}

func StringToXDPAttachMode(Mode string) link.XDPAttachFlags {
	switch Mode {
	case "generic":
		return link.XDPGenericMode
	case "native":
		return link.XDPDriverMode
	case "offload":
		return link.XDPOffloadMode
	default:
		return link.XDPGenericMode
	}
}
