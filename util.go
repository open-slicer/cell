package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const dayDuration = time.Hour * 24

type statusCode int

const (
	errorInternalError statusCode = iota
	errorExists
	errorPasswordInsecure
	errorTooLarge
	errorBindFailed
	errorMissingField
	errorDidntMatch
)

type tooLargeData struct {
	Offending []string `json:"offending"`
	Got       int      `json:"got"`
	Want      int      `json:"want"`
}

// response is a generic HTTP response. If HTTP is zeroed, Code should be used.
type response struct {
	Code    statusCode  `json:"code"`
	Message string      `json:"message"`
	HTTP    int         `json:"-"`
	Data    interface{} `json:"data,omitempty"`
}

func (r response) send(c *gin.Context) {
	if r.HTTP == 0 {
		c.JSON(int(r.Code), r)
		return
	}
	c.JSON(r.HTTP, r)
}

// someZero checks if a single value is zeroed. Supports: string, int.
func someZero(vals ...interface{}) bool {
	if vals[0] == nil {
		return true
	}
	check := func(val interface{}) bool {
		return true
	}

	switch vals[0].(type) {
	case string:
		check = func(val interface{}) bool {
			if val == "" {
				return true
			}
			return false
		}
	case int:
		check = func(val interface{}) bool {
			if val == 0 {
				return true
			}
			return false
		}
	}

	for _, val := range vals {
		if val == nil {
			return true
		}
		if check(val) {
			return true
		}
	}
	return false
}

type internalErrorData struct {
	EventID string `json:"event_id"`
}

func internalError(err error) response {
	eventID := string(*sentry.CaptureException(err))

	return response{
		Code:    errorInternalError,
		Message: "An internal server error was encountered and it was recorded",
		HTTP:    http.StatusInternalServerError,
		Data: internalErrorData{
			EventID: eventID,
		},
	}
}
