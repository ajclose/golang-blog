package controllers

import (
	"fmt"
	"net/http"

	"github.com/ajclose/golang-blog/config"
	"github.com/ajclose/golang-blog/models"
	"github.com/julienschmidt/httprouter"
)

type BlogController struct{}

func NewBlogController() *BlogController {
	return &BlogController{}
}

func (bc BlogController) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	blogs := models.FindBlogs()
	fmt.Println("in index")
	config.TPL.ExecuteTemplate(w, "blog_index.gohtml", blogs)
}

func (bc BlogController) New(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	config.TPL.ExecuteTemplate(w, "blog_new.gohtml", nil)
}

func (bc BlogController) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user := models.FindUserBySessionId(r)
	title := r.FormValue("title")
	body := r.FormValue("body")
	author_id := user.Id.Hex()
	models.CreateBlog(title, body, author_id)
	http.Redirect(w, r, "/blogs", http.StatusFound)
}

func (bc BlogController) Show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	blog := models.FindBlog(p.ByName("id"))
	config.TPL.ExecuteTemplate(w, "blog_show.gohtml", blog)
}

func (bc BlogController) Edit(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	blog := models.FindBlog(p.ByName("id"))
	config.TPL.ExecuteTemplate(w, "blog_edit.gohtml", blog)
}

func (bc BlogController) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	blog := models.FindBlog(p.ByName("id"))
	blog.Title = r.FormValue("title")
	blog.Body = r.FormValue("body")
	config.Blogs.UpdateId(p.ByName("id"), blog)
	http.Redirect(w, r, "/blogs/"+p.ByName("id"), http.StatusFound)
}

func (bc BlogController) Destroy(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := config.Blogs.RemoveId(p.ByName("id"))
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Redirect(w, r, "/blogs", http.StatusSeeOther)
}
