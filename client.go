package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/parnurzeal/gorequest"
)

type AuthToken struct {
	RefreshToken string `json:"refresh_token"`
	Token        string `json:"token"`
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
		ApiUrl: "https://api.usecanvas.com/",
		WebUrl: "https://beta.usecanvas.com/",
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
	case 401, 403:
		err = errors.New("Login invalid")
	case 200:
		err = json.Unmarshal([]byte(body), &account)
	}
	return
}

//Create a new canvas
func (c *Client) NewCanvas(collection string, data string) (canvas Canvas, err error) {
	newCanvasUrl := c.Url("canvases/" + collection)
	postBody, err := json.Marshal(ShareData{Data: data})
	check(err)

	agent := c.post(newCanvasUrl).Send(string(postBody))
	resp, body, errs := agent.End()

	if errs != nil {
		check(errs[0])
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

func (c *Client) DeleteCanvas(collection string, name string) (err error) {
	canvasUrl := c.Url("canvases/" + collection + "/" + name)
	agent := c.del(canvasUrl)
	resp, _, errs := agent.End()

	if errs != nil {
		err = errs[0]
		return
	}

	switch resp.StatusCode {
	case 204:
		err = nil
	case 404:
		err = errors.New("Not found")
	default:
		err = errors.New("Delete Canvas Failed")
	}
	return
}

func (c *Client) post(path string) (agent *gorequest.SuperAgent) {
	agent = gorequest.New()
	agent.Post(path).
		Type("json").
		Set("Accept", "application/json").
		Set("Authorization", "Bearer "+c.Auth.Token)

	return agent
}

func (c *Client) get(path string) (agent *gorequest.SuperAgent) {
	agent = gorequest.New()
	agent.Get(path).
		Set("Accept", "application/json").
		Set("Authorization", "Bearer "+c.Auth.Token)

	return agent
}

func (c *Client) del(path string) (agent *gorequest.SuperAgent) {
	agent = gorequest.New()
	agent.Delete(path).
		Set("Accept", "application/json").
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
