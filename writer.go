package elogx

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"crx_log/config"

	"crx_log/common"
)

type Writer interface {
	Debug(v ...interface{})
	DebugF(format string, fields ...interface{})
	DebugW(format string, fields ...LogField)

	Info(v ...interface{})
	InfoF(format string, fields ...interface{})
	InfoW(format string, fields ...LogField)

	Warn(v ...interface{})
	WarnF(format string, fields ...interface{})
	WarnW(format string, fields ...LogField)

	Error(v ...interface{})
	ErrorF(format string, fields ...interface{})
	ErrorW(format string, fields ...LogField)

	SetLevel(level string)
	GetLevel() int

	SetConfig(cfg *config.LogConfig)
	SetStackOffset(offset int)
	Close()
}

type LogWriter struct {
	writer Writer
}

func (l *LogWriter) SetWriter(writer Writer) {
	if l.writer != nil {
		return
	}
	l.writer = writer
}

func (l *LogWriter) GetWriter() Writer {
	return l.writer
}

type LogFileWrite interface {
	io.Writer
	Exit() error
}

type NormalLogFile struct {
	File *os.File
}

func NewNormalLogFile(dir, name string) (*NormalLogFile, error) {
	var file *os.File
	var err error
	newDir := common.ReplaceDir(dir)
	err = os.MkdirAll(newDir, 0755)
	if err != nil {
		return nil, err
	}
	fileName := filepath.Join(newDir, common.ReplaceName(name))
	if !strings.Contains(fileName, common.LogFormal) {
		fileName = fileName + common.LogTemp
	}
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(fileName)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		if fileInfo.IsDir() {
			return nil, errors.New("file is dir")
		}
		file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			return nil, err
		}
	}

	return &NormalLogFile{
		File: file,
	}, nil
}

func (l *NormalLogFile) Write(p []byte) (n int, err error) {
	return l.File.Write(p)
}

func (l *NormalLogFile) Exit() error {
	return l.File.Close()
}
