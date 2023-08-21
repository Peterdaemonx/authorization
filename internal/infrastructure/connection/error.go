package connection

import "errors"

var (
	ErrRequestCanceled   = errors.New("request canceled")
	ErrRequestTimeout    = errors.New("request timeout")
	ErrRequestAbandoned  = errors.New("request abandoned")
	ErrNoFreeConnections = errors.New("no free connections in the pool")
)
