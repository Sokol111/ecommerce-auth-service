package main

import (
	"context"
	"github.com/Sokol111/ecommerce-auth-service/internal/config"
	"github.com/Sokol111/ecommerce-auth-service/internal/handler"
	"github.com/Sokol111/ecommerce-auth-service/internal/repository"
	"github.com/Sokol111/ecommerce-auth-service/internal/service"
	"github.com/Sokol111/ecommerce-commons/pkg/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	conf := config.LoadConfig("./configs")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+conf.DB.Username+":"+conf.DB.Password+"@"+conf.DB.Host+":"+strconv.Itoa(conf.DB.Port)))
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = client.Ping(ctxTimeout, nil)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()

	usersCollection := client.Database(conf.DB.DBName).Collection("users")

	repository.CreateIndexes(ctx, usersCollection)
	ur := repository.NewUserMongoRepository(usersCollection)
	us := service.NewUserService(ur)
	as := service.NewAuthService(ur, conf.SecretKey)
	uh := handler.NewUserHandler(us)
	ah := handler.NewAuthHandler(as)
	server.NewServer(conf.Port, ctx, uh, ah).Start()
}
