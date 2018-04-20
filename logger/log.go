package logger

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

//Logger is used for all logging
type Logger interface {
	Open(string) (*Log, error)
	Trace(...interface{})
	Info(...interface{})
	Warning(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})
	Write([]byte) (int, error)
	Printf(format string, v ...interface{})
	Size() int64
	GzipClose() error
	Close() error
}

//Log is used to log information to one or several files
type Log struct {
	levels Levels
	file   *os.File
	path   string
}

//New returns a newly initialized log
func New(path string) (*Log, error) {
	if path == "" {
		return nil, errors.New("file path must be provided")
	}

	//Create log file and open for writing
	l, err := new(Log).Open(path)
	if err != nil {
		return nil, err
	}
	l.path = path

	if l.file == nil {
		return nil, errors.New("file could not be opened")
	}

	//Set log levels and default log level
	l.levels = DefaultLevels(l.file)
	return l, nil
}

//Open opens the specified log file
func (l *Log) Open(path string) (*Log, error) {
	//Check if directory exists
	if dir, _ := splitFilepath(path); dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0766)
		}
	}

	//Open log file
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0766)
	if err != nil {
		return l, err
	}
	l.file = f

	return l, nil
}

//Size returns size of log file
func (l *Log) Size() int64 {
	fi, err := l.file.Stat()
	if err != nil {
		l.Fatal(err)
	}

	return fi.Size()
}

//Trace level log entry
func (l *Log) Trace(p ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		_, file = splitFilepath(file)
		l.levels[TRACE].Println(fmt.Sprintf("%s %d:", file, line), p)
		return
	}

	l.levels[TRACE].Println(p)
}

//Info level log entry
func (l *Log) Info(p ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		_, file = splitFilepath(file)
		l.levels[INFO].Println(fmt.Sprintf("%s %d:", file, line), p)
		return
	}

	l.levels[INFO].Println(p)
}

//Warning level log entry
func (l *Log) Warning(p ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		_, file = splitFilepath(file)
		l.levels[WARN].Println(fmt.Sprintf("%s %d:", file, line), p)
		return
	}

	l.levels[WARN].Println(p)
}

//Error level log entry
func (l *Log) Error(p ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		_, file = splitFilepath(file)
		l.levels[ERR].Println(fmt.Sprintf("%s %d:", file, line), p)
		return
	}

	l.levels[ERR].Println(p)
}

//Fatal level log entry
func (l *Log) Fatal(p ...interface{}) {
	defer l.Close()

	_, file, line, ok := runtime.Caller(1)
	if ok {
		_, file = splitFilepath(file)
		l.levels[FATAL].Println(fmt.Sprintf("%s %d:", file, line), p)
		return
	}

	l.levels[FATAL].Fatalln(p)
}

//Panic level log entry
func (l *Log) Panic(p ...interface{}) {
	defer l.Close()

	_, file, line, ok := runtime.Caller(1)
	if ok {
		path := strings.Split(file, "/")
		file = path[len(path)-1]
		l.levels[PANIC].Println(fmt.Sprintf("%s %d:", file, line), p)
		return
	}

	l.levels[PANIC].Panicln(p)
}

//Printf is similar to fmt printf
func (l *Log) Printf(format string, v ...interface{}) {
	var msg = fmt.Sprintf(format, v...)
	l.levels[INFO].Println(msg)
}

//Write ...
func (l *Log) Write(b []byte) (int, error) {
	n, err := l.file.Write(b)
	if err != nil {
		return n, err
	}

	err = l.file.Sync()
	return n, err
}

//GzipClose zips the old data before closing
func (l *Log) GzipClose() error {
	err := l.file.Close()
	if err != nil {
		return err
	}

	fp, err := os.Open(l.path)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	w := gzip.NewWriter(&buff)
	_, err = bufio.NewReader(fp).WriteTo(w)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	fp, err = os.Create(fmt.Sprintf("%s.gz", l.path))
	if err != nil {
		return err
	}

	_, err = buff.WriteTo(fp)
	if err != nil {
		return err
	}

	err = fp.Sync()
	if err != nil {
		return err
	}

	return os.Remove(l.path)
}

//Close log file
func (l *Log) Close() error {
	return l.file.Close()
}

func splitFilepath(path string) (string, string) {
	var dir, file string

	spl := strings.Split(path, "/")
	if len(spl) > 1 {
		file = spl[len(spl)-1]
		spl = spl[:len(spl)-1]
		dir = strings.Join(spl, "/")
	} else {
		file = spl[0]
	}

	return dir, file
}
