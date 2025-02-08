package console

import "github.com/coranlabs/CORAN_LIB_OPENAPI/models"

type SubsData struct {
	PlmnID                            string                                     `json:"plmnID"`
	UeId                              string                                     `json:"ueId"`
	AuthenticationSubscription        models.AuthenticationSubscription          `json:"AuthenticationSubscription"`
	AccessAndMobilitySubscriptionData models.AccessAndMobilitySubscriptionData   `json:"AccessAndMobilitySubscriptionData"`
	SessionManagementSubscriptionData []models.SessionManagementSubscriptionData `json:"SessionManagementSubscriptionData"`
	SmfSelectionSubscriptionData      models.SmfSelectionSubscriptionData        `json:"SmfSelectionSubscriptionData"`
	AmPolicyData                      models.AmPolicyData                        `json:"AmPolicyData"`
	SmPolicyData                      models.SmPolicyData                        `json:"SmPolicyData"`
	FlowRules                         []FlowRule                                 `json:"FlowRules"`
	QosFlows                          []QosFlow                                  `json:"QosFlows"`
	ChargingDatas                     []ChargingData                             `json:"ChargingDatas"`
}
