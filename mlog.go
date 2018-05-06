package mlog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Log base struct
type Log struct {
	console bool
	writer  io.Writer
	mu      *sync.Mutex
}

var (
	log *Log
)

func init() {
	log = &Log{
		console: true,
		mu:      new(sync.Mutex),
	}
}

// Get returns singleton log instance
func Get() *Log {
	return log
}

// AddWriter adds a writer to log
func (l *Log) AddWriter(writer io.Writer) {
	if writer == nil {
		return
	}

	l.mu.Lock()
	if l.writer == nil {
		l.writer = writer
	} else {
		l.writer = io.MultiWriter(l.writer, writer)
	}
	l.mu.Unlock()
}

// SetFileWriter adds a *osFile as a writer
// path can be either relative or absolute
// if the file don't exist will be created
func (l *Log) SetFileWriter(path string) error {
	if path == "" {
		return errors.New("No file specified")
	}

	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}

	var file *os.File
	if _, err = os.Stat(path); err != nil {
		if _, err = os.Stat(filepath.Dir(path)); err != nil {
			err = os.MkdirAll(filepath.Dir(path), 0666)
			if err != nil {
				return err
			}
		}
		file, err = os.Create(path)
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
	}

	file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		return err
	}

	l.AddWriter(file)

	return nil
}

// EnableConsole enables the standard output
func (l *Log) EnableConsole() {
	l.mu.Lock()
	l.console = true
	l.mu.Unlock()
}

// DisableConsole disables the standard output
func (l *Log) DisableConsole() {
	l.mu.Lock()
	l.console = false
	l.mu.Unlock()
}

func (l *Log) out(level, format string, a []interface{}) {
	if l.console || l.writer != nil {
		nf := "%s [%s]: "
		s := ""
		date := time.Now().Format("02-01-2006 15:04:05")

		if len(a) == 0 {
			nf += "%v"
			s = fmt.Sprintf(nf, date, level, format)
		} else {
			nf += format
			na := append([]interface{}{date, level}, a...)

			s = fmt.Sprintf(nf, na...)
		}

		if l.console {
			l.mu.Lock()
			switch level {
			case "DEBG":
				color.HiCyan(s)
				break
			case "NOTI":
				color.HiGreen(s)
				break
			case "INFO":
				color.HiWhite(s)
				break
			case "WARN":
				color.HiYellow(s)
				break
			case "ERR ":
				color.HiRed(s)
				break
			case "CRIT":
				color.HiMagenta(s)
				break
			}
			l.mu.Unlock()
		}

		if !strings.HasSuffix(s, "\n") {
			s += "\n"
		}

		if l.writer != nil {
			l.mu.Lock()
			_, _ = fmt.Fprint(l.writer, s)
			l.mu.Unlock()
		}
	}
}

// Debug log
// Arguments are handled in the manner of fmt.Print
// a newline is appended if missing.
func (l *Log) Debug(format string, a ...interface{}) {
	l.out("DEBG", format, a)
}

// Notice log
// Arguments are handled in the manner of fmt.Print
// a newline is appended if missing.
func (l *Log) Notice(format string, a ...interface{}) {
	l.out("NOTI", format, a)
}

// Critical log and exit the program
// Arguments are handled in the manner of fmt.Print
// a newline is appended if missing.
func (l *Log) Critical(format string, a ...interface{}) {
	l.out("CRIT", format, a)

	time.Sleep(100 * time.Millisecond)

	l.Warning("The program will exit in:")

	for i := 5; i > 0; i-- {
		l.Warning("%v seconds", i)
		time.Sleep(1 * time.Second)
	}

	l.Info("Bye!")
	time.Sleep(1 * time.Second)
	os.Exit(1)
}

// Error log
// Arguments are handled in the manner of fmt.Print
// a newline is appended if missing.
func (l *Log) Error(format string, a ...interface{}) {
	l.out("ERR ", format, a)
}

// Info log
// Arguments are handled in the manner of fmt.Print
// a newline is appended if missing.
func (l *Log) Info(format string, a ...interface{}) {
	l.out("INFO", format, a)
}

// Warning log
// Arguments are handled in the manner of fmt.Print
// a newline is appended if missing.
func (l *Log) Warning(format string, a ...interface{}) {
	l.out("WARN", format, a)
}
