package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	feed_lib "loc-com.de/feed"
	"log"
	"net/http"
)

type dbFeed struct {
	Url   string
	Feed  *feed_lib.Feed
	Users []int
}

type DB struct {
	session *mgo.Session
}

var db = &DB{}

func init() {
	var err error
	db.session, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	db.session.SetMode(mgo.Monotonic, true)
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadFile("index.html")
		fmt.Fprint(w, string(body))
	})
	http.HandleFunc("/feed", getFeed)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/tmpl/", http.StripPrefix("/tmpl/", http.FileServer(http.Dir("tmpl"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	serverAddr := ":8080"

	log.Println("Starting Server on", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
	db.session.Close()
}

func getFeed(w http.ResponseWriter, r *http.Request) {
	feed := &dbFeed{}

	url := "http://loc-blog.de/rss.php?blog_id=5"

	session := db.session.Copy()
	c := session.DB("test").C("feeds")

	numResults, err := c.Find(bson.M{"url": url}).Count()
	if err != nil {
		panic(err)
	}

	if numResults == 0 {
		feed = insertFeed(url)
	} else {
		err = c.Find(bson.M{"url": url}).One(&feed)
		if err != nil {
			panic(err)
		}
	}

	respJSON, _ := json.Marshal(feed)
	fmt.Fprint(w, string(respJSON))
}

func insertFeed(url string) *dbFeed {

	session := db.session.Copy()
	c := session.DB("test").C("feeds")
	var feed dbFeed

	numResults, err := c.Find(bson.M{"url": url}).Count()
	if err != nil {
		panic(err)
	}

	if numResults == 0 {
		response, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		feed.Feed, _ = feed_lib.Decode([]byte(contents))
		feed.Url = url
		feed.Users = []int{1}

		for _, feed := range feed.Feed.Items {
			log.Println(feed)
		}

		err = c.Insert(feed)
		if err != nil {
			panic(err)
		}
	} else {
		// feed.feed = updateFeed(url)
	}
	return &feed
}
