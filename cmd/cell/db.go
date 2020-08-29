package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const callTimeout = 4 * time.Second

var rdb *redis.Client
var mng *mongoWrapper

type mongoWrapper struct {
	uri string

	client       *mongo.Client
	mainDatabase *mongo.Database
	users        *mongo.Collection
}

func (d *mongoWrapper) connect() error {
	log.Debug().Str("uri", mng.uri).Dur("timeout", callTimeout).Msg("Connecting to MongoDB")
	var err error

	ctx, _ := context.WithTimeout(context.Background(), callTimeout)
	d.client, err = mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
	if err != nil {
		return err
	}

	ctx, _ = context.WithTimeout(context.Background(), callTimeout)
	err = d.client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	d.mainDatabase = d.client.Database("cell")
	d.users = d.mainDatabase.Collection("users")
	return nil
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
