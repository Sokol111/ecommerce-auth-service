package repository

import (
	"context"
	"fmt"
	"github.com/Sokol111/ecommerce-auth-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Version          int32
	Login            string
	HashedPassword   string `bson:"hashed_password"`
	Enabled          bool
	Permissions      []string  `bson:"permissions,omitempty"`
	CreatedDate      time.Time `bson:"created_date"`
	LastModifiedDate time.Time `bson:"last_modified_date"`
}

type UserRepository interface {
	GetById(ctx context.Context, id string) (model.User, error)

	GetByLogin(ctx context.Context, login string) (model.User, error)

	Create(ctx context.Context, user model.User) (model.User, error)

	Update(ctx context.Context, user model.User) (model.User, error)

	GetUsers(ctx context.Context) ([]model.User, error)
}

type UserMongoRepository struct {
	collection *mongo.Collection
}

func NewUserMongoRepository(collection *mongo.Collection) *UserMongoRepository {
	return &UserMongoRepository{collection}
}

func (r *UserMongoRepository) GetById(ctx context.Context, id string) (model.User, error) {
	convertedId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to convert user id [%v], reason: %v", id, err)
	}
	result := r.collection.FindOne(ctx, bson.D{{"_id", convertedId}})
	var u User
	if err := result.Decode(&u); err == nil {
		return toDomain(u), nil
	} else {
		return model.User{}, fmt.Errorf("failed to find user by id [%v], reason: %v", id, err)
	}
}

func (r *UserMongoRepository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	result := r.collection.FindOne(ctx, bson.D{{"login", login}})
	var u User
	if err := result.Decode(&u); err == nil {
		return toDomain(u), nil
	} else {
		return model.User{}, fmt.Errorf("failed to find user by login [%v], reason: %v", login, err)
	}
}

func (r *UserMongoRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	if result, err := r.collection.InsertOne(ctx, fromDomain(user)); err == nil {
		if id, ok := result.InsertedID.(primitive.ObjectID); ok {
			user.ID = id.Hex()
			return user, nil
		} else {
			return model.User{}, fmt.Errorf("failed to convert user id")
		}
	} else {
		return model.User{}, fmt.Errorf("failed to create user [%v], reason: %v", user, err)
	}
}

func (r *UserMongoRepository) Update(ctx context.Context, user model.User) (model.User, error) {
	id, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to convert user id [%v], reason: %v", user.ID, err)
	}
	result, err := r.collection.UpdateOne(ctx,
		bson.D{{"_id", id}, {"version", user.Version}},
		bson.D{{"$set", bson.D{{"login", user.Login}, {"enabled", user.Enabled}}}, {"$inc", bson.D{{"version", 1}}}})

	if err != nil {
		return model.User{}, fmt.Errorf("failed to update user [%v], reason: %v", user.ID, err)
	}

	if result.ModifiedCount != 1 {
		return model.User{}, fmt.Errorf("failed to update user [%v], reason: %v", user.ID, "id not found or versions mismatch")
	}

	user.Version++
	return user, nil
}

func (r *UserMongoRepository) GetUsers(ctx context.Context) ([]model.User, error) {
	cursor, err := r.collection.Find(ctx, bson.D{{}})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var result []User
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	s := make([]model.User, 0, len(result))
	for _, u := range result {
		s = append(s, toDomain(u))
	}
	return s, nil
}

func toDomain(u User) model.User {
	return model.User{ID: u.ID.Hex(), Version: u.Version, Enabled: u.Enabled, HashedPassword: u.HashedPassword, Login: u.Login, CreatedDate: u.CreatedDate, LastModifiedDate: u.LastModifiedDate}
}

func fromDomain(u model.User) User {
	user := User{Version: u.Version, Enabled: u.Enabled, Login: u.Login, HashedPassword: u.HashedPassword, Permissions: u.Permissions, CreatedDate: u.CreatedDate, LastModifiedDate: u.LastModifiedDate}
	return user
}
