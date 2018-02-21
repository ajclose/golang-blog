package controllers

import (
	"fmt"
	"net/http"

	"github.com/ajclose/golang-blog/config"
	"github.com/julienschmidt/httprouter"
)

func NewSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	config.TPL.ExecuteTemplate(w, "session_new.gohtml", nil)
}

func CreateSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("CREATE SESSION")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
