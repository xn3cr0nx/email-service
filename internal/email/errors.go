package email

import "errors"

var (
	errServiceInitialization = errors.New("injected service not correcly initialized")
	errInvalidFrom           = errors.New("invalid from parameter")
	errInvalidTo             = errors.New("invalid to parameter")
	errInvalidSubject        = errors.New("invalid subject parameter")
	errInvalidName           = errors.New("invalid name parameter for welcome email")
	errInvalidURL            = errors.New("invalid URL parameter for welcome email")
	errTemplateNotFound      = errors.New("cannot find template path using task type")
)
