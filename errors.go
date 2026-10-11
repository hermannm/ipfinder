package ipfinder

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// wrapErrors wraps the given errors with the given message.
//
// If there are multiple errors, the returned error implements the Unwrap() []error interface from
// the standard [errors] package.
func wrapErrors(msg string, causes []error) error {
	switch len(causes) {
	case 0:
		return errors.New(msg)
	case 1:
		return fmt.Errorf("%s: %w", msg, causes[0])
	default:
		return wrappedErrors{msg: msg, causes: causes}
	}
}

//nolint:errname // This is a generic error type, so we don't want a specific error name
type wrappedErrors struct {
	msg    string
	causes []error
}

func (err wrappedErrors) Error() string {
	var msg strings.Builder
	msg.WriteString(err.msg)
	msg.WriteString(": ")

	for i, cause := range err.causes {
		if i != 0 {
			msg.WriteString(", ")
		}
		msg.WriteString(strconv.Itoa(i + 1))
		msg.WriteString(") ")
		msg.WriteString(cause.Error())
	}

	return msg.String()
}

// Unwrap implements the wrapper interface from the standard library [errors] package.
func (err wrappedErrors) Unwrap() []error {
	return err.causes
}
