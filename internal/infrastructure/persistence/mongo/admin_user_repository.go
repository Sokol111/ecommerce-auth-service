package mongo

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type adminUserRepository struct {
	*commonsmongo.GenericRepository[adminuser.AdminUser, adminUserEntity]
}

func newAdminUserRepository(db commonsmongo.Mongo, mapper *adminUserMapper) (adminuser.Repository, error) {
	genericRepo, err := commonsmongo.NewGenericRepository(
		db.GetCollection("admin_user"),
		mapper,
	)
	if err != nil {
		return nil, err
	}

	return &adminUserRepository{
		GenericRepository: genericRepo,
	}, nil
}

func (r *adminUserRepository) FindByEmail(ctx context.Context, email string) (*adminuser.AdminUser, error) {
	opts := commonsmongo.QueryOptions{
		Filter: bson.D{{Key: "email", Value: email}},
		Page:   1,
		Size:   1,
	}

	result, err := r.FindWithOptions(ctx, opts)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, persistence.ErrEntityNotFound
	}

	return result.Items[0], nil
}

func (r *adminUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.ExistsWithFilter(ctx, bson.D{{Key: "email", Value: email}})
}

func (r *adminUserRepository) FindList(ctx context.Context, query adminuser.ListQuery) (*commonsmongo.PageResult[adminuser.AdminUser], error) {
	opts := commonsmongo.QueryOptions{
		Filter: buildListFilter(query),
		Page:   query.Page,
		Size:   query.Size,
		Sort:   bson.D{{Key: "createdAt", Value: -1}},
	}

	return r.FindWithOptions(ctx, opts)
}

func buildListFilter(query adminuser.ListQuery) bson.D {
	filter := bson.D{}

	if query.Role != nil {
		filter = append(filter, bson.E{Key: "role", Value: *query.Role})
	}
	if query.Enabled != nil {
		filter = append(filter, bson.E{Key: "enabled", Value: *query.Enabled})
	}
	if query.Search != nil && *query.Search != "" {
		searchRegex := bson.M{"$regex": *query.Search, "$options": "i"}
		filter = append(filter, bson.E{
			Key: "$or",
			Value: bson.A{
				bson.M{"email": searchRegex},
				bson.M{"firstName": searchRegex},
				bson.M{"lastName": searchRegex},
			},
		})
	}

	return filter
}

func (r *adminUserRepository) Insert(ctx context.Context, u *adminuser.AdminUser) error {
	err := r.GenericRepository.Insert(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return adminuser.ErrEmailAlreadyExists
	}
	return err
}
