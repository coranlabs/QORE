package logger

import (
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var (
	Log         *logrus.Logger
	NfLog       *logrus.Entry
	MainLog     *logrus.Entry
	InitLog     *logrus.Entry
	CfgLog      *logrus.Entry
	CtxLog      *logrus.Entry
	GinLog      *logrus.Entry
	SBILog      *logrus.Entry
	ConsumerLog *logrus.Entry
	HttpLog     *logrus.Entry
	UeauLog     *logrus.Entry
	UecmLog     *logrus.Entry
	SdmLog      *logrus.Entry
	PpLog       *logrus.Entry
	EeLog       *logrus.Entry
	UtilLog     *logrus.Entry
	SuciLog     *logrus.Entry
	CallbackLog *logrus.Entry
	ProcLog     *logrus.Entry
	AppLog      *logrus.Entry
	
)


func init() {
	Log = logrus.New()
	Log.SetReportCaller(false)

	Log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	NfLog       = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "NF"})
	MainLog     = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "Main"})
	InitLog     = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "Init"})
	CfgLog      = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "CFG"})
	CtxLog      = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "CTX"})
	GinLog      = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "Gin"})
	SBILog      = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "SBI"})
	ConsumerLog = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "Consumer"})
	HttpLog     = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "HTTP"})
	UeauLog     = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "UEAU"})
	UecmLog     = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "UECM"})
	SdmLog      = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "SDM"})
	PpLog       = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "PP"})
	EeLog       = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "EE"})
	UtilLog     = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "Util"})
	SuciLog     = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "SUCI"})
	CallbackLog = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "Callback"})
	ProcLog     = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "Proc"})
	AppLog      = Log.WithFields(logrus.Fields{"component": "CORAN_UDM", "category": "App"})

}

func SetLogLevel(level logrus.Level) {
	Log.SetLevel(level)
}

// func SetReportCaller(set bool) {
// 	log.SetReportCaller(set)
// }