package pkg

import (
	"io"

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
	BasePath      string `yaml:"base-url"`
	EmbedderPath  string `yaml:"embedd-url"`
	EmbedderModel string `yaml:"model"`
	MaxResult     int    `yaml:"max-values"`
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
