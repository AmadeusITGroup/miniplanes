package mongo

type ConflictError struct {
}

func (a *ConflictError) Error() string {
	return "conflict"
}

type UnprocessableError struct {
}

func (u *UnprocessableError) Error() string {
	return "unprocessable object"
}
