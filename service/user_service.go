package service

import (
	"errors"

	"github.com/fsimic346/go-blog/model"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	UserRepository model.UserRepository
}

func CreateUserService(userRepository model.UserRepository) model.UserService {
	return &userService{
		UserRepository: userRepository,
	}
}

func (us *userService) GetById(userId string) (model.User, error) {
	return us.UserRepository.GetById(userId)
}

func (us *userService) Add(username, password string, isAdmin bool) (model.User, error) {
	return us.UserRepository.Add(username, password, isAdmin)
}

func (us *userService) Login(username string, password string) error {
	user, err := us.UserRepository.GetByUsername(username)
	if err != nil {
		return errors.New("user doesn't exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.New("incorrect password")
	}

	return nil
}
