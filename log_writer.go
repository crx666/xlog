package xlog

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/crx666/xlog/config"
)

type logWriter struct {
	logger *log.Logger
}

func (lw logWriter) Write(data []byte) (int, error) {
	lw.logger.Print(string(data))
	return len(data), nil
}

func newLogWriter(logger *log.Logger) logWriter {
	return logWriter{
		logger: logger,
	}
}

type concreteWriter struct {
	infoLog     io.Writer
	errorLog    io.Writer
	level       int
	encode      int
	stackOffset int
}

func NewWriter(w io.Writer) Writer {
	lw := newLogWriter(log.New(w, "", Flags))
	return &concreteWriter{
		infoLog:  lw,
		errorLog: lw,
		level:    DebugLevel,
		encode:   JsonEncodingType,
	}
}

func NewConsoleWriter(lv int, encode int) Writer {
	outLog := newLogWriter(log.New(os.Stderr, "", Flags))
	errLog := newLogWriter(log.New(os.Stderr, "", Flags))
	return &concreteWriter{
		infoLog:  outLog,
		errorLog: errLog,
		level:    lv,
		encode:   encode,
	}
}

func (w *concreteWriter) Close() {

}
func (w *concreteWriter) SetStackOffset(stackOffset int) {
	w.stackOffset = stackOffset
}

func (w *concreteWriter) SetConfig(config *config.LogConfig) {

}

func (w *concreteWriter) SetLevel(level string) {
	if lv, ok := LogLevel[level]; ok {
		w.level = lv
	}
}

func (w *concreteWriter) GetLevel() int {
	return w.level
}

func (w *concreteWriter) checkLevel(levle string) bool {
	if lv, ok := LogLevel[levle]; ok {
		if w.level > lv {
			return false
		}
		return true
	}
	return false

}

func (w *concreteWriter) Error(v ...interface{}) {
	w.output(w.errorLog, LevelError, fmt.Sprint(v...))
}

func (w *concreteWriter) ErrorF(format string, fields ...interface{}) {
	w.output(w.errorLog, LevelError, fmt.Sprintf(format, fields...))
}

func (w *concreteWriter) ErrorW(format string, fields ...LogField) {
	w.output(w.errorLog, LevelError, format, fields...)
}

func (w *concreteWriter) Info(v ...interface{}) {
	w.output(w.infoLog, LevelInfo, fmt.Sprint(v...))
}

func (w *concreteWriter) InfoF(format string, fields ...interface{}) {
	w.output(w.infoLog, LevelInfo, fmt.Sprintf(format, fields...))
}

func (w *concreteWriter) InfoW(format string, fields ...LogField) {
	w.output(w.infoLog, LevelInfo, format, fields...)
}

func (w *concreteWriter) Debug(v ...interface{}) {
	w.output(w.infoLog, LevelDebug, fmt.Sprint(v...))
}

func (w *concreteWriter) DebugF(format string, fields ...interface{}) {
	w.output(w.infoLog, LevelDebug, fmt.Sprintf(format, fields...))
}

func (w *concreteWriter) DebugW(format string, fields ...LogField) {
	w.output(w.infoLog, LevelDebug, format, fields...)
}

func (w *concreteWriter) Warn(v ...interface{}) {
	w.output(w.infoLog, LevelWarn, fmt.Sprint(v...))
}

func (w *concreteWriter) WarnF(format string, fields ...interface{}) {
	w.output(w.infoLog, LevelWarn, fmt.Sprintf(format, fields...))

}

func (w *concreteWriter) WarnW(format string, fields ...LogField) {
	w.output(w.infoLog, LevelWarn, format, fields...)
}

func (w *concreteWriter) SetEncoding(t int) {
	if t != TextEncodingType && t != JsonEncodingType {
		panic("unknow encoding type")
	}
	w.encode = t
}

func (w *concreteWriter) output(writer io.Writer, level string, val interface{}, fields ...LogField) {
	if level != LevelError && !w.checkLevel(level) {
		return
	}
	switch w.encode {
	case TextEncodingType:
		writePlainAny(writer, level, val, buildFields(fields...)...)
	default:
		entry := make(LogEntryWithFields)
		for _, field := range fields {
			entry[field.Key] = field.Value
		}
		entry[TimestampKey] = getTimestamp()
		entry[LevelKey] = level
		entry[ContentKey] = val
		writeJson(writer, entry)
	}
}

func buildFields(fields ...LogField) []string {
	var items []string

	for _, field := range fields {
		items = append(items, fmt.Sprintf("%s=%v", field.Key, field.Value))
	}

	return items
}

func getTimestamp() string {
	return time.Now().Format(TimeFormat)
}

func writePlainAny(writer io.Writer, level string, val interface{}, fields ...string) {
	switch v := val.(type) {
	case string:
		writePlainText(writer, level, v, fields...)
	case error:
		writePlainText(writer, level, v.Error(), fields...)
	case fmt.Stringer:
		writePlainText(writer, level, v.String(), fields...)
	default:
		var buf strings.Builder
		buf.WriteString(getTimestamp())
		buf.WriteByte(PlainEncodingSep)
		buf.WriteString(level)
		buf.WriteByte(PlainEncodingSep)
		if err := json.NewEncoder(&buf).Encode(val); err != nil {
			log.Println(err.Error())
			return
		}

		for _, item := range fields {
			buf.WriteByte(PlainEncodingSep)
			buf.WriteString(item)
		}
		buf.WriteByte('\n')
		if writer == nil {
			log.Println(buf.String())
			return
		}

		if _, err := fmt.Fprint(writer, buf.String()); err != nil {
			log.Println(err.Error())
		}
	}
}

func writePlainText(writer io.Writer, level, msg string, fields ...string) {
	var buf strings.Builder
	buf.WriteString(getTimestamp())
	buf.WriteByte(PlainEncodingSep)
	buf.WriteString(level)
	buf.WriteByte(PlainEncodingSep)
	buf.WriteString(msg)
	for _, item := range fields {
		buf.WriteByte(PlainEncodingSep)
		buf.WriteString(item)
	}
	buf.WriteByte('\n')
	if writer == nil {
		log.Println(buf.String())
		return
	}

	if _, err := fmt.Fprint(writer, buf.String()); err != nil {
		log.Println(err.Error())
	}
}

func writeJson(writer io.Writer, info interface{}) {
	if content, err := json.Marshal(info); err != nil {
		log.Println(err.Error())
	} else if writer == nil {
		log.Println(string(content))
	} else {
		writer.Write(append(content, '\n'))
	}
}
