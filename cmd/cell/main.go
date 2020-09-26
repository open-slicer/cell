package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var rdb *redis.Client
var pg *pgx.Conn

const epoch = 1577836800398

var useSentry = false

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	viper.SetConfigName("cell")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't read config")
	}

	environment := viper.GetString("environment")
	if environment != "release" {
		log.Logger = log.Level(zerolog.TraceLevel)
		log.Info().Msg("Environment isn't 'release'; using trace level")
	}

	if dsn := viper.GetString("sentry.dsn"); dsn != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: dsn,
		})
		if err != nil {
			log.Error().Err(err).Str("dsn", dsn).Msg("Initialising Sentry")
		} else {
			useSentry = true
		}
	}

	pgURI := viper.GetString("database.postgres")
	pg, err = pgx.Connect(context.Background(), pgURI)
	if err != nil {
		log.Fatal().Err(err).Str("uri", pgURI).Msg("Failed to connect to postgres")
	}

	redisAddr := viper.GetString("database.redis.address")
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: viper.GetString("database.redis.password"),
		DB:       viper.GetInt("database.redis.db"),
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal().Err(err).Str("address", redisAddr).Msg("Failed to connect to Redis")
	}
	gin.SetMode(environment)
	r := setupRouter()

	addr := viper.GetString("http.address")
	if certFile := viper.GetString("security.cert_file"); certFile != "" {
		log.Info().Bool("tls", true).Str("addr", addr).Msg("Starting HTTP server with TLS")

		// Let's assume key_file is present.
		keyFile := viper.GetString("security.key_file")
		go func() {
			if err := r.RunTLS(addr, certFile, keyFile); err != nil {
				log.Fatal().Bool("tls", true).Str("addr", addr).Err(err).Msg("Failed to start HTTP server")
			}
		}()
	}

	log.Info().Bool("tls", false).Str("addr", addr).Msg("Starting HTTP server")
	go func() {
		if err := r.Run(addr); err != nil {
			log.Fatal().Err(err).Msg("Starting HTTP server")
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs

	log.Info().Msg("Interrupt received, gracefully exiting")
	_ = pg.Close(context.Background())
	_ = rdb.Close()
	os.Exit(0)
}
