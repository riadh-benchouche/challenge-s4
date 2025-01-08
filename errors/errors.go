package errors

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserNotActive = errors.New("user not active or confirmed")
var ErrInternal = errors.New("internal error")
var ErrNotFound = errors.New("not found")

// ErrCodeDoesNotExist var ErrBadRequest = errors.New("bad request")
// var ErrForbidden = errors.New("forbidden")
var ErrCodeDoesNotExist = errors.New("code does not exist")
var ErrInvalidCode = errors.New("invalid code")

// ErrAlreadyJoined var ErrEmailAlreadyUsed = errors.New("email already used")
var ErrAlreadyJoined = errors.New("user already joined association")
var ErrInvalidULIDFormat = errors.New("invalid ULID format")
var ErrAssociationNotFound = errors.New("Association not found")
var ErrEmailNotVerified = errors.New("email not verified")
var ErrInvalidToken = errors.New("invalid or expired token")
var ErrUserNotFound = errors.New("user not found")
var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidPassword = errors.New("invalid password provided")
