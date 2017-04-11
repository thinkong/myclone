package main;

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"log"
)

type Entry struct {
	Title    string
	Contents string
	Upvote int
	Downvote int
}

var testEntries []Entry

func main() {
	for i := 0; i < 10; i++ {
		testEntries = append(testEntries, Entry{
			Title: "Hello",
			Contents: fmt.Sprintf("COntents %d", i),
		})
	}
	r := gin.Default();
	r.LoadHTMLGlob("templates/*")
	r.GET("/", indexpage)
	r.POST("/write", func(c *gin.Context){
		title := c.PostForm("title")
		contents := c.PostForm("contents")
		log.Println(title)
		log.Println(contents)
		testEntries = append(testEntries, Entry{
			Title: title,
			Contents:contents,
		})
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	log.Fatalln(r.Run(":8080"))
}

func indexpage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Items": testEntries,
	})
}