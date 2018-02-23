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
	config.TPL.ExecuteTemplate(w, "user_new.gohtml", nil)
}

func (uc UserController) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	u, err := models.CreateUser(email, username, password)
	if err != nil {
		config.TPL.ExecuteTemplate(w, "user_new.gohtml", err)
		return
	}
	models.CreateSession(w, r, u.Id.Hex())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (uc UserController) Show(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if models.IsLoggedIn(r) {
		user := models.FindUserBySessionId(r)
		config.TPL.ExecuteTemplate(w, "user_show.gohtml", user)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
