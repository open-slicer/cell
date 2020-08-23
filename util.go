package main

import "github.com/gin-gonic/gin"

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
