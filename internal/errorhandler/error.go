package errorhandler

import "errors"

var (
	ErrUnauthenticated      = errors.New("unauthenticated")
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrPasswordIncorrect    = errors.New("password incorrect")
	ErrUserNotFound         = errors.New("user not found")
	ErrMessageNotRequest    = errors.New("message is not request")
	ErrCommonTransferTo     = errors.New("common request transfer to unique request failed")
)
