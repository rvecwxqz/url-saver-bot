package storage

type AlreadyExistsError struct {
	text string
}

func (e *AlreadyExistsError) Error() string {
	return e.text
}

func NewAlreadyExistsError() *AlreadyExistsError {
	return &AlreadyExistsError{
		text: "url already exists",
	}
}

type NoResultError struct {
	text string
}

func (e *NoResultError) Error() string {
	return e.text
}

func NewNoResultError() *NoResultError {
	return &NoResultError{
		text: "no result",
	}
}
