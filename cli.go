package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

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
	return
}

func (cli *CLI) NewCanvas(collection string, filepath string) {
	cli.doAuth()
	var body string

	if filepath != "" {
		//read from file
		fileExists, err := exists(filepath)
		check(err)
		if !fileExists {
			fmt.Printf("File %s not found", filepath)
		}
		bytes, err := ioutil.ReadFile(filepath)
		check(err)
		body = string(bytes)
	} else if !terminal.IsTerminal(int(os.Stdin.Fd())) {
		// read from STDIN if not a terminal
		bytes, err := ioutil.ReadAll(os.Stdin)
		check(err)
		body = string(bytes)
	}

	// default collection to username
	if collection == "" {
		collection = cli.Account.Username
	}

	// make the canvas
	canvas, err := cli.Client.NewCanvas(collection, body)
	check(err)
	canvas.URL = cli.Client.JoinWebUrl(canvas.WebName())
	fmt.Println(canvas.URL)
}

func (cli *CLI) WhoAmI() {
	cli.doAuth()
	account := cli.Account
	fmt.Println("Username: ", account.Username)
	fmt.Println("Email:    ", account.Email)
}

func (cli *CLI) PullCanvas(id string, format string) {
	cli.doAuth()
	canvasText, err := cli.Client.GetCanvas(id, format)
	check(err)
	fmt.Println(canvasText)
}

func (cli *CLI) DeleteCanvas(id string) {
	cli.doAuth()
	err := cli.Client.DeleteCanvas(id)
	check(err)
	fmt.Println("Deleted: ", id)
}

func (cli *CLI) ListCanvases(collection string) {
	cli.doAuth()
	canvases, err := cli.Client.GetCanvases(collection)
	check(err)

	//TODO: have API return collection names
	collections, err := cli.Client.GetCollections()
	check(err)

	//make a map of canvas id to name
	cMap := make(map[string]string)
	for _, collection := range collections {
		cMap[collection.Id] = collection.Name
	}

	for _, canvas := range canvases {
		//pull in collectionName
		canvas.CollectionName = cMap[canvas.CollectionId]
		url := cli.Client.JoinWebUrl(canvas.WebName())
		fmt.Printf("%-30.30s # %s\n", canvas.Title(), url)
	}
}

//display configuration info about env.
func (cli *CLI) Env() {
	authTokenPath := home(authTokenFile)
	authTokenExists, err := exists(authTokenPath)
	check(err)
	if authTokenExists {
		fmt.Println("Auth token exists at:", authTokenPath)
	} else {
		fmt.Println("No auth token", authTokenPath)
	}
	fmt.Println("Canvas API Url:", cli.Client.ApiUrl)
	fmt.Println("Canvas Web Url:", cli.Client.WebUrl)
}

//Prompt user for login and auth with
//acquire auth token
func (cli *CLI) RefreshLogin() (err error) {
	refreshToken := RefreshToken{Id: cli.Client.Auth.RefreshToken}
	token, err := cli.Client.RefreshTokenLogin(refreshToken)
	if err != nil {
		return
	}

	cli.Client.Auth = token
	return
}

//Prompt user for login and auth with
//acquire auth token
func (cli *CLI) UserLogin() {
	//get username or password
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

	//determine if given username or password
	var auth User
	if strings.ContainsRune(identity, '@') {
		auth = User{Email: identity, Password: password}
	} else {
		auth = User{Username: identity, Password: password}
	}

	token, err := cli.Client.UserLogin(auth)
	check(err)

	cli.Client.Auth = token
	cli.save()
}

//use stored token or initiate login
func (cli *CLI) doAuth() {
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
		err = cli.RefreshLogin()
		if err != nil {
			cli.UserLogin()
		}
	} else {
		cli.UserLogin()
	}
	//we've acquired a token! save it
	cli.save()

	//TODO: maybe convert to error checking
	var account Account
	account, err = cli.Client.FetchAccount()
	check(err)
	cli.Account = account
}

// save login to `~/.canvas/auth-token.json`
func (cli *CLI) save() {
	dirExists, err := exists(home(""))
	check(err)

	if !dirExists {
		err := os.Mkdir(home(""), 0755)
		check(err)
	}

	authTokenPath := home(authTokenFile)
	authTokenJson, _ := json.Marshal(cli.Client.Auth)
	err = ioutil.WriteFile(authTokenPath, authTokenJson, 0600)
	check(err)
}

//helper functions
//TODO: configure canvas home via ENV
//TODO: make part of CLI struct
//TODO: use filepath
func home(path string) string {
	usr, _ := user.Current()
	canvasHome := usr.HomeDir + "/" + canvasDir
	if path == "" {
		return canvasHome
	}
	return canvasHome + "/" + path
}

//helper for dir or path existence
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
