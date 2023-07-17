package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fsimic346/go-blog/model"
	"github.com/fsimic346/go-blog/util"
	"github.com/go-chi/chi/v5"
)

type BlogController struct {
	BlogService model.BlogService
}

func (bc *BlogController) GetBlog(w http.ResponseWriter, r *http.Request) {
	blogId := chi.URLParam(r, "blogId")
	blog, err := bc.BlogService.GetById(blogId)

	if err != nil {
		log.Printf("Error while fetching blog: %v", err)
		util.RespondWithError(w, http.StatusNotFound, "Couldn't find blog")
		return
	}

	util.RespondWithJSON(w, 200, blog)
}

func (bc *BlogController) GetAll(w http.ResponseWriter, r *http.Request) {
	blogs, err := bc.BlogService.GetAll()

	if err != nil {
		log.Printf("Error while fetching blogs: %v", err)
		util.RespondWithError(w, http.StatusNotFound, "Couldn't fetch blogs")
		return
	}

	util.RespondWithJSON(w, 200, blogs)
}

func (bc *BlogController) AddBlog(w http.ResponseWriter, r *http.Request) {

	type reqParams struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var params reqParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid blog data")
		return
	}

	authorId, ok := r.Context().Value("userId").(string)

	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't convert user id")
		return
	}

	blog, err := bc.BlogService.Add(params.Title, params.Content, authorId)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, blog)
}
