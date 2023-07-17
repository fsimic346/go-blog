package routes

import (
	"database/sql"

	"github.com/fsimic346/go-blog/controller"
	"github.com/fsimic346/go-blog/repository"
	"github.com/fsimic346/go-blog/service"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func CreateUserRoutes(db *sql.DB, jwtKey string, redisDb *redis.Client) *chi.Mux {

	userRouter := chi.NewRouter()
	userRepository := repository.CreateUserRepository(db)
	authenticator := Authenticator{
		JWTKey:  jwtKey,
		RedisDB: redisDb,
	}

	userController := &controller.UserController{
		UserService: service.CreateUserService(userRepository),
		JWTKey:      jwtKey,
		RedisDB:     redisDb,
	}

	userRouter.Post("/login", authenticator.NotAuthenticated(userController.Login))

	userRouter.Get("/{userId}", authenticator.Authenticated(userController.GetUser))
	userRouter.Get("/logout", authenticator.Authenticated(userController.Logout))

	userRouter.Post("/", userController.AddUser)

	return userRouter
}
