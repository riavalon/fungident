package fungident

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
)

// Fungus Model that represents a single fungus
type Fungus struct {
	CommonNames []string
	SpeciesName string
	GenusName   string
	Edibility   string
	Traits      Traits
}

// TaxonomicName will return the genus and species name together as
// one string
func (f Fungus) TaxonomicName() string {
	return strings.Join([]string{f.GenusName, f.SpeciesName}, " ")
}

// Archive will act as a handler for interacting with a given
// database for tracking, querying, and working with fungi
// persisted in whatever database setup the user has
type Archive struct {
	db *sql.DB

	tablename string
}

// ArchiveOption used to modify the newly created archive. Meant
// to be used with the `NewArchive` function as functional options.
type ArchiveOption func(*Archive)

// UseTableName A functional option for configuring what name
// the archive will initialize the fungi with in the database.
func UseTableName(name string) ArchiveOption {
	return func(arc *Archive) {
		arc.tablename = name
	}
}

// NewArchive will return an archive handler for querying,
// creating, updating, and deleting fungal records.
func NewArchive(db *sql.DB, opts ...ArchiveOption) (*Archive, error) {
	ret := &Archive{db: db}

	for _, opt := range opts {
		opt(ret)
	}

	err := ret.initializeDB()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (arc *Archive) initializeDB() error {
	tablename := "fungident_fungi"

	if arc.tablename != "" {
		tablename = arc.tablename
	}

	// TODO: Make this not dumb
	// .. shouldn't be using raw SQL statements like this.
	// assumes too much about what database is being used.
	statement := fmt.Sprintf(`CREATE TABLE %s (
		fungusID int NOT NULL AUTO_INCREMENT,
		commonNames longtext,
		genusName varchar(255),
		speciesName varchar(255),
		edibility varchar(255),
		traits int FOREIGN KEY
	)`, tablename)

	_ = statement

	return nil
}

// Identify will take a Traits struct and find all fungi that match
// selected traits and place the results as a fungus slice
func (arc *Archive) Identify(opts Traits) []Fungus {
	result := rankResults(opts, fungi)
	return result
}

// Traits are used to categorize and search for fungi
type Traits struct {
	SporePrintColor string
	Hymenium        string
}

// Example fungi for use in testing at the moment.
// will eventually come from the database.
var fungi = []Fungus{
	Fungus{
		CommonNames: []string{"Poison Pie", "Fairy Cakes"},
		GenusName:   "Hebeloma",
		SpeciesName: "Crustuliniforme",
		Edibility:   "Poisonous",
		Traits: Traits{
			SporePrintColor: "brown",
			Hymenium:        "gills",
		},
	},
	Fungus{
		CommonNames: []string{"Destroying Angel", "Angel of Death"},
		GenusName:   "Amanita",
		SpeciesName: "Bisporigera",
		Edibility:   "Deadly Toxic",
		Traits: Traits{
			SporePrintColor: "white",
			Hymenium:        "gills",
		},
	},
	Fungus{
		CommonNames: []string{"King Stropharia", "Garden Giant", "Winecap Mushroom"},
		GenusName:   "Stropharia",
		SpeciesName: "Rugosoannulata",
		Edibility:   "Choice Edible",
		Traits: Traits{
			SporePrintColor: "purple",
			Hymenium:        "gills",
		},
	},
	Fungus{
		CommonNames: []string{"Violet-Toothed Polypore"},
		GenusName:   "Trichaptum",
		SpeciesName: "Biforme",
		Edibility:   "Inedible",
		Traits: Traits{
			SporePrintColor: "white",
			Hymenium:        "pores",
		},
	},
}

func rankResults(criteria Traits, f []Fungus) []Fungus {
	var matches []struct {
		rank   uint32
		result Fungus
	}
	var rankedResults []Fungus
	for _, fungus := range f {
		var rank uint32
		if criteria.SporePrintColor == fungus.Traits.SporePrintColor {
			rank++
		}

		if criteria.Hymenium == fungus.Traits.Hymenium {
			rank++
		}

		if rank > 0 {
			matches = append(matches, struct {
				rank   uint32
				result Fungus
			}{
				rank:   rank,
				result: fungus,
			})
		}
	}

	// greatest to least rank
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].rank > matches[j].rank
	})

	for _, rankedRes := range matches {
		rankedResults = append(rankedResults, rankedRes.result)
	}

	return rankedResults
}
