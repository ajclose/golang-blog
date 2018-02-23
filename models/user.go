package models

import (
	"errors"
	"net/http"
	"strings"

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

func CreateUser(email, username, password string) (User, error) {
	var u User
	err := validEmail(email)
	if err != nil {
		return u, err
	}
	err = validUsername(username)
	if err != nil {
		return u, err
	}
	err = validPassword(password)
	if err != nil {
		return u, err
	}
	u.Id = bson.NewObjectId()
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return u, errors.New("500. Internal Server Error." + err.Error())
	}
	u.Email = email
	u.Username = username
	u.Password = bs
	err = config.Users.Insert(u)
	if err != nil {
		return u, errors.New("500. Internal Server Error." + err.Error())
	}
	return u, nil
}

func FindUserBySessionId(r *http.Request) User {
	user := User{}
	session := Session{}
	c, _ := r.Cookie("session")
	err := config.Sessions.Find(bson.M{"session_id": c.Value}).Select(bson.M{"_id": 0, "user_id": 1}).One(&session)
	err = config.Users.FindId(bson.ObjectIdHex(session.User_ID)).One(&user)
	if err != nil {
		return user
	}
	return user
}

func validEmail(email string) error {
	emails := []string{}
	if len(email) == 0 {
		return errors.New("Email can't be blank")
	}
	err := config.Users.Find(bson.M{}).Distinct("email", &emails)
	if err != nil {
		return err
	}
	for _, v := range emails {
		if strings.ToLower(v) == strings.ToLower(email) {
			return errors.New("Email already exists")
		}
	}
	return nil
}

func validUsername(username string) error {
	usernames := []string{}
	if len(username) == 0 {
		return errors.New("Username can't be blank")
	}
	err := config.Users.Find(bson.M{}).Distinct("username", &usernames)
	if err != nil {
		return err
	}
	for _, v := range usernames {
		if strings.ToLower(v) == strings.ToLower(username) {
			return errors.New("Username already exists")
		}
	}
	return nil
}

func validPassword(password string) error {
	if len(password) == 0 {
		return errors.New("Password can't be blank")
	}
	return nil
}
