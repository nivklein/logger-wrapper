package logger

import (
	"errors"
	buildinlogger "log"
	"net/http"
	"os"
	"sync"
)

//Based on the good blogpost at https://www.mountedthoughts.com/golang-logger-interface/
// A global variable so that log functions can be directly accessed
var log Logger
var doOnce sync.Once

//Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

const (
	//Debug has verbose message
	Debug = "debug"
	//Info is default log level
	Info = "info"
	//Warn is for logging messages about possible issues
	Warn = "warn"
	//Error is for logging errors
	Error = "error"
	//Fatal is for logging fatal messages. The sytem shutsdown after logging the message.
	Fatal = "fatal"
)

const (
	InstanceZapLogger int = iota
)

const EnvKeyEnv = "env"

var (
	errInvalidLoggerInstance = errors.New("Invalid logger instance")
	DefaultConfig            = LoggerConfig{
		EnableConsole:     true,
		ConsoleLevel:      Debug,
		ConsoleJSONFormat: false,
		EnableFile:        false,
	}
)

//Logger is our contract for the logger
type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Panicf(format string, args ...interface{})

	ChangeLogLevel(w http.ResponseWriter, r *http.Request)

	WithFields(keyValues Fields) Logger
}

// LoggerConfig stores the config for the logger
// For some loggers there can only be one level across writers, for such the level of Console is picked by default
type LoggerConfig struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
}

//NewLogger returns an instance of logger
func NewLogger(config LoggerConfig, loggerInstance int) error {
	switch loggerInstance {
	case InstanceZapLogger:
		logger, err := newZapLogger(config)
		if err != nil {
			return err
		}
		log = logger
		return nil
	default:
		return errInvalidLoggerInstance
	}
}

func Debugf(format string, args ...interface{}) {
	if log == nil {
		initLogger()
	}
	log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	if log == nil {
		initLogger()
	}
	log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	if log == nil {
		initLogger()
	}
	log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	if log == nil {
		initLogger()
	}
	log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	if log == nil {
		initLogger()
	}
	log.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	if log == nil {
		initLogger()
	}
	log.Panicf(format, args...)
}

func ChangeLogLevel(w http.ResponseWriter, r *http.Request) {
	if log == nil {
		initLogger()
	}
	log.ChangeLogLevel(w, r)
}

func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}

func initLogger() {
	doOnce.Do(func() {
		if os.Getenv(EnvKeyEnv) == "prod" {
			DefaultConfig.ConsoleLevel = Warn
		}
		err := NewLogger(DefaultConfig, InstanceZapLogger)
		if err != nil {
			buildinlogger.Fatalf("Could not instantiate log %s", err.Error())
		}
		Infof("Logger created successfuly")
	})
}
