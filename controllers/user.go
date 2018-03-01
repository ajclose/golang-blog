package controllers

import (
	"net/http"

	"github.com/ajclose/golang-blog/config"
	"github.com/ajclose/golang-blog/models"
	"github.com/julienschmidt/httprouter"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (uc UserController) New(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if models.IsLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var data interface{}
	user := models.User{}
	vd := models.ViewData{user, data}
	config.CreateView("user_new.gohtml")
	config.Base.ExecuteTemplate(w, "Base", vd)
}

func (uc UserController) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if models.IsLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	u, err := models.CreateUser(email, username, password)
	if err != nil {
		u := models.User{}
		config.CreateView("user_new.gohtml")
		vd := models.ViewData{u, err}
		config.Base.ExecuteTemplate(w, "Base", vd)
		return
	}
	models.CreateSession(w, r, u.Id.Hex())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (uc UserController) Show(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user := models.FindUserBySessionId(r)
	blogs := models.FindUserBlogs(user.Id.Hex())
	vd := models.ViewData{user, blogs}
	config.CreateView("user_show.gohtml")
	config.Base.ExecuteTemplate(w, "Base", vd)
}
