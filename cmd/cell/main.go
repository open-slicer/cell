package main

import (
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var r = gin.Default()

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	environment := viper.GetString("environment")
	if environment != "release" {
		log.Logger = log.Level(zerolog.TraceLevel)
		log.Info().Msg("Environment isn't 'release'; using trace level")
	}

	viper.SetConfigName("cell")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't read config")
	}

	if dsn := viper.GetString("sentry.dsn"); dsn != "" {
		log.Debug().Err(err).Str("dsn", dsn).Msg("Initialising Sentry")
		err := sentry.Init(sentry.ClientOptions{
			Dsn: dsn,
		})
		if err != nil {
			log.Error().Err(err).Str("dsn", dsn).Msg("Initialising Sentry")
		}
	}

	mng = &mongoWrapper{
		uri: viper.GetString("database.mongodb"),
	}
	if err := mng.connect(); err != nil {
		log.Fatal().Err(err).Str("uri", mng.uri).Msg("Connecting to MongoDB")
	}

	redisAddr := viper.GetString("database.redis.address")
	if rdb, err = redisConnect(
		redisAddr, viper.GetString("database.redis.password"), viper.GetInt("database.redis.db"),
	); err != nil {
		log.Fatal().Err(err).Str("address", redisAddr).Msg("Failed to connect to Redis")
	}

	gin.SetMode(environment)
	setupRouter()

	addr := viper.GetString("http.address")
	if certFile := viper.GetString("security.cert_file"); certFile != "" {
		log.Info().Bool("tls", true).Str("addr", addr).Msg("Starting HTTP server with TLS")

		// Let's assume key_file is present.
		keyFile := viper.GetString("security.key_file")
		err = r.RunTLS(addr, certFile, keyFile)
		if err != nil {
			log.Fatal().Bool("tls", true).Str("addr", addr).Err(err).Msg("Failed to start HTTP server")
		}
	}

	log.Info().Bool("tls", false).Str("addr", addr).Msg("Starting HTTP server")
	err = r.Run(addr)
	if err != nil {
		log.Fatal().Err(err).Msg("Starting HTTP server")
	}
}
