package main

import "fmt"

type CLI struct{}

func (cli *CLI) New() {
	client := NewClient()
	client.DoAuth()
	url := client.NewCanvas()
	fmt.Println(url)
}
