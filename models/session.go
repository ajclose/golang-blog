package models

import (
	"fmt"
	"net/http"

	"github.com/ajclose/golang-blog/config"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

type Session struct {
	Session_ID string
	User_ID    string
}

func CreateSession(w http.ResponseWriter, r *http.Request, id string) {
	var s Session
	sID, _ := uuid.NewV4()
	c := &http.Cookie{
		Name:  "session",
		Value: sID.String(),
	}
	http.SetCookie(w, c)
	s = Session{c.Value, id}
	err := config.Sessions.Insert(s)
	if err != nil {
		fmt.Println(err)
	}
}

func IsLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	u := config.Sessions.Find(bson.M{"session_id": c.Value})
	if u != nil {
		return true
	}
	return false
}
