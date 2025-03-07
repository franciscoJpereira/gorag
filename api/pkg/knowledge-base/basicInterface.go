package knowledgebase

import (
	"errors"
)

type BasicBase struct {
	collections []string
	data        map[string][]string
}

func NewBasicBase() *BasicBase {
	return &BasicBase{
		make([]string, 0),
		make(map[string][]string),
	}
}

func (b *BasicBase) ListCollections() ([]string, error) {
	return b.collections, nil
}

func (b *BasicBase) CreateCollection(collectionName string) error {
	if _, ok := b.data[collectionName]; ok {
		return errors.New("already exists")
	}
	b.collections = append(b.collections, collectionName)
	b.data[collectionName] = make([]string, 0)
	return nil
}

func (b *BasicBase) AddDataToCollection(collection string, data []string) error {
	b.data[collection] = append(b.data[collection], data...)
	return nil
}

func (b *BasicBase) Retrieve(from string, query string) []string {
	return b.data[from]
}
