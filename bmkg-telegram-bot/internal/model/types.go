package model

import "encoding/xml"

// RSSFeed represents the root structure of the BMKG RSS feed
type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel represents the channel information in the RSS feed
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Item represents a single alert item in the RSS feed
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Author      string `xml:"author"`
	Category    string `xml:"category"` // Met
	GUID        GUID   `xml:"guid"`
	PubDate     string `xml:"pubDate"`
}

// GUID represents the unique identifier for an item
type GUID struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

// CAPAlert represents the detailed CAP XML structure
type CAPAlert struct {
	XMLName xml.Name `xml:"alert"`
	Info    CAPInfo  `xml:"info"`
}

type CAPInfo struct {
	Headline    string `xml:"headline"`
	Description string `xml:"description"`
	Web         string `xml:"web"` // This contains the infographic URL
	Area        struct {
		AreaDesc string `xml:"areaDesc"`
	} `xml:"area"`
}
