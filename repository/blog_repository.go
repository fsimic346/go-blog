package repository

import (
	"database/sql"
	"time"

	"github.com/fsimic346/go-blog/model"
	"github.com/google/uuid"
)

type blogRepository struct {
	db *sql.DB
}

func CreateBlogRepository(db *sql.DB) model.BlogRepository {
	return &blogRepository{db: db}
}

func (br *blogRepository) GetById(id string) (model.DBBlog, error) {
	var blog model.DBBlog
	rows, err := br.db.Query("SELECT * FROM blogs WHERE id=$1", id)
	if err != nil {
		return model.DBBlog{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&blog.Id, &blog.Title, &blog.Content, &blog.AuthorId, &blog.CreatedAt, &blog.UpdatedAt)
		if err != nil {
			return model.DBBlog{}, err
		}
	}

	return blog, nil
}

func (br *blogRepository) Add(title, content, authorId string) (model.DBBlog, error) {

	blog := model.DBBlog{
		Id:        uuid.NewString(),
		Title:     title,
		Content:   content,
		AuthorId:  authorId,
		CreatedAt: time.Now().UTC().GoString(),
		UpdatedAt: time.Now().UTC().GoString(),
	}

	_, err := br.db.Exec("INSERT INTO blogs VALUES($1,$2,$3,$4,$5,$6)", blog.Id, blog.Title, blog.Content, blog.AuthorId, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		return model.DBBlog{}, err
	}
	return blog, nil
}

func (br *blogRepository) GetAll() ([]model.DBBlog, error) {
	var blogs []model.DBBlog
	rows, err := br.db.Query("SELECT * FROM blogs")
	if err != nil {
		return []model.DBBlog{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var blog model.DBBlog
		err = rows.Scan(&blog.Id, &blog.Title, &blog.Content, &blog.AuthorId, &blog.CreatedAt, &blog.UpdatedAt)
		if err != nil {
			return []model.DBBlog{}, err
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}
