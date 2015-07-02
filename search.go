package main

import "fmt"
import "github.com/algolia/algoliasearch-client-go/algoliasearch"

func (cli *CLI) Search(collection string) {
	cli.doAuth()

	// default collection to username
	if collection == "" {
		collection = cli.Account.Username
	}

	searchToken, err := cli.Client.GetSearchToken(collection)
	check(err)

	client := algoliasearch.NewClient(searchToken.ApplicationId, searchToken.SearchKey)

	facet := "collection_id:" + searchToken.CollectionId
	filter := "facetFilters=(" + facet + ")"
	client.SetExtraHeader("X-Algolia-QueryParameters", filter)

	index := client.InitIndex("canvases")
	params := make(map[string]interface{})
	params["facetFilters"] = facet
  results, err := index.Search("hello", params)
	check(err)

	hits := results.(map[string]interface{})["hits"]
	fmt.Println(hits)
}
