package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
)

type (
	UserRepository interface {
		GetByUserName(c context.Context, username string) (*model.User, error)
		IsUserExists(c context.Context, username string) (bool, error)
		RegisterUser(c context.Context, user *model.User) error
		Create(c context.Context, user *model.User) error
	}
	// UserRepository provides database operations with User model
	userRepository struct {
		db *pg.DB
	}
)

func NewUserRepository(db *pg.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) GetByUserName(c context.Context, username string) (*model.User, error) {
	var user model.User
	if err := u.db.WithContext(c).Model(&user).Where("user_name = ?", username).First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find user")
	}
	return &user, nil
}

func (u *userRepository) IsUserExists(c context.Context, username string) (bool, error) {
	exists, err := u.db.WithContext(c).Model(&model.User{}).Where("user_name = ?", username).Exists()
	if err != nil {
		return false, utils.NotFound.Wrap(err, "can not find user")
	}
	return exists, err
}

func (u *userRepository) Create(c context.Context, user *model.User) error {
	if _, err := u.db.WithContext(c).Model(user).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not create user")
	}
	return nil

}
func (u *userRepository) RegisterUser(c context.Context, user *model.User) error {
	return u.db.RunInTransaction(c, func(tx *pg.Tx) error {
		tenant := model.Tenant{
			ID:            user.TenantID,
			Name:          user.FirstName,
			Configuration: nil,
		}
		if _, err := tx.Model(&tenant).Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not create tenant")
		}
		if _, err := tx.Model(user).Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not create user")
		}
		return nil
	})
}
