package xlog

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
	"time"

	"github.com/crx666/xlog/config"

	"github.com/crx666/xlog/common"

	"errors"
)

type Mgs struct {
	Name string
	Age  int
}

var mgs = &Mgs{
	Name: "xxx", Age: 18,
}

func TestConfigs(t *testing.T) {
	c := new(config.RepeateConfig)
	common.ParserYamlData("./config/repeate_log.yaml", c)
	fmt.Println(c.Configs[0], c.Configs[1])
}

func TestZapSetLevel(t *testing.T) {
	c := common.GetBaseLogConfig()
	writer, err := NewZapWriter(JsonEncodingType)
	if err != nil {
		panic(err)
	}
	writer.SetConfig(c)
	SetWriter(writer)
	go func() {
		time.Sleep(3 * time.Second)
		writer.SetLevel("info")
		time.Sleep(3 * time.Second)
		writer.SetLevel("debug")
	}()
	for i := 0; i < 10; i++ {
		WarnW("zap warn", LogFields{"mgs": *mgs, "i": i})
		DebugW("zap debug", LogFields{"mgs": *mgs, "i": i})
		InfoW("zap info", LogFields{"mgs": *mgs, "i": i})
		Error("zap error", LogFields{"err": errors.New("err msg"), "i": i})
		time.Sleep(1 * time.Second)
	}
	writer.Close()
}

func TestZapLogger(t *testing.T) {
	c := common.GetBaseLogConfig()
	writer, err := NewZapWriter(JsonEncodingType)
	if err != nil {
		panic(err)
	}
	writer.SetConfig(c)
	SetWriter(writer)
	Debug("zap debug  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	DebugF("zap debugf. i am %s. i am %d years old", "super man", 18)
	DebugW("zap debugw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	Info("zap info  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	InfoF("zap infof. i am %s. i am %d years old", "super man", 18)
	InfoW("zap infow", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	Warn("zap warn  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	WarnF("zap warnf . i am %s. i am %d years old", "super man", 18)
	WarnW("zap warnfw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	Error("zap error  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	ErrorF("zap errorf. i am %s. i am %d years old", "super man", 18)
	ErrorW("zap errorw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	writer.Close()
}

func TestZapOffsetLogger(t *testing.T) {
	c := common.GetBaseLogConfig()
	w, err := NewZapWriter(JsonEncodingType)
	if err != nil {
		panic(err)
	}
	w.SetStackOffset(ThirdSkipOffset)
	w.SetConfig(c)
	w.Debug("zap debug  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	w.DebugF("zap debugf. i am %s. i am %d years old", "super man", 18)
	w.DebugW("zap debugw", Fields(LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})...)
	w.Info("zap info  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	w.InfoF("zap infof. i am %s. i am %d years old", "super man", 18)
	w.InfoW("zap infow", Fields(LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})...)
	w.Warn("zap warn  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	w.WarnF("zap warnf . i am %s. i am %d years old", "super man", 18)
	w.WarnW("zap warnfw", Fields(LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})...)
	w.Error("zap error  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	w.ErrorF("zap errorf. i am %s. i am %d years old", "super man", 18)
	w.ErrorW("zap errorw", Fields(LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})...)
	w.Close()
}

func TestLogrusSetLevel(t *testing.T) {
	c := common.GetBaseLogConfig()
	writer := NewLogrusWriter(func(logger *logrus.Logger) {
		logger.SetFormatter(JsonFormatter) //自定义输出格式
	})
	writer.SetConfig(c)
	SetWriter(writer)
	go func() {
		time.Sleep(3 * time.Second)
		writer.SetLevel("info")
		time.Sleep(3 * time.Second)
		writer.SetLevel("debug")
	}()
	for i := 0; i < 10; i++ {
		WarnW("zap warn", LogFields{"mgs": *mgs, "i": i})
		DebugW("zap debug", LogFields{"mgs": *mgs, "i": i})
		InfoW("zap info", LogFields{"mgs": *mgs, "i": i})
		ErrorW("zap error", LogFields{"err": errors.New("err msg"), "i": i})
		time.Sleep(1 * time.Second)
	}
	writer.Close()
}

func TestLogrusLogger(t *testing.T) {
	c := common.GetBaseLogConfig()
	//writer := NewLogrusWriter(c, func(logger *logrus.Logger) {
	//	logger.SetFormatter(JsonFormatter) //自定义输出格式
	//})

	writer := NewLogrusWriter(func(logger *logrus.Logger) { //默认 logurs 格式输出
		logger.SetFormatter(CustomizeFormatter) //自定义输出格式
	})
	writer.SetConfig(c)
	SetWriter(writer)
	Debug("logrus debug  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	DebugF("logrus debugf. i am %s. i am %d years old", "super man", 18)
	DebugW("logrus debugw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	Info("logrus info  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	InfoF("logrus infof. i am %s. i am %d years old", "super man", 18)
	InfoW("logrus infow", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	Warn("logrus warn  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	WarnF("logrus warnf . i am %s. i am %d years old", "super man", 18)
	WarnW("logrus warnfw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	Error("logrus error  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	ErrorF("logrus errorf. i am %s. i am %d years old", "super man", 18)
	ErrorW("logrus errorw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	writer.Close()
}

func TestLogSetLevel(t *testing.T) {
	writer := GetWriter()
	go func() {
		time.Sleep(3 * time.Second)
		writer.SetLevel("info")
		time.Sleep(3 * time.Second)
		writer.SetLevel("debug")
	}()
	for i := 0; i < 10; i++ {
		WarnW("zap warn", LogFields{"mgs": *mgs, "i": i})
		DebugW("zap debug", LogFields{"mgs": *mgs, "i": i})
		InfoW("zap info", LogFields{"mgs": *mgs, "i": i})
		ErrorW("zap error", LogFields{"err": errors.New("err msg"), "i": i})
		time.Sleep(1 * time.Second)
	}
	writer.Close()
}

func TestLogger(t *testing.T) { //golang 自带log包
	//SetLevel("info") //设置打印等级
	Debug("go debug  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	DebugF("go debugf. i am %s. i am %d years old", "super man", 18)
	DebugW("go debugw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	Info("go info  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	InfoF("go infof. i am %s. i am %d years old", "super man", 18)
	InfoW("go infow", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	Warn("go warn  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	WarnF("go warnf . i am %s. i am %d years old", "super man", 18)
	WarnW("go warnfw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	Error("go error  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
	ErrorF("go errorf. i am %s. i am %d years old", "super man", 18)
	ErrorW("go errorw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
}

func TestNewLoggerManager(t *testing.T) {
	c := common.GetBaseLogConfig()
	//w, err := NewZapWriter(c, JsonEncodingType,DefaultSkipOffset)
	//if err != nil {
	//	panic(err)
	//}
	//RegisterWriter("zap", w)

	c.LogName = c.LogName + "_logrus"
	c.ErrLogName = c.ErrLogName + "_logrus"
	w1 := NewLogrusWriter(func(logger *logrus.Logger) {
		logger.SetFormatter(CustomizeFormatter) //自定义输出格式
	})
	w1.SetConfig(c)
	RegisterWriter("logrus", w1)

	//w.DebugW("zap debug", Fields(LogFields{"mgs": mgs})...)
	//w.InfoW("zap info", Fields(LogFields{"mgs": mgs})...)
	//w.ErrorW("zap error", Fields(LogFields{"err": errors.New("err msg")})...)
	//CloseWrite("zap")

	//=============
	w1.Debug("logrus debug", mgs)
	w1.Info("logrus info", mgs)
	w1.Error("logrus error", errors.New("err msg"))

	w1.DebugW("logrus debug", Fields(LogFields{"mgs": mgs})...)
	w1.InfoW("logrus info", Fields(LogFields{"mgs": mgs})...)
	w1.ErrorW("logrus error", Fields(LogFields{"err": errors.New("err msg")})...)

	w1.DebugF("go debugf. i am %s. i am %d years old", "super man", 18)
	w1.InfoF("go infof. i am %s. i am %d years old", "super man", 18)
	w1.ErrorF("go errorf. i am %s. i am %d years old", "super man", 18)

	CloseAllWrite()

}

func TestLumberjack(t *testing.T) {
	c := common.GetBaseLogConfig() //配置设置为lumberjack
	w, err := NewZapWriter(JsonEncodingType)
	if err != nil {
		panic(err)
	}
	w.SetConfig(c)
	SetWriter(w)
	for i := 0; i < 70; i++ {
		for j := 0; j < 200; j++ {
			Debug("zap debug  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs, "i": i, "j": j})
			DebugF("zap debugf. i am %s. i am %d years old, i: %d,j: %d", "super man", 18, i, j)
			DebugW("zap debugw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs, "i": i, "j": j})
			Info("zap info  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs, "i": i, "j": j})
			InfoF("zap infof.i am %s. i am %d years old, i: %d,j: %d", "super man", 18, i, j)
			InfoW("zap infow", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs, "i": i, "j": j})
			Warn("zap warn  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs, "i": i, "j": j})
			WarnF("zap warnf . i am %s. i am %d years old, i: %d,j: %d", "super man", 18, i, j)
			WarnW("zap warnfw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs, "i": i, "j": j})
			Error("zap error  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs, "i": i, "j": j})
			ErrorF("zap errorf. i am %s. i am %d years old, i: %d,j: %d", "super man", 18, i, j)
			ErrorW("zap errorw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs, "i": i, "j": j})
		}
		time.Sleep(time.Second)
	}
	w.Close()
}

func TestRotateLog(t *testing.T) {
	c := common.GetBaseLogConfig() //配置设置为rotatelog
	w, err := NewZapWriter(JsonEncodingType)
	if err != nil {
		panic(err)
	}
	w.SetConfig(c)
	SetWriter(w)
	for i := 0; i < 70; i++ {
		Debug("zap debug  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
		DebugF("zap debugf. i am %s. i am %d years old", "super man", 18)
		DebugW("zap debugw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
		Info("zap info  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
		InfoF("zap infof. i am %s. i am %d years old", "super man", 18)
		InfoW("zap infow", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
		Warn("zap warn  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
		WarnF("zap warnf . i am %s. i am %d years old", "super man", 18)
		WarnW("zap warnfw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
		Error("zap error  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
		ErrorF("zap errorf. i am %s. i am %d years old", "super man", 18)
		ErrorW("zap errorw", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
		time.Sleep(time.Second)
	}
	w.Close()
}

func BenchmarkZap(b *testing.B) {
	c := common.GetBaseLogConfig()
	writer, err := NewZapWriter(JsonEncodingType)
	if err != nil {
		panic(err)
	}
	writer.SetConfig(c)
	//SetWriter(writer)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		writer.DebugW("zap debug", Fields(LogFields{"xxx": 1})...)
	}

}

func BenchmarkLogrus(b *testing.B) {
	c := common.GetBaseLogConfig()
	writer := NewLogrusWriter(func(logger *logrus.Logger) {
		logger.SetFormatter(JsonFormatter) //自定义输出格式
	})
	writer.SetConfig(c)
	//SetWriter(writer)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		writer.DebugW("zap debug  ", Fields(LogFields{"xxx": 1})...)
	} //writer.Close()
}

func BenchmarkLogger(b *testing.B) {
	c := common.GetBaseLogConfig()
	b.Run("Zap", func(b *testing.B) {
		writer, err := NewZapWriter(JsonEncodingType)
		if err != nil {
			panic(err)
		}
		writer.SetConfig(c)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				writer.Debug("zap debug  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
			}
		})
	})

	b.Run("Logrus", func(b *testing.B) {
		writer := NewLogrusWriter(func(logger *logrus.Logger) {
			logger.SetFormatter(JsonFormatter) //自定义输出格式
		})
		writer.SetConfig(c)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				writer.Debug("logrus debug  ", LogFields{"xxx": 1, "yyyy": 2, "zzz": mgs})
			}
		})
	})
}
