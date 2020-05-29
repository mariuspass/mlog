// +build linux windows
// +build amd64 arm

package main

import (
	"github.com/mariuspass/mlog"
)

var (
	log *mlog.Log
)

func main() {
	// get the log instance
	log = mlog.Get()

	// log to file log.txt
	log.SetFileWriter("log.txt")

	// Messages without formatting
	println("Messages without formatting:")
	log.Debug("Debug Message")
	log.Notice("Notice Message")
	log.Error("Error Message")
	log.Info("Info Message")
	log.Warning("Warning Message")
	// log.Critical("Critical Message")


	// Messages with formatting (using format 'verbs' from: https://golang.org/pkg/fmt/)
	println("")
	println("Messages with formatting (using format 'verbs' from: https://golang.org/pkg/fmt/):")
	log.Debug("%s", "Debug Message")
	log.Notice("%s", "Notice Message")
	log.Error("%s", "Error Message")
	log.Info("%s", "Info Message")
	log.Warning("%s", "Warning Message")
	log.Critical("%s", "Critical Error, Program will exit after this call")

}