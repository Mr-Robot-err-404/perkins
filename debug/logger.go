package debug

import (
	"log"
	"os"

	"github.com/Mr-Robot-err-404/perkins/core"
)

const logfile = "debug.log"

var logger *log.Logger

func Init() error {
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	logger = log.New(f, "", log.Ltime|log.Lshortfile)
	return nil
}

func Pos(pos core.Pos) {
	Logf("%d:%d\n", pos.Row, pos.Col)
}

func Log(msg string) {
	if logger == nil {
		return
	}
	logger.Println(msg)
}

func Logf(format string, args ...any) {
	if logger == nil {
		return
	}
	logger.Printf(format, args...)
}
