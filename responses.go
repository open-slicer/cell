package main

type errorCode int

const (
	errorInternalError errorCode = iota
	errorExists
)

type errorResponse struct {
	Code    errorCode `json:"code"`
	Message string    `json:"message"`
	HTTP    int       `json:"-"`
}
