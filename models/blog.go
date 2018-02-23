package models

import (
	"fmt"

	"github.com/ajclose/golang-blog/config"
	"gopkg.in/mgo.v2/bson"
)

type Author struct {
	Author_id string
	Username  string
}

type Blog struct {
	Id    string `json:"id" bson:"_id"`
	Title string
	Body  string
	Author
}

func FindBlogs() []Blog {
	blogs := []Blog{}
	err := config.Blogs.Find(bson.M{}).All(&blogs)
	if err != nil {
		fmt.Println("error")
		return blogs
	}
	return blogs
}

func CreateBlog(title string, body string, user User) {
	var blog Blog
	blog.Id = bson.NewObjectId().Hex()
	blog.Title = title
	blog.Body = body
	blog.Author = Author{user.Id.Hex(), user.Username}
	err := config.Blogs.Insert(blog)
	if err != nil {
		fmt.Println(err)
	}
}

func FindBlog(id string) Blog {
	blog := Blog{}
	err := config.Blogs.FindId(id).One(&blog)
	if err != nil {
		return blog
	}
	return blog
}
