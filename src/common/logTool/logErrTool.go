package logTool

import "fmt"

type LogErrorT struct {
	Code    string
	Message string
	Err     error
}

func (err *LogErrorT) Error() string {
	var result string
	if err.Err != nil {
		result = fmt.Sprintf(" [%s] %s %s", err.Code, err.Message, err.Err.Error())
	} else {
		result = fmt.Sprintf("[%s] %s", err.Code, err.Message)
	}
	return result
}

func NewLogError(code, message string, err error) error {
	return &LogErrorT{Code: code, Message: message, Err: err}
}
func NewLogErrorAndLog(code, message string, err error) error {
	NewErr := &LogErrorT{Code: code, Message: message, Err: err}
	ErrorAJ(NewErr)
	return NewErr
}

const (
	NIL_ERROR_CODE = ""
)
