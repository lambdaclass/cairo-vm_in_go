package memory

type SubReloctableError struct {
	Msg string
}

func (e *SubReloctableError) Error() string {
	return e.Msg
}
