package elogx

import (
	"sync"

	"errors"
)

var loggerMgr = NewLoggerManager()

type LoggerManager struct {
	sync.RWMutex
	LoggerInfo map[string]Writer
}

func NewLoggerManager() *LoggerManager {
	return &LoggerManager{
		LoggerInfo: make(map[string]Writer),
	}
}

func (l *LoggerManager) addWriter(mark string, w Writer) error {
	l.Lock()
	defer l.Unlock()
	if _, ok := l.LoggerInfo[mark]; !ok {
		l.LoggerInfo[mark] = w
		return nil
	}
	return errors.New("repeated logger name!")
}

func (l *LoggerManager) getWriter(mark string) Writer {
	l.RLock()
	defer l.RUnlock()
	return l.LoggerInfo[mark]
}

func (l *LoggerManager) close(mark string) {
	l.Lock()
	defer l.Unlock()
	if w, ok := l.LoggerInfo[mark]; ok {
		w.Close()
		delete(l.LoggerInfo, mark)
	}
}

func (l *LoggerManager) closeAll() {
	l.Lock()
	defer l.Unlock()
	for _, w := range l.LoggerInfo {
		w.Close()
	}
	l.LoggerInfo = make(map[string]Writer)
}

func RegisterWriter(mark string, w Writer) {
	loggerMgr.addWriter(mark, w)
}

func GetWriterInstance(mark string) Writer {
	return loggerMgr.getWriter(mark)
}

func CloseWrite(mark string) {
	loggerMgr.close(mark)
}

func CloseAllWrite() {
	loggerMgr.closeAll()
}
