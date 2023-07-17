package repository

import (
	"database/sql"
	"time"

	"github.com/fsimic346/go-blog/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	db *sql.DB
}

func CreateUserRepository(db *sql.DB) model.UserRepository {
	return &userRepository{
		db: db}
}

func (ur *userRepository) GetById(id string) (model.User, error) {
	var dbUser model.DBUser
	rows, err := ur.db.Query("SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		return model.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&dbUser.Id, &dbUser.Username, &dbUser.Password, &dbUser.IsAdmin, &dbUser.CreatedAt, &dbUser.UpdatedAt)
		if err != nil {
			return model.User{}, err
		}
	}

	user := model.ConvertDBUserToUser(dbUser)

	return user, nil
}

func (ur *userRepository) GetByUsername(username string) (model.User, error) {
	var dbUser model.DBUser
	rows, err := ur.db.Query("SELECT * FROM users WHERE username=$1", username)
	if err != nil {
		return model.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&dbUser.Id, &dbUser.Username, &dbUser.Password, &dbUser.IsAdmin, &dbUser.CreatedAt, &dbUser.UpdatedAt)
		if err != nil {
			return model.User{}, err
		}
	}

	user := model.ConvertDBUserToUser(dbUser)

	return user, nil
}

func (ur *userRepository) Add(username, password string, isAdmin bool) (model.User, error) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		Id:       uuid.NewString(),
		Username: username,
		Password: string(hashedPw),
		IsAdmin:  isAdmin,
	}

	_, err = ur.db.Exec("INSERT INTO users VALUES($1,$2,$3,$4,$5,$6)", user.Id, user.Username, user.Password, user.IsAdmin, time.Now().UTC(), time.Now().UTC())

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
