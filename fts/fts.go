package fts

import (
	"os"

	"github.com/blevesearch/bleve/v2"
)

var Index bleve.Index

func Load(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return initIndex(path)
	}
	Index, err = bleve.Open(path)
	return err
}

func Close() {
	Index.Close()
}
