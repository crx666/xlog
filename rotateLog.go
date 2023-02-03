package elogx

import (
	"strings"
	"time"

	"crx_log/config"

	"crx_log/common"

	"crx_log/rotatelogs"
)

func GetRotateLogWriter(dir, file string, cfg *config.Rotatelog) LogFileWrite {
	var ti time.Duration
	if !strings.Contains(file, "$ti") && !strings.Contains(file, "$day") && !strings.Contains(file, "$hour") && !strings.Contains(file, "$minute") {
		panic("split time log file need time format. please use $ti or $day or $hour or $minute of logName")
	}
	if !strings.Contains(file, common.LogFormal) {
		file = file + common.LogTemp
	}

	if cfg.SplitDay > 0 { //按天切分
		ti = time.Duration(cfg.SplitDay*24) * time.Hour
	} else if cfg.SplitHour > 0 { //按小时切分
		ti = time.Duration(cfg.SplitHour) * time.Hour
	} else {
		ti = time.Duration(cfg.SplitMinute) * time.Minute
	}

	options := []rotatelogs.Option{
		rotatelogs.WithRotationCount(cfg.MaxSave), // 文件最大保存份数
		rotatelogs.WithRotationTime(ti),           // 日志切割时间间隔
	}
	if cfg.LinkName != "" {
		options = append(options, rotatelogs.WithLinkName(cfg.LinkName)) // 生成软链，指向最新日志文件
	}
	hook, err := rotatelogs.New(
		file,
		dir,
		options...,
	)
	if err != nil {
		panic(err)
	}
	return hook
}
