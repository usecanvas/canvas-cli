package main

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)

func (cli *CLI) Search(collection string, query string) {
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
	results, err := index.Search(query, params)
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
