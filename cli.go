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

// Interatcts with a user via STDIN/STDOUT/STDERR,
// filesystem, and HTTP via Client
type CLI struct {
	Client
	Account
}

var canvasDir = ".canvas"
var authTokenFile = "auth-token.json"

func NewCLI() (cli *CLI) {
	client := NewClient()
	cli = &CLI{Client: *client}
	cli.DoAuth()
	return
}

func (cli *CLI) NewCanvas() {
	canvas, err := cli.Client.NewCanvas(cli.Account.Username)
	check(err)
	canvas.URL = cli.Client.JoinWebUrl(canvas.Collection + "/untitled/" + canvas.Name)
	fmt.Println(canvas.URL)
}

func (cli *CLI) WhoAmI() {
	account := cli.Account
	fmt.Println("Username: ", account.Username)
	fmt.Println("Email:    ", account.Email)
}

func (cli *CLI) List() {
	canvases, err := cli.Client.Canvases(cli.Account.Username)
	check(err)
	for _, canvas := range canvases {
		fmt.Println(canvas)
	}
}

//use stored token or initiate login
func (cli *CLI) DoAuth() {
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
		cli.Client.Auth = token
	} else {
		cli.Login()
	}

	//TODO: maybe convert to error checking
	account, err := cli.Client.FetchAccount()
	if err != nil {
		cli.DoAuth()
	} else {
		//prompt again
		cli.Account = account
		cli.Save()
	}
}

//Prompt user for login and auth with
//acquire auth token
func (cli *CLI) Login() {
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
	resp, body, errs := request.Post(cli.Client.Url("tokens")).
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
		cli.Client.Auth = token
		cli.Save()
	default:
		log.Fatal("Login not valid: ", resp.StatusCode)
	}
}

// save login to `~/.canvas/auth-token.json`
func (cli *CLI) Save() {
	dirExists, err := exists(home(""))
	check(err)

	if !dirExists {
		err := os.Mkdir(home(""), 0755)
		check(err)
	}

	authTokenPath := home(authTokenFile)
	authTokenJson, _ := json.Marshal(cli.Client.Auth)
	err = ioutil.WriteFile(authTokenPath, authTokenJson, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

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
