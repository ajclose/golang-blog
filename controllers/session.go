package controllers

import (
	"fmt"
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
		config.TPL.ExecuteTemplate(w, "session_new.gohtml", "Username or password is incorrect")
		return
	}
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
	if err != nil {
		config.TPL.ExecuteTemplate(w, "session_new.gohtml", "Username or password is incorrect")
		return
	}
	models.CreateSession(w, r, u.Id.Hex())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (sc SessionController) Destroy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c, _ := r.Cookie("session")
	err := config.Sessions.Remove(bson.M{"session_id": c.Value})
	if err != nil {
		fmt.Println(err)
	}
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
