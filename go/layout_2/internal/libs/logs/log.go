package logs

import (
	"time"

	"layout_2/configs"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

func SetLogEnv() error {
	level, err := log.ParseLevel(configs.App.GetLogLevel())
	if err != nil {
		return err
	}

	logFormat := &log.JSONFormatter{}

	log.SetLevel(level)
	log.SetFormatter(logFormat)
	log.SetReportCaller(true)

	if err := logAddHookRotatelogs(logFormat); err != nil {
		return err
	}

	return nil
}

func logAddHookRotatelogs(formatter log.Formatter) error {
	writer, err := rotatelogs.New(
		configs.App.GetLogPath()+".%Y-%m-%d-%H-%M",

		rotatelogs.WithLinkName(configs.App.GetLogPath()),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		return err
	}

	log.AddHook(lfshook.NewHook(lfshook.WriterMap{
		log.PanicLevel: writer,
		log.FatalLevel: writer,
		log.ErrorLevel: writer,
		log.WarnLevel:  writer,
		log.InfoLevel:  writer,
		log.DebugLevel: writer,
		log.TraceLevel: writer,
	}, formatter))

	return nil
}

func ReloadSetLogLevel() {
	level, err := log.ParseLevel(configs.App.GetLogLevel())
	if err != nil {
		log.WithFields(log.Fields{
			"newLogLevel": configs.App.GetLogLevel(),
			"oldLogLevel": log.GetLevel().String(),
			"err":         err,
		}).Error("ReloadSetLogLevel")
		return
	}

	log.SetLevel(level)
}
