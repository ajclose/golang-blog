package models

import (
	"fmt"
	"time"

	"github.com/ajclose/golang-blog/config"
	"gopkg.in/mgo.v2/bson"
)

type Creator struct {
	Creator_id string `bson:"creator_id"`
	Username   string `bson:"username"`
}

type Blog struct {
	Id         string `json:"id" bson:"_id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Creator    `bson:"creator"`
	Published  bool    `json:"published"`
	Author     string  `json:"author"`
	Images     []Image `json:"images"`
	Created_at time.Time
	Updated_at time.Time
}

type DefaultText struct {
	ParagraphTitle string `bson:"paragraphTitle"`
	ParagraphBody  string `bson:"paragraphBody"`
}

type Image struct {
	Img string `json:'img'`
}

func FindBlogs(published bool, query string) []Blog {
	blogs := []Blog{}
	err := config.Blogs.Find(bson.M{"$and": []bson.M{bson.M{"body": bson.M{"$regex": query, "$options": "i"}}, bson.M{"published": published}}}).Sort("-updated_at").All(&blogs)
	if err != nil {
		fmt.Println("error")
	}
	return blogs
}

func FindUserBlogs(user_id string) []Blog {
	blogs := []Blog{}
	err := config.Blogs.Find(bson.M{"creator.creator_id": user_id}).All(&blogs)
	if err != nil {
		fmt.Println(err)
		return blogs
	}
	return blogs
}

func CreateBlog(blog Blog, user User) {
	blog.Id = bson.NewObjectId().Hex()
	blog.Creator = Creator{user.Id.Hex(), user.Username}
	blog.Created_at = time.Now()
	blog.Updated_at = time.Now()
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
