// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0
//

//go:build !debug
// +build !debug

package util

import (
	"github.com/omec-project/path_util"
)

var (
	AusfLogPath           = path_util.Free5gcPath("free5gc/ausfsslkey.log")
	AusfPemPath           = path_util.Free5gcPath("free5gc/support/TLS/ausf.pem")
	AusfKeyPath           = path_util.Free5gcPath("free5gc/support/TLS/ausf.key")
	DefaultAusfConfigPath = path_util.Free5gcPath("free5gc/config/ausfcfg.yaml")
)
