package PFCP_server

import (
	"net"

	"github.com/coranlabs/CORAN_GO_PFCP/message"
	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/logger"
	"github.com/pkg/errors"
)

func (s *Pfcp_Link) Handle_request(msg message.Message, addr net.Addr) error {
	switch req := msg.(type) {
	case *message.HeartbeatRequest:{
		//TS 29.244 6.2.2 PFCP Heartbeat Procedure
		s.Heartbeat_Procedure(req, addr.String())}
	case *message.AssociationSetupRequest:{
		//TS 29.244 6.2.6 PFCP Association Setup Procedure
		s.PFCP_Association_Setup_Procedure(req, addr.String())}
	case *message.AssociationUpdateRequest:{
		//TS 29.244 6.2.7 PFCP Association Update Procedure
		s.PFCP_Association_Update_Procedure(req, addr)}
	case *message.AssociationReleaseRequest:{
		//TS 29.244 6.2.8 PFCP Association Release Procedure
		s.PFCP_Association_Release_Procedure(req, addr)}
	case *message.SessionEstablishmentRequest:{
		//TS 29.244 6.3.2 PFCP Session Establishment Procedure
		s.PFCP_Session_Establishment_Procedure(req, addr.String())}
	case *message.SessionModificationRequest:{
		//TS 29.244 6.3.3 PFCP Session Modification Procedure
		s.PFCP_Session_Modification_Procedure(req, addr.String())}
	case *message.SessionDeletionRequest:{
		//TS 29.244 6.3.4 PFCP Session Deletion Procedure
		s.PFCP_Session_Deletion_Procedure(req, addr.String())}
	default:
		return errors.Errorf("pfcp reqDispacher unknown msg type: %d", msg.MessageType())
	}
	return nil
}
func (s *Pfcp_Link) Handle_response(msg message.Message, addr net.Addr, req message.Message) error {
	switch rsp := msg.(type) {
	case *message.SessionReportResponse:{
		//TS 29.244 6.2.2.3 Heartbeat Response 
		s.Heartbeat_Response(rsp, addr.String())}
	default:
		return errors.Errorf("pfcp rspDispacher unknown msg type: %d", msg.MessageType())
	}
	return nil
}

func (s *Pfcp_Link) PFCP_Association_Update_Procedure(
	req *message.AssociationUpdateRequest,
	addr net.Addr,
) {
	logger.PfcpAssocLog.Infoln("handleAssociationUpdateRequest not supported")
}

func (s *Pfcp_Link) PFCP_Association_Release_Procedure(
	req *message.AssociationReleaseRequest,
	addr net.Addr,
) {
	logger.PfcpAssocLog.Infoln("handleAssociationReleaseRequest not supported")
}
