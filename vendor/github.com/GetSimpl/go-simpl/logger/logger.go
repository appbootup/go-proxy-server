package logger

import (
	"fmt"
	"os"
	"errors"
	"strconv"
	"net/http"

	"github.com/Sirupsen/logrus"
	"gopkg.in/polds/logrus-papertrail-hook.v2"
	"gopkg.in/airbrake/gobrake.v2"
	"github.com/GetSimpl/go-simpl/db"
)

var log *logrus.Logger
var env string
var airbrake *gobrake.Notifier
var airbrakeEnabled = false

const (
	DEBUG = 0
	INFO = 1
	WARNING = 2
	ERROR = 3
)

func init() {
	env = os.Getenv("ENV")
	if env == "" {
		env = "dev"
		return
	}
}

func Init(mode int) {
	log = logrus.New()
	log.Formatter = new(logrus.TextFormatter)
	switch mode{
	case DEBUG:
		log.Level = logrus.DebugLevel
	case INFO:
		log.Level = logrus.InfoLevel
	case WARNING:
		log.Level = logrus.WarnLevel
	case ERROR:
		log.Level = logrus.ErrorLevel
	}
}

func AddAirbrakeHook(projectId int64, projectKey string) {
	airbrake = gobrake.NewNotifier(projectId, projectKey)
	airbrakeEnabled = true
}

func AddPapertrailHook(serviceName string, papertrailHost string, papertrailPort int) {
	switch env{
	case "dev":
		fallthrough
	case "test":
		log.Error("Can not add papertrail hook for dev/test env")
		return
	}
	log.Debug("Adding papertrail hook")
	papertrailHook := &logrus_papertrail.Hook{}
	papertrailHook.Host = os.Getenv("PAPERTRAIL_HOST")
	port, err := strconv.Atoi(os.Getenv("PAPERTRAIL_PORT"))
	if err != nil {
		panic(err)
	}
	papertrailHook.Port = port
	papertrailHook.Hostname = fmt.Sprintf("%s-%s", serviceName, env)
	papertrailHook.Appname = serviceName
	log.Println(papertrailHook.Port)
	hook, err := logrus_papertrail.NewPapertrailHook(papertrailHook)
	if err != nil {
		panic(err)
	}
	log.Hooks.Add(hook)
}

func Get() *logrus.Logger {
	return log
}

func I(args ...interface{}) {
	log.Info(args)
}

func Df(format string, args ...interface{}) {
	log.Debugf(format, args)
}

func D(args ...interface{}) {
	log.Debug(args)
}

func W(args ...interface{}) {
	log.Warn(args)
}

func E(args ...interface{}) {
	log.Error(args)
	if airbrakeEnabled == false {
		return
	}
	switch env {
	case "dev":
		fallthrough
	case "test":
		log.Error("Can not add airbrake hook for dev/test env")
		return
	case "staging":
		fallthrough
	case "production":
		var err error
		var req *http.Request
		var rawData []interface{}
		for _, arg := range args {
			switch arg.(type) {
			case *http.Request:
				req = arg.(*http.Request)
			case error:
				err = arg.(error)
				log.Println("error", err)
				continue
			default:
				rawData = append(rawData, arg)
			}
		}
		if err == nil {
			err = errors.New(fmt.Sprintf("%v", rawData))
		}
		// skipping record not found error
		if !shouldPushToAirbrake(err) {
			return
		}
		notice := airbrake.Notice(err, req, 1)
		notice.Context["environment"] = env
		_, err = airbrake.SendNotice(notice)
		if err != nil {
			log.Error("While sendning notice to the airbrake", err)
		}
	}
}

func shouldPushToAirbrake(err error) bool {
	if err == db.RecordNotFound {
		return false
	}
	return true
}
