package main

import (
	"KeyMouseSimulation/pkg/common/logTool"
	ui "KeyMouseSimulation/pkg/module/UI"
	_ "KeyMouseSimulation/pkg/module/server"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	// todo 整理锁关系
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
