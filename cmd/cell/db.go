package main

import (
	"context"
	"github.com/spf13/viper"

	"github.com/go-redis/redis/v8"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var rdb *redis.Client
var mng *mongoWrapper

type mongoWrapper struct {
	uri string

	client       *mongo.Client
	mainDatabase *mongo.Database
	users        *mongo.Collection
}

func (d *mongoWrapper) connect() error {
	log.Debug().Str("uri", mng.uri).Msg("Connecting to MongoDB")
	var err error

	d.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(d.uri))
	if err != nil {
		return err
	}
	err = d.client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}

	d.mainDatabase = d.client.Database("cell")
	d.users = d.mainDatabase.Collection("users")
	return nil
}

func redisConnect(address, password string, db int) (*redis.Client, error) {
	log.Debug().Str("address", address).Msg("Connecting to Redis")

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	_, err := client.Ping(context.Background()).Result()

	return client, err
}

func dbConnect() {
	var err error

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
}
