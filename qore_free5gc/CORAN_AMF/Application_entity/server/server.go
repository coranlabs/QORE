package server

import (
	"context"
	"os/signal"
	"syscall"

	"os"

	// "github.com/coranlabs/CORAN_AMF/Application_entity/config"
	//"github.com/coranlabs/CORAN_AMF/Messages_controller/context"
	// version "github.com/coranlabs/CORAN_LIB_FSM"
	//"github.com/coranlabs/CORAN_AMF/Messages_handling_entity/ngap_handler"
	"github.com/coranlabs/CORAN_AMF/Application_entity/config/factory"
	"github.com/coranlabs/CORAN_AMF/Application_entity/logger"
	"github.com/coranlabs/CORAN_AMF/Application_entity/pkg/service"
)

func Action() error {
	tlsKeyLogPath := ""

	// logger.MainLog.Infoln("AMF version: ", version.GetVersion())

	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh  // Wait for interrupt signal to gracefully shutdown
		cancel() // Notify each goroutine and wait them stopped
	}()
	logger.MainLog.Infoln("after ")

	cfg, err := factory.ReadConfig("./config/CORAN_AMF.yaml")
	logger.MainLog.Infoln("after 2")
	if err != nil {
		return err
	}
	factory.AmfConfig = cfg

	amf, err := service.NewApp(ctx, cfg, tlsKeyLogPath)
	if err != nil {
		return err
	}

	amf.Start()

	return nil
}

// func initLogFile(logNfPath []string) (string, error) {
// 	logTlsKeyPath := ""

// 	for _, path := range logNfPath {
// 		if err := logger_util.LogFileHook(logger.Log, path); err != nil {
// 			return "", err
// 		}

// 		if logTlsKeyPath != "" {
// 			continue
// 		}

// 		nfDir, _ := filepath.Split(path)
// 		tmpDir := filepath.Join(nfDir, "key")
// 		if err := os.MkdirAll(tmpDir, 0o775); err != nil {
// 			logger.InitLog.Errorf("Make directory %s failed: %+v", tmpDir, err)
// 			return "", err
// 		}
// 		_, name := filepath.Split(factory.AmfDefaultTLSKeyLogPath)
// 		logTlsKeyPath = filepath.Join(tmpDir, name)
// 	}

// 	return logTlsKeyPath, nil
// }
// type NGAPHandler struct {
//     HandleMessage         func(conn net.Conn, msg []byte)
//     HandleNotification    func(conn net.Conn)
//     HandleConnectionError func(conn net.Conn)
// }
// var connections  sync.Map
// var Ranpool  sync.Map

// type Server struct{
// 	SctpListener *sctp.SCTPListener
// 	// Nf_server *gin.Engine

// }

// const readBufSize uint32 = 262144

// func NewServer() (*Server, error) {
// 	var s *Server

// 	return s , nil
// }
// func Start_amf_server(c *config.Amf_config) error{

// 	fmt.Println("starting amf server")

// 	amf, err := NewServer()
// 	if err != nil{
// 		fmt.Println(" oops error occured")
// 		os.Exit(1)
// 	}
// 	fmt.Printf("amf %v",amf)
// 	return nil
// }

// func NewSctpConfig(cfg *config.Amf_config) *sctp.SocketConfig {
// 	sctpConfig := &sctp.SocketConfig{
// 		InitMsg: sctp.InitMsg{
// 			NumOstreams:    uint16(cfg.NumOstreams),
// 			MaxInstreams:   uint16(cfg.MaxInstreams),
// 			MaxAttempts:    uint16(cfg.MaxAttempts),
// 			MaxInitTimeout: uint16(cfg.MaxInitTimeout),
// 		},
// 		// RtoInfo:   &sctp.RtoInfo{SrtoAssocID: 0, SrtoInitial: 500, SrtoMax: 1500, StroMin: 100},
// 		// AssocInfo: &sctp.AssocInfo{AsocMaxRxt: 4},
// 	}
// 	return sctpConfig
// }

// func Sctp_handler(conf *config.Amf_config) string {

//     addresses := []string{"192.168.4.19"}

// 	ips := []net.IPAddr{}
// 	port := 38412
// 	sctpConfig := NewSctpConfig(conf)
// 	for _, addr := range addresses {
// 		if netAddr, err := net.ResolveIPAddr("ip", addr); err != nil {
// 		} else {
// 			ips = append(ips, *netAddr)
// 		}
// 	}

// 	addr := &sctp.SCTPAddr{
// 		IPAddrs: ips,
// 		Port:    port,
// 	}
// 	var wg sync.WaitGroup
// 		wg.Add(1)

// 	go func() {
// 		defer wg.Done() // Decrement the counter when the goroutine completes
// 		Serving_sctp(addr,sctpConfig)
// 	}()

// 	// Wait for the goroutine to finish
// 	wg.Wait()

// 		return "success"

// }

// func Serving_sctp(addr *sctp.SCTPAddr, sctpConfig *sctp.SocketConfig){

// 	listener, err := sctpConfig.Listen("sctp", addr)
// 	if  err != nil {
// 		fmt.Printf("Failed to start server on : %v\n", err)
//     	os.Exit(1)
// 		return
// 	}
// 	sctpListener := listener

// 	fmt.Printf("Listen on %s", sctpListener.Addr())

// 	for{
// 		newConn, err := sctpListener.AcceptSCTP()

// 		if err != nil {
// 			switch err {
// 			case syscall.EINTR, syscall.EAGAIN:
// 				fmt.Printf("AcceptSCTP: %+v", err)
// 			default:
// 				fmt.Printf("Failed to accept: %+v", err)
// 			}
// 			continue
// 		}
// 		//buffer := make([]byte, 2046) // Adjust the size based on your expected message size
// 		var info *sctp.SndRcvInfo
// 		if infoTmp, errGetDefaultSentParam := newConn.GetDefaultSentParam(); errGetDefaultSentParam != nil {
// 			fmt.Printf("Get default sent param error: %+v, accept failed", errGetDefaultSentParam)
// 			if errGetDefaultSentParam = newConn.Close(); errGetDefaultSentParam != nil {
// 				fmt.Printf("Close error: %+v", errGetDefaultSentParam)
// 			}
// 			continue
// 		} else {
// 			info = infoTmp
// 			fmt.Printf("Get default sent param[value: %+v]\n", info)
// 		}
// 		info.PPID = ngap.PPID
// 		if errSetDefaultSentParam := newConn.SetDefaultSentParam(info); errSetDefaultSentParam != nil {
// 			fmt.Printf("Set default sent param error: %+v, accept failed", errSetDefaultSentParam)
// 			if errSetDefaultSentParam = newConn.Close(); errSetDefaultSentParam != nil {
// 				fmt.Printf("Close error: %+v", errSetDefaultSentParam)
// 			}
// 			continue
// 		} else {
// 			fmt.Printf("Set default sent param[value: %+v]\n", info)
// 		}
// 		events := sctp.SCTP_EVENT_DATA_IO | sctp.SCTP_EVENT_SHUTDOWN | sctp.SCTP_EVENT_ASSOCIATION
// 		if errSubscribeEvents := newConn.SubscribeEvents(events); errSubscribeEvents != nil {
// 			fmt.Printf("Failed to accept: %+v", errSubscribeEvents)
// 			if errSubscribeEvents = newConn.Close(); errSubscribeEvents != nil {
// 				fmt.Printf("Close error: %+v", errSubscribeEvents)
// 			}
// 			continue
// 		} else {
// 			fmt.Printf("Subscribe SCTP event[DATA_IO, SHUTDOWN_EVENT, ASSOCIATION_CHANGE]")
// 		}
// 		if errSetReadBuffer := newConn.SetReadBuffer(int(readBufSize)); errSetReadBuffer != nil {
// 			fmt.Printf("Set read buffer error: %+v, accept failed", errSetReadBuffer)
// 			if errSetReadBuffer = newConn.Close(); errSetReadBuffer != nil {
// 				fmt.Printf("Close error: %+v", errSetReadBuffer)
// 			}
// 			continue
// 		} else {
// 			fmt.Printf("Set read buffer to %d bytes\n", readBufSize)
// 		}

// 		connections.Store(newConn,newConn)

// 		 go Sctp_type_handler(newConn,readBufSize)

// }
// }

// func Sctp_type_handler(conn *sctp.SCTPConn, bufsize uint32) {
// 	defer func() {
// 		if p := recover(); p != nil {
// 			fmt.Printf("this is panic log")
// 			// Print stack for panic to log. Fatalf() will let program exit.

// 		}

// 		if err := conn.Close(); err != nil && err != syscall.EBADF {
// 			fmt.Printf("close connection error: %+v", err)
// 		}
// 		connections.Delete(conn)
// 	}()

// 	for {
// 		buf := make([]byte, bufsize)

// 		n, info,  err := conn.SCTPRead(buf)
// 		if err != nil {
// 			switch err {
// 			case io.EOF, io.ErrUnexpectedEOF:
// 				fmt.Printf("Read EOF from client")
// 				// handler.HandleConnectionError(conn)
// 				return
// 			case syscall.EAGAIN:
// 				fmt.Println("SCTP read timeout")
// 				continue
// 			case syscall.EINTR:
// 				fmt.Printf("SCTPRead: %+v", err)
// 				continue
// 			default:
// 				fmt.Printf(
// 					"Handle connection[addr: %+v] error: %+v",
// 					conn.RemoteAddr(),
// 					err,
// 				)
// 				// HandleConnectionError(conn)
// 				return
// 			}
// 		}else {
// 			if info == nil || info.PPID != ngap.PPID {
// 				fmt.Printf("Received SCTP PPID != 60, discard this packet")
// 				continue
// 			}

// 			fmt.Printf("Read %d bytes", n)
// 			fmt.Printf("Packet content:\n%+v", hex.Dump(buf[:n]))

// 			// TODO: concurrent on per-UE message
// 			internal_ngap.HandleMessage(conn, buf[:n])
// 		}
// 	}
// }
