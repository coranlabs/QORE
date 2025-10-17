package rating

import (
	"fmt"
	"strconv"
	"time"

	chf_context "github.com/coranlabs/CORAN_CHF/Application_entity/internal/context"
	"github.com/coranlabs/CORAN_CHF/Application_entity/pkg/factory"

	"github.com/coranlabs/CORAN_CHF/Application_entity/internal/logger"

	charging_code "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/ccs_diameter/code"
	charging_datatype "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/ccs_diameter/datatype"
	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/fiorix/go-diameter/diam/sm/smpeer"
)

func SendServiceUsageRequest(
	ue *chf_context.ChfUe, sur *charging_datatype.ServiceUsageRequest,
) (*charging_datatype.ServiceUsageResponse, error) {
	ue.RatingMux.Handle("SUA", HandleSUA(ue.RatingChan))
	rfDiameter := factory.ChfConfig.Configuration.RfDiameter
	addr := rfDiameter.HostIPv4 + ":" + strconv.Itoa(rfDiameter.Port)
	conn, err := ue.RatingClient.DialNetworkTLS(rfDiameter.Protocol, addr, rfDiameter.Tls.Pem, rfDiameter.Tls.Key)
	if err != nil {
		return nil, err
	}

	meta, ok := smpeer.FromContext(conn.Context())
	if !ok {
		return nil, fmt.Errorf("peer metadata unavailable")
	}

	sur.DestinationRealm = datatype.DiameterIdentity(meta.OriginRealm)
	sur.DestinationHost = datatype.DiameterIdentity(meta.OriginHost)

	msg := diam.NewRequest(charging_code.ServiceUsageMessage, charging_code.Re_interface, dict.Default)
	err = msg.Marshal(sur)
	if err != nil {
		return nil, fmt.Errorf("Marshal SUR Failed: %s\n", err)
	}

	_, err = msg.WriteTo(conn)
	if err != nil {
		return nil, fmt.Errorf("Failed to send message from %s: %s\n",
			conn.RemoteAddr(), err)
	}

	select {
	case m := <-ue.RatingChan:
		var sua charging_datatype.ServiceUsageResponse
		if err := m.Unmarshal(&sua); err != nil {
			return nil, fmt.Errorf("Failed to parse message from %v", err)
		}
		return &sua, nil
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout: no rate answer received")
	}
}

func HandleSUA(rgChan chan *diam.Message) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		logger.RatingLog.Tracef("Received SUA from %s", c.RemoteAddr())

		rgChan <- m
	}
}
