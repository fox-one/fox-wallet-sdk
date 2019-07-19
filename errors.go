package sdk

import (
	"errors"
	"fmt"
)

var (
	// ErrInternalServerError code 2, internal server error
	ErrInternalServerError = errors.New("internal server error")
	// ErrInvalidBroker code 1500, invalid broker
	ErrInvalidBroker = errors.New("invalid broker")
	// ErrReplayRequest code 1501, replay request
	ErrReplayRequest = errors.New("replay request")
	// ErrInvalidUser code 1502, invalid user
	ErrInvalidUser = errors.New("invalid user")
	// ErrForbidden code 1503, forbiddenn
	ErrForbidden = errors.New("forbidden")
	// ErrInvalidPINToken code 1504, invalid pin token
	ErrInvalidPINToken = errors.New("invalid pin token")
	// ErrIncorrectPIN code 1505, incorrect pin
	ErrIncorrectPIN = errors.New("incorrect pin")
	// ErrIncorrectPINNonce code 1506, invalid pin nonce
	ErrIncorrectPINNonce = errors.New("invalid pin nonce")
	// ErrNotFound code 1507, not found
	ErrNotFound = errors.New("not found")
	// ErrInvalidRequest code 1508, invalid request
	ErrInvalidRequest = errors.New("invalid request")
	// ErrTooManyRequests code 1509, too many requests
	ErrTooManyRequests = errors.New("too many requests")
	// ErrInsufficientBalance code 1510, insufficient balance
	ErrInsufficientBalance = errors.New("insufficient balance")
	// ErrInsufficientFee code 1511, insufficient fee
	ErrInsufficientFee = errors.New("insufficient fee")
	// ErrAmountTooSmall code 1512, amount too small
	ErrAmountTooSmall = errors.New("amount too small")
	// ErrInvalidTrace code 1513, invalid trace id
	ErrInvalidTrace = errors.New("invalid trace id")
	// ErrAuthFailed auth failed
	ErrAuthFailed = errors.New("auth failed")
	// ErrInvalidAddress code 1515, invalid address
	ErrInvalidAddress = errors.New("invalid address")
)

func errorWithWalletError(err *Error) error {
	switch err.Code {
	case 2:
		return ErrInternalServerError
	case 1500:
		return ErrInvalidBroker
	case 1501:
		return ErrReplayRequest
	case 1502:
		return ErrInvalidUser
	case 1503:
		return ErrForbidden
	case 1504:
		return ErrInvalidPINToken
	case 1505:
		return ErrIncorrectPIN
	case 1506:
		return ErrIncorrectPINNonce
	case 1507:
		return ErrNotFound
	case 1508:
		return ErrInvalidRequest
	case 1509:
		return ErrTooManyRequests
	case 1510:
		return ErrInsufficientBalance
	case 1511:
		return ErrInsufficientFee
	case 1512:
		return ErrAmountTooSmall
	case 1513:
		return ErrInvalidTrace
	case 1514:
		return ErrAuthFailed
	case 1515:
		return ErrInvalidAddress
	}

	return err
}

const (
	// RequestFailed request failed
	RequestFailed = 3000000
)

// Error error
type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Hint string `json:"hint,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Msg)
}
