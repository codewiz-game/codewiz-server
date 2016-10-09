package log

import (
	"github.com/Sirupsen/logrus"
)

var theInstance *Logger
var theLevel Level

type Logger struct {
	*logrus.Logger
}

type Fields logrus.Fields

type Level uint8

const (
    // PanicLevel level, highest level of severity. Logs and then calls panic with the
    // message passed to Debug, Info, ...
    PanicLevel = Level(logrus.PanicLevel)
    // FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
    // logging level is set to Panic.
    FatalLevel =  Level(logrus.FatalLevel)
    // ErrorLevel level. Logs. Used for errors that should definitely be noted.
    // Commonly used for hooks to send errors to an error tracking service.
    ErrorLevel = Level(logrus.ErrorLevel)
    // WarnLevel level. Non-critical entries that deserve eyes.
    WarnLevel =  Level(logrus.WarnLevel)
    // InfoLevel level. General operational entries about what's going on inside the
    // application.
    InfoLevel = Level(logrus.InfoLevel)
    // DebugLevel level. Usually only enabled when debugging. Very verbose logging.
    DebugLevel = Level(logrus.DebugLevel)
)

func init() { 	
	theInstance = NewLogger()
}

func NewLogger() *Logger {
	logger := logrus.New()

	// TODO: these should be args
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.InfoLevel
	return &Logger{logger}
}

func (logger *Logger) SetLevel(level Level) {
	logger.Logger.Level = logrus.Level(level)
}

func (logger *Logger) Panic(message interface{}, fields ...Fields) {
	var log logrus.FieldLogger = logger.Logger
	if len(fields) > 0 {
		log = log.WithFields(logrus.Fields(fields[0]))
	}

	log.Panic(message)
}

func (logger *Logger) Fatal(message interface{}, fields ...Fields) {
	var log logrus.FieldLogger = logger.Logger
	if len(fields) > 0 {
		log = log.WithFields(logrus.Fields(fields[0]))
	}

	log.Fatal(message)
}

func (logger *Logger) Error(message interface{}, fields ...Fields) {
	var log logrus.FieldLogger = logger.Logger
	if len(fields) > 0 {
		log = log.WithFields(logrus.Fields(fields[0]))
	}

	log.Error(message)
}

func (logger *Logger) Warn(message interface{}, fields ...Fields) {
	var log logrus.FieldLogger = logger.Logger
	if len(fields) > 0 {
		log = log.WithFields(logrus.Fields(fields[0]))
	}

	log.Warn(message)
}

func (logger *Logger) Info(message interface{}, fields ...Fields) {
	var log logrus.FieldLogger = logger.Logger
	if len(fields) > 0 {
		log = log.WithFields(logrus.Fields(fields[0]))
	}

	log.Info(message)
}

func (logger *Logger) Debug(message interface{}, fields ...Fields) {
	var log logrus.FieldLogger = logger.Logger
	if len(fields) > 0 {
		log = log.WithFields(logrus.Fields(fields[0]))
	}

	log.Debug(message)
}

func SetLevel(level Level) {
	theInstance.SetLevel(level)
}

func ParseLevel(levelName string) (Level, error) {
	level, err := logrus.ParseLevel(levelName)
	return Level(level), err
}

func Panic(message interface{}, fields ...Fields) {
	theInstance.Panic(message, fields...)
}

func Fatal(message interface{}, fields ...Fields) {
	theInstance.Fatal(message, fields...)
}

func Error(message interface{}, fields ...Fields) {
	theInstance.Error(message, fields...)
}

func Warn(message interface{}, fields ...Fields) {
	theInstance.Warn(message, fields...)
}

func Info(message interface{}, fields ...Fields) {
	theInstance.Info(message, fields...)
}

func Debug(message interface{}, fields ...Fields) {
	theInstance.Debug(message, fields...)	
}