package ebpf_datapath

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"unsafe"

	config "github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/config"
	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/logger"

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

// func PreprocessPdrWithSdf(lookup func(interface{}, interface{}) error, key interface{}, pdrInfo PdrInfo) (Coran_ebpf_datapathPDR, error) {
// 	var defaultPdr Coran_ebpf_datapathPDR
// 	if err := lookup(key, &defaultPdr); err != nil {
// 		return CombinePdrWithSdf(nil, pdrInfo), nil
// 	}

// 	return CombinePdrWithSdf(&defaultPdr, pdrInfo), nil
// }

func (EBPF_controller *EBPF_entity) PutPdrUplink(teid uint32, pdrInfo PdrInfo) error {
	logger.EBPF_Datapath.Infof("EBPF: Put PDR Uplink: teid=%d, pdrInfo=%+v", teid, pdrInfo)
	var pdrToStore Coran_ebpf_datapathPDR
	
	if pdrInfo.SdfFilter != nil {
		return fmt.Errorf("Can't apply SDF PDR")
	} else {
		pdrToStore = ToCoran_ebpf_datapathPDR(pdrInfo)
	}
	return EBPF_controller.PDR_uplinkMap.Put(teid, unsafe.Pointer(&pdrToStore))
}

func (EBPF_controller *EBPF_entity) PutPdrDownlink(ipv4 net.IP, pdrInfo PdrInfo) error {
	logger.EBPF_Datapath.Infof("EBPF: Put PDR Downlink: ipv4=%s, pdrInfo=%+v", ipv4, pdrInfo)
	var pdrToStore Coran_ebpf_datapathPDR
	if pdrInfo.SdfFilter != nil {
		return fmt.Errorf("Can't apply SDF PDR")

		
	} else {
		pdrToStore = ToCoran_ebpf_datapathPDR(pdrInfo)
	}
	return EBPF_controller.PDR_downlinkMap.Put(ipv4, unsafe.Pointer(&pdrToStore))
}

func (EBPF_controller *EBPF_entity) UpdatePdrUplink(teid uint32, pdrInfo PdrInfo) error {
	logger.EBPF_Datapath.Infof("EBPF: Update PDR Uplink: teid=%d, pdrInfo=%+v", teid, pdrInfo)
	var pdrToStore Coran_ebpf_datapathPDR
	if pdrInfo.SdfFilter != nil {
		return fmt.Errorf("Can't apply SDF PDR")

	} else {
		pdrToStore = ToCoran_ebpf_datapathPDR(pdrInfo)
	}
	return EBPF_controller.PDR_uplinkMap.Update(teid, unsafe.Pointer(&pdrToStore), ebpf.UpdateExist)
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
func (EBPF_controller *EBPF_entity) Update_m_arpUplink() error {
	gatewayip := config.Conf.Gatewayip
	// if err1 != nil {
	// 	fmt.Printf("Error fetching uplink config: %v\n", err1)

	// }
	ipAddr := gatewayip

	// Get the MAC address
	mac, err := get_UL_MACAddress(ipAddr)
	if err != nil {
		fmt.Printf("Error fetching uplink MAC address: %v\n for ip : %v", err, ipAddr)

	}
	return EBPF_controller.ArpMap.Put(N6_INTERFACE, unsafe.Pointer(&mac))
}

func (EBPF_controller *EBPF_entity) Update_m_arpDownlink(gnbip string) error {
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
	return EBPF_controller.ArpMap.Put(N3_INTERFACE, unsafe.Pointer(&mac))
}
func (EBPF_controller *EBPF_entity) UpdatePdrDownlink(ipv4 net.IP, pdrInfo PdrInfo) error {
	logger.EBPF_Datapath.Infof("EBPF: Update PDR Downlink: ipv4=%s, pdrInfo=%+v", ipv4, pdrInfo)
	var pdrToStore Coran_ebpf_datapathPDR
	if pdrInfo.SdfFilter != nil {
		return fmt.Errorf("Can't apply SDF PDR")

	} else {
		pdrToStore = ToCoran_ebpf_datapathPDR(pdrInfo)
	}
	return EBPF_controller.PDR_downlinkMap.Update(ipv4, unsafe.Pointer(&pdrToStore), ebpf.UpdateExist)
}

func (EBPF_controller *EBPF_entity) DeletePdrUplink(teid uint32) error {
	logger.EBPF_Datapath.Infof("EBPF: Delete PDR Uplink: teid=%d", teid)
	return EBPF_controller.PDR_uplinkMap.Delete(teid)
}

func (EBPF_controller *EBPF_entity) DeletePdrDownlink(ipv4 net.IP) error {
	logger.EBPF_Datapath.Infof("EBPF: Delete PDR Downlink: ipv4=%s", ipv4)
	return EBPF_controller.PDR_downlinkMap.Delete(ipv4)
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

func (EBPF_controller *EBPF_entity) NewFar(farid uint32, farInfo FarInfo) (uint32, error) {
	//internalId, err := EBPF_controller.FarIdTracker.GetNext()
	internalId := farid
	// if err != nil {
	// 	return 0, err
	// }
	logger.EBPF_Datapath.Infof("EBPF: Put FAR: internalId=%d, qerInfo=%+v", internalId, farInfo)
	return internalId, EBPF_controller.FAR_map.Put(internalId, unsafe.Pointer(&farInfo))
}

func (EBPF_controller *EBPF_entity) UpdateFar(internalId uint32, farInfo FarInfo) error {
	logger.EBPF_Datapath.Infof("EBPF: Update FAR: internalId=%d, farInfo=%+v", internalId, farInfo)
	return EBPF_controller.FAR_map.Update(internalId, unsafe.Pointer(&farInfo), ebpf.UpdateExist)
}

func (EBPF_controller *EBPF_entity) DeleteFar(intenalId uint32) error {
	logger.EBPF_Datapath.Infof("EBPF: Delete FAR: intenalId=%d", intenalId)
	EBPF_controller.FarIdTracker.Release(intenalId)
	return EBPF_controller.FAR_map.Update(intenalId, unsafe.Pointer(&FarInfo{}), ebpf.UpdateExist)
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



type EBPFMapInterface interface {
	Update_m_arpUplink() error
	Update_m_arpDownlink(gnbip string) error
	PutPdrUplink(teid uint32, pdrInfo PdrInfo) error
	PutPdrDownlink(ipv4 net.IP, pdrInfo PdrInfo) error
	UpdatePdrUplink(teid uint32, pdrInfo PdrInfo) error
	UpdatePdrDownlink(ipv4 net.IP, pdrInfo PdrInfo) error
	DeletePdrUplink(teid uint32) error
	DeletePdrDownlink(ipv4 net.IP) error
	NewFar(farid uint32, farInfo FarInfo) (uint32, error)
	UpdateFar(internalId uint32, farInfo FarInfo) error
	DeleteFar(internalId uint32) error
}

// func CombinePdrWithSdf(defaultPdr *Coran_ebpf_datapathPDR, sdfPdr PdrInfo) Coran_ebpf_datapathPDR {
// 	var pdrToStore Coran_ebpf_datapathPDR
// 	// Default mapping options.
// 	if defaultPdr != nil {
// 		pdrToStore.OHR = defaultPdr.OHR
// 		pdrToStore.FarId = defaultPdr.FarId
// 		pdrToStore.QerId = defaultPdr.QerId
// 	}

// 	// SDF mapping options.

// 	return pdrToStore
// }

func ToCoran_ebpf_datapathPDR(defaultPdr PdrInfo) Coran_ebpf_datapathPDR {
	var pdrToStore Coran_ebpf_datapathPDR
	pdrToStore.OHR = defaultPdr.OuterHeaderRemoval
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
