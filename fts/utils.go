package fts

import "github.com/blevesearch/bleve/v2"

var err error

func initIndex(path string) error {
	mapping := bleve.NewIndexMapping()
	Index, err = bleve.New(path, mapping)
	return err
}
