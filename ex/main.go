package main

import (
	"fmt"
	"strings"

	"riavalon.com/fungident"
)

func main() {
	var searchOptions fungident.Traits
	searchOptions.SporePrintColor = "white"

	// no db handler for now. Also ignoring error for now
	archive, _ := fungident.NewArchive(nil)
	results := archive.Identify(searchOptions)

	for _, res := range results {
		var name string
		var commonNames string
		if len(res.CommonNames) > 0 {
			name = res.CommonNames[0]
			commonNames = strings.Join(res.CommonNames, ", ")
		}
		fmt.Printf("%s:\nCommon Names: %s\nTaxonomic Name: %s\nEdibility: %s\n",
			name, commonNames, res.TaxonomicName(), res.Edibility)
	}
}
