package main

type CLI struct{}

func (cli *CLI) New() {
	client := NewClient()
	client.Auth()
	client.NewCanvas()
}
