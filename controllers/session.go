package controllers

import (
	"net/http"

	"github.com/ajclose/golang-blog/config"
	"github.com/ajclose/golang-blog/models"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type SessionController struct{}

func NewSessionController() *SessionController {
	return &SessionController{}
}

func (sc SessionController) New(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	config.TPL.ExecuteTemplate(w, "session_new.gohtml", nil)
}

func (sc SessionController) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	un := r.FormValue("username")
	p := r.FormValue("password")
	u := models.User{}
	ok := config.Users.Find(bson.M{"username": un}).One(&u)
	if ok != nil {
		http.Error(w, "Username or password is not correct", http.StatusForbidden)
		return
	}
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
	if err != nil {
		http.Error(w, "Username or password is not correct", http.StatusForbidden)
		return
	}
	models.CreateSession(w, r, u.Id.Hex())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
