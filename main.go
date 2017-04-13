package myclone;

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"sync/atomic"
	"sync"
	"github.com/satori/go.uuid"
	. "github.com/ahmetb/go-linq"
)

// Entry
type Entry struct {
	Title    string
	Contents string
	Upvote   int
	Downvote int
	Uid      string
}
// DataStore.
// atomic.Value makes sure that it would be atomic.
var dataStore atomic.Value
var lock sync.Mutex

type DBMap map[string]*Entry

// InsertValue
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
//downvote
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
//upvote
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

var r = gin.Default();
func init() {
	dataStore.Store(make(DBMap))

	r.LoadHTMLGlob("templates/*")
	r.GET("/", indexpage)
	r.POST("/write", WritePost)
	r.GET("/up/:title", func(c *gin.Context){
		uid := c.Param("title")
		Upvote(uid)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})
	r.GET("/down/:title", func(c *gin.Context){
		uid := c.Param("title")
		Downvote(uid)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

}

func Start(){
	log.Fatalln(r.Run(":8080"))
}

func WritePost (c *gin.Context) {
	title := c.PostForm("title")
	contents := c.PostForm("contents")
	// make sure there is data in contents and title..
	// this should be in the tmpl file
	// lazyness to implement it in the tmpl file
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
	var items []*Entry

	db1 := dataStore.Load().(DBMap)
	From(db1).Select(func(e interface{})interface{}{
		return e.(KeyValue).Value
	}).ToSlice(&items)
	From(items).OrderByDescending(func(e interface{}) interface{}{
		return e.(*Entry).Upvote
	}).Take(20).ToSlice(&items)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Items": items,
	})
}