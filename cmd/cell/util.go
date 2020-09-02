package main

import (
	"net/http"
	"regexp"
	"time"

	"github.com/spf13/viper"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const dayDuration = time.Hour * 24

type statusCode int

const (
	errorInternalError statusCode = iota
	errorExists
	errorPasswordInsecure
	errorBindFailed
	errorNotFound
	errorInvalidConfigToken
	errorDomainFailedLookup
	errorDomainDidntMatch
	errorParentNotExists
	errorNotCommonName
)

var commonNameRegex = regexp.MustCompile("^[A-Za-z0-9]+(?:[_-][A-Za-z0-9]+)*$")

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

type internalErrorData struct {
	EventID string `json:"event_id"`
}

func internalError(err error) response {
	eventID := captureException(err)
	return response{
		Code:    errorInternalError,
		Message: "An internal server error was encountered and it was recorded",
		HTTP:    http.StatusInternalServerError,
		Data: internalErrorData{
			EventID: eventID,
		},
	}
}

func captureException(err error) string {
	eventID := string(*sentry.CaptureException(err))
	log.Error().Str("id", eventID).Err(err).Msg("Exception captured")
	return eventID
}

func configAuthMiddleware(configPath string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != viper.GetString(configPath) {
			response{
				Code:    errorInvalidConfigToken,
				Message: "Authorization header didn't contain the value of config." + configPath,
				HTTP:    http.StatusUnauthorized,
			}.send(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
