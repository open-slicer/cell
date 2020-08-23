package main

type statusCode int

const (
	errorInternalError statusCode = iota
	errorExists
)

// response is a generic HTTP response. If HTTP is zeroed, Code will be used.
type response struct {
	Code    statusCode `json:"code"`
	Message string     `json:"message"`
	HTTP    int        `json:"-"`
}
