package errors

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInternal = errors.New("internal error")
var ErrNotFound = errors.New("not found")
var ErrBadRequest = errors.New("bad request")
var ErrForbidden = errors.New("forbidden")
var ErrCodeDoesNotExist = errors.New("code does not exist")
var ErrInvalidCode = errors.New("invalid code")
var ErrEmailAlreadyUsed = errors.New("email already used")
var ErrAlreadyJoined = errors.New("user already joined association")
var ErrInvalidULIDFormat = errors.New("invalid ULID format")
var ErrAssociationNotFound = errors.New("Association not found")
