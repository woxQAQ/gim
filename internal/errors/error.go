package errors

import "errors"

var (
	ErrUnauthenticated      = errors.New("unauthenticated")
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrPasswordIncorrect    = errors.New("password incorrect")
	ErrUserNotFound         = errors.New("user not found")
	ErrMessageNotRequest    = errors.New("message is not request")
	ErrCommonTransferTo     = errors.New("common request transfer to unique request failed")
	ErrTokenInvalid         = errors.New("token invalid")
	ErrAuthMessageNotFound  = errors.New("auth message not found")
	ErrTemp                 = errors.New("temp error")
)
