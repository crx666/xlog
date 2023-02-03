package xlog

import (
	"fmt"
	"time"
)

var (
	TimeFormat = "2006-01-02 15:04:05"
	writer     = new(LogWriter)
	ConsoleLog = NewConsoleWriter(DebugLevel, TextEncodingType)
)

type LogEntryWithFields map[string]interface{}
type LogFields map[string]interface{}

func SetWriter(w Writer) {
	writer.SetWriter(w)
}

func GetWriter() Writer {
	w := writer.GetWriter()
	if w == nil {
		w = ConsoleLog
		SetWriter(w)
	}
	return w
}

func SetLevel(level string) {
	if writer.GetWriter() != nil {
		GetWriter().SetLevel(level)
	}
}

func Close() {
	w := GetWriter()
	if w != nil {
		w.Close()
	}
	CloseAllWrite()
}

func GetLevel() int {
	if writer.GetWriter() != nil {
		return GetWriter().GetLevel()
	}
	return DebugLevel
}

func Debug(v ...interface{}) {
	GetWriter().Debug(v...)
}

func DebugF(format string, fields ...interface{}) {
	GetWriter().DebugF(format, fields...)
}

func DebugW(format string, fields ...LogFields) {
	GetWriter().DebugW(format, Fields(fields...)...)
}

func Info(v ...interface{}) {
	GetWriter().Info(v...)
}

func InfoF(format string, fields ...interface{}) {
	GetWriter().InfoF(format, fields...)
}

func InfoW(format string, v ...LogFields) {
	GetWriter().InfoW(format, Fields(v...)...)
}

func Warn(v ...interface{}) {
	GetWriter().Warn(v...)
}

func WarnF(format string, fields ...interface{}) {
	GetWriter().WarnF(format, fields...)
}

func WarnW(format string, v ...LogFields) {
	GetWriter().WarnW(format, Fields(v...)...)
}

func Error(v ...interface{}) {
	GetWriter().Error(v...)
}

func ErrorF(format string, fields ...interface{}) {
	GetWriter().ErrorF(format, fields...)
}

func ErrorW(format string, v ...LogFields) {
	GetWriter().ErrorW(format, Fields(v...)...)
}

type LogField struct {
	Key   string
	Value interface{}
}

func Field(key string, value interface{}) LogField {
	switch val := value.(type) {
	case error:
		return LogField{Key: key, Value: val.Error()}
	case []error:
		var errs []string
		for _, err := range val {
			errs = append(errs, err.Error())
		}
		return LogField{Key: key, Value: errs}
	case time.Duration:
		return LogField{Key: key, Value: fmt.Sprint(val)}
	case []time.Duration:
		var durs []string
		for _, dur := range val {
			durs = append(durs, fmt.Sprint(dur))
		}
		return LogField{Key: key, Value: durs}
	case []time.Time:
		var times []string
		for _, t := range val {
			times = append(times, fmt.Sprint(t))
		}
		return LogField{Key: key, Value: times}
	case fmt.Stringer:
		return LogField{Key: key, Value: val.String()}
	case []fmt.Stringer:
		var strs []string
		for _, str := range val {
			strs = append(strs, str.String())
		}
		return LogField{Key: key, Value: strs}
	default:
		return LogField{Key: key, Value: val}
	}
}

func Fields(fields ...LogFields) []LogField {
	if len(fields) <= 0 {
		return nil
	}
	var values []LogField
	for _, l := range fields {
		for k, v := range l {
			field := Field(k, v)
			values = append(values, field)
		}
	}
	return values
}
