package main

import (
	"net/http"

	"github.com/ajclose/golang-blog/controllers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController()
	sc := controllers.NewSessionController()
	r.GET("/", uc.Show)
	r.GET("/signup", uc.New)
	r.POST("/signup", uc.Create)
	r.GET("/login", sc.New)
	r.POST("/login", sc.Create)
	http.ListenAndServe("localhost:8080", r)
}
