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
	DataRepoLog *logrus.Entry
	UtilLog     *logrus.Entry
	HttpLog     *logrus.Entry
	ConsumerLog *logrus.Entry
	GinLog      *logrus.Entry
	ProcLog     *logrus.Entry
	SBILog      *logrus.Entry
	DbLog       *logrus.Entry
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
	
	NfLog       = Log.WithFields(logrus.Fields{"component": "CORAN_UDR", "category": "NF"})
	MainLog     = Log.WithFields(logrus.Fields{"component": "CORAN_UDR", "category": "Main"})
	InitLog     = Log.WithFields(logrus.Fields{"component": "CORAN_UDR", "category": "Init"})
	CfgLog      = Log.WithFields(logrus.Fields{"component": "CORAN_UDR", "category": "Config"})
	CtxLog      = Log.WithFields(logrus.Fields{"component": "CORAN_UDR", "category": "Context"})
	DataRepoLog = Log.WithFields(logrus.Fields{"component": "CORAN_UDR", "category": "DataRepo"})
	UtilLog     = Log.WithFields(logrus.Fields{"component": "CORAN_UDR", "category": "Util"})
	HttpLog     = Log.WithFields(logrus.Fields{"component": "CORAN_UDR", "category": "HTTP"})
	ConsumerLog = Log.WithFields(logrus.Fields{"component": "CORAN_UDR", "category": "Consumer"})
	GinLog      = Log.WithFields(logrus.Fields{"component": "CORAN_UDR", "category": "Gin"})
	ProcLog     = Log.WithFields(logrus.Fields{"component": "CORAN_UDR", "category": "Processor"})
	SBILog      = Log.WithFields(logrus.Fields{"component": "CORAN_UDR", "category": "SBI"})
	DbLog       = Log.WithFields(logrus.Fields{"component": "CORAN_UDR", "category": "Database"})


}


func SetLogLevel(level logrus.Level) {
	Log.SetLevel(level)
}
