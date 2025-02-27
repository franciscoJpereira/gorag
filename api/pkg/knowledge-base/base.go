package knowledgebase

import (
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/collection"
	"github.com/amikos-tech/chroma-go/openai"
	"github.com/amikos-tech/chroma-go/types"
)

//I'll implement everything on the same file
//because I only plan to use ChromaDB

type BaseInterface interface {
	//Creates a new collection
	CreateColletion(collectionName string) error
	//Adds data to a collection
	AddDataToCollection(collection string, data []string) error
	//Retrieves documents that are mos relevant to a certain query
	Retrieve(query string) []string
}

type ChromaKBOptions struct {
	BasePath      string `json:"ChromaURL"`
	EmbedderPath  string `json:"ModelURL"`
	EmbedderModel string `json:"EmbeddModel"`
	OpenAIAPIKey  string `json:"OpenAIKey"`
}

type ChromaKB struct {
	ctx     context.Context
	client  *chromago.Client
	options ChromaKBOptions
}

func NewChromaKB(
	ctx context.Context,
	options ChromaKBOptions,
) (*ChromaKB, error) {
	client, err := chromago.NewClient(options.BasePath)
	if err != nil {
		return nil, err
	}
	return &ChromaKB{
		ctx,
		client,
		options,
	}, nil
}

func (c *ChromaKB) EmbeddFunction() (types.EmbeddingFunction, error) {
	return openai.NewOpenAIEmbeddingFunction(
		c.options.OpenAIAPIKey,
		openai.WithBaseURL(c.options.EmbedderPath),
		openai.WithModel(openai.EmbeddingModel(c.options.EmbedderModel)),
	)
}

func (c *ChromaKB) CreateColletion(collectionName string) error {
	_, err := c.client.NewCollection(
		c.ctx,
		collection.WithName(collectionName),
	)
	return err
}

func (c *ChromaKB) AddDataToCollection(collection string, data []string) error {
	embeddingFunction, err := c.EmbeddFunction()
	if err != nil {
		return err
	}
	db, err := c.client.GetCollection(c.ctx, collection, embeddingFunction)
	if err != nil {
		return err
	}
	rs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(db.EmbeddingFunction),
		types.WithIDGenerator(types.NewULIDGenerator()),
	)
	if err != nil {
		return err
	}
	for _, entry := range data {
		rs.WithRecord(types.WithDocument(entry))
	}
	if _, err = rs.BuildAndValidate(c.ctx); err != nil {
		return err
	}
	if _, err = db.AddRecords(c.ctx, rs); err != nil {
		return err
	}
	return nil
}

func (c *ChromaKB) Retrieve(query string) []string {
	return nil
}
