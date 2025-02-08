package ebpf

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"unsafe"

	"github.com/coranlabs/CORAN_UPF_eBPF/logger"
	config "github.com/coranlabs/CORAN_UPF_eBPF/config"
	

	"github.com/cilium/ebpf"
)

// The BPF_ARRAY map type has no delete operation. The only way to delete an element is to replace it with a new one.

type PdrInfo struct {
	OuterHeaderRemoval uint8
	FarId              uint32
	QerId              uint32
	SdfFilter          *SdfFilter
}

type SdfFilter struct {
	Protocol     uint8 // 0: icmp, 1: ip, 2: tcp, 3: udp, 4: icmp6
	SrcAddress   IpWMask
	SrcPortRange PortRange
	DstAddress   IpWMask
	DstPortRange PortRange
}

type IpWMask struct {
	Type uint8 // 0: any, 1: ip4, 2: ip6
	Ip   net.IP
	Mask net.IPMask
}

type EReferencePoint int32

const (
	N3_INTERFACE  EReferencePoint = iota // 0
	N6_INTERFACE                         // 1
	N4_INTERFACE                         // 2
	N9_INTERFACE                         // 3
	N19_INTERFACE                        // 4
)

type PortRange struct {
	LowerBound uint16
	UpperBound uint16
}

func PreprocessPdrWithSdf(lookup func(interface{}, interface{}) error, key interface{}, pdrInfo PdrInfo) (IpEntrypointPdrInfo, error) {
	var defaultPdr IpEntrypointPdrInfo
	if err := lookup(key, &defaultPdr); err != nil {
		return CombinePdrWithSdf(nil, pdrInfo), nil
	}

	return CombinePdrWithSdf(&defaultPdr, pdrInfo), nil
}

func (bpfObjects *BpfObjects) PutPdrUplink(teid uint32, pdrInfo PdrInfo) error {
	logger.Pfcplog.Infof("EBPF: Put PDR Uplink: teid=%d, pdrInfo=%+v", teid, pdrInfo)
	var pdrToStore IpEntrypointPdrInfo
	var err error
	if pdrInfo.SdfFilter != nil {
		if pdrToStore, err = PreprocessPdrWithSdf(bpfObjects.PdrMapUplinkIp4.Lookup, teid, pdrInfo); err != nil {
			return err
		}
	} else {
		pdrToStore = ToIpEntrypointPdrInfo(pdrInfo)
	}
	return bpfObjects.PdrMapUplinkIp4.Put(teid, unsafe.Pointer(&pdrToStore))
}

func (bpfObjects *BpfObjects) PutPdrDownlink(ipv4 net.IP, pdrInfo PdrInfo) error {
	logger.Pfcplog.Infof("EBPF: Put PDR Downlink: ipv4=%s, pdrInfo=%+v", ipv4, pdrInfo)
	var pdrToStore IpEntrypointPdrInfo
	var err error
	if pdrInfo.SdfFilter != nil {
		if pdrToStore, err = PreprocessPdrWithSdf(bpfObjects.PdrMapDownlinkIp4.Lookup, ipv4, pdrInfo); err != nil {
			return err
		}
	} else {
		pdrToStore = ToIpEntrypointPdrInfo(pdrInfo)
	}
	return bpfObjects.PdrMapDownlinkIp4.Put(ipv4, unsafe.Pointer(&pdrToStore))
}

func (bpfObjects *BpfObjects) UpdatePdrUplink(teid uint32, pdrInfo PdrInfo) error {
	logger.Pfcplog.Infof("EBPF: Update PDR Uplink: teid=%d, pdrInfo=%+v", teid, pdrInfo)
	var pdrToStore IpEntrypointPdrInfo
	var err error
	if pdrInfo.SdfFilter != nil {
		if pdrToStore, err = PreprocessPdrWithSdf(bpfObjects.PdrMapUplinkIp4.Lookup, teid, pdrInfo); err != nil {
			return err
		}
	} else {
		pdrToStore = ToIpEntrypointPdrInfo(pdrInfo)
	}
	return bpfObjects.PdrMapUplinkIp4.Update(teid, unsafe.Pointer(&pdrToStore), ebpf.UpdateExist)
}

func parseMAC(macStr string) ([6]byte, error) {
	var mac [6]byte
	macStr = strings.ReplaceAll(macStr, ":", "")
	if len(macStr) != 12 {
		return mac, fmt.Errorf("invalid MAC address length")
	}

	for i := 0; i < 6; i++ {
		var b byte
		fmt.Sscanf(macStr[2*i:2*i+2], "%02x", &b)
		mac[i] = b
	}
	return mac, nil
}
func get_DL_MACAddress(ipAddr string) ([6]byte, error) {
	var mac [6]byte

	// Execute the command to get the MAC address
	cmd := exec.Command("sh", "-c", fmt.Sprintf("arping -I n3 -c 1 %s | awk -F'[][]' '/Unicast reply from/ {print $2}' | head -n 1", ipAddr))
	//cmd := exec.Command("sh", "-c", fmt.Sprintf("ip neigh show %s | awk '/lladdr/ {print $5; exit}'", ipAddr))

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return mac, fmt.Errorf("error running command: %v", err)
	}

	// Parse the output to extract the MAC address
	macStr := strings.TrimSpace(out.String())
	if macStr == "" {
		return mac, fmt.Errorf("MAC address not found")
	}

	// Convert the MAC string to [6]byte format
	return parseMAC(macStr)
}
func get_UL_MACAddress(ipAddr string) ([6]byte, error) {
	var mac [6]byte

	// Execute the command to get the MAC address
	cmd := exec.Command("sh", "-c", fmt.Sprintf("arping -I n6 -c 1 %s | awk -F'[][]' '/Unicast reply from/ {print $2}' | head -n 1", ipAddr))
	//cmd := exec.Command("sh", "-c", fmt.Sprintf("ip neigh show %s | awk '/lladdr/ {print $5; exit}'", ipAddr))

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return mac, fmt.Errorf("error running command: %v", err)
	}

	// Parse the output to extract the MAC address
	macStr := strings.TrimSpace(out.String())
	if macStr == "" {
		return mac, fmt.Errorf("MAC address not found")
	}

	// Convert the MAC string to [6]byte format
	return parseMAC(macStr)
}
func (bpfObjects *BpfObjects) Update_m_arpUplink() error {
	gatewayip:= config.Conf.Gatewayip
	// if err1 != nil {
	// 	fmt.Printf("Error fetching uplink config: %v\n", err1)

	// }
	ipAddr := gatewayip

	// Get the MAC address
	mac, err := get_UL_MACAddress(ipAddr)
	if err != nil {
		fmt.Printf("Error fetching uplink MAC address: %v\n for ip : %v", err, ipAddr)

	}
	return bpfObjects.M_arpTable.Put(N6_INTERFACE, unsafe.Pointer(&mac))
}

func (bpfObjects *BpfObjects) Update_m_arpDownlink(gnbip string) error {
	//cfg, err1 := factory.ReadConfig("./config/upfcfg.yaml")
	if gnbip == "" {
		fmt.Printf("Error fetching downlink config address: \n")

	}

	ipAddr := gnbip

	// Get the MAC address
	mac, err := get_DL_MACAddress(ipAddr)
	if err != nil {
		fmt.Printf("Error fetching downlink MAC address: %v\n for ip %v", err, ipAddr)

	}
	return bpfObjects.M_arpTable.Put(N3_INTERFACE, unsafe.Pointer(&mac))
}
func (bpfObjects *BpfObjects) UpdatePdrDownlink(ipv4 net.IP, pdrInfo PdrInfo) error {
	logger.Pfcplog.Infof("EBPF: Update PDR Downlink: ipv4=%s, pdrInfo=%+v", ipv4, pdrInfo)
	var pdrToStore IpEntrypointPdrInfo
	var err error
	if pdrInfo.SdfFilter != nil {
		if pdrToStore, err = PreprocessPdrWithSdf(bpfObjects.PdrMapDownlinkIp4.Lookup, ipv4, pdrInfo); err != nil {
			return err
		}
	} else {
		pdrToStore = ToIpEntrypointPdrInfo(pdrInfo)
	}
	return bpfObjects.PdrMapDownlinkIp4.Update(ipv4, unsafe.Pointer(&pdrToStore), ebpf.UpdateExist)
}

func (bpfObjects *BpfObjects) DeletePdrUplink(teid uint32) error {
	logger.Pfcplog.Infof("EBPF: Delete PDR Uplink: teid=%d", teid)
	return bpfObjects.PdrMapUplinkIp4.Delete(teid)
}

func (bpfObjects *BpfObjects) DeletePdrDownlink(ipv4 net.IP) error {
	logger.Pfcplog.Infof("EBPF: Delete PDR Downlink: ipv4=%s", ipv4)
	return bpfObjects.PdrMapDownlinkIp4.Delete(ipv4)
}

func (bpfObjects *BpfObjects) PutDownlinkPdrIp6(ipv6 net.IP, pdrInfo PdrInfo) error {
	logger.Pfcplog.Infof("EBPF: Put PDR Ipv6 Downlink: ipv6=%s, pdrInfo=%+v", ipv6, pdrInfo)
	var pdrToStore IpEntrypointPdrInfo
	var err error
	if pdrInfo.SdfFilter != nil {
		if pdrToStore, err = PreprocessPdrWithSdf(bpfObjects.PdrMapDownlinkIp6.Lookup, ipv6, pdrInfo); err != nil {
			return err
		}
	} else {
		pdrToStore = ToIpEntrypointPdrInfo(pdrInfo)
	}
	return bpfObjects.PdrMapDownlinkIp6.Put(ipv6, unsafe.Pointer(&pdrToStore))
}

func (bpfObjects *BpfObjects) UpdateDownlinkPdrIp6(ipv6 net.IP, pdrInfo PdrInfo) error {
	logger.Pfcplog.Infof("EBPF: Update PDR Ipv6 Downlink: ipv6=%s, pdrInfo=%+v", ipv6, pdrInfo)
	var pdrToStore IpEntrypointPdrInfo
	var err error
	if pdrInfo.SdfFilter != nil {
		if pdrToStore, err = PreprocessPdrWithSdf(bpfObjects.PdrMapDownlinkIp6.Lookup, ipv6, pdrInfo); err != nil {
			return err
		}
	} else {
		pdrToStore = ToIpEntrypointPdrInfo(pdrInfo)
	}
	return bpfObjects.PdrMapDownlinkIp6.Update(ipv6, unsafe.Pointer(&pdrToStore), ebpf.UpdateExist)
}

func (bpfObjects *BpfObjects) DeleteDownlinkPdrIp6(ipv6 net.IP) error {
	logger.Pfcplog.Infof("EBPF: Delete PDR Ipv6 Downlink: ipv6=%s", ipv6)
	return bpfObjects.PdrMapDownlinkIp6.Delete(ipv6)
}

type FarInfo struct {
	Action                uint8
	OuterHeaderCreation   uint8
	Teid                  uint32
	RemoteIP              uint32
	LocalIP               uint32
	TransportLevelMarking uint16
}

func (f FarInfo) MarshalJSON() ([]byte, error) {
	remoteIP := make(net.IP, 4)
	localIP := make(net.IP, 4)
	binary.LittleEndian.PutUint32(remoteIP, f.RemoteIP)
	binary.LittleEndian.PutUint32(localIP, f.LocalIP)
	data := map[string]interface{}{
		"action":                  f.Action,
		"outer_header_creation":   f.OuterHeaderCreation,
		"teid":                    f.Teid,
		"remote_ip":               remoteIP.String(),
		"local_ip":                localIP.String(),
		"transport_level_marking": f.TransportLevelMarking,
	}
	return json.Marshal(data)
}

func (bpfObjects *BpfObjects) NewFar(farid uint32, farInfo FarInfo) (uint32, error) {
	//internalId, err := bpfObjects.FarIdTracker.GetNext()
	internalId := farid
	// if err != nil {
	// 	return 0, err
	// }
	logger.Pfcplog.Infof("EBPF: Put FAR: internalId=%d, qerInfo=%+v", internalId, farInfo)
	return internalId, bpfObjects.FarMap.Put(internalId, unsafe.Pointer(&farInfo))
}

func (bpfObjects *BpfObjects) UpdateFar(internalId uint32, farInfo FarInfo) error {
	logger.Pfcplog.Infof("EBPF: Update FAR: internalId=%d, farInfo=%+v", internalId, farInfo)
	return bpfObjects.FarMap.Update(internalId, unsafe.Pointer(&farInfo), ebpf.UpdateExist)
}

func (bpfObjects *BpfObjects) DeleteFar(intenalId uint32) error {
	logger.Pfcplog.Infof("EBPF: Delete FAR: intenalId=%d", intenalId)
	bpfObjects.FarIdTracker.Release(intenalId)
	return bpfObjects.FarMap.Update(intenalId, unsafe.Pointer(&FarInfo{}), ebpf.UpdateExist)
}

type QerInfo struct {
	GateStatusUL uint8
	GateStatusDL uint8
	Qfi          uint8
	MaxBitrateUL uint32
	MaxBitrateDL uint32
	StartUL      uint64
	StartDL      uint64
}

func (bpfObjects *BpfObjects) NewQer(qerInfo QerInfo) (uint32, error) {
	internalId, err := bpfObjects.QerIdTracker.GetNext()
	if err != nil {
		return 0, err
	}
	logger.Pfcplog.Infof("EBPF: Put QER: internalId=%d, qerInfo=%+v", internalId, qerInfo)
	return internalId, bpfObjects.QerMap.Put(internalId, unsafe.Pointer(&qerInfo))
}

func (bpfObjects *BpfObjects) UpdateQer(internalId uint32, qerInfo QerInfo) error {
	logger.Pfcplog.Infof("EBPF: Update QER: internalId=%d, qerInfo=%+v", internalId, qerInfo)
	return bpfObjects.QerMap.Update(internalId, unsafe.Pointer(&qerInfo), ebpf.UpdateExist)
}

func (bpfObjects *BpfObjects) DeleteQer(internalId uint32) error {
	logger.Pfcplog.Infof("EBPF: Delete QER: internalId=%d", internalId)
	bpfObjects.QerIdTracker.Release(internalId)
	return bpfObjects.QerMap.Update(internalId, unsafe.Pointer(&QerInfo{}), ebpf.UpdateExist)
}

type ForwardingPlaneController interface {
	Update_m_arpUplink() error
	Update_m_arpDownlink(gnbip string) error
	PutPdrUplink(teid uint32, pdrInfo PdrInfo) error
	PutPdrDownlink(ipv4 net.IP, pdrInfo PdrInfo) error
	UpdatePdrUplink(teid uint32, pdrInfo PdrInfo) error
	UpdatePdrDownlink(ipv4 net.IP, pdrInfo PdrInfo) error
	DeletePdrUplink(teid uint32) error
	DeletePdrDownlink(ipv4 net.IP) error
	PutDownlinkPdrIp6(ipv6 net.IP, pdrInfo PdrInfo) error
	UpdateDownlinkPdrIp6(ipv6 net.IP, pdrInfo PdrInfo) error
	DeleteDownlinkPdrIp6(ipv6 net.IP) error
	NewFar(farid uint32, farInfo FarInfo) (uint32, error)
	UpdateFar(internalId uint32, farInfo FarInfo) error
	DeleteFar(internalId uint32) error
	NewQer(qerInfo QerInfo) (uint32, error)
	UpdateQer(internalId uint32, qerInfo QerInfo) error
	DeleteQer(internalId uint32) error
}

func CombinePdrWithSdf(defaultPdr *IpEntrypointPdrInfo, sdfPdr PdrInfo) IpEntrypointPdrInfo {
	var pdrToStore IpEntrypointPdrInfo
	// Default mapping options.
	if defaultPdr != nil {
		pdrToStore.OuterHeaderRemoval = defaultPdr.OuterHeaderRemoval
		pdrToStore.FarId = defaultPdr.FarId
		pdrToStore.QerId = defaultPdr.QerId
		pdrToStore.SdfMode = 2
	} else {
		pdrToStore.SdfMode = 1
	}

	// SDF mapping options.
	pdrToStore.SdfRules.SdfFilter.Protocol = sdfPdr.SdfFilter.Protocol
	pdrToStore.SdfRules.SdfFilter.SrcAddr.Type = sdfPdr.SdfFilter.SrcAddress.Type
	pdrToStore.SdfRules.SdfFilter.SrcAddr.Ip = Copy16Ip(sdfPdr.SdfFilter.SrcAddress.Ip)
	pdrToStore.SdfRules.SdfFilter.SrcAddr.Mask = Copy16Ip(sdfPdr.SdfFilter.SrcAddress.Mask)
	pdrToStore.SdfRules.SdfFilter.SrcPort.LowerBound = sdfPdr.SdfFilter.SrcPortRange.LowerBound
	pdrToStore.SdfRules.SdfFilter.SrcPort.UpperBound = sdfPdr.SdfFilter.SrcPortRange.UpperBound
	pdrToStore.SdfRules.SdfFilter.DstAddr.Type = sdfPdr.SdfFilter.DstAddress.Type
	pdrToStore.SdfRules.SdfFilter.DstAddr.Ip = Copy16Ip(sdfPdr.SdfFilter.DstAddress.Ip)
	pdrToStore.SdfRules.SdfFilter.DstAddr.Mask = Copy16Ip(sdfPdr.SdfFilter.DstAddress.Mask)
	pdrToStore.SdfRules.SdfFilter.DstPort.LowerBound = sdfPdr.SdfFilter.DstPortRange.LowerBound
	pdrToStore.SdfRules.SdfFilter.DstPort.UpperBound = sdfPdr.SdfFilter.DstPortRange.UpperBound
	pdrToStore.SdfRules.OuterHeaderRemoval = sdfPdr.OuterHeaderRemoval
	pdrToStore.SdfRules.FarId = sdfPdr.FarId
	pdrToStore.SdfRules.QerId = sdfPdr.QerId
	return pdrToStore
}

func ToIpEntrypointPdrInfo(defaultPdr PdrInfo) IpEntrypointPdrInfo {
	var pdrToStore IpEntrypointPdrInfo
	pdrToStore.OuterHeaderRemoval = defaultPdr.OuterHeaderRemoval
	pdrToStore.FarId = defaultPdr.FarId
	pdrToStore.QerId = defaultPdr.QerId
	return pdrToStore
}

func Copy16Ip[T ~[]byte](arr T) [16]byte {
	const Ipv4len = 4
	const Ipv6len = 16
	var c [Ipv6len]byte
	var arrLen int
	if len(arr) == Ipv4len {
		arrLen = Ipv4len
	} else if len(arr) == Ipv6len {
		arrLen = Ipv6len
	} else if len(arr) == 0 || arr == nil {
		return c
	}
	for i := 0; i < arrLen; i++ {
		c[i] = (arr)[arrLen-1-i]
	}
	return c
}

func (sdfFilter *SdfFilter) String() string {
	return fmt.Sprintf("%+v", *sdfFilter)
}
