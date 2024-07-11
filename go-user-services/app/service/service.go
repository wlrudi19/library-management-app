package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/wlrudi19/library-management-app/go-user-services/app/model"
	"github.com/wlrudi19/library-management-app/go-user-services/app/repository"
	"github.com/wlrudi19/library-management-app/go-user-services/infrastructure/middlewares"
	logger "github.com/wlrudi19/library-management-app/go-user-services/utils/log"
	"golang.org/x/crypto/bcrypt"
)

type UserLogic interface {
	FindUser(ctx context.Context, email string) (model.UserResponse, error)
	LoginUser(ctx context.Context, email string, password string) (model.LoginResponse, error)
}

type userlogic struct {
	UserRepository repository.UserRepository
}

func NewUserLogic(userRepository repository.UserRepository) UserLogic {
	return &userlogic{
		UserRepository: userRepository,
	}
}

func (l *userlogic) FindUser(ctx context.Context, email string) (model.UserResponse, error) {
	logCtx := "FindUser"
	logger.GetRequestLogEntry(ctx, logCtx, email).Info("FindUser Invoked...")

	var user model.UserResponse

	user, err := l.UserRepository.GetUserRedis(ctx, email)
	if err != nil {
		user, err := l.UserRepository.FindUser(ctx, email)
		if err != nil {
			logger.GetRequestLogEntry(ctx, logCtx, email).Error(fmt.Sprintf("error get redis: %v", err))
			return user, err
		}
		return user, nil
	}

	return user, nil
}

func (l *userlogic) LoginUser(ctx context.Context, email string, password string) (model.LoginResponse, error) {
	logCtx := "LoginUser"
	logger.GetRequestLogEntry(ctx, logCtx, email).Info("LoginUser Invoked...")

	var login model.LoginResponse

	user, err := l.FindUser(ctx, email)
	if err != nil {
		logger.GetRequestLogEntry(ctx, logCtx, email).Error(fmt.Sprintf("error find user: %v", err))
		return login, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.GetRequestLogEntry(ctx, logCtx, email).Error(fmt.Sprintf("error hash: %v", err))
		return login, errors.New("invalid password")
	}

	token, err := middlewares.GenerateAccessToken(user.Id, email)
	if err != nil {
		logger.GetRequestLogEntry(ctx, logCtx, email).Error(fmt.Sprintf("error gen token: %v", err))
		return login, err
	}

	login = model.LoginResponse{
		Id:          user.Id,
		AccessToken: token,
	}

	return login, nil
}
