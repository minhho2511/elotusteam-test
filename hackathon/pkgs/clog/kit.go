package clog

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"os"
	"time"
)

const (
	FileLogger  = "file"
	DepthCaller = 5
)

type gokitLog struct {
	logger log.Logger
}

func NewGoKitLog() Writer {
	l := gokitLog{}
	l.logger = l.stdLogger()
	l.logger = level.NewFilter(l.logger, level.AllowAll())
	l.logger = log.With(l.logger, "ts", l.timeNow(), "caller", log.Caller(DepthCaller))
	return &l
}

func (l *gokitLog) fileLogger(filePath string) log.Logger {
	logfile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	defer func(logfile *os.File) {
		_ = logfile.Close()
	}(logfile)

	return log.NewLogfmtLogger(log.NewSyncWriter(logfile))
}

func (l *gokitLog) stdLogger() log.Logger {
	return log.NewLogfmtLogger(os.Stderr)
}

func (l *gokitLog) timeNow() log.Valuer {
	return log.TimestampFormat(func() time.Time { return time.Now() }, time.DateTime)
}

func (l *gokitLog) Printf(log *LogCollection) {
	switch log.Level {
	case DebugLevel:
		_ = level.Debug(l.logger).Log("debug", log.Message)
	case InfoLevel:
		_ = level.Info(l.logger).Log("info", log.Message)
	case WarnLevel:
		_ = level.Warn(l.logger).Log("warn", log.Message)
	case ErrorLevel:
		_ = level.Error(l.logger).Log("err", log.Message)
	}
}

func (l *gokitLog) Log(keyvals ...interface{}) error {
	return l.logger.Log(keyvals...)
}
