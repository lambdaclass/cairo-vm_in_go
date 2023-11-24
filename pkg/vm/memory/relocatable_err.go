package memory

type SubRelocatableError struct {
	Msg string
}

func (e *SubRelocatableError) Error() string {
	return e.Msg
}
