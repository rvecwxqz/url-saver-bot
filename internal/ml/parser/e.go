package parser

type NoDataError struct {
}

func (e *NoDataError) Error() string {
	return "no data"
}

func NewNoDataError() error {
	return &NoDataError{}
}
