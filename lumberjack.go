package elogx

import (
	"strings"

	"xlog/common"

	"xlog/config"

	"xlog/lumberjack"
)

func GetLumberjackLogWriter(dir, file string, cfg *config.Lumberjack) LogFileWrite {
	if !strings.Contains(file, common.LogFormal) {
		file = file + common.LogTemp
	}
	logger := lumberjack.NewLumberjack(file, dir, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.SplitTime, cfg.Compress, true)
	return logger
}
