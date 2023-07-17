package service

import "github.com/fsimic346/go-blog/model"

type blogService struct {
	userRepository model.UserRepository
	blogRepository model.BlogRepository
}

func CreateBlogService(userRepository model.UserRepository, blogRepository model.BlogRepository) model.BlogService {
	return &blogService{
		userRepository: userRepository,
		blogRepository: blogRepository,
	}
}

func (bs *blogService) GetById(blogId string) (model.Blog, error) {
	dbBlog, err := bs.blogRepository.GetById(blogId)
	if err != nil {
		return model.Blog{}, err
	}

	author, err := bs.userRepository.GetById(dbBlog.AuthorId)
	if err != nil {
		return model.Blog{}, err
	}

	blog := model.ConvertDBBlogToBlog(dbBlog, author)

	return blog, nil

}

func (bs *blogService) Add(title, content, authorId string) (model.Blog, error) {
	dbBlog, err := bs.blogRepository.Add(title, content, authorId)
	if err != nil {
		return model.Blog{}, err
	}

	author, err := bs.userRepository.GetById(dbBlog.AuthorId)
	if err != nil {
		return model.Blog{}, err
	}

	return model.ConvertDBBlogToBlog(dbBlog, author), nil
}

func (bs *blogService) GetAll() ([]model.Blog, error) {
	dbBlogs, err := bs.blogRepository.GetAll()
	if err != nil {
		return []model.Blog{}, err
	}

	var blogs []model.Blog
	for _, dbBlog := range dbBlogs {
		author, err := bs.userRepository.GetById(dbBlog.AuthorId)
		if err != nil {
			return []model.Blog{}, err
		}
		blogs = append(blogs, model.ConvertDBBlogToBlog(dbBlog, author))
	}

	return blogs, nil
}
