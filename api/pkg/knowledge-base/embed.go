package knowledgebase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/amikos-tech/chroma-go/types"
)

type LmEmbeddFunction struct {
	url   string
	model string
}

func NewLmEmbeddFunction(url string, model string) *LmEmbeddFunction {
	return &LmEmbeddFunction{
		url,
		model,
	}
}

func (l *LmEmbeddFunction) Embed(text string) ([]float32, error) {
	rBody := map[string]string{
		"model": l.model,
		"input": text,
	}
	bodyReader, err := json.Marshal(rBody)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(l.url, "application/json", bytes.NewReader(bodyReader))
	if err != nil {
		return nil, err
	}
	var decoded map[string]any
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&decoded); err != nil {
		return nil, err
	}
	data := decoded["data"].([]interface{})[0].(map[string]interface{})["embedding"].([]interface{})
	result := make([]float32, len(data))
	for index, value := range data {
		result[index] = float32(value.(float64))
	}
	return result, nil
}

func (l *LmEmbeddFunction) EmbedDocuments(ctx context.Context, texts []string) ([]*types.Embedding, error) {
	embeddings := make([]*types.Embedding, len(texts))
	for index, text := range texts {
		embedding, err := l.EmbedQuery(ctx, text)
		if err != nil {
			return nil, err
		}
		embeddings[index] = embedding
	}
	return embeddings, nil
}

type queryResponse struct {
	Embedding *types.Embedding
	Err       error
}

func (l *LmEmbeddFunction) embedQuery(text string) <-chan queryResponse {
	channel := make(chan queryResponse)
	runQuery := func() {
		response, err := l.Embed(text)
		channel <- queryResponse{
			Embedding: &types.Embedding{ArrayOfFloat32: &response},
			Err:       err,
		}
	}
	go runQuery()
	return channel
}

func (l *LmEmbeddFunction) EmbedQuery(ctx context.Context, text string) (*types.Embedding, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("context Done")
	case result := <-l.embedQuery(text):
		return result.Embedding, result.Err
	}
}

func (l *LmEmbeddFunction) EmbedRecords(ctx context.Context, records []*types.Record, force bool) error {
	for _, record := range records {
		if !force && (record.Embedding.ArrayOfFloat32 != nil || record.Embedding.ArrayOfInt32 != nil) {
			continue
		}
		embedding, err := l.EmbedQuery(ctx, record.Document)
		if err != nil {
			return err
		}
		record.Embedding = *embedding
	}
	return nil
}
