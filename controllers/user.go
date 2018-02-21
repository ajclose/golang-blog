package controllers

import (
	"log"
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
	u, err := models.CreateUser(r)
	if err != nil {
		log.Fatalln(err)
	}
	models.CreateSession(w, r, u.Id.Hex())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (uc UserController) Show(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	config.TPL.ExecuteTemplate(w, "user_show.gohtml", nil)
}
