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
	if models.IsLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var data interface{}
	user := models.User{}
	vd := models.ViewData{user, data}
	config.CreateView("session_new.gohtml")
	config.Base.ExecuteTemplate(w, "Base", vd)
}

func (sc SessionController) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if models.IsLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u := models.User{}
	un := r.FormValue("username")
	p := r.FormValue("password")
	ok := config.Users.Find(bson.M{"username": un}).One(&u)
	if ok != nil {
		u := models.User{}
		config.CreateView("session_new.gohtml")
		vd := models.ViewData{u, "Username or password is incorrect"}
		config.Base.ExecuteTemplate(w, "Base", vd)
		return
	}
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
	if err != nil {
		u := models.User{}
		config.CreateView("session_new.gohtml")
		vd := models.ViewData{u, "Username or password is incorrect"}
		config.Base.ExecuteTemplate(w, "Base", vd)
		return
	}
	models.CreateSession(w, r, u.Id.Hex())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (sc SessionController) Destroy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
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
