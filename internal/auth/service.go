package auth

import (
	"errors"
	"go/adv-demo/internal/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepo *user.UserRepository) *AuthService {
	return &AuthService{UserRepository: userRepo}
}

func (service *AuthService) RegisterUser(email, password, name string) (string, error) {
	existedUser, _ := service.UserRepository.GetByEmail(email)

	if existedUser != nil {
		return "", errors.New(ErrUserAlreadyExists)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	user := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	result, err := service.UserRepository.Create(user)

	if err != nil {
		return "", err
	}

	return result.Email, nil
}

func (service *AuthService) AuthenticateUser(email, password string) (string, error) {
	user, err := service.UserRepository.GetByEmail(email)

	if err != nil {
		return "", errors.New(ErrUserAlreadyExists)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}

	return user.Email, nil
}
