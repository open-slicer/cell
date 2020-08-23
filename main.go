package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/gitlab"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var dbConn *pgx.Conn

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't read config")
	}

	databaseURL := viper.GetString("database.url")
	mgr, err := migrate.New(viper.GetString("database.migrations.source"), databaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't create migration instance")
	}

	err = mgr.Up()
	if err != nil {
		log.Warn().Err(err).Msg("Couldn't run migrations")
	}

	dbConn, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Opening database")
	}
	defer log.Err(dbConn.Close(context.Background())).Msg("Closing database")

	r := gin.Default()
	// Gin defaults to debug mode.
	if !viper.GetBool("debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	addr := viper.GetString("http.address")
	if certFile := viper.GetString("security.cert_file"); certFile != "" {
		// Let's assume key_file is present.
		keyFile := viper.GetString("security.key_file")
		log.Err(r.RunTLS(addr, certFile, keyFile)).Msg("Running HTTP server (tls)")
	}
	log.Err(r.Run(addr)).Msg("Running HTTP server")
}
