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
	bc := controllers.NewBlogController()
	r.GET("/", uc.Show)
	r.GET("/signup", uc.New)
	r.POST("/signup", uc.Create)
	r.GET("/login", sc.New)
	r.POST("/login", sc.Create)
	r.GET("/logout", sc.Destroy)
	r.GET("/blogs", bc.Index)
	r.GET("/api/blogs/:id", bc.APIShow)
	r.GET("/blogs/:id", bc.Show)
	r.GET("/blog/new", bc.New)
	r.POST("/blog/new", bc.Create)
	r.GET("/blogs/:id/edit", bc.Edit)
	r.POST("/blogs/:id/edit", bc.Update)
	r.DELETE("/blogs/:id", bc.Destroy)
	r.ServeFiles("/dist/*filepath", http.Dir("dist"))
	http.ListenAndServe("localhost:8080", r)
}
