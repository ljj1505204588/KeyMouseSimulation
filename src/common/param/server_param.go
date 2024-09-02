package param

var ServerParam serverParamT

//func init() {
//	ServerParam.defaultParam()
//	Manager.Register(&ServerParam)
//}

type serverParamT struct {
	ServerPort int
	OpenPProf  bool
	PProfPort  int
}

func (s *serverParamT) defaultParam() {
	s.ServerPort = 21002
	s.PProfPort = 6061
}

func (s *serverParamT) ParamName() (name string) {
	return "serverParam"
}
