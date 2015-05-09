package main

import "regexp"

type Canvas struct {
	Id             string `json:"id"`
	CollectionId   string `json:"collection_id"`
	CollectionName string `json:"collection_name"`
	URL            string
	Data           ShareData
}

type ShareData struct {
	Version int    `json:"v"`
	Data    string `json:"data"`
	Type    string `json:"type"`
	ModTime `json:"m"`
}

type ModTime struct {
	Mtime int
	Ctime int
}

var titleRegexp = regexp.MustCompile(`^# ([^\n])+`)

func (c *Canvas) WebName() string {
	return c.CollectionName + "/-/" + c.Id
}

func (c *Canvas) Body() string {
	return c.Data.Data
}

func (c *Canvas) Title() (title string) {
	body := c.Body()
	title = "Untitled"
	match := titleRegexp.FindString(body)
	if match != "" {
		title = match[2:]
	}
	return
}
