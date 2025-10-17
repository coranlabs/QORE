package consumer

import (
	"context"
	"testing"

	openapi "github.com/coranlabs/CORAN_LIB_OPENAPI"
	"github.com/coranlabs/CORAN_UDM/Application_entity/pkg/app"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	udm_context "github.com/coranlabs/CORAN_UDM/Message_controller/context"
)

func TestSendRegisterNFInstance(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.InterceptClient(openapi.GetHttpClient())
	defer gock.RestoreClient(openapi.GetHttpClient())

	gock.New("http://127.0.0.10:8000").
		Put("/nnrf-nfm/v1/nf-instances/1").
		Reply(200).
		JSON(map[string]string{})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApp := app.NewMockApp(ctrl)
	consumer, err := NewConsumer(mockApp)
	require.NoError(t, err)

	mockApp.EXPECT().Context().Times(1).Return(
		&udm_context.UDMContext{
			NrfUri: "http://127.0.0.10:8000",
			NfId:   "1",
		},
	)

	_, _, err = consumer.RegisterNFInstance(context.TODO())
	require.NoError(t, err)
}
