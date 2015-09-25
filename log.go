package httputil2

import (
    "github.com/albertjin/log2"
)

var critical = log2.Critical

var Logger = log2.NewStdLogger(nil)
var LogDebug = false

func log(a... interface{}) {
    Logger.Output(LogDebug, 1, a)
}

var stack = log2.StackLog
