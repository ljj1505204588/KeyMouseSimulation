package main

import (
	uiWindows "KeyMouseSimulation/internal/ui/windows"
	_ "KeyMouseSimulation/pkg/windows"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	go func() {
		err := http.ListenAndServe(":6060", nil)
		if err != nil {
			panic(err.Error())
		}
	}()

	uiWindows.MainWindows()
	time.Sleep(10 * time.Second)
}
