package main;

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"sync/atomic"
	"sync"
	"github.com/satori/go.uuid"
)

type Entry struct {
	Title    string
	Contents string
	Upvote   int
	Downvote int
	Uid      string
}
var dataStore atomic.Value
var lock sync.Mutex
type DBMap map[string]*Entry

func InsertValue(newEntry *Entry){
	lock.Lock()
	defer lock.Unlock()
	db1 := dataStore.Load().(DBMap)
	db2 := make(DBMap)
	for k, v := range db1{
		db2[k] = v
	}
	db2[newEntry.Uid] = newEntry
	dataStore.Store(db2)
}

func Downvote(uid string){
	lock.Lock()
	defer lock.Unlock()
	db1 := dataStore.Load().(DBMap)
	db2 := make(DBMap)
	for k, v := range db1{
		db2[k] = v
	}
	db2[uid].Downvote++
	dataStore.Store(db2)
}

func Upvote(uid string){
	lock.Lock()
	defer lock.Unlock()
	db1 := dataStore.Load().(DBMap)
	db2 := make(DBMap)
	for k, v := range db1{
		db2[k] = v
	}
	db2[uid].Upvote++
	dataStore.Store(db2)
}

func main() {
	dataStore.Store(make(DBMap))

	r := gin.Default();
	r.LoadHTMLGlob("templates/*")
	r.GET("/", indexpage)

	r.POST("/write", WritePost)

	log.Fatalln(r.Run(":8080"))
}

func WritePost (c *gin.Context) {
	title := c.PostForm("title")
	contents := c.PostForm("contents")
	if contents == ""{
		c.Redirect(http.StatusMovedPermanently, "/")
	}
	if title == ""{
		c.Redirect(http.StatusMovedPermanently, "/")
	}
	newEntry := Entry{
		Title:title,
		Contents:contents,
		Uid:uuid.NewV4().String(),
	}
	InsertValue(&newEntry)
	c.Redirect(http.StatusMovedPermanently, "/")
}

func indexpage(c *gin.Context) {
	items := dataStore.Load().(DBMap)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Items": items,
	})
}