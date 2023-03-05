package cmd

import "errors"

var (
	ErrMultipleOutputFlag = errors.New("multiple output flags are specified at once")
)
