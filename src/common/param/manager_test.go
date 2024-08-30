package param

import (
	"fmt"
	"testing"
	"time"
)

type serverParam2T struct {
	ServerPort int  //`xml:"ServerPort"`
	OpenPProf  bool //`xml:"OpenPProf"`
	PProfPort  int  //`xml:"PProfPort"`
}

func TestParam(t *testing.T) {
	var s = &serverParam2T{}
	var s2 = &esParamT{}
	Manager.Register(s)
	Manager.Register(s2)

	var times int
	for {
		fmt.Println(times, *s)
		fmt.Println(times, *s2)
		times++
		time.Sleep(3 * time.Second)
	}

}

func (s *serverParam2T) ParamName() (name string) {
	return "server"
}

type esParamT struct {
	EsPort     []int
	UserName   string
	PassWord   string
	RecordBody bool
}

func (s *esParamT) ParamName() (name string) {
	return "esParamT"
}
