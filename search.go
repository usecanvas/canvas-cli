package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"golang.org/x/crypto/ssh/terminal"
)

func (cli *CLI) SearchInteractive(collection string) {
	cli.doAuth()

	// default collection to username
	if collection == "" {
		collection = cli.Account.Username
	}

	searchToken, err := cli.Client.GetSearchToken(collection)
	check(err)

	client := algoliasearch.NewClient(searchToken.ApplicationId, searchToken.SearchKey)
	facet := "collection_id:" + searchToken.CollectionId

	oldState, err := terminal.MakeRaw(0)
	check(err)
	defer terminal.Restore(0, oldState)

	var screen = struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	term := terminal.NewTerminal(screen, "Canvas Search > ")
	term.AutoCompleteCallback = func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		query := line + string(key)
		results, err := search(client, query, facet)
		check(err)

		hits := results.Hits
		for _, hit := range hits {
			canvas := &Canvas{}
			canvas.Id = hit["id"].(string)
			canvas.CollectionName = collection
			title := hit["title"].(string)
			url := cli.Client.JoinWebUrl(canvas.WebName())
			fmt.Printf("%-30.30s # %s\n", title, url)
		}

		//newLine = query
		//newPos = pos
		//ok = true
		return
	}

	errc := make(chan error, 2)
	go func() { errc <- readConsole(term) }()
	<-errc
}

func (cli *CLI) SearchUnix(collection string, query string) {
	cli.doAuth()

	// default collection to username
	if collection == "" {
		collection = cli.Account.Username
	}

	searchToken, err := cli.Client.GetSearchToken(collection)
	check(err)

	client := algoliasearch.NewClient(searchToken.ApplicationId, searchToken.SearchKey)
	facet := "collection_id:" + searchToken.CollectionId
	results, err := search(client, query, facet)
	check(err)

	hits := results.Hits
	for _, hit := range hits {
		canvas := &Canvas{}
		canvas.Id = hit["id"].(string)
		canvas.CollectionName = collection
		title := hit["title"].(string)
		url := cli.Client.JoinWebUrl(canvas.WebName())
		fmt.Printf("%-30.30s # %s\n", title, url)
	}
}

func search(client *algoliasearch.Client, query, facet string) (result algoliasearch.SearchResult, err error) {
	filter := "facetFilters=(" + facet + ")"
	client.SetExtraHeader("X-Algolia-QueryParameters", filter)

	index := client.InitIndex("canvases")
	params := make(map[string]interface{})
	params["facetFilters"] = facet
	return index.Search(query, params)
}

func readConsole(term *terminal.Terminal) error {
	for {
		line, err := term.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("terminal.ReadLine: %v", err)
		}
		f := strings.Fields(line)
		if len(f) == 0 {
			continue
		}

		fmt.Println(f)
		if err != nil {
			return err
		}
	}
}
