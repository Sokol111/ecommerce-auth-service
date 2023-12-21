package repository

//import (
//	"context"
//	"errors"
//	"github.com/Sokol111/auth-service/internal/mock"
//	"github.com/Sokol111/auth-service/internal/model"
//	"github.com/Sokol111/auth-service/internal/repository/db"
//	"github.com/jackc/pgx/v5/pgtype"
//	"github.com/stretchr/testify/assert"
//	"strings"
//	"testing"
//	"time"
//)
//
//func TestUserPostgresRepository_GetById(t *testing.T) {
//	now := time.Now()
//	tests := map[string]struct {
//		id            int64
//		userResult    db.User
//		errorResult   error
//		expectedUser  model.User
//		expectedError error
//	}{
//		"with result": {
//			id: int64(1),
//			userResult: db.User{
//				ID:               1,
//				Version:          1,
//				Enabled:          true,
//				Name:             "name",
//				CreatedDate:      pgtype.Timestamp{Time: now},
//				LastModifiedDate: pgtype.Timestamp{Time: now},
//			},
//			errorResult: nil,
//			expectedUser: model.User{
//				ID:               1,
//				Version:          1,
//				Enabled:          true,
//				Name:             "name",
//				CreatedDate:      now,
//				LastModifiedDate: now,
//			},
//			expectedError: nil,
//		},
//		"with error": {
//			id:            int64(1),
//			userResult:    db.User{},
//			errorResult:   errors.New("error"),
//			expectedUser:  model.User{},
//			expectedError: errors.New("error"),
//		},
//	}
//
//	for name, tc := range tests {
//		t.Run(name, func(t *testing.T) {
//			querier := mock.NewMockQuerier(t)
//			r := NewUserPostgresRepository(querier)
//			querier.EXPECT().
//				GetUserById(context.Background(), tc.id).
//				Return(tc.userResult, tc.errorResult)
//			result, err := r.GetById(context.Background(), tc.id)
//			assert.Equal(t, tc.expectedUser, result)
//			if tc.expectedError == nil {
//				assert.Equal(t, tc.expectedError, err)
//			} else {
//				assert.True(t, strings.Contains(err.Error(), tc.expectedError.Error()))
//			}
//		})
//	}
//}
