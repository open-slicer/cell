package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var locketToken string

type locketInsertion struct {
}

func handleLocketPost(c *gin.Context) {
	locket := locketInsertion{}
	if err := c.ShouldBindJSON(&locket); err != nil {
		response{
			Code:    errorBindFailed,
			Message: "Failed to bind JSON",
			HTTP:    http.StatusBadRequest,
			Data:    err.Error(),
		}.send(c)
		return
	}
}

func locketAuthMiddleware(c *gin.Context) {
	if c.GetHeader("Authorization") != locketToken {
		response{
			Code:    errorInvalidLocketAuth,
			Message: "Authorization header didn't contain the required token (config.locket.token)",
			HTTP:    http.StatusUnauthorized,
		}.send(c)
		return
	}

	c.Next()
}
