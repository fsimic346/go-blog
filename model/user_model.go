package model

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Id       string
	Username string
	Password string
	IsAdmin  bool
}

type DBUser struct {
	Id        string
	Username  string
	Password  string
	IsAdmin   bool
	CreatedAt string
	UpdatedAt string
}

type UserClaim struct {
	jwt.RegisteredClaims
	Username string
}

type UserRepository interface {
	GetById(id string) (User, error)
	GetByUsername(username string) (User, error)
	Add(username string, password string, isAdmin bool) (User, error)
}

type UserService interface {
	GetById(id string) (User, error)
	Add(username string, password string, isAdmin bool) (User, error)
	Login(username string, password string) error
}

func ConvertDBUserToUser(dbUser DBUser) User {
	return User{
		Id:       dbUser.Id,
		Username: dbUser.Username,
		Password: dbUser.Password,
		IsAdmin:  dbUser.IsAdmin,
	}
}
