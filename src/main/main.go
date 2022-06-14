package main

import (
	"KeyMouseSimulation/common/logTool"
	ui "KeyMouseSimulation/module/UI"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	_ = logTool.SetLogLevel(logTool.LEVEL_DEBUG_STR)
	go func() {
		err := http.ListenAndServe(":6060", nil)
		if err != nil {
			panic(err.Error())
		}
	}()
	ui.MainWindows()
	time.Sleep(10 * time.Second)

}
