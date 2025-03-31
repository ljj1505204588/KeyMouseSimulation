package main

import (
	uiWindows "KeyMouseSimulation/internal/ui/windows"
	"fmt"
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

	if err := uiWindows.MainWindows(); err != nil {
		fmt.Println(err.Error())
	}
	time.Sleep(10 * time.Second)
}
