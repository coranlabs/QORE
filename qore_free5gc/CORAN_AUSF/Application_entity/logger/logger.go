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
	Auth5gAkaLog *logrus.Entry
	UeAuthLog    *logrus.Entry
	AuthELog   	 *logrus.Entry
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
	
	NfLog       = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "NF"})
	MainLog     = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "Main"})
	InitLog     = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "Init"})
	CfgLog      = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "Config"})
	CtxLog      = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "Context"})
	DataRepoLog = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "DataRepo"})
	UtilLog     = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "Util"})
	HttpLog     = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "HTTP"})
	ConsumerLog = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "Consumer"})
	GinLog      = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "Gin"})
	ProcLog     = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "Processor"})
	SBILog      = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "SBI"})
	DbLog       = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "Database"})
	Auth5gAkaLog = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "Auth5g"})
	UeAuthLog   = Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "UEAuth"})
	AuthELog   =  Log.WithFields(logrus.Fields{"component": "CORAN_AUSF", "category": "AuthE"})
}


func SetLogLevel(level logrus.Level) {
	Log.SetLevel(level)
}