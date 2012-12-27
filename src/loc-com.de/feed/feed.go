// Package loc-com.de/feed provides functions and types to decode Atom-/RSS-/RSS2-Feeds
package feed

import (
	"encoding/xml"
	"fmt"
	"time"
)

type rssItem struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	Author      string   `xml:"author"`
	Categories  []string `xml:"category"`
	Comments    string   `xml:"comments"`
	Guid        string   `xml:"guid"`
	Content     string   `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	Enclosure   struct {
		Url    string `xml:"url,attr"`
		Length string `xml:"length,attr"`
		Type   string `xml:"type,attr"`
	} `xml:"enclosure"`
	Source struct {
		Name string `xml:",innerxml"`
		Url  string `xml:"url,attr"`
	} `xml:"source"`
}
type rssChannel struct {
	Title          string    `xml:"title"`
	Link           string    `xml:"link"`
	Description    string    `xml:"description"`
	Language       string    `xml:"language"`
	Copyright      string    `xml:"copyright"`
	ManagingEditor string    `xml:"managingEditor"`
	WebMaster      string    `xml:"webMaster"`
	PubDate        string    `xml:"pubDate"`
	LastBuildDate  string    `xml:"lastBuildDate"`
	Categories     []string  `xml:"category"`
	Generator      string    `xml:"generator"`
	Docs           string    `xml:"docs"`
	Ttl            string    `xml:"ttl"`
	Rating         string    `xml:"rating"`
	SkipHours      string    `xml:"skipHours"`
	SkipDays       string    `xml:"skipDays"`
	Items          []rssItem `xml:"item"`
	Image          struct {
		Url   string `xml:"url,attr"`
		Title string `xml:"title,attr"`
		Link  string `xml:"link,attr"`
	} `xml:"image"`
}
type rss2Feed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}
type rssFeed struct {
	XMLName xml.Name   `xml:"rdf"`
	Channel rssChannel `xml:"channel"`
}

type atomEntry struct {
	Author       string   `xml:"author"`
	Categories   []string `xml:"category"`
	Content      string   `xml:"content"`
	Contributors []string `xml:"contributor"`
	Id           string   `xml:"id"`
	Link         string   `xml:"link"`
	Published    string   `xml:"published"`
	Rights       string   `xml:"rights"`
	Source       string   `xml:"source"`
	Summary      string   `xml:"summary"`
	Title        string   `xml:"title"`
	Updated      string   `xml:"updated"`
}
type atomFeed struct {
	XMLName      xml.Name    `xml:"feed"`
	Author       string      `xml:"author"`
	Categories   []string    `xml:"category"`
	Contributors []string    `xml:"contributor"`
	Generator    string      `xml:"generator"`
	Logo         string      `xml:"logo"`
	Icon         string      `xml:"icon"`
	Link         string      `xml:"link"`
	Rights       string      `xml:"rights"`
	Subtitle     string      `xml:"subtitle"`
	Title        string      `xml:"title"`
	Updated      string      `xml:"updated"`
	Entries      []atomEntry `xml:"entry"`
}

// FeedItem represents Items in a Feed
type FeedItem struct {
	Author       string
	Categories   []string
	Comments     string
	Content      string
	Contributors []string
	Guid         string
	Id           string
	Link         string
	Published    time.Time // == rss.PubDate
	Received     time.Time
	Rights       string
	Summary      string // == rss.Description
	Title        string
	Updated      time.Time
	Enclosure    struct {
		Url    string
		Length string
		Type   string
	}
	Source struct {
		Title string
		Link  string
	}
}

// Feed represents an RSS/Atom-Feed
type Feed struct {
	Author         string
	Categories     []string
	Contributors   []string
	Copyright      string // == atom.Rights
	Generator      string
	Icon           string
	Logo           string
	Language       string
	Link           string
	ManagingEditor string
	PubDate        string
	Rating         string
	Rights         string
	SkipHours      string
	SkipDays       string
	Subtitle       string // == rss.Description
	Title          string
	Ttl            string
	Updated        string // == rss.LastBuildDate
	WebMaster      string
	Items          []FeedItem
	Image          struct {
		Url   string
		Title string
		Link  string
	}
}

func getAtomItems(atom *atomFeed) []FeedItem {
	feedItems := make([]FeedItem, len(atom.Entries))
	for i := 0; i < len(atom.Entries); i++ {
		feedItems[i].Author = atom.Entries[i].Author
		feedItems[i].Categories = atom.Entries[i].Categories
		feedItems[i].Content = atom.Entries[i].Content
		feedItems[i].Contributors = atom.Entries[i].Contributors
		feedItems[i].Id = atom.Entries[i].Id
		feedItems[i].Link = atom.Entries[i].Link
		feedItems[i].Published, _ = time.Parse("RFC3339", atom.Entries[i].Published)
		feedItems[i].Rights = atom.Entries[i].Rights
		feedItems[i].Source.Link = atom.Entries[i].Source
		feedItems[i].Summary = atom.Entries[i].Summary
		feedItems[i].Title = atom.Entries[i].Title
		feedItems[i].Updated, _ = time.Parse("RFC3339", atom.Entries[i].Updated)
		sanitizeItem(&feedItems[i])
	}
	return feedItems
}

func getRSSItems(channel *rssChannel) []FeedItem {
	feedItems := make([]FeedItem, len(channel.Items))
	for i := 0; i < len(channel.Items); i++ {
		feedItems[i].Author = channel.Items[i].Author
		feedItems[i].Link = channel.Items[i].Link
		feedItems[i].Summary = channel.Items[i].Description
		feedItems[i].Published, _ = time.Parse("RFC3339", channel.Items[i].PubDate)
		feedItems[i].Categories = channel.Items[i].Categories
		feedItems[i].Guid = channel.Items[i].Guid
		feedItems[i].Content = channel.Items[i].Content
		feedItems[i].Title = channel.Items[i].Title
		feedItems[i].Enclosure.Url = channel.Items[i].Enclosure.Url
		feedItems[i].Enclosure.Length = channel.Items[i].Enclosure.Length
		feedItems[i].Enclosure.Type = channel.Items[i].Enclosure.Type
		feedItems[i].Source.Link = channel.Items[i].Source.Url
		feedItems[i].Source.Title = channel.Items[i].Source.Name
		sanitizeItem(&feedItems[i])
	}
	return feedItems
}

func sanitizeItem(feedItem *FeedItem) {
	if len(feedItem.Content) < len(feedItem.Summary) {
		feedItem.Content = feedItem.Summary
	}
	feedItem.Received = time.Now()
}

// Decode decodes an XML byte slice to a Feed struct
func Decode(feedXML []byte) (*Feed, error) {
	feedType := ""
	rss2 := rss2Feed{}
	err := xml.Unmarshal(feedXML, &rss2)
	if err == nil {
		feedType = "rss2"
	}
	atom := atomFeed{}
	err = xml.Unmarshal(feedXML, &atom)
	if err == nil {
		feedType = "atom"
	}
	rss := rssFeed{}
	err = xml.Unmarshal(feedXML, &rss)
	if err != nil && feedType == "" {
		fmt.Println(err)
	} else if err == nil {
		feedType = "rss"
	}
	feed := Feed{}
	if feedType == "rss" || feedType == "rss2" {
		channel := new(rssChannel)
		if feedType == "rss" {
			channel = &rss.Channel
		} else {
			channel = &rss2.Channel
		}
		feed.Author = channel.ManagingEditor
		feed.Categories = channel.Categories
		feed.Copyright = channel.Copyright
		feed.Generator = channel.Generator
		feed.Image.Url = channel.Image.Url
		feed.Image.Title = channel.Image.Title
		feed.Image.Link = channel.Image.Link
		feed.Language = channel.Language
		feed.Link = channel.Link
		feed.PubDate = channel.PubDate
		feed.Rating = channel.Rating
		feed.SkipDays = channel.SkipDays
		feed.SkipHours = channel.SkipHours
		feed.Subtitle = channel.Description
		feed.Title = channel.Title
		feed.Ttl = channel.Ttl
		feed.WebMaster = channel.WebMaster
		feed.Updated = channel.LastBuildDate
		feed.Items = getRSSItems(channel)
	} else if feedType == "atom" {
		feed.Author = atom.Author
		feed.Categories = atom.Categories
		feed.Contributors = atom.Contributors
		feed.Copyright = atom.Rights
		feed.Generator = atom.Generator
		feed.Icon = atom.Icon
		feed.Logo = atom.Logo
		feed.Link = atom.Link
		feed.Title = atom.Title
		feed.Subtitle = atom.Subtitle
		feed.Updated = atom.Updated
		feed.Items = getAtomItems(&atom)
	}
	return &feed, err
}
