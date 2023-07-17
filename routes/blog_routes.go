package routes

import (
	"database/sql"

	"github.com/fsimic346/go-blog/controller"
	"github.com/fsimic346/go-blog/repository"
	"github.com/fsimic346/go-blog/service"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func CreateBlogRoutes(db *sql.DB, jwtKey string, redisDb *redis.Client) *chi.Mux {

	blogRouter := chi.NewRouter()
	userRepository := repository.CreateUserRepository(db)
	blogRepository := repository.CreateBlogRepository(db)
	authenticator := Authenticator{
		JWTKey:  jwtKey,
		RedisDB: redisDb,
	}

	blogController := &controller.BlogController{
		BlogService: service.CreateBlogService(userRepository, blogRepository),
	}

	blogRouter.Post("/", authenticator.Authenticated(authenticator.IsAdmin(blogController.AddBlog, userRepository)))

	blogRouter.Get("/{blogId}", blogController.GetBlog)
	blogRouter.Get("/all", blogController.GetAll)

	return blogRouter
}
