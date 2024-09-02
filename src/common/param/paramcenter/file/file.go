package paramcenter

import (
	"KeyMouseSimulation/common/commonTool"
	"encoding/json"
	"os"
	"strings"
)

type modData struct {
	Name string
	Data any
}

type fileDealT struct {
	path string
}

func (fd *fileDealT) init() {
	var err error
	if fd.path, err = os.Getwd(); err != nil {
		panic(err)
	}

	configPath := fd.path + commonTool.GetSysPthSep() + "config"
	_ = os.Mkdir(configPath, 0764)

	fd.path = strings.Join([]string{configPath, "config.json"}, commonTool.GetSysPthSep())
}

func (fd *fileDealT) Read() (data []modData, err error) {
	var text []byte
	if text, err = os.ReadFile(fd.path); err != nil {
		return
	}

	if err = json.Unmarshal(text, &data); err != nil {
		return
	}

	return
}

func (fd *fileDealT) WriteBack(data []modData) (err error) {
	var file *os.File
	if file, err = os.OpenFile(fd.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0762); err != nil {
		return
	}

	defer file.Close()

	var test, _ = json.MarshalIndent(data, "", "   ")

	if _, err = file.Write(test); err != nil {
		return
	}
	return
}
