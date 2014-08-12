package log5go

import "bytes"

type compositeError struct {
	errs []error
}

func (e *compositeError) Error() string {
	if !e.hasErrors() {
		return ""
	} else if len(e.errs) == 1 {
		return e.errs[0].Error()
	}

	var buffer bytes.Buffer

	buffer.WriteString("Composite Error:\n")
	for _, err := range e.errs {
		buffer.WriteString(err.Error())
		buffer.WriteRune('\n')
	}

	return buffer.String()
}

func (e *compositeError) append(err error) {
	e.errs = append(e.errs, err)
}

func (e *compositeError) hasErrors() bool {
	return len(e.errs) > 0
}

func newCompositeError() *compositeError {
	return &compositeError{make([]error, 0, 0)}
}
