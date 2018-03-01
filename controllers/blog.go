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
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

type BlogController struct{}

func NewBlogController() *BlogController {
	return &BlogController{}
}

func (bc BlogController) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	blogs := models.FindBlogs(true)
	user := models.FindUserBySessionId(r)
	vd := models.ViewData{user, blogs}
	config.CreateView("blog_index.gohtml")
	config.Base.ExecuteTemplate(w, "Base", vd)
}

func (bc BlogController) Drafts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	blogs := models.FindBlogs(false)
	user := models.FindUserBySessionId(r)
	vd := models.ViewData{user, blogs}
	config.CreateView("blog_index.gohtml")
	config.Base.ExecuteTemplate(w, "Base", vd)
}

func (bc BlogController) New(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user := models.FindUserBySessionId(r)
	vd := models.ViewData{user, nil}
	config.CreateView("blog_new.gohtml")
	config.Base.ExecuteTemplate(w, "Base", vd)
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
	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var blog interface{}
	user := models.FindUserBySessionId(r)
	vd := models.ViewData{user, blog}
	config.CreateView("blog_show.gohtml")
	config.Base.ExecuteTemplate(w, "Base", vd)
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
	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var blog interface{}
	user := models.FindUserBySessionId(r)
	vd := models.ViewData{user, blog}
	config.CreateView("blog_edit.gohtml")
	config.Base.ExecuteTemplate(w, "Base", vd)
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
	blog := models.FindBlog(p.ByName("id"))
	images := blog.Images
	for _, v := range images {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		path := filepath.Join(wd, "public", "images", v.Img)
		os.Remove(path)
	}
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
	uid, _ := uuid.NewV4()
	fname := uid.String() + "." + ext
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
func (bc BlogController) DeleteImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println(r)
	img := models.Image{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(b, &img)
	if err != nil {
		fmt.Println(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	path := filepath.Join(wd, "public", "images", img.Img)
	os.Remove(path)
}

func (bc BlogController) APIDefaultText(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	text := models.DefaultText{}
	err := config.DefaultText.Find(bson.M{}).One(&text)
	fmt.Println(text)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json, err := json.Marshal(text)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(json)
}
