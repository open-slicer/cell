package main

import (
	"context"
	"github.com/JakeMakesStuff/structuredhttp"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/url"
	"os"
	"path"
	"time"
)

var callTimeout = time.Second * 4

var rdb *redis.Client

type registration struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type genericResponse struct {
	Data interface{} `json:"data"`
}

type registrationResponse struct {
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
	if rdb, err = redisConnect(
		regData.Address, regData.Password, regData.DB,
	); err != nil {
		log.Fatal().Err(err).Str("address", regData.Address).Msg("Failed to connect to Redis")
	}
}

func register() registrationResponse {
	apiURL, err := url.Parse(viper.GetString("url"))
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid URL (config.url)")
	}
	apiURL.Path = path.Join(apiURL.Path, "api", "v2", "lockets")

	response, err := structuredhttp.PUT(apiURL.String()).JSON(registration{
		Port: viper.GetInt("registration.port"),
		Host: viper.GetString("registration.host"),
	}).Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Registration request threw an error")
	}

	err = response.RaiseForStatus()
	if err != nil {
		log.Fatal().Err(err).Msg("Received bad status code when registering")
	}
	respData, err := response.JSON()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't unmarshal response data")
	}

	return respData.(genericResponse).Data.(registrationResponse)
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
