package errors

type BaseError struct {
	msg string
}

func NewError(msg string) BaseError {
	return BaseError{
		msg: msg,
	}
}
func (e BaseError) Error() string {
	return e.msg
}
