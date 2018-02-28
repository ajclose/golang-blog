package controllers

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ajclose/golang-blog/config"
	"github.com/ajclose/golang-blog/models"
	"github.com/julienschmidt/httprouter"
)

type BlogController struct{}

func NewBlogController() *BlogController {
	return &BlogController{}
}

func (bc BlogController) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	blogs := models.FindBlogs(true)
	config.TPL.ExecuteTemplate(w, "blog_index.gohtml", blogs)
}
func (bc BlogController) Drafts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	blogs := models.FindBlogs(false)
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

func (bc BlogController) APIEdit(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	blog := models.FindBlog(p.ByName("id"))
	w.Header().Set("Content-Type", "application/json")
	json, err := json.Marshal(blog)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(json)
}

func (bc BlogController) Edit(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	config.TPL.ExecuteTemplate(w, "blog_edit.gohtml", nil)
}

func (bc BlogController) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	blog := models.FindBlog(p.ByName("id"))
	err = json.Unmarshal(b, &blog)
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

func (bc BlogController) UploadImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	mf, fh, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
	}
	defer mf.Close()
	ext := strings.Split(fh.Filename, ".")[1]
	h := sha1.New()
	io.Copy(h, mf)
	fname := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	path := filepath.Join(wd, "public", "images", fname)
	nf, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	defer nf.Close()
	mf.Seek(0, 0)
	io.Copy(nf, mf)
	res := make(map[string]string)
	res["link"] = "http://localhost:8080/public/images/" + fname
	json, err := json.Marshal(res)
	w.Write(json)
}
