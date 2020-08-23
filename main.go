package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

var r = gin.Default()

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	environment := viper.GetString("environment")
	if environment != "release" {
		log.Logger = log.Level(zerolog.TraceLevel)
		log.Debug().Msg("Environment isn't 'release'; using trace level")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't read config")
	}

	db = &database{
		uri: viper.GetString("database.uri"),
	}
	if err := db.connect(); err != nil {
		log.Fatal().Err(err).Str("uri", db.uri).Msg("Connecting to MongoDB")
	}

	gin.SetMode(environment)

	addr := viper.GetString("http.address")
	if certFile := viper.GetString("security.cert_file"); certFile != "" {
		// Let's assume key_file is present.
		keyFile := viper.GetString("security.key_file")
		log.Err(r.RunTLS(addr, certFile, keyFile)).Msg("Running HTTP server (tls)")
	}
	log.Err(r.Run(addr)).Msg("Running HTTP server")
}
