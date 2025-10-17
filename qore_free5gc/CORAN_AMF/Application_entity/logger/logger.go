package logger

import (
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	logger_util "github.com/coranlabs/CORAN_LIB_UTIL/logger"
	"github.com/sirupsen/logrus"
)

var (
	log         *logrus.Logger
	AppLog      *logrus.Entry
	InitLog     *logrus.Entry
	CfgLog      *logrus.Entry
	XdpLog      *logrus.Entry
	Pfcplog     *logrus.Entry
	CtxLog      *logrus.Entry
	gmmlog      *logrus.Entry
	NgapLog     *logrus.Entry
	Log         *logrus.Logger
	NfLog       *logrus.Entry
	MainLog     *logrus.Entry
	GinLog      *logrus.Entry
	HandlerLog  *logrus.Entry
	HttpLog     *logrus.Entry
	GmmLog      *logrus.Entry
	MtLog       *logrus.Entry
	ProducerLog *logrus.Entry
	SBILog      *logrus.Entry
	LocationLog *logrus.Entry
	CommLog     *logrus.Entry
	CallbackLog *logrus.Entry
	UtilLog     *logrus.Entry
	NasLog      *logrus.Entry
	ConsumerLog *logrus.Entry
	EeLog       *logrus.Entry
)

const (
	FieldRanAddr     string = "ran_addr"
	FieldAmfUeNgapID string = "amf_ue_ngap_id"
	FieldSupi        string = "supi"
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)
	fieldsOrder := []string{
		logger_util.FieldNF,
		logger_util.FieldCategory,
	}

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	AppLog = log.WithFields(logrus.Fields{"component": "HEXA_UPF", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "HEXA_UPF", "category": "Init"})
	CfgLog = log.WithFields(logrus.Fields{"component": "HEXA_UPF", "category": "CFG"})
	UtilLog = log.WithFields(logrus.Fields{"component": "HEXA_UPF", "category": "UTL"})

	CtxLog = log.WithFields(logrus.Fields{"component": "HEXA_UPF", "category": "CTX"})
	XdpLog = log.WithFields(logrus.Fields{"component": "HEXA_UPF", "category": "XDP"})
	Pfcplog = log.WithFields(logrus.Fields{"component": "HEXA_UPF", "category": "PFCP"})
	Log = logger_util.New(fieldsOrder)
	NfLog = Log.WithField(logger_util.FieldNF, "AMF")
	MainLog = NfLog.WithField(logger_util.FieldCategory, "Main")
	InitLog = NfLog.WithField(logger_util.FieldCategory, "Init")
	CfgLog = NfLog.WithField(logger_util.FieldCategory, "CFG")
	CtxLog = NfLog.WithField(logger_util.FieldCategory, "CTX")
	GinLog = NfLog.WithField(logger_util.FieldCategory, "GIN")
	NgapLog = NfLog.WithField(logger_util.FieldCategory, "Ngap")
	HandlerLog = NfLog.WithField(logger_util.FieldCategory, "Handler")
	HttpLog = NfLog.WithField(logger_util.FieldCategory, "Http")
	GmmLog = NfLog.WithField(logger_util.FieldCategory, "Gmm")
	MtLog = NfLog.WithField(logger_util.FieldCategory, "Mt")
	ProducerLog = NfLog.WithField(logger_util.FieldCategory, "Producer")
	SBILog = NfLog.WithField(logger_util.FieldCategory, "SBI")
	LocationLog = NfLog.WithField(logger_util.FieldCategory, "Location")
	CommLog = NfLog.WithField(logger_util.FieldCategory, "Comm")
	CallbackLog = NfLog.WithField(logger_util.FieldCategory, "Callback")
	UtilLog = NfLog.WithField(logger_util.FieldCategory, "Util")
	NasLog = NfLog.WithField(logger_util.FieldCategory, "Nas")
	ConsumerLog = NfLog.WithField(logger_util.FieldCategory, "Consumer")
	EeLog = NfLog.WithField(logger_util.FieldCategory, "Ee")
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

// func SetReportCaller(set bool) {
// 	log.SetReportCaller(set)
// }
