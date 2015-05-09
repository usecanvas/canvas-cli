package main

import (
	"encoding/json"
	"os"

	"github.com/parnurzeal/gorequest"
)

type ApiPayload struct {
	Data interface{} `json:"data"`
}

type AuthToken struct {
	RefreshToken string `json:"refresh_token"`
	Token        string `json:"token"`
}

type UserResource struct {
	User `json:"user"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
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

func (c *Client) UserLogin(auth User) (token AuthToken, err error) {
	payload := &ApiPayload{Data: &UserResource{auth}}
	//build auth body
	authJson, err := json.Marshal(payload)
	check(err)
	reqBody := string(authJson)

	request := gorequest.New()
	resp, body, errs := request.Post(c.Url("tokens")).
		Type("json").
		Send(reqBody).
		End()

	if errs != nil {
		check(errs[0])
	}

	switch resp.StatusCode {
	case 201:
		var resp struct {
			AuthToken `json:"data"`
		}
		err = json.Unmarshal([]byte(body), &resp)
		token = resp.AuthToken
	default:
		err = decodeError(body)
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
	case 200:
		var resp struct {
			Account `json:"data"`
		}
		err = json.Unmarshal([]byte(body), &resp)
		account = resp.Account
	default:
		err = decodeError(body)
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
		err = decodeError(body)
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
	default:
		err = decodeError(body)
	}
	return
}

func (c *Client) GetCollections() (collections []Collection, err error) {
	canvasesUrl := c.Url("collections")
	resp, body, errs := c.get(canvasesUrl).End()

	if errs != nil {
		err = errs[0]
		return
	}

	switch resp.StatusCode {
	case 200:
		var resp struct {
			Collections []Collection `json:"data"`
		}
		err = json.Unmarshal([]byte(body), &resp)
		collections = resp.Collections
	default:
		err = decodeError(body)
	}
	return
}

func (c *Client) GetCanvases(collection string) (canvases []Canvas, err error) {
	canvasesUrl := c.Url("canvases/" + collection)
	resp, body, errs := c.get(canvasesUrl).End()

	if errs != nil {
		err = errs[0]
		return
	}

	switch resp.StatusCode {
	case 200:
		var resp struct {
			Canvases []Canvas `json:"data"`
		}
		err = json.Unmarshal([]byte(body), &resp)
		canvases = resp.Canvases
	default:
		err = decodeError(body)
	}

	// map out the collection names.
	// TODO: get the client to return this
	collections, err := c.GetCollections()
	check(err)
	cMap := make(map[string]string)
	for _, collection := range collections {
		cMap[collection.Id] = collection.Name
	}
	for i, c := range canvases {
		canvases[i].CollectionName = cMap[c.CollectionId]
	}

	return
}

func (c *Client) DeleteCanvas(collection string, name string) (err error) {
	canvasUrl := c.Url("canvases/" + collection + "/" + name)
	agent := c.del(canvasUrl)
	resp, body, errs := agent.End()

	if errs != nil {
		err = errs[0]
		return
	}

	switch resp.StatusCode {
	case 204:
		err = nil
	default:
		err = decodeError(body)
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

func decodeError(body string) error {
	var errorRes ErrorPayload
	err := json.Unmarshal([]byte(body), &errorRes)
	check(err)
	return errorRes
}
