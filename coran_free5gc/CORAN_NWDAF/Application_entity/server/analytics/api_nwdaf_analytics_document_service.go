/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the OAI Public License, Version 1.1  (the "License"); you may not use this
 * file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 *      http://www.openairinterface.org/?page_id=698
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*
 * Author: Abdelkader Mekrache <mekrache@eurecom.fr>
 * Author: Karim Boutiba 	   <boutiba@eurecom.fr>
 * Author: Arina Prostakova    <prostako@eurecom.fr>
 * Description: Functions of the analytics nbi service.
 */

package analytics

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// ------------------------------------------------------------------------------
type NWDAFAnalyticsDocumentApiService struct {
}

// ------------------------------------------------------------------------------
// Type of num_of_ue data to request engine
type EngineReqData struct {
	StartTs time.Time `json:"startTs,omitempty"`
	EndTs   time.Time `json:"endTs,omitempty"`
	Tais    []Tai     `json:"networkArea,omitempty"`
	Dnns    []string  `json:"dnns,omitempty"`
	Snssaia []Snssai  `json:"snssaia,omitempty"`
	Supi    string    `json:"supi,omitempty"`
}

// ------------------------------------------------------------------------------
// Type of network_performance response from engine.
type NwPerfResp struct {
	RelativeRatio int32 `json:"relativeRatio,omitempty"`
	AbsoluteNum   int32 `json:"absoluteNum,omitempty"`
	Confidence    int32 `json:"confidence,omitempty"`
}

// ------------------------------------------------------------------------------
// Type of Ue_communication response from engine.
type UeCommResp struct {
	CommDur       int32     `json:"commDur"`
	Ts            time.Time `json:"ts,omitempty"`
	UlVol         int64     `json:"ulVol,omitempty"`
	UlVolVariance float32   `json:"ulVolVariance,omitempty"`
	DlVol         int64     `json:"dlVol,omitempty"`
	DlVolVariance float32   `json:"dlVolVariance,omitempty"`
}

// ------------------------------------------------------------------------------
// Type of Ue_mobility response from engine.
type UeMobResp struct {
	Loc []UserLocation `json:"loc"`
}

// ------------------------------------------------------------------------------
// NewNWDAFAnalyticsDocumentApiService - create a default api service
func NewNWDAFAnalyticsDocumentApiService() NWDAFAnalyticsDocumentApiServicer {
	return &NWDAFAnalyticsDocumentApiService{}
}

// ------------------------------------------------------------------------------
// GetNWDAFAnalytics - read a NWDAF Analytics
func (s *NWDAFAnalyticsDocumentApiService) GetNWDAFAnalytics(
	ctx context.Context,
	eventId EventIdAnyOf,
	anaReq EventReportingRequirement,
	eventFilter EventFilter,
	supportedFeatures string,
	tgtUe TargetUeInformation,
) (ImplResponse, error) {
	// Create AnalyticsData Report
	analyticsData := AnalyticsData{}
	if eventId == "" {
		return Response(http.StatusBadRequest, ProblemDetails{}),
			errors.New("missing Event Id param")
	}
	switch eventId {
	case EVENTIDANYOF_NETWORK_PERFORMANCE:
		// get list of NetworkPerfInfo
		nwPerfAnalyticsData, err := getNwPerfAnalytics(anaReq, eventFilter)
		if err != nil {
			return Response(http.StatusBadRequest, ProblemDetails{}), err
		}
		// check if list is empty
		if len(nwPerfAnalyticsData) == 0 {
			return Response(http.StatusBadRequest, ProblemDetails{}),
				errors.New("missing Network Performance Data")
		}
		analyticsData.NwPerfs = nwPerfAnalyticsData

	case EVENTIDANYOF_UE_COMMUNICATION:
		// get list of UeCommunication
		ueCommsAnalyticsData, err := getUeCommsAnalytics(anaReq, eventFilter)
		if err != nil {
			return Response(http.StatusBadRequest, ProblemDetails{}), err
		}
		// check if list is empty
		if len(ueCommsAnalyticsData) == 0 {
			return Response(http.StatusBadRequest, ProblemDetails{}),
				errors.New("missing UE Communications Data")
		}
		analyticsData.UeComms = ueCommsAnalyticsData

	case EVENTIDANYOF_UE_MOBILITY:
		// get list of UeCommunication
		ueMobAnalyticsData, err := getUeMobAnalytics(anaReq, eventFilter, tgtUe)
		if err != nil {
			return Response(http.StatusBadRequest, ProblemDetails{}), err
		}
		// check if list is empty
		if len(ueMobAnalyticsData) == 0 {
			return Response(http.StatusBadRequest, ProblemDetails{}),
				errors.New("missing UE Communications Data")
		}
		analyticsData.UeMobs = ueMobAnalyticsData

	default:
		return Response(http.StatusBadRequest, ProblemDetails{}),
			errors.New("invalid Event Id param")
	}
	analyticsData.AnaMetaInfo.DataWindow.StartTime = anaReq.StartTs
	analyticsData.AnaMetaInfo.DataWindow.StopTime = anaReq.EndTs
	log.Printf("Returning NWDAF Analytics Data")
	return Response(http.StatusOK, analyticsData), nil
}

// ------------------------------------------------------------------------------
// getNwPerfAnalytics - Get list of NetworkPerfInfo
func getNwPerfAnalytics(
	anaReq EventReportingRequirement,
	eventFilter EventFilter,
) ([]NetworkPerfInfo, error) {
	log.Printf("Getting NW Performance Analytics")
	var nwPerfList []NetworkPerfInfo
	// for each NwPerfType, request the engine
	for _, nwPerfType := range eventFilter.NwPerfTypes {
		var nwPerfInfo NetworkPerfInfo
		var err error
		switch nwPerfType {
		case NETWORKPERFTYPEANYOF_NUM_OF_UE:
			nwPerfInfo, err = requestNwPerfEngine(
				anaReq,
				eventFilter,
				config.Engine.Uri+config.Routes.NumOfUe,
			)
			if err != nil {
				return nwPerfList, err
			}

		case NETWORKPERFTYPEANYOF_SESS_SUCC_RATIO:
			nwPerfInfo, err = requestNwPerfEngine(
				anaReq,
				eventFilter,
				config.Engine.Uri+config.Routes.SessSuccRatio,
			)
			if err != nil {
				return nwPerfList, err
			}

		default:
			return nil, errors.New("invalid Network Performance Type")
		}
		nwPerfInfo.NwPerfType = nwPerfType
		nwPerfList = append(nwPerfList, nwPerfInfo)
	}
	return nwPerfList, nil
}

// ------------------------------------------------------------------------------
// getUeCommNotifData - Get list Ue Communication
func getUeCommsAnalytics(
	anaReq EventReportingRequirement,
	eventFilter EventFilter,
) ([]UeCommunication, error) {

	log.Printf("Getting UE Communications Notification Data")
	var ueCommList []UeCommunication
	// this treat just one type of UE_COMMUNICATION
	var ueCommInfo UeCommunication
	var err error
	ueCommInfo, err = requestUeCommEngine(
		anaReq,
		eventFilter,
		config.Engine.Uri+config.Routes.UeComm,
	)
	if err != nil {
		return ueCommList, err
	}
	ueCommList = append(ueCommList, ueCommInfo)
	return ueCommList, nil
}

// ------------------------------------------------------------------------------
// getUeCommNotifData - Get list Ue Communication
func getUeMobAnalytics(
	anaReq EventReportingRequirement,
	eventFilter EventFilter,
	tgtUe TargetUeInformation,
) ([]UeMobility, error) {
	log.Printf("Getting UE Mobility Notification Data")
	var ueMobList []UeMobility
	// check supis not empty
	if len(tgtUe.Supis) == 0 {
		return ueMobList, errors.New("missing supis param in TgtUe")
	}
	// for each User imsi, request the engine to get location.
	for _, supi := range tgtUe.Supis {
		var ueMobInfo UeMobility
		var err error
		ueMobInfo, err = requestUeMobEngine(
			anaReq,
			eventFilter,
			supi,
			config.Engine.Uri+config.Routes.UeMob,
		)
		if err != nil {
			return ueMobList, err
		}
		// TODO we need to add supi to the response ? ueMobInfo.Supi = supi
		ueMobList = append(ueMobList, ueMobInfo)
	}
	return ueMobList, nil
}

// ------------------------------------------------------------------------------
func requestNwPerfEngine(
	anaReq EventReportingRequirement,
	eventFilter EventFilter,
	enginePath string,
) (NetworkPerfInfo, error) {
	log.Printf("Reaching engine to get Network Performance Info from DB")
	var engineReqData EngineReqData
	engineReqData.StartTs = anaReq.StartTs
	engineReqData.EndTs = anaReq.EndTs
	// for num_of_ue
	engineReqData.Tais = eventFilter.NetworkArea.Tais
	// for sess_succ_ratio request
	engineReqData.Dnns = eventFilter.Dnns
	engineReqData.Snssaia = eventFilter.Snssais
	// Convert the data to a JSON byte array
	engineReqJsonData, err := json.Marshal(engineReqData)
	if err != nil {
		return NetworkPerfInfo{}, err
	}
	// Create a POST request with the JSON data in the body
	req, err := http.NewRequest(
		http.MethodGet,
		enginePath, bytes.NewBuffer(engineReqJsonData))
	if err != nil {
		return NetworkPerfInfo{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	// Send the request and print the response body
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return NetworkPerfInfo{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	var nwPerfResp NwPerfResp
	err = json.Unmarshal(body, &nwPerfResp)
	if err != nil {
		return NetworkPerfInfo{}, err
	}
	nwPerfInfo := NetworkPerfInfo{
		NetworkArea:   eventFilter.NetworkArea,
		AbsoluteNum:   &nwPerfResp.AbsoluteNum,
		RelativeRatio: &nwPerfResp.RelativeRatio,
		Confidence:    &nwPerfResp.Confidence,
	}
	return nwPerfInfo, nil
}

// ------------------------------------------------------------------------------
func requestUeCommEngine(
	anaReq EventReportingRequirement,
	eventFilter EventFilter,
	enginePath string,
) (UeCommunication, error) {
	log.Printf("Reaching engine to get UE Communication Info from DB")
	var engineReqData EngineReqData
	engineReqData.StartTs = anaReq.StartTs
	engineReqData.EndTs = anaReq.EndTs
	// Convert the data to a JSON byte array
	engineReqJsonData, err := json.Marshal(engineReqData)
	if err != nil {
		return UeCommunication{}, err
	}
	// Create a POST request with the JSON data in the body
	req, err := http.NewRequest(
		http.MethodGet,
		enginePath,
		bytes.NewBuffer(engineReqJsonData),
	)
	if err != nil {
		return UeCommunication{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	// Send the request and print the response body
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return UeCommunication{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	var ueCommResp UeCommResp
	err = json.Unmarshal(body, &ueCommResp)
	if err != nil {
		return UeCommunication{}, err
	}
	trafChar := TrafficCharacterization{
		UlVol:         &ueCommResp.UlVol,
		UlVolVariance: &ueCommResp.UlVolVariance,
		DlVol:         &ueCommResp.DlVol,
		DlVolVariance: &ueCommResp.DlVolVariance,
	}
	ueCommunication := UeCommunication{
		CommDur:  ueCommResp.CommDur,
		Ts:       anaReq.StartTs,
		TrafChar: trafChar,
	}
	return ueCommunication, nil
}

// ------------------------------------------------------------------------------
func requestUeMobEngine(
	anaReq EventReportingRequirement,
	eventFilter EventFilter,
	supi string,
	enginePath string,
) (UeMobility, error) {
	log.Printf("Reaching engine to get UE Mobility Info from DB")
	log.Printf("Supi : %s", supi)
	var engineReqData EngineReqData
	engineReqData.StartTs = anaReq.StartTs
	engineReqData.EndTs = anaReq.EndTs
	engineReqData.Supi = supi
	// Convert the data to a JSON byte array
	engineReqJsonData, err := json.Marshal(engineReqData)
	if err != nil {
		return UeMobility{}, err
	}
	// Create a POST request with the JSON data in the body
	req, err := http.NewRequest(
		http.MethodGet,
		enginePath,
		bytes.NewBuffer(engineReqJsonData),
	)
	if err != nil {
		return UeMobility{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	// Send the request and print the response body
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return UeMobility{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var ueMobResp UeMobResp
	err = json.Unmarshal(body, &ueMobResp)
	if err != nil {
		return UeMobility{}, err
	}
	// Create a variable of type UeMobility
	var ueMobility UeMobility
	// Fill the Ts field with the current time, and duration
	ueMobility.Ts = time.Now()
	ueMobility.Duration = 10
	// Iterate over the Loc slice in UeMobResp
	for _, userLocation := range ueMobResp.Loc {
		locationInfo := LocationInfo{
			Loc:        userLocation,
			Ratio:      100, // Set the ratio to 100 as an example, you can change it as needed
			Confidence: 0,   // Set the confidence to 0 as an example, you can change it as needed
		}
		ueMobility.LocInfos = append(ueMobility.LocInfos, locationInfo)
	}
	return ueMobility, nil
}
