package elegantExit

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var exitFuncList []func()

func AddElegantExit(h func()) {
	exitFuncList = append(exitFuncList, h)
}

func init() {
	// 监听系统退出信号
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigChan

		for _, h := range exitFuncList {
			h()
		}

		fmt.Println("优雅退出执行完毕.")
		os.Exit(0)
	}()
}
