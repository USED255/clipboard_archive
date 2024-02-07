package fts

import (
	"os"

	"github.com/blevesearch/bleve/v2"
)

var index bleve.Index
var err error

func Load(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return initIndex(path)
	}
	index, err = bleve.Open(path)
	return err
}

func initIndex(path string) error {
	mapping := bleve.NewIndexMapping()
	index, err = bleve.New(path, mapping)
	return err
}
