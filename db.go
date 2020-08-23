package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const callTimeout = 4 * time.Second

var db *database

type database struct {
	uri string

	c            *mongo.Client
	mainDatabase *mongo.Database
	users        *mongo.Collection
}

func (d *database) connect() error {
	log.Debug().Str("uri", db.uri).Dur("timeout", callTimeout).Msg("Connecting to MongoDB")
	var err error

	ctx, _ := context.WithTimeout(context.Background(), callTimeout)
	d.c, err = mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
	if err != nil {
		return err
	}

	ctx, _ = context.WithTimeout(context.Background(), callTimeout)
	err = d.c.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	d.mainDatabase = d.c.Database("cell")
	d.users = d.mainDatabase.Collection("users")
	return nil
}
