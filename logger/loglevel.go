package logfile

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

//Loglevel types
const (
	TRACE = "TRACE"
	INFO  = "INFO"
	WARN  = "WARNING"
	ERR   = "ERROR"
	FATAL = "FATAL"
	PANIC = "PANIC"
)

//Levels contain each log level
type Levels map[string]*log.Logger

//LogLevel contains info about a log level
type LogLevel struct {
	color     func(...interface{}) string
	colorText string
	prefix    string
}

//DefaultLevels initializes and returns the available log levels
func DefaultLevels(file *os.File) Levels {
	var levels = make(Levels)

	levels[TRACE] = trace().Build(file)
	levels[INFO] = info().Build(file)
	levels[WARN] = warning().Build(file)
	levels[ERR] = err().Build(file)
	levels[FATAL] = fatal().Build(file)
	levels[PANIC] = panic().Build(file)

	return levels
}

//Build creates individual logger
func (ll *LogLevel) Build(file *os.File) *log.Logger {
	return log.New(file,
		ll.Prefix(),
		log.LstdFlags)
}

//Prefix ...
func (ll *LogLevel) Prefix() string {
	if ll.prefix == "" {
		ll.prefix = fmt.Sprintf("[%s] ", ll.color(ll.colorText))
	}
	return ll.prefix
}

func trace() *LogLevel {
	return &LogLevel{
		color:     color.New(color.FgGreen).SprintFunc(),
		colorText: TRACE,
		prefix:    "",
	}
}

func info() *LogLevel {
	return &LogLevel{
		color:     color.New(color.FgBlue).SprintFunc(),
		colorText: INFO,
		prefix:    "",
	}
}

func warning() *LogLevel {
	return &LogLevel{
		color:     color.New(color.FgYellow).SprintFunc(),
		colorText: WARN,
		prefix:    "",
	}
}

func err() *LogLevel {
	return &LogLevel{
		color:     color.New(color.FgRed).SprintFunc(),
		colorText: ERR,
		prefix:    "",
	}
}

func fatal() *LogLevel {
	return &LogLevel{
		color:     color.New(color.FgMagenta).SprintFunc(),
		colorText: FATAL,
		prefix:    "",
	}
}

func panic() *LogLevel {
	return &LogLevel{
		color:     color.New(color.FgBlack, color.BgWhite).SprintFunc(),
		colorText: PANIC,
		prefix:    "",
	}
}
