package pkg

import (
	"fmt"
	"io"
	knowledgebase "ragAPI/pkg/knowledge-base"

	"gopkg.in/yaml.v3"
)

/*
	Config Reader
*/

type EchoConfig struct {
	Port  string `yaml:"port"`
	Local bool   `yaml:"local"`
}

type ChromaOptions struct {
	BasePath        string `yaml:"base-url"`
	EmbedderPath    string `yaml:"embedd-url"`
	EmbedderModel   string `yaml:"model"`
	MaxResult       int    `yaml:"max-values"`
	DefaultEmbedder bool   `yaml:"use-default"`
}

type ModelOptions struct {
	Model     string `yaml:"model"`
	ModelPath string `yaml:"model-url"`
}

type StoreOptions struct {
	StorePath string `yaml:"store-path"`
}

// Main configuration options
type T struct {
	Server EchoConfig    `yaml:"echo"`
	Chroma ChromaOptions `yaml:"chroma"`
	Model  ModelOptions  `yaml:"model"`
	Store  StoreOptions  `yaml:"store"`
}

func GetConfiguration(f io.Reader) (t T, err error) {
	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&t); err != nil {
		return
	}
	return
}

func (t T) GetChromaConfig() knowledgebase.ChromaKBOptions {
	return knowledgebase.ChromaKBOptions{
		BasePath:         t.Chroma.BasePath,
		EmbedderPath:     t.Chroma.EmbedderPath,
		EmbedderModel:    t.Chroma.EmbedderModel,
		MaxResults:       t.Chroma.MaxResult,
		DefaultEmbedding: t.Chroma.DefaultEmbedder,
	}
}

func (t T) GetServerConfig() string {
	url := ""
	if t.Server.Local {
		url = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%s", url, t.Server.Port)
}

func (t T) GetStoreConfig() string {
	return t.Store.StorePath
}

func (t T) GetModelConfig() (model string, modelurl string) {
	model, modelurl = t.Model.Model, t.Model.ModelPath
	return
}
