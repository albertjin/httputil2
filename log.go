package httputil2

import (
    "os"

    "github.com/albertjin/log2"
)

var critical = log2.Critical

var Logger = log2.NewStdLogger(nil)
var LogDebug = (os.Getenv("httputil2_log") == "1")

func log(a... interface{}) {
    Logger.Output(LogDebug, 1, a)
}

var stack = log2.StackLog
