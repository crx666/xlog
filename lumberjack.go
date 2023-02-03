package elogx

import (
	"strings"

	"crx_log/common"

	"crx_log/config"

	"crx_log/lumberjack"
)

func GetLumberjackLogWriter(dir, file string, cfg *config.Lumberjack) LogFileWrite {
	if !strings.Contains(file, common.LogFormal) {
		file = file + common.LogTemp
	}
	logger := lumberjack.NewLumberjack(file, dir, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.SplitTime, cfg.Compress, true)
	return logger
}
