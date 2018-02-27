package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	fmt.Println(r.URL)
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
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	var blog models.Blog
	err = json.Unmarshal(b, &blog)
	fmt.Println(blog)
	user := models.FindUserBySessionId(r)
	models.CreateBlog(blog, user)
	http.Redirect(w, r, "/blogs", http.StatusFound)
}

func (bc BlogController) APIShow(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	blog := models.FindBlog(p.ByName("id"))
	w.Header().Set("Content-Type", "application/json")
	json, err := json.Marshal(blog)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(json)
}

func (bc BlogController) Show(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	config.TPL.ExecuteTemplate(w, "blog_show.gohtml", nil)
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
