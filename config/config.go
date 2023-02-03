package config

type Rotatelog struct {
	MaxSave     int    `json:"max_save" yaml:"max_save"`         //文件最大保存分数
	SplitDay    int    `json:"split_day" yaml:"split_day"`       //日志切割时间  单位:天
	SplitHour   int    `json:"split_hour" yaml:"split_hour"`     //日志切割时间  单位:小时
	SplitMinute int    `json:"split_minute" yaml:"split_minute"` //日志切割时间  单位:分钟
	LinkName    string `json:"link_name" yaml:"link_name"`       //是否需要软连接
}

type Lumberjack struct {
	MaxSize    int  `json:"max_size" yaml:"max_size"`       //在进行切割之前，日志文件的最大大小（以 MB 为单位）
	MaxBackups int  `json:"max_backups" yaml:"max_backups"` //保留旧文件的最大个数
	MaxAge     int  `json:"max_age" yaml:"max_age"`         //保留旧文件的最大天数
	Compress   bool `json:"compress" yaml:"compress"`       //是否压缩 / 归档旧文件
	SplitTime  int  `json:"split_time" yaml:"split_time"`   //定时分割  单位:分钟
}

type LogConfig struct {
	LogDir     string      `json:"log_dir" yaml:"log_dir"`           //日志路径
	LogName    string      `json:"log_name" yaml:"log_name"`         //正常打印日志文件名字
	ErrLogName string      `json:"err_log_name" yaml:"err_log_name"` //错误日志文件名字  为空时代表 正常打印和错误打印在同一个文件
	LogLevel   string      `json:"log_level" yaml:"log_level"`       //日志打印等级 debug info
	IsProd     bool        `json:"is_prod" yaml:"is_prod"`           //是否正式服
	IsConsole  bool        `json:"is_console" yaml:"is_console"`     //是否控制台打印
	IsCall     bool        `json:"is_call" yaml:"is_call"`           //是否需要调用行数打印
	Rotatelog  *Rotatelog  `json:"rotatelog" yaml:"rotatelog"`       //按时间切分日志
	Lumberjack *Lumberjack `json:"lumberjack" yaml:"lumberjack"`     //按日志大小切分日志
	LogMark    string      `json:"log_mark" yaml:"log_mark"`         //日志标记
}

type RepeateConfig struct {
	Configs []*LogConfig `json:"configs" yaml:"configs"`
}
