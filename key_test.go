package fungident

import (
	"testing"
)

var testShroom1 = Fungus{
	GenusName:   "Test Shroom",
	SpeciesName: "One",
	Traits: Traits{
		Hymenium:        "gills",
		SporePrintColor: "nope",
	},
}

var testShroom2 = Fungus{
	GenusName:   "Test Shroom",
	SpeciesName: "Two",
	Traits: Traits{
		Hymenium:        "gills",
		SporePrintColor: "white",
	},
}

var testShroom3 = Fungus{
	GenusName:   "Test Shroom",
	SpeciesName: "Three",
	Traits: Traits{
		Hymenium: "ridges",
	},
}

func TestIdentify(t *testing.T) {
	var searchOpts Traits
	archive, _ := NewArchive(nil)

	searchOpts.SporePrintColor = "brown"
	results := archive.Identify(searchOpts)

	expectedStr := "Poison Pie"
	if results[0].CommonNames[0] != expectedStr {
		t.Errorf("Expected first item in returned results to be %s. Got %s", expectedStr, results[0].CommonNames[0])
	}

	searchOpts.SporePrintColor = "white"
	results = archive.Identify(searchOpts)

	expectedStr = "Destroying Angel"
	if results[0].CommonNames[0] != expectedStr {
		t.Errorf("Expected first item in returned results to be %s. Got %s", expectedStr, results[0].CommonNames[0])
	}

	searchOpts.SporePrintColor = "purple"
	results = archive.Identify(searchOpts)

	expectedStr = "King Stropharia"
	if results[0].CommonNames[0] != expectedStr {
		t.Errorf("Expected first item in returned results to be %s. Got %s", expectedStr, results[0].CommonNames[0])
	}
}

func TestRankResults(t *testing.T) {
	archive, _ := NewArchive(nil)
	result := archive.Identify(testShroom1.Traits)

	if result[0].TaxonomicName() != "Hebeloma Crustuliniforme" {
		t.Errorf("Expected gilled fungi in default order. Got: %+v", result)
	}
	if len(result) > 3 {
		t.Errorf("Should only have gilled fungi in result (3 total), but found %d total results", len(result))
	}

	result = archive.Identify(testShroom2.Traits)

	if result[0].TaxonomicName() != "Amanita Bisporigera" {
		t.Errorf("Expected more specific trait options to return destroying angel first. Got %s", result[0].TaxonomicName())
	}

	result = archive.Identify(testShroom3.Traits)

	if len(result) != 0 {
		t.Errorf("Should not find any fungi with ridges, but got results: %+v", result)
	}
}
