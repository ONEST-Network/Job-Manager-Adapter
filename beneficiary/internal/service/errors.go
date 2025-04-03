package service

import "errors"

var (
	ErrInvalidFilter       = errors.New("invalid filter parameters")
	ErrInvalidApplication  = errors.New("invalid application data")
	ErrSchemeNotFound      = errors.New("scheme not found")
	ErrApplicationNotFound = errors.New("application not found")
	ErrSchemeInactive      = errors.New("scheme is not active")
	ErrSchemeExpired       = errors.New("scheme has expired")
)
