package PFCP_server

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"

	UPF_config "github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/config"
	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/service"

	"github.com/pkg/errors"

	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/logger"
	ebpf_datapath "github.com/coranlabs/CORAN_UPF_eBPF/eBPF_Datapath_entity"

	"github.com/coranlabs/CORAN_GO_PFCP/message"
)

// type Pfcp_Message_function func(conn *Pfcp_Link, msg message.Message, addr string) (message.Message, error)

// type Pfcp_Messages_handler map[uint8]Pfcp_Message_function
type ReceivePacket struct {
	RemoteAddr net.Addr
	Buf        []byte
}

const (
	RECEIVE_CHANNEL_LEN       = 512
	REPORT_CHANNEL_LEN        = 128
	TRANS_TIMEOUT_CHANNEL_LEN = 64
	MAX_PFCP_MSG_LEN          = 65536
)

type TransType int

const (
	TX TransType = iota
	RX
)

type TransactionTimeout struct {
	TrType TransType
	TrID   string
}

type Pfcp_Link struct {
	UdpConn           *net.UDPConn
	NodeId            string
	NodeAddrV4        net.IP
	N3Address         net.IP
	rxTrans           map[string]*RxTransaction // key: RemoteAddr-Sequence
	txTrans           map[string]*TxTransaction // key: RemoteAddr-Sequence
	EBPFMapManager    ebpf_datapath.EBPFMapInterface
	RecoveryTimestamp time.Time
	rcvCh             chan ReceivePacket
	NodeAssociations  map[string]*NodeAssociation
	trToCh            chan TransactionTimeout
	ResourceManager   *service.ResourceManager
	cfg               *UPF_config.UpfConfig
}

func Setup_PFCP_server(config *UPF_config.UpfConfig, EBPFMapManager ebpf_datapath.EBPFMapInterface, resourceManager *service.ResourceManager) (*Pfcp_Link, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", config.PfcpAddress)
	if err != nil {
		logger.PFCPSessLog.Warnf("Can't resolve UDP address: %s", err.Error())
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.PFCPSessLog.Warnf("Can't listen UDP address: %s", err.Error())
		return nil, err
	}

	n3Addr := net.ParseIP(config.N3Address)
	logger.InitLog.Infof("Starting n3 IP: %v ", config.N3Address)
	if n3Addr == nil {
		return nil, fmt.Errorf("failed to parse N3 IP address ID: %s", config.N3Address)
	}
	logger.InitLog.Infof("Starting PFCP connection: %v with Node ID: %v and N3 address: %v", udpAddr, config.PfcpNodeId, n3Addr)

	PFCP_controller := &Pfcp_Link{
		UdpConn:           udpConn,
		NodeId:            config.PfcpNodeId,
		NodeAddrV4:        udpAddr.IP,
		N3Address:         net.IP(config.N3Address),
		EBPFMapManager:    EBPFMapManager,
		rcvCh:             make(chan ReceivePacket, RECEIVE_CHANNEL_LEN),
		RecoveryTimestamp: time.Now(),
		cfg:               config,
		rxTrans:           make(map[string]*RxTransaction),
		txTrans:           make(map[string]*TxTransaction),
		NodeAssociations:  map[string]*NodeAssociation{},
		ResourceManager:   resourceManager,
	}

	go PFCP_controller.Run()

	return PFCP_controller, nil
}

// func hexDump(data []byte) {
// 	fmt.Printf("[%s] [DEBUG] [ SCTP ] Message received: Length: %d\n",
// 		time.Now().Format("2006-01-02T15:04:05"), len(data))
// 	fmt.Println("[DEBUG] [ SCTP ] Message Bytes (Hex):")

// 	for i := 0; i < len(data); i += 16 {
// 		fmt.Printf("%04X: ", i) // Print memory offset
// 		for j := 0; j < 16 && i+j < len(data); j++ {
// 			fmt.Printf("%02X ", data[i+j]) // Print hex values
// 		}
// 		fmt.Println()
// 	}
// }

func structuredHexDump(data []byte) {
	fmt.Printf("[%s] Message received: Length: %d\n",
		time.Now().Format("2006-01-02T15:04:05"), len(data))
	fmt.Println("Message Bytes (Hex):")

	for i := 0; i < len(data); i += 16 {
		fmt.Printf("%04X: ", i) // Print memory offset

		// Print data in a 4x4 block format
		for j := 0; j < 16 && i+j < len(data); j++ {
			fmt.Printf("%02X ", data[i+j])
			if (j+1)%4 == 0 { // Add spacing every 4 bytes
				fmt.Print(" ")
			}
		}

		fmt.Println() // Newline after each row
	}
}
func (s *Pfcp_Link) Run() {
	// go func() {
	// 	for {
	// 		connection.RefreshAssociations()
	// 		time.Sleep(time.Duration(config.Conf.HeartbeatInterval) * time.Second)
	// 	}
	// }()
	go s.receiver()

	for {
		select {

		case rcvPkt := <-s.rcvCh:
			logger.InitLog.Infof("receive buf(len=%d) from rcvCh", len(rcvPkt.Buf))
			logger.CtxLog.Infof("received buffer: ")
			structuredHexDump(rcvPkt.Buf)
			if len(rcvPkt.Buf) == 0 {
				// receiver closed
				return
			}
			msg, err := message.Parse(rcvPkt.Buf)
			if err != nil {
				logger.InitLog.Errorf("parse error %v", err)
				logger.InitLog.Errorf("ignored undecodable message:\n%+v", hex.Dump(rcvPkt.Buf))
				continue
			}

			trID := fmt.Sprintf("%s-%d", rcvPkt.RemoteAddr, msg.Sequence())
			if isRequest(msg) {
				logger.InitLog.Debugf("receive req pkt from %s", trID)
				rx, ok := s.rxTrans[trID]
				if !ok {
					rx = NewRxTransaction(s, rcvPkt.RemoteAddr, msg.Sequence())
					s.rxTrans[trID] = rx // nil pointer dereference
				}
				needDispatch, err1 := rx.recv(msg, ok)
				if err1 != nil {
					logger.InitLog.Debugf("rcvCh: %v", err1)
					continue
				} else if !needDispatch {
					logger.InitLog.Debugf("rcvCh: rxtr[%s] req no need to dispatch", trID)
					continue
				}
				err := s.Handle_request(msg, rcvPkt.RemoteAddr)
				if err != nil {
					logger.InitLog.Errorf("%v", err)
					logger.InitLog.Errorf("ignored undecodable message:\n%+v", hex.Dump(rcvPkt.Buf))
				}
			} else if isResponse(msg) {
				logger.InitLog.Errorf("receive rsp pkt from %s", trID)
				tx, ok := s.txTrans[trID]
				if !ok {
					logger.InitLog.Debugf("rcvCh: No txtr[%s] found for rsp", trID)
					continue
				}
				req := tx.recv(msg)
				err = s.Handle_response(msg, rcvPkt.RemoteAddr, req)
				if err != nil {
					logger.InitLog.Errorf("%v", err)
					logger.InitLog.Errorf("ignored undecodable message:\n%+v", hex.Dump(rcvPkt.Buf))
				}
			}
		case trTo := <-s.trToCh:
			logger.InitLog.Errorf("receive tr timeout (%v) from trToCh", trTo)
			if trTo.TrType == TX {
				tx, ok := s.txTrans[trTo.TrID]
				if !ok {
					logger.InitLog.Debugf("trToCh: txtr[%s] not found", trTo.TrID)
					continue
				}
				logger.InitLog.Debugf("trToCh: txtr[%v] timeout", tx)
				//tx.handleTimeout()
			} else { // RX
				rx, ok := s.rxTrans[trTo.TrID]
				if !ok {
					logger.InitLog.Debugf("trToCh: rxtr[%s] not found", trTo.TrID)
					continue
				}
				rx.handleTimeout()
			}
		}
	}
	// buf := make([]byte, 1500)
	// for {
	// 	n, addr, err := connection.Receive(buf)
	// 	if err != nil {
	// 		logger.PFCPSessLog.Warnf("Error reading from UDP socket: %s", err.Error())
	// 		time.Sleep(1 * time.Second)
	// 		continue
	// 	}
	// 	logger.InitLog.Debugf("Received %d bytes from %s", n, addr)
	// 	connection.Handle(buf[:n], addr)
	// }
}
func (s *Pfcp_Link) NotifyTransTimeout(trType TransType, trID string) {
	s.trToCh <- TransactionTimeout{TrType: trType, TrID: trID}
}
func (s *Pfcp_Link) sendRspTo(msg message.Message, addr net.Addr) error {
	if !isResponse(msg) {
		return errors.Errorf("sendRspTo: invalid rsp type(%d)", msg.MessageType())
	}

	// find transaction
	trID := fmt.Sprintf("%s-%d", addr, msg.Sequence())
	rxtr, ok := s.rxTrans[trID]
	if !ok {
		return errors.Errorf("sendRspTo: rxtr(%s) not found", trID)
	}

	return rxtr.send(msg)
}
func isRequest(msg message.Message) bool {
	switch msg.MessageType() {
	case message.MsgTypeHeartbeatRequest:
		return true
	case message.MsgTypePFDManagementRequest:
		return true
	case message.MsgTypeAssociationSetupRequest:
		return true
	case message.MsgTypeAssociationUpdateRequest:
		return true
	case message.MsgTypeAssociationReleaseRequest:
		return true
	case message.MsgTypeNodeReportRequest:
		return true
	case message.MsgTypeSessionSetDeletionRequest:
		return true
	case message.MsgTypeSessionEstablishmentRequest:
		return true
	case message.MsgTypeSessionModificationRequest:
		return true
	case message.MsgTypeSessionDeletionRequest:
		return true
	case message.MsgTypeSessionReportRequest:
		return true
	default:
	}
	return false
}

func isResponse(msg message.Message) bool {
	switch msg.MessageType() {
	case message.MsgTypeHeartbeatResponse:
		return true
	case message.MsgTypePFDManagementResponse:
		return true
	case message.MsgTypeAssociationSetupResponse:
		return true
	case message.MsgTypeAssociationUpdateResponse:
		return true
	case message.MsgTypeAssociationReleaseResponse:
		return true
	case message.MsgTypeNodeReportResponse:
		return true
	case message.MsgTypeSessionSetDeletionResponse:
		return true
	case message.MsgTypeSessionEstablishmentResponse:
		return true
	case message.MsgTypeSessionModificationResponse:
		return true
	case message.MsgTypeSessionDeletionResponse:
		return true
	case message.MsgTypeSessionReportResponse:
		return true
	default:
	}
	return false
}
func (s *Pfcp_Link) receiver() {

	buf := make([]byte, MAX_PFCP_MSG_LEN)
	for {
		logger.InitLog.Debugf("receiver starts to read...")
		n, addr, err := s.UdpConn.ReadFrom(buf)
		if err != nil {
			logger.InitLog.Errorf("%+v", err)
			s.rcvCh <- ReceivePacket{}
			break
		}

		logger.InitLog.Debugf("receiver reads message(len=%d)", n)
		msgBuf := make([]byte, n)
		copy(msgBuf, buf)
		s.rcvCh <- ReceivePacket{
			RemoteAddr: addr,
			Buf:        msgBuf,
		}
	}
}
func (connection *Pfcp_Link) Close() {
	connection.UdpConn.Close()
}

func (connection *Pfcp_Link) Receive(b []byte) (n int, addr *net.UDPAddr, err error) {
	return connection.UdpConn.ReadFromUDP(b)
}

func (connection *Pfcp_Link) Send(b []byte, addr *net.UDPAddr) (int, error) {
	return connection.UdpConn.WriteTo(b, addr)
}

func (connection *Pfcp_Link) SendMessage(msg message.Message, addr *net.UDPAddr) error {
	responseBytes := make([]byte, msg.MarshalLen())
	if err := msg.MarshalTo(responseBytes); err != nil {
		logger.PFCPSessLog.Warn(err.Error())
		return err
	}
	if _, err := connection.Send(responseBytes, addr); err != nil {
		logger.PFCPSessLog.Warn(err.Error())
		return err
	}
	return nil
}

func (connection *Pfcp_Link) RefreshAssociations() {
	for _, assoc := range connection.NodeAssociations {
		if !assoc.HeartbeatsActive {
			go assoc.ScheduleHeartbeat(connection)
		}
	}
}

// DeleteAssociation deletes an association and all sessions associated with it.
func (connection *Pfcp_Link) DeleteAssociation(assocAddr string) {
	assoc := connection.GetAssociation(assocAddr)
	logger.InitLog.Infof("Pruning expired node association: %s", assocAddr)
	for sessionId, session := range assoc.Sessions {
		logger.InitLog.Infof("Deleting session: %d", sessionId)
		connection.DeleteSession(session)
	}
	delete(connection.NodeAssociations, assocAddr)
}

// DeleteSession deletes a session and all PDRs, FARs and QERs associated with it.
func (connection *Pfcp_Link) DeleteSession(session *Session) {
	for _, far := range session.FARs {
		_ = connection.EBPFMapManager.DeleteFar(far.GlobalId)
	}
	pdrContext := NewPDRCreationContext(session, connection.ResourceManager)
	for _, PDR := range session.PDRs {
		_ = pdrContext.deletePDR(PDR, connection.EBPFMapManager)
	}
}

func (connection *Pfcp_Link) GetSessionCount() int {
	count := 0
	for _, assoc := range connection.NodeAssociations {
		count += len(assoc.Sessions)
	}
	return count
}

func (connection *Pfcp_Link) GetAssiciationCount() int {
	return len(connection.NodeAssociations)
}

func (connection *Pfcp_Link) ReleaseResources(seID uint64) {
	if connection.ResourceManager == nil {
		return
	}

	if connection.ResourceManager.IPAM != nil {
		connection.ResourceManager.IPAM.ReleaseIP(seID)
	}

	if connection.ResourceManager.FTEIDM != nil {
		connection.ResourceManager.FTEIDM.ReleaseTEID(seID)
	}
}
