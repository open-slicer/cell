package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var previousLockets = map[string]bool{}

type locketInterface struct {
	Port int    `json:"port" binding:"required"`
	Host string `json:"host"`
}

type locketPUTResponse struct {
	Address  string          `json:"address"`
	Password string          `json:"password"`
	DB       int             `json:"db"`
	Locket   locketInterface `json:"locket"`
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
		context.Background(),
		"lockets",
		map[string]interface{}{
			ipAddr: fmt.Sprintf("%s:%d", locket.Host, locket.Port),
		},
	).Err()
	if err != nil {
		return internalError(err)
	}
	return response{
		Code:    http.StatusCreated,
		Message: "Locket added, expected to be ready",
		Data: locketPUTResponse{
			Address:  viper.GetString("database.redis.address"),
			Password: viper.GetString("database.redis.password"),
			DB:       viper.GetInt("database.redis.db"),
			Locket:   *locket,
		},
	}
}

func handleLocketsPUT(c *gin.Context) {
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

func (locket *locketInterface) getHostname() (string, error) {
	res, err := rdb.HGetAll(context.Background(), "lockets").Result()
	if err != nil {
		return "", err
	}

	for _, hostname := range res {
		if used, ok := previousLockets[hostname]; !ok || !used {
			previousLockets[hostname] = true
			return hostname, nil
		}
	}

	// If we haven't returned yet, all lockets were used.
	previousLockets = map[string]bool{}

	var genesisLocket string
	for _, hostname := range res {
		genesisLocket = hostname
		break
	}
	return genesisLocket, nil
}

func handleLocketsGET(c *gin.Context) {
	locket := locketInterface{}
	hostname, err := locket.getHostname()
	if err != nil {
		internalError(err).send(c)
		return
	}

	response{
		Code:    http.StatusOK,
		Message: "Locket rotated",
		Data:    hostname,
	}.send(c)
}
