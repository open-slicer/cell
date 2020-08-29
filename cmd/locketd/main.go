package main

import (
	"context"
	"fmt"
	"github.com/JakeMakesStuff/structuredhttp"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"time"
)

var callTimeout = time.Second * 4

var rdb *redis.Client

type registration struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type registrationResponse struct {
	Data registrationResponseData `json:"data"`
}

type registrationResponseData struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	viper.SetConfigName("locketd")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't read config")
	}

	regData := register()
	if rdb, err = redisConnect(regData.Address, regData.Password, regData.DB); err != nil {
		log.Fatal().Err(err).Str("address", regData.Address).Msg("Failed to connect to Redis")
	}

	addr := fmt.Sprintf(":%d", viper.GetInt("port"))
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Str("address", addr).Msg("Failed to listen")
	}
	log.Info().Str("address", l.Addr().String()).Msg("Listening (tcp)")

	s := &http.Server{
		Handler:      websocketServer{},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Error().Err(err).Msg("Failed to serve")
	case <-sigs:
		log.Info().Msg("Interrupt received, gracefully exiting")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_ = s.Shutdown(ctx)
}

func register() registrationResponseData {
	apiURL, err := url.Parse(viper.GetString("registration.home"))
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid home URL")
	}
	apiURL.Path = path.Join(apiURL.Path, "api", "v2", "lockets")

	response, err := structuredhttp.
		PUT(apiURL.String()).
		Header("Authorization", viper.GetString("registration.token")).
		JSON(registration{
			Port: viper.GetInt("port"),
			Host: viper.GetString("registration.host"),
		}).Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Registration request threw an error")
	}
	err = response.RaiseForStatus()
	if err != nil {
		log.Fatal().Err(err).Msg("Received bad status code when registering")
	}

	var respData registrationResponse
	if err = response.JSONToPointer(&respData); err != nil {
		log.Fatal().Err(err).Msg("Couldn't unmarshal response data")
	}

	return respData.Data
}

func redisConnect(address, password string, db int) (*redis.Client, error) {
	log.Debug().Str("address", address).Dur("timeout", callTimeout).Msg("Connecting to Redis")

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	ctx, _ := context.WithTimeout(context.Background(), callTimeout)
	_, err := client.Ping(ctx).Result()

	return client, err
}
