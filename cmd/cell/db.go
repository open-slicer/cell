package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var rdb *redis.Client
var pg *pgx.Conn

func dbConnect() {
	var err error

	pgURI := viper.GetString("database.postgres")
	pg, err = pgx.Connect(context.Background(), pgURI)
	if err != nil {
		log.Fatal().Err(err).Str("uri", pgURI).Msg("Failed to connect to Redis")
	}

	redisAddr := viper.GetString("database.redis.address")
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: viper.GetString("database.redis.password"),
		DB:       viper.GetInt("database.redis.db"),
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal().Err(err).Str("address", redisAddr).Msg("Failed to connect to Redis")
	}
}
