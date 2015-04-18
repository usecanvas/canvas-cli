package main

import "fmt"

type CLI struct{}

func (cli *CLI) New() {
	client := NewClient()
	client.DoAuth()
	url := client.NewCanvas()
	fmt.Println(url)
}

func (cli *CLI) Account() {
	client := NewClient()
	client.DoAuth()
	account := client.Account
	fmt.Println("Username: ", account.Username)
	fmt.Println("Email:    ", account.Email)
}
