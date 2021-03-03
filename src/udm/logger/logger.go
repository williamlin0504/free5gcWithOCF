package logger

import (
	"os"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"

	"free5gcWithOCF/lib/logger_conf"
	"free5gcWithOCF/lib/logger_util"
)

var log *logrus.Logger
var AppLog *logrus.Entry
var InitLog *logrus.Entry
var Handlelog *logrus.Entry
var HttpLog *logrus.Entry
var UeauLog *logrus.Entry
var UecmLog *logrus.Entry
var SdmLog *logrus.Entry
var PpLop *logrus.Entry
var EeLog *logrus.Entry
var UtilLog *logrus.Entry
var CallbackLog *logrus.Entry
var ContextLog *logrus.Entry
var GinLog *logrus.Entry

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	free5gcWithOCFLogHook, err := logger_util.NewFileHook(logger_conf.Free5gcLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		log.Hooks.Add(free5gcWithOCFLogHook)
	}

	selfLogHook, err := logger_util.NewFileHook(logger_conf.NfLogDir+"udm.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		log.Hooks.Add(selfLogHook)
	}

	AppLog = log.WithFields(logrus.Fields{"component": "UDM", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "UDM", "category": "Init"})
	Handlelog = log.WithFields(logrus.Fields{"component": "UDM", "category": "Handler"})
	HttpLog = log.WithFields(logrus.Fields{"component": "UDM", "category": "HTTP"})
	UeauLog = log.WithFields(logrus.Fields{"component": "UDM", "category": "UEAU"})
	UecmLog = log.WithFields(logrus.Fields{"component": "UDM", "category": "UECM"})
	SdmLog = log.WithFields(logrus.Fields{"component": "UDM", "category": "SDM"})
	PpLop = log.WithFields(logrus.Fields{"component": "UDM", "category": "PP"})
	EeLog = log.WithFields(logrus.Fields{"component": "UDM", "category": "EE"})
	UtilLog = log.WithFields(logrus.Fields{"component": "UDM", "category": "Util"})
	CallbackLog = log.WithFields(logrus.Fields{"component": "UDM", "category": "Callback"})
	ContextLog = log.WithFields(logrus.Fields{"component": "UDM", "category": "Context"})
	GinLog = log.WithFields(logrus.Fields{"component": "UDM", "category": "GIN"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}
