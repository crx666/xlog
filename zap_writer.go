package xlog

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"time"

	"github.com/crx666/xlog/config"

	"github.com/crx666/xlog/common"

	"go.uber.org/zap/zapcore"
)

const (
	ThirdSkipOffset   = -1
	DefaultSkipOffset = 0
	CallerSkipOffset  = 2
)

var normalLevel = zap.NewAtomicLevel()
var errLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(fmt.Sprintf("[%s]", TimeFormat)))
}

func customJsonTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(fmt.Sprintf("%s", TimeFormat)))
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getJsonEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customJsonTimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

type ZapWriter struct {
	stackOffset int //默认输出为0
	encodeType  int
	opts        []zap.Option
	logger      *zap.Logger
}

func NewZapWriter(encodeType int, opts ...zap.Option) (Writer, error) {
	var logger *zap.Logger
	var err error

	if encodeType == JsonEncodingType {
		logger, err = zap.NewProduction(opts...)
		if err != nil {
			return nil, err
		}
	} else {
		logger, err = zap.NewDevelopment(opts...)
		if err != nil {
			return nil, err
		}
	}

	w := &ZapWriter{
		logger:      logger,
		encodeType:  encodeType,
		opts:        opts,
		stackOffset: DefaultSkipOffset,
	}
	return w, nil
}

func (w *ZapWriter) SetStackOffset(stackOffset int) {
	w.stackOffset = stackOffset
}

func (w *ZapWriter) SetLevel(level string) {
	lv, err := zapcore.ParseLevel(level)
	if err != nil {
		lv = zapcore.DebugLevel
	}
	normalLevel.SetLevel(lv)
}

func (w *ZapWriter) GetLevel() int {
	return LogLevel[normalLevel.String()]
}

func (w *ZapWriter) SetConfig(config *config.LogConfig) {
	if config == nil {
		return
	}
	err := common.LogConfigCheck(config)
	if err != nil {
		panic(err)
	}
	level, err := zapcore.ParseLevel(config.LogLevel)
	if err != nil {
		level = zapcore.DebugLevel
	}
	normalLevel.SetLevel(level)
	var encoder zapcore.Encoder
	if w.encodeType == JsonEncodingType {
		encoder = getJsonEncoder()
	} else {
		encoder = getEncoder()
	}

	cores := []zapcore.Core{}
	if config.IsConsole {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), level))
	}
	if config.LogDir != "" && config.LogName != "" {
		var info, warn io.Writer
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
		} else { //没有就默认创建一个日志文件
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

		var infoLevel zap.LevelEnablerFunc
		if config.ErrLogName != "" {
			infoLevel = func(lvl zapcore.Level) bool {
				return normalLevel.Level() <= lvl && lvl <= zapcore.WarnLevel
			}
		} else {
			infoLevel = func(lvl zapcore.Level) bool {
				return normalLevel.Level() <= lvl && lvl <= zapcore.FatalLevel
			}
		}

		if info != nil {
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(info), infoLevel)) //输出到日志文件
		}
		if warn != nil {
			//warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			//	return lvl >= zapcore.ErrorLevel
			//})
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(warn), errLevel)) //错误输出到日志文件
		}
	}
	core := zapcore.NewTee(
		cores...,
	)
	if !config.IsProd {
		w.opts = append(w.opts, zap.AddStacktrace(zap.WarnLevel))
	}
	if config.IsCall {
		w.opts = append(w.opts, zap.AddCaller())
		if w.stackOffset != 0 {
			w.opts = append(w.opts, zap.AddCallerSkip(CallerSkipOffset+w.stackOffset))
		} else {
			w.opts = append(w.opts, zap.AddCallerSkip(CallerSkipOffset))
		}
	}
	logger := zap.New(core, w.opts...)
	w.logger = logger
}

func (w *ZapWriter) Close() {
	w.logger.Sync()
}

func (w *ZapWriter) Error(v ...interface{}) {
	w.logger.Error(fmt.Sprint(v...))
}

func (w *ZapWriter) ErrorF(format string, fields ...interface{}) {
	w.logger.Error(fmt.Sprintf(format, fields...))
}

func (w *ZapWriter) ErrorW(format string, fields ...LogField) {
	w.logger.Error(format, toZapFields(fields...)...)
}

func (w *ZapWriter) Debug(v ...interface{}) {
	w.logger.Debug(fmt.Sprint(v...))
}

func (w *ZapWriter) DebugF(format string, fields ...interface{}) {
	w.logger.Debug(fmt.Sprintf(format, fields...))
}

func (w *ZapWriter) DebugW(format string, fields ...LogField) {
	w.logger.Debug(format, toZapFields(fields...)...)
}

func (w *ZapWriter) Info(v ...interface{}) {
	w.logger.Info(fmt.Sprint(v...))
}

func (w *ZapWriter) InfoF(format string, fields ...interface{}) {
	w.logger.Info(fmt.Sprintf(format, fields...))
}

func (w *ZapWriter) InfoW(format string, fields ...LogField) {
	w.logger.Info(format, toZapFields(fields...)...)
}

func (w *ZapWriter) Warn(v ...interface{}) {
	w.logger.Warn(fmt.Sprint(v...))
}

func (w *ZapWriter) WarnF(format string, fields ...interface{}) {
	w.logger.Warn(fmt.Sprintf(format, fields...))
}

func (w *ZapWriter) WarnW(format string, fields ...LogField) {
	w.logger.Warn(format, toZapFields(fields...)...)
}

func toZapFields(fields ...LogField) []zap.Field {
	if len(fields) <= 0 {
		return nil
	}
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
