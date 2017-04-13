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
func insertValue(newEntry *Entry) {
	// lock the store to make sure 2 people don't update at the same time.
	lock.Lock()
	defer lock.Unlock()
	db1 := dataStore.Load().(DBMap)
	db2 := make(DBMap)
	for k, v := range db1 {
		db2[k] = v
	}
	// add the entry
	db2[newEntry.Uid] = newEntry
	// when this is called all readers will start using the new one
	dataStore.Store(db2)
}
//downvote
func downvote(uid string) {
	lock.Lock()
	defer lock.Unlock()
	db1 := dataStore.Load().(DBMap)
	db2 := make(DBMap)
	for k, v := range db1 {
		db2[k] = v
	}
	db2[uid].Downvote++
	dataStore.Store(db2)
}
//upvote
func upvote(uid string) {
	lock.Lock()
	defer lock.Unlock()
	db1 := dataStore.Load().(DBMap)
	db2 := make(DBMap)
	for k, v := range db1 {
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
	r.POST("/write", writePost)
	r.GET("/up/:title", upvotePost)
	r.GET("/down/:title", downvotePost)
}

func Start() {
	log.Fatalln(r.Run(":8080"))
}

func upvotePost(c *gin.Context) {
	uid := c.Param("title")
	upvote(uid)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func downvotePost(c *gin.Context) {
	uid := c.Param("title")
	downvote(uid)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func writePost(c *gin.Context) {
	title := c.PostForm("title")
	contents := c.PostForm("contents")
	// make sure there is data in contents and title..
	// this should be in the tmpl file
	// lazyness to implement it in the tmpl file
	if contents == "" {
		c.Redirect(http.StatusMovedPermanently, "/")
	}
	if title == "" {
		c.Redirect(http.StatusMovedPermanently, "/")
	}
	newEntry := Entry{
		Title:title,
		Contents:contents,
		Uid:uuid.NewV4().String(),
	}
	insertValue(&newEntry)
	c.Redirect(http.StatusMovedPermanently, "/")
}

func indexpage(c *gin.Context) {
	var items []*Entry

	db1 := dataStore.Load().(DBMap)
	From(db1).Select(func(e interface{}) interface{} {
		return e.(KeyValue).Value
	}).ToSlice(&items)
	From(items).OrderByDescending(func(e interface{}) interface{} {
		return e.(*Entry).Upvote
	}).Take(20).ToSlice(&items)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Items": items,
	})
}