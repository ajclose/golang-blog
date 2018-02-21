package models

import (
	"errors"
	"net/http"

	"github.com/ajclose/golang-blog/config"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string
	Username string
	Password []byte
}

func CreateUser(r *http.Request) (User, error) {
	var u User
	un := r.FormValue("username")
	e := r.FormValue("email")
	p := r.FormValue("password")

	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		return u, errors.New("500. Internal Server Error." + err.Error())
	}
	u = User{e, un, bs}
	err = config.Users.Insert(u)
	if err != nil {
		return u, errors.New("500. Internal Server Error." + err.Error())
	}
	return u, nil
}
