package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/parnurzeal/gorequest"
)

type AuthToken struct {
	RefreshToken string
	Token        string
}

type Login struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

//TODO: get orgs list from API
type Account struct {
	Id       string
	Username string
	Email    string
}

type Client struct {
	ApiUrl string
	WebUrl string
	Auth   AuthToken
}

//TODO: return InvalidLogin error so CLI can re-auth?

var contentJSON = "application/json"

func NewClient() *Client {
	client := Client{
		ApiUrl: "http://localhost:5001/",
		WebUrl: "http://localhost:5000/",
	}
	url := os.Getenv("CANVAS_API_URL")
	if url != "" {
		client.ApiUrl = url
	}
	url = os.Getenv("CANVAS_WEB_URL")
	if url != "" {
		client.WebUrl = url
	}
	return &client
}

func (c *Client) TokenLogin(auth Login) (token AuthToken, err error) {
	//build auth body
	authJson, err := json.Marshal(auth)
	check(err)

	body := string(authJson)
	request := gorequest.New()
	resp, body, errs := request.Post(c.Url("tokens")).
		Type("json").
		Send(body).
		End()

	if errs != nil {
		check(errs[0])
	}

	switch resp.StatusCode {
	case 201:
		err = json.Unmarshal([]byte(body), &token)
	default:
		err = errors.New("Login not valid")
	}
	return
}

func (c *Client) FetchAccount() (account Account, err error) {
	resp, body, errs := c.get(c.Url("account")).End()

	if errs != nil {
		err = errs[0]
		return
	}

	switch resp.StatusCode {
	case 404:
		err = errors.New("Account not found")
	case 403:
		err = errors.New("Login invalid")
	case 200:
		err = json.Unmarshal([]byte(body), &account)
	}
	return
}

//Create a new canvas
func (c *Client) NewCanvas(collection string) (canvas Canvas, err error) {
	request := gorequest.New()
	newCanvasUrl := c.Url("canvases/" + collection)
	resp, body, errs := request.Post(newCanvasUrl).
		Set("Authorization", "Bearer "+c.Auth.Token).
		End()

	if errs != nil {
		err = errs[0]
		return
	}

	switch resp.StatusCode {
	case 201:
		err = json.Unmarshal([]byte(body), &canvas)
	default:
		err = errors.New("New canvas failed")
	}

	return
}

func (c *Client) GetCanvas(collection string, name string) (canvas Canvas, err error) {
	canvasUrl := c.Url("canvases/" + collection + "/" + name)
	agent := c.get(canvasUrl)
	resp, body, errs := agent.End()

	if errs != nil {
		err = errs[0]
		return
	}

	switch resp.StatusCode {
	case 200:
		err = json.Unmarshal([]byte(body), &canvas)
	case 404:
		err = errors.New("Not found")
	default:
		err = errors.New("Get Canvas Failed")
	}
	return
}

func (c *Client) GetCanvases(collection string) (canvases []Canvas, err error) {
	canvasesUrl := c.Url("canvases/" + collection)
	agent := c.get(canvasesUrl)
	resp, body, errs := agent.End()

	if errs != nil {
		err = errs[0]
		return
	}

	switch resp.StatusCode {
	case 200:
		err = json.Unmarshal([]byte(body), &canvases)
	case 404:
		err = errors.New("Not found")
	default:
		err = errors.New("Get Canvases failed")
	}

	return
}

func (c *Client) get(path string) (agent *gorequest.SuperAgent) {
	agent = gorequest.New()
	agent.Get(path).
		Set("Authorization", "Bearer "+c.Auth.Token)

	return agent
}

func (c *Client) Url(path string) string {
	return c.ApiUrl + path
}

func (c *Client) JoinWebUrl(path string) string {
	return c.WebUrl + path
}

func (c *Client) LoggedIn() bool {
	return c.Auth.Token != ""
}
