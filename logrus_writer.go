package xlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"runtime"
	"strings"

	"github.com/crx666/xlog/config"

	"github.com/crx666/xlog/common"
)

const (
	Skip     = 8
	HookSkip = 3
)

var (
	logFieldMap *logrus.FieldMap = &logrus.FieldMap{
		logrus.FieldKeyTime:  "@time",
		logrus.FieldKeyLevel: "@lv",
		logrus.FieldKeyMsg:   "@msg",
		logrus.FieldKeyFunc:  "@func",
	}

	JsonFormatter *logrus.JSONFormatter = &logrus.JSONFormatter{
		FieldMap:        *logFieldMap,
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     false,
		DataKey:         "field",
	}

	TextFormatter *logrus.TextFormatter = &logrus.TextFormatter{
		FieldMap:        *logFieldMap,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableQuote:    true,
		DisableColors:   true,
		DisableSorting:  true,
	}

	CustomizeFormatter   = NewSimpleFormatter(Skip)
	PureMessageFormatter = new(pureMessageFormatter)
	PureFieldsFormatter  = new(pureFieldsFormatter)
)

type SimpleFormatter struct {
	Skip int
}

func (m *SimpleFormatter) getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}
	n := 0
	//获取执行代码的文件名
	for i := len(file) - 1; i > 0; i-- {
		if string(file[i]) == "/" {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}

func (m *SimpleFormatter) findCaller() string {
	file := ""
	line := 0
	skip := m.GetSkip()
	for i := 0; i < 10; i++ {
		file, line = m.getCaller(skip + i)
		if !strings.Contains(file, "logrus") && !strings.Contains(file, "log.go") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func NewSimpleFormatter(skip int) *SimpleFormatter {
	return &SimpleFormatter{Skip: skip}
}

func (m *SimpleFormatter) GetSkip() int {
	return m.Skip
}

func (m *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format(TimeFormat)
	var newLog string
	var value []byte
	if len(entry.Data) > 0 {
		value, _ = json.Marshal(entry.Data)
	}
	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := m.findCaller()
		if len(value) > 0 {
			newLog = fmt.Sprintf("[%s] [%s] %s %s %s\n",
				timestamp, entry.Level, fName, entry.Message, common.Bytes2String(value))
		} else {
			newLog = fmt.Sprintf("[%s] [%s] %s %s\n",
				timestamp, entry.Level, fName, entry.Message)
		}
	} else {
		if len(value) > 0 {
			newLog = fmt.Sprintf("[%s] [%s] %s %s\n", timestamp, entry.Level, entry.Message, common.Bytes2String(value))
		} else {
			newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
		}

	}
	b.WriteString(newLog)
	return b.Bytes(), nil
}

type pureMessageFormatter struct{}

func (p *pureMessageFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(entry.Message)
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

type pureFieldsFormatter struct{}

func (p *pureFieldsFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, err
	}
	return append(serialized, '\n'), nil
}

type LogrusWriter struct {
	logger      *logrus.Logger
	stackOffset int //默认输出为0
}

func NewLogrusWriter(opts ...func(logger *logrus.Logger)) Writer {
	logger := logrus.New()
	logger.ExitFunc = func(i int) {}
	for _, opt := range opts {
		opt(logger)
	}
	w := &LogrusWriter{
		logger:      logger,
		stackOffset: HookSkip,
	}
	return w
}

func (w *LogrusWriter) SetStackOffset(offset int) {
	w.stackOffset = offset
}

func (w *LogrusWriter) SetConfig(config *config.LogConfig) {
	if config == nil {
		return
	}
	err := common.LogConfigCheck(config)
	if err != nil {
		panic(err)
	}
	lv, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		lv = logrus.DebugLevel
	}
	w.logger.SetLevel(lv)
	formatter := w.logger.Formatter
	if !config.IsConsole {
		w.logger.SetOutput(io.Discard)
		w.logger.SetFormatter(TextFormatter) //不输出控制台 改成text格式 减少json内存分配
	}
	if config.IsCall {
		w.logger.SetReportCaller(true)
	}

	if config.LogDir != "" && config.LogName != "" {
		var info, warn LogFileWrite
		if config.Rotatelog != nil {
			info = GetRotateLogWriter(config.LogDir, config.LogName, config.Rotatelog)
			if config.ErrLogName != "" {
				warn = GetRotateLogWriter(config.LogDir, config.ErrLogName, config.Rotatelog)
			}
		} else if config.Lumberjack != nil {
			info = GetLumberjackLogWriter(config.LogDir, config.LogName, config.Lumberjack)
			if config.ErrLogName != "" {
				warn = GetLumberjackLogWriter(config.LogDir, config.ErrLogName, config.Lumberjack)
			}
		} else {
			info, err = NewNormalLogFile(config.LogDir, config.LogName)
			if err != nil {
				panic(err)
			}
			if config.ErrLogName != "" {
				warn, err = NewNormalLogFile(config.LogDir, config.ErrLogName)
				if err != nil {
					panic(err)
				}
			}
		}

		if warn == nil {
			warn = info
		}

		if info != nil {
			if _, ok := formatter.(*SimpleFormatter); ok {
				formatter = NewSimpleFormatter(Skip + w.stackOffset)
			}
			hook := lfshook.NewHook(lfshook.WriterMap{
				logrus.DebugLevel: info, // 为不同级别设置不同的输出目的
				logrus.InfoLevel:  info,
				logrus.WarnLevel:  info,
				logrus.ErrorLevel: warn,
				logrus.FatalLevel: warn,
				logrus.PanicLevel: warn,
			}, formatter)
			w.logger.AddHook(hook)
			logrus.RegisterExitHandler(func() {
				defer func() {
					if r := recover(); r != nil {
						log.Println("logrus.RegisterExitHandler error", r)
					}
				}()
				err := info.Exit()
				if err != nil {
					panic(err)
				}
				if config.ErrLogName != "" {
					err = warn.Exit()
					if err != nil {
						panic(err)
					}
				}
			})
		}
	}
}

func (w *LogrusWriter) Close() {
	w.logger.Exit(1)
}

func (w *LogrusWriter) SetLevel(level string) {
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		lv = logrus.DebugLevel
	}
	w.logger.SetLevel(lv)
}

func (w *LogrusWriter) GetLevel() int {
	return LogLevel[w.logger.GetLevel().String()]
}

func (w *LogrusWriter) Error(v ...interface{}) {
	w.logger.Error(fmt.Sprint(v...))
}

func (w *LogrusWriter) ErrorF(format string, fields ...interface{}) {
	w.logger.Errorf(format, fields...)
}

func (w *LogrusWriter) ErrorW(format string, fields ...LogField) {
	w.logger.WithFields(toLogrusFields(fields...)).Error(format)
}

func (w *LogrusWriter) Debug(v ...interface{}) {
	w.logger.Debug(fmt.Sprint(v...))
}

func (w *LogrusWriter) DebugF(format string, fields ...interface{}) {
	w.logger.Debugf(format, fields...)
}

func (w *LogrusWriter) DebugW(format string, fields ...LogField) {
	w.logger.WithFields(toLogrusFields(fields...)).Debug(format)
}

func (w *LogrusWriter) Info(v ...interface{}) {
	w.logger.Info(fmt.Sprint(v...))
}

func (w *LogrusWriter) InfoF(format string, fields ...interface{}) {
	w.logger.Infof(format, fields...)
}

func (w *LogrusWriter) InfoW(format string, fields ...LogField) {
	w.logger.WithFields(toLogrusFields(fields...)).Info(format)
}

func (w *LogrusWriter) Warn(v ...interface{}) {
	w.logger.Warn(fmt.Sprint(v...))
}

func (w *LogrusWriter) WarnF(format string, fields ...interface{}) {
	w.logger.Warnf(format, fields...)
}

func (w *LogrusWriter) WarnW(format string, fields ...LogField) {
	w.logger.WithFields(toLogrusFields(fields...)).Warn(format)
}

func toLogrusFields(fields ...LogField) logrus.Fields {
	if len(fields) <= 0 {
		return nil
	}
	logrusFields := make(logrus.Fields)
	for _, field := range fields {
		logrusFields[field.Key] = field.Value
	}
	return logrusFields
}
