package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

var locketToken string

var previousLockets = map[string]bool{}

type locketInterface struct {
	Port int    `json:"port" binding:"required"`
	Host string `json:"host"`
}

func (locket *locketInterface) insert(ipAddr string) response {
	if locket.Host != "" {
		resolved := false

		ips, err := net.LookupIP(locket.Host)
		if err != nil {
			return response{
				Code:    errorDomainFailedLookup,
				Message: "The domain provided couldn't be looked up",
				HTTP:    http.StatusBadRequest,
			}
		}

		for _, v := range ips {
			if v.String() == ipAddr {
				resolved = true
				break
			}
		}
		if !resolved {
			return response{
				Code:    errorDomainDidntMatch,
				Message: "The domain provided didn't resolve to the client IP",
				HTTP:    http.StatusBadRequest,
			}
		}
	} else {
		locket.Host = ipAddr
	}

	err := rdb.HSet(
		context.Background(), "lockets", ipAddr, fmt.Sprintf("%s:%d", locket.Host, locket.Port),
	).Err()
	if err != nil {
		return internalError(err)
	}
	return response{
		Code:    http.StatusCreated,
		Message: "Locket added, expected to be ready",
		Data:    locket,
	}
}

func handleLocketPut(c *gin.Context) {
	locket := locketInterface{}
	if err := c.ShouldBindJSON(&locket); err != nil {
		response{
			Code:    errorBindFailed,
			Message: "Failed to bind JSON",
			HTTP:    http.StatusBadRequest,
			Data:    err.Error(),
		}.send(c)
		return
	}

	locket.insert(c.ClientIP()).send(c)
}

func (locket *locketInterface) get() response {
	res, err := rdb.HGetAll(context.Background(), "").Result()
	if err != nil {
		return internalError(err)
	}

	for _, hostname := range res {
		if used, ok := previousLockets[hostname]; !ok || !used {
			previousLockets[hostname] = true
			return response{
				Code:    http.StatusOK,
				Message: "Locket found",
				Data:    hostname,
			}
		}
	}

	// If we haven't returned yet, all lockets were used.
	previousLockets = map[string]bool{}

	var genesisLocket string
	for _, hostname := range res {
		genesisLocket = hostname
		break
	}
	return response{
		Code:    http.StatusOK,
		Message: "Locket found (reset)",
		Data:    genesisLocket,
	}
}

func handleLocketGet(c *gin.Context) {
	locket := locketInterface{}
	locket.get().send(c)
}

func locketAuthMiddleware(c *gin.Context) {
	if c.GetHeader("Authorization") != locketToken {
		response{
			Code:    errorInvalidLocketAuth,
			Message: "Authorization header didn't contain the required token (config.locket.token)",
			HTTP:    http.StatusUnauthorized,
		}.send(c)
		c.Abort()
		return
	}

	c.Next()
}
