package config

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

var DB *mgo.Database

var Users *mgo.Collection
var Sessions *mgo.Collection

func init() {
	s, err := mgo.Dial("mongodb://localhost/golang-blog")
	if err != nil {
		panic(err)
	}

	if err = s.Ping(); err != nil {
		panic(err)
	}

	DB = s.DB("golang-blog")
	Users = DB.C("users")
	Sessions = DB.C("sessions")

	fmt.Println("Connected to db")
}
