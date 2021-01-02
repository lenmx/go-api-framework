package errno

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "ok"}
	InternalServerError = &Errno{Code: 10001, Message: "服务器内部错误"}
	ErrBind             = &Errno{Code: 10002, Message: "参数绑定失败"}
	ErrToken = &Errno{Code: 10003, Message: "身份鉴权失败"}
)
