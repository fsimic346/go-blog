package model

type Blog struct {
	Id      string
	Title   string
	Content string
	Author  User
}

type DBBlog struct {
	Id        string
	Title     string
	Content   string
	AuthorId  string
	CreatedAt string
	UpdatedAt string
}

type BlogRepository interface {
	GetAll() ([]DBBlog, error)
	GetById(id string) (DBBlog, error)
	Add(title string, content string, authorId string) (DBBlog, error)
}

type BlogService interface {
	GetAll() ([]Blog, error)
	GetById(id string) (Blog, error)
	Add(title string, content string, authorId string) (Blog, error)
}

func ConvertDBBlogToBlog(dbBlog DBBlog, author User) Blog {
	return Blog{
		Id:      dbBlog.Id,
		Title:   dbBlog.Title,
		Content: dbBlog.Content,
		Author:  author,
	}
}
