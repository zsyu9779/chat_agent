package error

type ServerError struct {
	err error
}

func (e ServerError) Error() string {
	return e.err.Error()
}

func NewServerError(err error) error {
	if err == nil {
		return nil
	}
	return ServerError{
		err: err,
	}
}
