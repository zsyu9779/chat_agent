package logger

import (
	"chat_agent/config"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

// 日志文件配置（日志切割文件）
type LogConfig struct {
	FileNamePrefix string
	LogFilePath    string
}

var Logger *logrus.Logger
var AccessLogger *logrus.Logger
var FuncLogger *logrus.Logger

func Init() {
	//文件初始化
	var logCnf config.LogStruct
	logCnf.Init()
	logConfig := LogConfig{
		FileNamePrefix: logCnf.Prefix,
		LogFilePath:    logCnf.FilePath,
	}
	//实例化--根据配置生成logrus-日志实例
	Logger = logrus.New()
	AccessLogger = logrus.New()
	FuncLogger = logrus.New()
	//logging.SetLevel(logrus.TraceLevel)
	switch config.Server.Mode {
	case config.ReleaseMode:
		Logger.SetLevel(logrus.InfoLevel)
		AccessLogger.SetLevel(logrus.InfoLevel)
	default:
		Logger.SetLevel(logrus.TraceLevel)
		AccessLogger.SetLevel(logrus.TraceLevel)
	}
	//设置日志格式
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	AccessLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	FuncLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	//切割写入日志
	Logger.AddHook(rotateLogs(logConfig, Logger))
	AccessLogger.AddHook(rotateAccessLogs(logConfig, AccessLogger))
	FuncLogger.AddHook(rotateFuncLogs(logConfig, FuncLogger))
}

// 为了少一层同时可以进行一些封装
func Info(args ...interface{}) {
	Logger.WithField("namespace", config.NamespaceMode).Info(args...)
}

func Infof(format string, args ...interface{}) {
	Logger.WithField("namespace", config.NamespaceMode).Infof(format, args...)
}
func Panic(args ...interface{}) {
	Logger.WithField("namespace", config.NamespaceMode).Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	Logger.WithField("namespace", config.NamespaceMode).Panicf(format, args...)
}

func Error(args ...interface{}) {
	function, location := caller()
	Logger.WithField("function", function).
		WithField("location", location).
		WithField("namespace", config.NamespaceMode).Error(args...)
}

func Errorf(format string, args ...interface{}) {
	function, location := caller()
	Logger.WithField("function", function).
		WithField("location", location).
		WithField("namespace", config.NamespaceMode).Errorf(format, args...)
}

func Warn(args ...interface{}) {
	Logger.WithField("namespace", config.NamespaceMode).Warn(args...)
}

func Debug(args ...interface{}) {
	Logger.WithField("namespace", config.NamespaceMode).Debug(args...)
}

func Trace(args ...interface{}) {
	Logger.WithField("namespace", config.NamespaceMode).Trace(args...)
}

func Fatal(args ...interface{}) {
	Logger.WithField("namespace", config.NamespaceMode).Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Logger.WithField("namespace", config.NamespaceMode).Fatalf(format, args...)
}

func caller() (string, string) {
	pc, file, line, _ := runtime.Caller(2)
	function := runtime.FuncForPC(pc).Name()                    // 获取函数名
	location := fmt.Sprintf("%s:%d", filepath.Base(file), line) // 获取文件名
	return function, location
}

// 日志切割配置
func rotateLogs(config LogConfig, logger *logrus.Logger) logrus.Hook {
	getWriterAndRotateLogs := func(level logrus.Level) io.Writer {
		file := path.Join(config.LogFilePath, fmt.Sprintf("%s_%s", config.FileNamePrefix, level.String()))
		writer, err := rotatelogs.New(
			file+"_%Y-%m-%d.log",
			// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
			rotatelogs.WithLinkName(file+".log"),
			// WithRotationTime设置日志分割的时间,这里设置为24小时分割一次
			rotatelogs.WithRotationTime(time.Hour*24),
			// WithMaxAge设置文件清理前的最长保存时间,
			rotatelogs.WithMaxAge(time.Hour*24*7),
		)
		if err != nil {
			logrus.Errorf("config local file system for logging error: %v", err)
		}
		return writer
	}
	//hook
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: getWriterAndRotateLogs(logrus.DebugLevel),
		logrus.InfoLevel:  getWriterAndRotateLogs(logrus.InfoLevel),
		logrus.WarnLevel:  getWriterAndRotateLogs(logrus.WarnLevel),
		logrus.ErrorLevel: getWriterAndRotateLogs(logrus.ErrorLevel),
		logrus.FatalLevel: getWriterAndRotateLogs(logrus.FatalLevel),
		logrus.PanicLevel: getWriterAndRotateLogs(logrus.PanicLevel),
	}, logger.Formatter)
	return lfsHook
}

func rotateAccessLogs(config LogConfig, logger *logrus.Logger) logrus.Hook {
	//access日志切分规则
	getWriterAndRotateLogs := func() io.Writer {
		file := path.Join(config.LogFilePath, fmt.Sprintf("%s_%s", config.FileNamePrefix, "access"))
		writer, err := rotatelogs.New(
			file+"_%Y-%m-%d.log",
			// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
			rotatelogs.WithLinkName(file+".log"),
			// WithRotationTime设置日志分割的时间,这里设置为24小时分割一次
			rotatelogs.WithRotationTime(time.Hour*24),
			// WithMaxAge设置文件清理前的最长保存时间,
			rotatelogs.WithMaxAge(time.Hour*24*7),
		)
		if err != nil {
			logrus.Errorf("config local file system for logging error: %v", err)
		}
		return writer
	}
	//hook
	lfsHook := lfshook.NewHook(getWriterAndRotateLogs(), logger.Formatter)
	return lfsHook
}

func rotateFuncLogs(config LogConfig, logger *logrus.Logger) logrus.Hook {
	//access日志切分规则
	getWriterAndRotateLogs := func() io.Writer {
		file := path.Join(config.LogFilePath, fmt.Sprintf("%s_%s", config.FileNamePrefix, "func"))
		writer, err := rotatelogs.New(
			file+"_%Y-%m-%d.log",
			// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
			rotatelogs.WithLinkName(file+".log"),
			// WithRotationTime设置日志分割的时间,这里设置为24小时分割一次
			rotatelogs.WithRotationTime(time.Hour*24),
			// WithMaxAge设置文件清理前的最长保存时间,
			rotatelogs.WithMaxAge(time.Hour*24*7),
		)
		if err != nil {
			logrus.Errorf("config local file system for logging error: %v", err)
		}
		return writer
	}
	//hook
	lfsHook := lfshook.NewHook(getWriterAndRotateLogs(), logger.Formatter)
	return lfsHook
}
