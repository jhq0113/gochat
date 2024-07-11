package protocol

import "errors"

var (
	ErrForbidden = errors.New(`forbidden`)
	ErrProtocol  = errors.New(`invalid protocol`)
)
