package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var logger = &rotatingLogger{
	fileMux: new(sync.RWMutex),
}

type rotatingLogger struct {
	prefix      string
	file        *os.File
	fileMux     *sync.RWMutex
	formatter   *log.JSONFormatter
	masterTimer *time.Timer
}

func newRotatingLogger(prefix string) (*rotatingLogger, error) {
	fPath := filepath.Join(prefix, time.Now().Format("2006-01-02")+".log")
	f, err := os.OpenFile(fPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	l := &rotatingLogger{
		prefix:    prefix,
		file:      f,
		fileMux:   new(sync.RWMutex),
		formatter: new(log.JSONFormatter),
	}
	l.masterTimer = time.AfterFunc(untilMidnight(), func() {
		l.rotate(time.Now())
		l.masterTimer.Reset(untilMidnight())
	})
	return l, nil
}

func untilMidnight() time.Duration {
	ts := time.Now()
	midnight := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location()).AddDate(0, 0, 1)
	return midnight.Sub(ts)
}

func (l *rotatingLogger) rotate(now time.Time) {
	l.fileMux.Lock()
	defer l.fileMux.Unlock()
	if l.file == nil {
		// closed already
		return
	}

	fPath := filepath.Join(l.prefix, now.Format("2006-01-02")+".log")
	if f, err := os.Create(fPath); err == nil {
		if err := l.file.Close(); err != nil {
			log.Warningf("failed to close previous log file: %v", err)
		}
		l.file = f
	} else {
		log.Warningf("failed to open new log file: %v", err)
	}

}

func (l *rotatingLogger) Levels() []log.Level {
	return []log.Level{
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
		log.PanicLevel,
	}
}

func (l *rotatingLogger) Fire(entry *log.Entry) error {
	l.fileMux.Lock()
	defer l.fileMux.Unlock()
	if l.file == nil {
		// closed already
		return nil
	}
	line, err := l.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(l.file, string(line))
	return err
}

func (r *rotatingLogger) Write(p []byte) (n int, err error) {
	r.fileMux.Lock()
	defer r.fileMux.Unlock()
	if r.file == nil {
		return 0, io.EOF
	}
	return r.file.Write(p)
}

func (l *rotatingLogger) Close() error {
	l.fileMux.Lock()
	defer l.fileMux.Unlock()
	if l.file == nil {
		// closed already
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}
