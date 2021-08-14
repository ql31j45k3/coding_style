package driver

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/ql31j45k3/coding_style/go/layout/di/configs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

func SetLogEnv() {
	level, err := log.ParseLevel(configs.Env.GetLogLevel())
	if err != nil {
		panic(err)
	}

	logFormat := &log.JSONFormatter{}

	log.SetLevel(level)
	log.SetFormatter(logFormat)

	logAddHookRotatelogs(logFormat)
}

func logAddHookRotatelogs(formatter log.Formatter) {
	writer, err := rotatelogs.New(
		configs.Env.GetLogPath()+".%Y-%m-%d-%H-%M",

		rotatelogs.WithLinkName(configs.Env.GetLogPath()),
		rotatelogs.WithMaxAge(time.Duration(24)*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
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
}

func ReloadSetLogLevel() {
	level, err := log.ParseLevel(configs.Env.GetLogLevel())
	if err != nil {
		log.WithFields(log.Fields{
			"newLogLevel": configs.Env.GetLogLevel(),
			"oldLogLevel": log.GetLevel().String(),
			"err":         err,
		}).Error("SetLogLevel")
		return
	}

	log.SetLevel(level)
}
