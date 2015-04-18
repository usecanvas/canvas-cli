package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/parnurzeal/gorequest"
	"golang.org/x/crypto/ssh/terminal"
)

type Canvas struct {
	Name       string
	Collection string
}

type AuthToken struct {
	RefreshToken string
	Token        string
}

type Login struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type Account struct {
	Id       string
	Username string
	Email    string
}

type Client struct {
	ApiUrl  string
	WebUrl  string
	auth    AuthToken
	account Account
}

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

func (c *Client) Url(path string) string {
	return c.ApiUrl + path
}

func (c *Client) JoinWebUrl(path string) string {
	return c.WebUrl + path
}

//Prompt user for login and auth with
//acquire auth token
func (c *Client) Login() {
	var identity string
	fmt.Fprintf(os.Stderr, "Please enter your username or email: ")
	_, err := fmt.Scanln(&identity)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "Please enter your password: ")
	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	password := string(pass)
	auth := Login{identity, password}
	loginJson, err := json.Marshal(auth)
	if err != nil {
		log.Fatal(err)
	}
	authString := string(loginJson)
	request := gorequest.New()
	resp, body, errs := request.Post(c.Url("tokens")).
		Type("json").
		Send(authString).
		End()

	if errs != nil {
		log.Fatal(errs)
	}
	if resp.StatusCode >= 400 {
		log.Fatal("Login not valid")
	} else {
		var token AuthToken
		err = json.Unmarshal([]byte(body), &token)
		if err != nil {
			log.Fatal(err)
		}
		c.auth = token
	}
}

func (c *Client) LoggedIn() bool {
	return c.auth.Token != ""
}

//use stored token or initiate login
func (c *Client) Auth() {
	//TODO: check for stored creds
	if !c.LoggedIn() {
		c.Login()
	}
	c.FetchAccount()
}

func (c *Client) FetchAccount() {
	if !c.LoggedIn() {
		log.Fatal("Must be logged in")
	}

	request := gorequest.New()
	resp, body, errs := request.Get(c.Url("account")).
		Set("Authorization", "Bearer "+c.auth.Token).
		End()

	if errs != nil {
		log.Fatal(errs)
	}

	if resp.StatusCode >= 400 {
		log.Fatal("Account not found")
	} else {
		var account Account
		err := json.Unmarshal([]byte(body), &account)
		if err != nil {
			log.Fatal(err)
		}
		c.account = account
	}
}

func (c *Client) NewCanvas() {
	if !c.LoggedIn() {
		log.Fatal("Must be logged in")
	}

	request := gorequest.New()
	newCanvasUrl := c.Url("canvases/" + c.account.Username)
	resp, body, errs := request.Post(newCanvasUrl).
		Set("Authorization", "Bearer "+c.auth.Token).
		End()

	if errs != nil {
		log.Fatal(errs)
	}

	if resp.StatusCode >= 400 {
		log.Fatal(body)
	} else {
		var canvas Canvas
		err := json.Unmarshal([]byte(body), &canvas)
		if err != nil {
			log.Fatal(err)
		}
		url := c.JoinWebUrl(canvas.Collection + "/untitled/" + canvas.Name)
		fmt.Println(url)
	}

}
