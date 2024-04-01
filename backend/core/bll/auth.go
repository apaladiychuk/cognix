package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/google/uuid"
)

type (
	AuthBL interface {
		Login(ctx context.Context, userName string) (*model.User, error)
		SignUp(ctx context.Context, identity *oauth.IdentityResponse) (*model.User, error)
	}
	authBL struct {
		userRepo repository.UserRepository
	}
)

func NewAuthBL(userRepo repository.UserRepository) AuthBL {
	return &authBL{
		userRepo: userRepo,
	}
}

func (a *authBL) Login(ctx context.Context, userName string) (*model.User, error) {
	user, err := a.userRepo.GetByUserName(ctx, userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *authBL) SignUp(ctx context.Context, identity *oauth.IdentityResponse) (*model.User, error) {
	exists, err := a.userRepo.IsUserExists(ctx, identity.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.InvalidInput.New("user already exists")
	}
	user := model.User{
		ID:         uuid.New(),
		TenantID:   uuid.New(),
		UserName:   identity.Email,
		FirstName:  identity.GivenName,
		LastName:   identity.FamilyName,
		ExternalID: identity.ID,
		Roles:      model.StringSlice{model.RoleSuperAdmin},
	}
	if user.FirstName == "" {
		user.FirstName = identity.Name
	}
	if err = a.userRepo.RegisterUser(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
