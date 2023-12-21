package repository

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func CreateIndexes(ctx context.Context, usersCollection *mongo.Collection) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	indexModel := mongo.IndexModel{
		Keys: bson.D{{"login", 1}},
	}
	name, err := usersCollection.Indexes().CreateOne(ctxTimeout, indexModel)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create index")
	}
	log.Info().Msg("Index Created: " + name)
}
