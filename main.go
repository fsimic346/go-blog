package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/fsimic346/go-blog/routes"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {

	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	connStr := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := chi.NewRouter()

	setMiddleware(router)

	jwtKey := os.Getenv("JWT_KEY")

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PW"),
		DB:       0,
	})

	userRouter := routes.CreateUserRoutes(db, jwtKey, rdb)
	blogRouter := routes.CreateBlogRoutes(db, jwtKey, rdb)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	router.Mount("/user", userRouter)
	router.Mount("/blog", blogRouter)

	log.Println("Starting server on port: " + port)
	http.ListenAndServe(":"+port, router)
}
