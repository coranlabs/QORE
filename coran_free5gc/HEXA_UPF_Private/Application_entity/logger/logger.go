// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package logger

// import (
// 	"time"

// 	formatter "github.com/antonfisher/nested-logrus-formatter"
// 	"github.com/sirupsen/logrus"
// )

// var (
// 	log       *logrus.Logger
// 	AppLog    *logrus.Entry
// 	InitLog   *logrus.Entry
// 	ConfigLog *logrus.Entry
// 	XdpLog    *logrus.Entry
// 	Pfcplog   *logrus.Entry
// )

// func init() {
// 	log = logrus.New()
// 	log.SetReportCaller(false)

// 	log.Formatter = &formatter.Formatter{
// 		TimestampFormat: time.RFC3339,
// 		TrimMessages:    true,
// 		NoFieldsSpace:   true,
// 		HideKeys:        true,
// 		FieldsOrder:     []string{"component", "category"},
// 	}

// 	AppLog = log.WithFields(logrus.Fields{"component": "CORAN_UPF_eBPF", "category": "App"})
// 	InitLog = log.WithFields(logrus.Fields{"component": "CORAN_UPF_eBPF", "category": "Init"})
// 	ConfigLog = log.WithFields(logrus.Fields{"component": "CORAN_UPF_eBPF", "category": "CFG"})
// 	XdpLog = log.WithFields(logrus.Fields{"component": "CORAN_UPF_eBPF", "category": "XDP"})
// 	Pfcplog = log.WithFields(logrus.Fields{"component": "CORAN_UPF_eBPF", "category": "PFCP"})
// }

// func SetLogLevel(level logrus.Level) {
// 	log.SetLevel(level)
// }

// // func SetReportCaller(set bool) {
// // 	log.SetReportCaller(set)
// // }
