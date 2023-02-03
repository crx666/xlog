package elogx

const (
	JsonEncodingType = iota
	TextEncodingType
)

const (
	Flags            = 0x0
	PlainEncodingSep = '\t'
)

const (
	DebugLevel int = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

const (
	ContentKey   = "content"
	LevelKey     = "level"
	TimestampKey = "@timestamp"

	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelDebug = "debug"
	LevelPanic = "panic"
	LevelFatal = "fatal"
)

var LogLevel = map[string]int{
	LevelDebug: DebugLevel,
	LevelInfo:  InfoLevel,
	LevelWarn:  WarnLevel,
	LevelError: ErrorLevel,
	LevelPanic: PanicLevel,
	LevelFatal: FatalLevel,
}
