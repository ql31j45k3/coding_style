package driver

import (
	"net"
	"time"

	"github.com/ql31j45k3/coding_style/go/layout/configs"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
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

	logAddServerIPHock()
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

type serverIPHock struct {
	_ struct{}
}

func (h *serverIPHock) Levels() []log.Level {
	return log.AllLevels
}

func (h *serverIPHock) Fire(entry *log.Entry) error {
	serverIP, err := h.getServerIP()
	if err != nil {
		return err
	}

	entry.Data["server_ip"] = serverIP

	return nil
}

func (h *serverIPHock) getServerIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, v := range addrs {
		ipNet, ok := v.(*net.IPNet)
		if !ok {
			continue
		}

		if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}

	return "", err
}

func logAddServerIPHock() {
	hock := &serverIPHock{}
	log.AddHook(hock)
}
