package utils

import (
	"github.com/pkg/errors"
)

func ErrorWithCause(err error, cause string) error {
	if errors.Cause(err) != err {
		return errors.WithMessage(err, cause)
	}
	return errors.Wrap(err, cause)
}
