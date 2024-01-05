package log

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"hypier.fun/hdwallet/hdwallet-go-sdk/utils"
	"io"
	"os"
	"path/filepath"
	"time"
)

var logger *Logger
var conSoleRule *string

const (
	CONSOLE_FILE = "CONSOLE_FILE"
	CONSOLE      = "CONSOLE"
	FILE         = "FILE"
	CLOSE        = "CLOSE"
)

type Logger struct {
	*logrus.Logger
	Pid int
}

func InitLog(logFile string, con string) {
	conSoleRule = &con
	logger = loggerInit(logFile, 6)
}

func defaultLog() {
	InitLog("../logs", CONSOLE_FILE)
}

func NewPrivateLog(dir string, logLevel uint32) {
	logger = loggerInit(dir, logLevel)
}

func GetLogger() *Logger {
	return logger
}

type UTCFormatter struct {
	logrus.Formatter
}

func (u UTCFormatter) Format(e *logrus.Entry) ([]byte, error) {

	e.Time = e.Time.In(time.UTC)
	return u.Formatter.Format(e)
}

func loggerInit(dir string, logLevel uint32) *Logger {
	var logger = logrus.New()
	logger.SetLevel(logrus.Level(logLevel))

	var formatter logrus.Formatter = &nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        true,
		FieldsOrder:     []string{"FilePath"},
	}

	formatter = UTCFormatter{formatter}
	logger.SetFormatter(formatter)
	logger.AddHook(newFileHook())
	if *conSoleRule == CONSOLE_FILE {
		logger.SetOutput(os.Stdout)
		logger.AddHook(NewLfsHook(time.Duration(1)*time.Hour, 24, dir))
	} else if *conSoleRule == CONSOLE {
		logger.SetOutput(os.Stdout)
	} else if *conSoleRule == FILE {
		logger.SetOutput(io.Discard)
		logger.AddHook(NewLfsHook(time.Duration(1)*time.Hour, 24, dir))
	} else if *conSoleRule == CLOSE {
		logger.SetOutput(io.Discard)
	} else {
		logger.SetOutput(os.Stdout)
	}
	return &Logger{
		logger,
		os.Getpid(),
	}
}

func NewLfsHook(rotationTime time.Duration, maxRemainNum uint, dir string) logrus.Hook {
	var formatter logrus.Formatter = &nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        true,
		FieldsOrder:     []string{"FilePath"},
		TrimMessages:    true,
		NoFieldsSpace:   true,
		NoColors:        true,
	}
	formatter = UTCFormatter{formatter}
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: initRotateLogs(rotationTime, maxRemainNum, "all", dir),
		logrus.InfoLevel:  initRotateLogs(rotationTime, maxRemainNum, "all", dir),
		logrus.WarnLevel:  initRotateLogs(rotationTime, maxRemainNum, "all", dir),
		logrus.ErrorLevel: initRotateLogs(rotationTime, maxRemainNum, "all", dir),
	}, formatter)
	return lfsHook
}

func initRotateLogs(rotationTime time.Duration, maxRemainNum uint, level string, dir string) *rotatelogs.RotateLogs {
	writer, err := rotatelogs.New(
		filepath.Join(dir, level+"."+"%Y-%m-%d-%H-%M"+".log"),
		rotatelogs.WithRotationTime(rotationTime),
		rotatelogs.WithRotationCount(maxRemainNum),
	)
	if err != nil {
		panic(err.Error())
	} else {
		return writer
	}
}

func Info(args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Infoln(args)
}
func Error(args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Errorln(args)
}
func Debug(args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Debugln(args)
}
func Warn(args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Warnln(args)
}
func Debugf(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Infof(format, args...)
}
func Errorf(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Errorf(format, args...)
}
func Warnf(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Warnf(format, args...)
}

func GetNetLogger() *logrus.Logger {
	return logger.Logger
}

func WithError(err error, message ...string) error {
	var u *utils.Error

	if !errors.As(err, &u) {
		u = utils.NewSysError(err)
	}

	if logger == nil {
		defaultLog()
	}

	logger.WithFields(logrus.Fields{}).Errorln(u, message)

	if !u.IsHasStack {
		var err error
		if len(message) == 0 {
			err = errors.WithStack(u)
		} else {
			err = errors.Wrap(u, message[0])
		}

		tmpErr := &utils.Error{
			ErrCode:    u.ErrCode,
			ErrMsg:     err,
			IsHasStack: true,
		}

		Errorf("this is a error stack: %+v", tmpErr.ErrMsg)
		u = tmpErr
	}

	return u
}
