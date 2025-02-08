package processor

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/coranlabs/CORAN_LIB_NAS/nasMessage"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	pfcp "github.com/coranlabs/CORAN_LIB_PFCP"
	"github.com/coranlabs/CORAN_LIB_PFCP/pfcpType"
	"github.com/coranlabs/CORAN_SMF/Application_entity/logger"
	smf_context "github.com/coranlabs/CORAN_SMF/Messages_handling_entity/context"
	"github.com/coranlabs/CORAN_SMF/Messages_handling_entity/pfcp/message"
)

func (p *Processor) ToBeAssociatedWithUPF(ctx context.Context, upf *smf_context.UPF) {
	var upfStr string
	if upf.NodeID.NodeIdType == pfcpType.NodeIdTypeFqdn {
		upfStr = fmt.Sprintf("[%s](%s)", upf.NodeID.FQDN, upf.NodeID.ResolveNodeIdToIp().String())
	} else {
		upfStr = fmt.Sprintf("[%s]", upf.NodeID.ResolveNodeIdToIp().String())
	}

	for {
		ensureSetupPfcpAssociation(ctx, upf, upfStr)
		if isDone(ctx, upf) {
			break
		}

		if smf_context.GetSelf().PfcpHeartbeatInterval == 0 {
			return
		}

		keepHeartbeatTo(ctx, upf, upfStr ,p)
		// return when UPF heartbeat lost is detected or association is canceled
		if isDone(ctx, upf) {
			break
		}

		p.releaseAllResourcesOfUPF(upf, upfStr)
		if isDone(ctx, upf) {
			break
		}
	}
}

func (p *Processor) ReleaseAllResourcesOfUPF(upf *smf_context.UPF) {
	var upfStr string
	if upf.NodeID.NodeIdType == pfcpType.NodeIdTypeFqdn {
		upfStr = fmt.Sprintf("[%s](%s)", upf.NodeID.FQDN, upf.NodeID.ResolveNodeIdToIp().String())
	} else {
		upfStr = fmt.Sprintf("[%s]", upf.NodeID.ResolveNodeIdToIp().String())
	}
	p.releaseAllResourcesOfUPF(upf, upfStr)
}

func isDone(ctx context.Context, upf *smf_context.UPF) bool {
	select {
	case <-ctx.Done():
		return true
	case <-upf.Ctx.Done():
		return true
	default:
		return false
	}
}

func ensureSetupPfcpAssociation(ctx context.Context, upf *smf_context.UPF, upfStr string) {
	alertTime := time.Now()
	alertInterval := smf_context.GetSelf().AssocFailAlertInterval
	retryInterval := smf_context.GetSelf().AssocFailRetryInterval
	for {
		timer := time.After(retryInterval)
		err := setupPfcpAssociation(upf, upfStr)
		if err == nil {
			return
		}
		logger.MainLog.Warnf("Failed to setup an association with UPF%s, error:%+v", upfStr, err)
		now := time.Now()
		logger.MainLog.Debugf("now %+v, alertTime %+v", now, alertTime)
		if now.After(alertTime.Add(alertInterval)) {
			logger.MainLog.Errorf("ALERT for UPF%s", upfStr)
			alertTime = now
		}
		logger.MainLog.Debugf("Wait %+v (or less) until next retry attempt", retryInterval)
		select {
		case <-ctx.Done():
			logger.MainLog.Infof("Canceled association request to UPF%s", upfStr)
			return
		case <-upf.Ctx.Done():
			logger.MainLog.Infof("Canceled association request to this UPF%s only", upfStr)
			return
		case <-timer:
			continue
		}
	}
}

func setupPfcpAssociation(upf *smf_context.UPF, upfStr string) error {
	logger.MainLog.Infof("Sending PFCP Association Request to UPF%s", upfStr)

	resMsg, err := message.SendPfcpAssociationSetupRequest(upf.NodeID)
	if err != nil {
		return err
	}

	rsp := resMsg.PfcpMessage.Body.(pfcp.PFCPAssociationSetupResponse)

	if rsp.Cause == nil || rsp.Cause.CauseValue != pfcpType.CauseRequestAccepted {
		return fmt.Errorf("received PFCP Association Setup Not Accepted Response from UPF%s", upfStr)
	}

	nodeID := rsp.NodeID
	if nodeID == nil {
		return fmt.Errorf("pfcp association needs NodeID")
	}

	logger.MainLog.Infof("Received PFCP Association Setup Accepted Response from UPF%s", upfStr)

	upf.UPFStatus = smf_context.AssociatedSetUpSuccess

	if rsp.UserPlaneIPResourceInformation != nil {
		upf.UPIPInfo = *rsp.UserPlaneIPResourceInformation

		logger.MainLog.Infof("UPF(%s)[%s] setup association",
			upf.NodeID.ResolveNodeIdToIp().String(), upf.UPIPInfo.NetworkInstance.NetworkInstance)
	}

	return nil
}
const heartbeatFailureThreshold = 3

func keepHeartbeatTo(ctx context.Context, upf *smf_context.UPF, upfStr string,p *Processor ) {
	unansweredHeartbeats := 0

	for {
		err := doPfcpHeartbeat(upf, upfStr)
		if err != nil {
			logger.MainLog.Errorf("PFCP Heartbeat error: %v", err)
			unansweredHeartbeats++
			if unansweredHeartbeats >= heartbeatFailureThreshold {
				logger.MainLog.Warnf("Heartbeat failed %d times. Releasing UPF resources and terminating association with UPF%s", unansweredHeartbeats, upfStr)
				p.releaseAllResourcesOfUPF(upf, upfStr)
				return
			}
		} else {
			// Reset the failure count if heartbeat succeeds
			unansweredHeartbeats = 0
		}

		timer := time.After(smf_context.GetSelf().PfcpHeartbeatInterval)
		select {
		case <-ctx.Done():
			logger.MainLog.Infof("Canceled Heartbeat with UPF%s", upfStr)
			return
		case <-upf.Ctx.Done():
			logger.MainLog.Infof("Canceled Heartbeat to this UPF%s only", upfStr)
			return
		case <-timer:
			continue
		}
	}
}

func doPfcpHeartbeat(upf *smf_context.UPF, upfStr string) error {
	if upf.UPFStatus != smf_context.AssociatedSetUpSuccess {
		return fmt.Errorf("invalid status of UPF%s: %d", upfStr, upf.UPFStatus)
	}

	logger.MainLog.Debugf("Sending PFCP Heartbeat Request to UPF%s", upfStr)

	resMsg, err := message.SendPfcpHeartbeatRequest(upf)
	if err != nil {
		upf.UPFStatus = smf_context.NotAssociated
		upf.RecoveryTimeStamp = time.Time{}
		return fmt.Errorf("SendPfcpHeartbeatRequest error: %w", err)
	}

	rsp := resMsg.PfcpMessage.Body.(pfcp.HeartbeatResponse)
	if rsp.RecoveryTimeStamp == nil {
		logger.MainLog.Warnf("Received PFCP Heartbeat Response without timestamp from UPF%s", upfStr)
		return nil
	}

	logger.MainLog.Debugf("Received PFCP Heartbeat Response from UPF%s", upfStr)
	if upf.RecoveryTimeStamp.IsZero() {
		// first receive
		upf.RecoveryTimeStamp = rsp.RecoveryTimeStamp.RecoveryTimeStamp
	} else if upf.RecoveryTimeStamp.Before(rsp.RecoveryTimeStamp.RecoveryTimeStamp) {
		// received a newer recovery timestamp
		upf.UPFStatus = smf_context.NotAssociated
		upf.RecoveryTimeStamp = time.Time{}
		return fmt.Errorf("received PFCP Heartbeat Response RecoveryTimeStamp has been updated")
	}
	return nil
}

func (p *Processor) releaseAllResourcesOfUPF(upf *smf_context.UPF, upfStr string) {
	logger.MainLog.Infof("Release all resources of UPF %s", upfStr)

	upf.ProcEachSMContext(func(smContext *smf_context.SMContext) {
		smContext.SMLock.Lock()
		defer smContext.SMLock.Unlock()
		switch smContext.State() {
		case smf_context.Active, smf_context.ModificationPending, smf_context.PFCPModification:
			needToSendNotify, removeContext := p.requestAMFToReleasePDUResources(smContext)
			if needToSendNotify {
				p.SendReleaseNotification(smContext)
			}
			if removeContext {
				// Notification has already been sent, if it is needed
				p.RemoveSMContextFromAllNF(smContext, false)
			}
		}
	})
}

func (p *Processor) requestAMFToReleasePDUResources(
	smContext *smf_context.SMContext,
) (sendNotify bool, releaseContext bool) {
	n1n2Request := models.N1N2MessageTransferRequest{}
	// TS 23.502 4.3.4.2 3b. Send Namf_Communication_N1N2MessageTransfer Request, SMF->AMF
	n1n2Request.JsonData = &models.N1N2MessageTransferReqData{
		PduSessionId: smContext.PDUSessionID,
		SkipInd:      true,
	}
	cause := nasMessage.Cause5GSMNetworkFailure
	if buf, err := smf_context.BuildGSMPDUSessionReleaseCommand(smContext, cause, false); err != nil {
		logger.MainLog.Errorf("Build GSM PDUSessionReleaseCommand failed: %+v", err)
	} else {
		n1n2Request.BinaryDataN1Message = buf
		n1n2Request.JsonData.N1MessageContainer = &models.N1MessageContainer{
			N1MessageClass:   "SM",
			N1MessageContent: &models.RefToBinaryData{ContentId: "GSM_NAS"},
		}
	}
	if smContext.UpCnxState != models.UpCnxState_DEACTIVATED {
		if buf, err := smf_context.BuildPDUSessionResourceReleaseCommandTransfer(smContext); err != nil {
			logger.MainLog.Errorf("Build PDUSessionResourceReleaseCommandTransfer failed: %+v", err)
		} else {
			n1n2Request.BinaryDataN2Information = buf
			n1n2Request.JsonData.N2InfoContainer = &models.N2InfoContainer{
				N2InformationClass: models.N2InformationClass_SM,
				SmInfo: &models.N2SmInformation{
					PduSessionId: smContext.PDUSessionID,
					N2InfoContent: &models.N2InfoContent{
						NgapIeType: models.NgapIeType_PDU_RES_REL_CMD,
						NgapData: &models.RefToBinaryData{
							ContentId: "N2SmInformation",
						},
					},
					SNssai: smContext.SNssai,
				},
			}
		}
	}

	ctx, _, errToken := smf_context.GetSelf().GetTokenCtx(models.ServiceName_NAMF_COMM, models.NfType_AMF)
	if errToken != nil {
		return false, false
	}

	rspData, statusCode, err := p.Consumer().
		N1N2MessageTransfer(ctx, smContext.Supi, n1n2Request, smContext.CommunicationClientApiPrefix)
	if err != nil {
		logger.ConsumerLog.Warnf("N1N2MessageTransfer for RequestAMFToReleasePDUResources failed: %+v", err)
	}

	switch *statusCode {
	case http.StatusOK:
		if rspData.Cause == models.N1N2MessageTransferCause_N1_MSG_NOT_TRANSFERRED {
			// the PDU Session Release Command was not transferred to the UE since it is in CM-IDLE state.
			//   ref. step3b of "4.3.4.2 UE or network requested PDU Session Release for Non-Roaming and
			//        Roaming with Local Breakout" in TS23.502
			// it is needed to remove both AMF's and SMF's SM Contexts immediately
			smContext.SetState(smf_context.InActive)
			return true, true
		} else if rspData.Cause == models.N1N2MessageTransferCause_N1_N2_TRANSFER_INITIATED {
			// wait for N2 PDU Session Release Response
			smContext.SetState(smf_context.InActivePending)
		} else {
			// other causes are unexpected.
			// keep SM Context to avoid inconsistency with AMF
			smContext.SetState(smf_context.InActive)
		}
	case http.StatusNotFound:
		// it is not needed to notify AMF, but needed to remove SM Context in SMF immediately
		smContext.SetState(smf_context.InActive)
		return false, true
	default:
		// keep SM Context to avoid inconsistency with AMF
		smContext.SetState(smf_context.InActive)
	}
	return false, false
}
