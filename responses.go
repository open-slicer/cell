package main

type statusCode int

const (
	errorInternalError statusCode = iota
	errorExists
)

type response struct {
	Code    statusCode `json:"code"`
	Message string     `json:"message"`
	HTTP    int        `json:"-"`
}
