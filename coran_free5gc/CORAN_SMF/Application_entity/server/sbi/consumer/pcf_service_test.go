package consumer_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/coranlabs/CORAN_LIB_NAS/nasType"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/factory"
	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/service"
	"github.com/coranlabs/CORAN_SMF/Application_entity/server/sbi/consumer"
	smf_context "github.com/coranlabs/CORAN_SMF/Messages_handling_entity/context"
)

var testConfig = factory.Config{
	Info: &factory.Info{
		Version:     "1.0.0",
		Description: "SMF procdeure test configuration",
	},
	Configuration: &factory.Configuration{
		Sbi: &factory.Sbi{
			Scheme:       "http",
			RegisterIPv4: "127.0.0.1",
			BindingIPv4:  "127.0.0.1",
			Port:         8000,
		},
	},
}

func TestSendSMPolicyAssociationUpdateByUERequestModification(t *testing.T) {
	smf_context.InitSmfContext(&testConfig)

	testCases := []struct {
		name         string
		smContext    *smf_context.SMContext
		qosRules     nasType.QoSRules
		qosFlowDescs nasType.QoSFlowDescs

		smPolicyDecision *models.SmPolicyDecision
		responseErr      error
	}{
		{
			name:             "QoSRules is nil",
			smContext:        smf_context.NewSMContext("imsi-208930000000001", 10),
			qosRules:         nasType.QoSRules{},
			qosFlowDescs:     nasType.QoSFlowDescs{nasType.QoSFlowDesc{}},
			smPolicyDecision: nil,
			responseErr:      fmt.Errorf("QoS Rule not found"),
		},
		{
			name:             "QoSFlowDescs is nil",
			smContext:        smf_context.NewSMContext("imsi-208930000000001", 10),
			qosRules:         nasType.QoSRules{nasType.QoSRule{}},
			qosFlowDescs:     nasType.QoSFlowDescs{},
			smPolicyDecision: nil,
			responseErr:      fmt.Errorf("QoS Flow Description not found"),
		},
	}

	mockSmf := service.NewMockSmfAppInterface(gomock.NewController(t))
	consumer, errNewConsumer := consumer.NewConsumer(mockSmf)
	if errNewConsumer != nil {
		t.Fatalf("Failed to create consumer: %+v", errNewConsumer)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			smPolicyDecision, err := consumer.SendSMPolicyAssociationUpdateByUERequestModification(
				tc.smContext, tc.qosRules, tc.qosFlowDescs)

			require.Equal(t, tc.smPolicyDecision, smPolicyDecision)
			require.Equal(t, tc.responseErr.Error(), err.Error())
		})
	}
}
