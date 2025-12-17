package clients

import "errors"

var (
	ErrRequestFailure = errors.New("failed to form request to service")
	ErrSendingFailure = errors.New("failed to send request to service")
	ErrInvalidInput   = errors.New("invalid input")

	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrUserServiceUnavailable = errors.New("user service unavailable")
	ErrUserServiceTimeout     = errors.New("user service timeout")

	ErrProducerServiceUnavailable = errors.New("producer service unavailable")
	ErrConsumerServiceUnavailable = errors.New("producer service unavailable")
)
