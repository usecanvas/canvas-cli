package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

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
	Auth    AuthToken
	Account Account
}

var contentJSON = "application/json"
var canvasDir = ".canvas"
var authTokenFile = "auth-token.json"

func home(path string) string {
	usr, _ := user.Current()
	canvasHome := usr.HomeDir + "/" + canvasDir
	if path == "" {
		return canvasHome
	}
	return canvasHome + "/" + path
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

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
	//get username
	var identity string
	fmt.Fprintf(os.Stderr, "Please enter your username or email: ")
	_, err := fmt.Scanln(&identity)
	check(err)

	//get password
	fmt.Fprintf(os.Stderr, "Please enter your password: ")
	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintf(os.Stderr, "\n")
	check(err)
	password := string(pass)

	//build auth body
	auth := Login{identity, password}
	authJson, err := json.Marshal(auth)
	check(err)
	body := string(authJson)

	//make request
	request := gorequest.New()
	resp, body, errs := request.Post(c.Url("tokens")).
		Type("json").
		Send(body).
		End()
	if errs != nil {
		log.Fatal(errs)
	}

	switch resp.StatusCode {
	case 201:
		var token AuthToken
		err = json.Unmarshal([]byte(body), &token)
		check(err)
		c.Auth = token
		c.Save()
	default:
		log.Fatal("Login not valid: ", resp.StatusCode)
	}
}

// save login to `~/.canvas/auth-token.json`
func (c *Client) Save() {
	dirExists, err := exists(home(""))
	check(err)

	if !dirExists {
		err := os.Mkdir(home(""), 0755)
		check(err)
	}

	authTokenPath := home(authTokenFile)
	authTokenJson, _ := json.Marshal(c.Auth)
	err = ioutil.WriteFile(authTokenPath, authTokenJson, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) LoggedIn() bool {
	return c.Auth.Token != ""
}

//use stored token or initiate login
func (c *Client) DoAuth() {
	//check for stored creds
	authPath := home(authTokenFile)
	authExists, err := exists(authPath)
	check(err)

	if authExists {
		authTokenJson, err := ioutil.ReadFile(authPath)
		check(err)
		var token AuthToken
		err = json.Unmarshal(authTokenJson, &token)
		check(err)
		c.Auth = token
	} else {
		c.Login()
	}

	//this will trigger login if account is stale
	c.FetchAccount()
}

func (c *Client) FetchAccount() {
	//TODO: change to HasAuth()
	if !c.LoggedIn() {
		log.Fatal("Must be logged in")
	}

	request := gorequest.New()
	resp, body, errs := request.Get(c.Url("account")).
		Set("Authorization", "Bearer "+c.Auth.Token).
		End()

	if errs != nil {
		log.Fatal(errs)
	}

	switch resp.StatusCode {
	case 404:
		log.Fatal("Account not found")
	case 403:
		fmt.Println("Login invalid")
		c.DoAuth()
	case 200:
		var account Account
		err := json.Unmarshal([]byte(body), &account)
		check(err)
		c.Account = account
	}
}

//Create a new canvas
func (c *Client) NewCanvas() (url string) {
	if !c.LoggedIn() {
		log.Fatal("Must be logged in")
	}

	request := gorequest.New()
	newCanvasUrl := c.Url("canvases/" + c.Account.Username)
	resp, body, errs := request.Post(newCanvasUrl).
		Set("Authorization", "Bearer "+c.Auth.Token).
		End()

	if errs != nil {
		log.Fatal(errs)
	}

	switch resp.StatusCode {
	case 201:
		var canvas Canvas
		err := json.Unmarshal([]byte(body), &canvas)
		check(err)
		url = c.JoinWebUrl(canvas.Collection + "/untitled/" + canvas.Name)
	default:
		log.Fatal(body)
	}

	return
}
