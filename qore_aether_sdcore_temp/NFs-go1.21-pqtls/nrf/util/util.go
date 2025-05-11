// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0

//go:build !debug
// +build !debug

package util

import (
	"encoding/json"
	"reflect"

	"github.com/omec-project/nrf/logger"
	"github.com/omec-project/path_util"
)

// Path of HTTP2 key and log file

var (
	NrfLogPath = path_util.Free5gcPath("free5gc/nrfsslkey.log")
	NrfPemPath = path_util.Free5gcPath("free5gc/support/TLS/nrfpq.pem")
	NrfKeyPath = path_util.Free5gcPath("free5gc/support/TLS/keypq.pem")
)

func MarshToJsonString(v interface{}) (result []string) {
	types := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if types.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			tmp, err := json.Marshal(val.Index(i).Interface())
			if err != nil {
				logger.UtilLog.Errorf("Marshal error: %+v", err)
			}

			result = append(result, string(tmp))
		}
	} else {
		tmp, err := json.Marshal(v)
		if err != nil {
			logger.UtilLog.Errorf("Marshal error: %+v", err)
		}

		result = append(result, string(tmp))
	}
	return
}
