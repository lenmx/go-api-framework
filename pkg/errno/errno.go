package errno

type Errno struct {
	Code    int
	Message string
}

func (err Errno) Error() string {
	return err.Message
}

func ConvertErr(err error) *Errno {
	code, message := DecodeErr(err)
	return &Errno{
		Code:    code,
		Message: message,
	}
}

func DecodeErr(err error) (code int, message string) {
	if err == nil {
		code = OK.Code
		message = OK.Message
		return
	}

	switch typed := err.(type) {
	case Errno:
		code = typed.Code
		message = typed.Message
		return
	default:
		code = InternalServerError.Code
		message = err.Error()
		return
	}
}
