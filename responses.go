package main

type errorCode int

const (
	errorExists errorCode = iota
)

type errorResponse struct {
	Code    errorCode `json:"code"`
	Message string    `json:"message"`
}
