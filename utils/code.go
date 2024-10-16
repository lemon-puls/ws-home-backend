package utils

type Code int

const (
	CodeSuccess = iota
	CodeServerInternalError
	CodeInvalidParams
	CodeNotFound
)

var MsgMap = map[Code]string{
	CodeSuccess:             "success",
	CodeServerInternalError: "server internal error",
	CodeInvalidParams:       "invalid params",
	CodeNotFound:            "not found",
}

// code to msg
func (code Code) ToMsg() string {
	msg, ok := MsgMap[code]
	if !ok {
		msg = MsgMap[CodeServerInternalError]
	}
	return msg
}
