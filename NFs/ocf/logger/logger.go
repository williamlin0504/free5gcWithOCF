package logger

import (
	"os"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"

	"github.com/free5gc/logger_conf"
	"github.com/free5gc/logger_util"
)

var (
	log         *logrus.Logger
	AppLog      *logrus.Entry
	InitLog     *logrus.Entry
	CfgLog      *logrus.Entry
	ContextLog  *logrus.Entry
	NgapLog     *logrus.Entry
	HandlerLog  *logrus.Entry
	HttpLog     *logrus.Entry
	GmmLog      *logrus.Entry
	MtLog       *logrus.Entry
	ProducerLog *logrus.Entry
	LocationLog *logrus.Entry
	CommLog     *logrus.Entry
	CallbackLog *logrus.Entry
	UtilLog     *logrus.Entry
	NasLog      *logrus.Entry
	ConsumerLog *logrus.Entry
	EeLog       *logrus.Entry
	GinLog      *logrus.Entry
)

const (
	FieldRanAddr     string = "ran_addr"
	FieldOcfUeNgapID string = "ocf_ue_ngap_id"
	FieldSupi        string = "supi"
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category", FieldRanAddr, FieldOcfUeNgapID, FieldSupi},
	}

	free5gcLogHook, err := logger_util.NewFileHook(logger_conf.Free5gcLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err == nil {
		log.Hooks.Add(free5gcLogHook)
	}

	selfLogHook, err := logger_util.NewFileHook(logger_conf.NfLogDir+"ocf.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err == nil {
		log.Hooks.Add(selfLogHook)
	}

	AppLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "Init"})
	CfgLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "CFG"})
	ContextLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "Context"})
	NgapLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "NGAP"})
	HandlerLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "Handler"})
	HttpLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "HTTP"})
	GmmLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "GMM"})
	MtLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "MT"})
	ProducerLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "Producer"})
	LocationLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "LocInfo"})
	CommLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "Comm"})
	CallbackLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "Callback"})
	UtilLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "Util"})
	NasLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "NAS"})
	ConsumerLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "Consumer"})
	EeLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "EventExposure"})
	GinLog = log.WithFields(logrus.Fields{"component": "OCF", "category": "GIN"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(set bool) {
	log.SetReportCaller(set)
}
