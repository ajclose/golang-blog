package models

import (
	"errors"
	"net/http"

	"github.com/ajclose/golang-blog/config"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Email    string        `json:"email" bson:"email"`
	Username string        `json:"username" bson:"username"`
	Password []byte
}

func CreateUser(r *http.Request) (User, error) {
	var u User
	u.Id = bson.NewObjectId()
	u.Username = r.FormValue("username")
	u.Email = r.FormValue("email")
	p := r.FormValue("password")
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		return u, errors.New("500. Internal Server Error." + err.Error())
	}
	u.Password = bs
	err = config.Users.Insert(u)
	if err != nil {
		return u, errors.New("500. Internal Server Error." + err.Error())
	}
	return u, nil
}
