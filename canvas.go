package main

type Canvas struct {
	Name       string `json:"name"`
	Collection string `json:"collection"`
	URL        string
	Share      ShareData `json:"data"`
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

func (c *Canvas) Body() string {
	return c.Share.Data
}
