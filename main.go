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

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
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

	db = &database{
		uri: viper.GetString("database.uri"),
	}
	if err := db.connect(); err != nil {
		log.Fatal().Err(err).Str("uri", db.uri).Msg("Connecting to MongoDB")
	}

	gin.SetMode(environment)
	setupRoutes()

	addr := viper.GetString("http.address")
	if certFile := viper.GetString("security.cert_file"); certFile != "" {
		log.Debug().Bool("tls", true).Str("addr", addr).Msg("Running HTTP server")

		// Let's assume key_file is present.
		keyFile := viper.GetString("security.key_file")
		err = r.RunTLS(addr, certFile, keyFile)
		if err != nil {
			log.Fatal().Bool("tls", true).Str("addr", addr).Err(err).Msg("Running HTTP server")
		}
	}

	log.Debug().Bool("tls", false).Str("addr", addr).Msg("Running HTTP server")
	err = r.Run(addr)
	if err != nil {
		log.Fatal().Err(err).Msg("Running HTTP server")
	}
}
