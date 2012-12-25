package main

import (
	"bytes"
	"encoding/json"
	"exp/html"
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	feed_lib "loc-com.de/feed"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type dbFeed struct {
	Url   string
	Feed  *feed_lib.Feed
	Users []int
	Id    bson.ObjectId `bson:"_id"`
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

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/feed/", feed)
	http.HandleFunc("/feeds/", getFeeds)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/tmpl/", http.StripPrefix("/tmpl/", http.FileServer(http.Dir("tmpl"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	serverAddr := ":8000"

	log.Println("Starting Server on", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
	db.session.Close()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadFile("index.html")
	fmt.Fprint(w, string(body))
}

func getFeeds(w http.ResponseWriter, r *http.Request) {
	user, _ := strconv.Atoi(r.URL.Path[len("/feeds/"):])

	session := db.session.Copy()
	c := session.DB("test").C("feeds")

	numResults, err := c.Find(bson.M{"users": bson.M{"$in": []int{user}}}).Count()
	if err != nil {
		panic(err)
	}

	feeds := make([]dbFeed, numResults)

	err = c.Find(bson.M{"users": bson.M{"$in": []int{user}}}).All(&feeds)
	if err != nil {
		panic(err)
	}

	respJSON, _ := json.Marshal(feeds)
	fmt.Fprint(w, string(respJSON)) //string(respJSON))
}

func feed(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if strings.HasPrefix(r.Header.Get("Accept"), "application/json") {
			feedId := r.URL.Path[len("/feed/"):]
			feed := &dbFeed{}

			// u := "http://loc-blog.de/rss.php?blog_id=5"

			session := db.session.Copy()
			c := session.DB("test").C("feeds")

			numResults, err := c.FindId(bson.ObjectIdHex(feedId)).Count()
			if err != nil {
				panic(err)
			}

			if numResults == 0 {
				// feed = insertFeed(u)
			} else {
				err = c.FindId(bson.ObjectIdHex(feedId)).One(&feed)
				if err != nil {
					panic(err)
				}
			}

			u, _ := url.Parse(feed.Url)

			for i, _ := range feed.Feed.Items {
				doc, err := html.Parse(strings.NewReader(feed.Feed.Items[i].Content))
				if err != nil {
					log.Fatal(err)
				}
				var f func(*html.Node, *url.URL)
				f = func(n *html.Node, u *url.URL) {
					if n.Type == html.ElementNode && n.Data == "img" {
						for i, _ := range n.Attr {
							if n.Attr[i].Key == "src" {
								u2, _ := url.Parse(n.Attr[i].Val)
								if !u2.IsAbs() {
									u2.Scheme = u.Scheme
									u2.Host = u.Host
								}
								if !strings.HasPrefix(u2.Path, "/") {
									u2.Path = "/" + u2.Path
								}
								n.Attr[i].Val = u2.String()
								break
							}
						}
					}
					if n.Type == html.ElementNode && n.Data == "a" {
						found := false
						for i, _ := range n.Attr {
							if n.Attr[i].Key == "target" {
								n.Attr[i].Val = "_blank"
								found = true
								break
							}
						}
						if !found {
							attr := new(html.Attribute)
							attr.Key = "target"
							attr.Val = "_blank"
							n.Attr = append(n.Attr, *attr)
						}
					}
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						f(c, u)
					}
				}
				f(doc, u)
				var wr bytes.Buffer
				html.Render(&wr, doc)
				feed.Feed.Items[i].Content = wr.String()
			}

			respJSON, _ := json.Marshal(feed)
			fmt.Fprint(w, string(respJSON))
		} else {
			indexHandler(w, r)
		}
	} else if r.Method == "POST" {
		respJSON, _ := json.Marshal(insertFeed(r.FormValue("feed_url")))
		fmt.Fprint(w, string(respJSON))
	}
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
		feed.Id = bson.NewObjectId()

		// for _, feed := range feed.Feed.Items {
		// log.Println(feed)
		// }

		err = c.Insert(feed)
		if err != nil {
			panic(err)
		}
	} else {
		// feed.feed = updateFeed(url)
	}
	return &feed
}
