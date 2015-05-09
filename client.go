package main

import (
	"encoding/json"
	"errors"
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

type RefreshToken struct {
	Id string `json:"id"`
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

func (c *Client) RefreshTokenLogin(refreshToken RefreshToken) (token AuthToken, err error) {

	var tokenData struct {
		RefreshToken `json:"refresh_token"`
	}
	tokenData.RefreshToken = refreshToken
	payload := &ApiPayload{Data: &tokenData}
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

func (c *Client) UserLogin(user User) (token AuthToken, err error) {
	var userData struct {
		User `json:"user"`
	}
	userData.User = user
	payload := &ApiPayload{Data: &userData}
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

func (c *Client) CollectionNameToId() (cMap map[string]string) {
	collections, err := c.GetCollections()
	check(err)

	//make a map of canvas id to name
	cMap = make(map[string]string)
	for _, collection := range collections {
		cMap[collection.Name] = collection.Id
	}

	return
}

//Create a new canvas
func (c *Client) NewCanvas(collection string, data string) (canvas Canvas, err error) {
	newCanvasUrl := c.Url("canvases")
	var postData struct {
		Collection `json:"collection"`
		Text       string `json:"text"`
	}
	// figure out the collection id from the name
	cMap := c.CollectionNameToId()
	if cMap[collection] == "" {
		err = errors.New("Collection \"" + collection + "\" not found")
		return
	}

	//set and serialize the post body
	postData.Collection.Id = cMap[collection]
	postData.Text = data
	postBody, err := json.Marshal(ApiPayload{Data: &postData})
	check(err)

	agent := c.post(newCanvasUrl).Send(string(postBody))
	resp, body, errs := agent.End()

	if errs != nil {
		check(errs[0])
	}

	switch resp.StatusCode {
	case 201:
		var resp struct {
			Canvas `json:"data"`
		}
		err = json.Unmarshal([]byte(body), &resp)
		canvas = resp.Canvas
		canvas.CollectionName = collection
	default:
		err = decodeError(body)
	}

	return
}

//Get canvas by id
func (c *Client) GetCanvas(id string, format string) (canvasText string, err error) {
	canvasUrl := c.Url("canvas/" + id)

	var mimeType string
	switch format {
	case "md":
		mimeType = "text/plain"
	case "html":
		mimeType = "text/html"
	case "json":
		mimeType = "application/vnd.canvas.doc"
	case "":
		mimeType = "text/plain"
	}

	agent := c.get(canvasUrl).Set("Accept", mimeType)
	resp, body, errs := agent.End()

	if errs != nil {
		err = errs[0]
		return
	}

	switch resp.StatusCode {
	case 200:
		canvasText = body
	case 404:
		err = errors.New("Canvas not found")
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
