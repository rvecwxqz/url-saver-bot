package telegram

type UnknownTypeError struct {
	text string
}

func (e *UnknownTypeError) Error() string {
	return e.text
}

func NewUnknownTypeError() *UnknownTypeError {
	return &UnknownTypeError{
		text: "unknown message type",
	}
}

type UnknownMetaTypeError struct {
	text string
}

func (e *UnknownMetaTypeError) Error() string {
	return e.text
}

func NewUnknownMetaTypeError() *UnknownMetaTypeError {
	return &UnknownMetaTypeError{
		text: "unknown meta type",
	}
}
