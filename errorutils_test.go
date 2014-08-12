package log4go

import (
  "fmt"
  "testing"
)

func TestCompositeError(t *testing.T) {
  cerr := newCompositeError()

  if len(cerr.errs) != 0 {
    t.Error("new composite error has non-empty error slice")
  }
  if cerr.hasErrors() {
    t.Error("new compositer error reporting errors. expecting none.")
  }

  cerr.append(fmt.Errorf("test error"))
  if !cerr.hasErrors() {
    t.Error("expected hasErrors() to be true after appending error")
  }
  if len(cerr.errs) != 1 {
    t.Errorf("expected errs to hold 1 error but holds %d", len(cerr.errs))
  }

  if cerr.Error() != "test error" {
    t.Errorf("expected error message: 'test error', actual: '%s'", cerr.Error())
  }

  cerr.append(fmt.Errorf("another error"))
  if !cerr.hasErrors() {
    t.Error("expected hasErrors() to be true after appending 2 errors")
  }
  if len(cerr.errs) != 2 {
    t.Errorf("expected errs to hold 2 error but holds %d", len(cerr.errs))
  }

  if cerr.Error() != "Composite Error:\ntest error\nanother error\n" {
    t.Errorf("expected error message not present. actual: '%s'", cerr.Error())
  }
}
