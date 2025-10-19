package errs

import "errors"

var ErrNotFound = errors.New("not found")

var ErrAlreadyExists = errors.New("entry already exists")

var ErrEmpty = errors.New("no entries exist")