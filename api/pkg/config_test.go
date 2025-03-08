package pkg_test

import (
	"os"
	"path/filepath"
	"ragAPI/pkg"
	"testing"
)

const (
	CONFIG_TEST_PATH = "../test/configtest.yaml"
)

func TestConfig(t *testing.T) {
	p, _ := filepath.Abs(CONFIG_TEST_PATH)
	f, _ := os.Open(p)
	c, err := pkg.GetConfiguration(f)
	if err != nil {
		t.Fatalf("Failed to load configuration: %s\n", err)
	}
	if !c.Server.Local {
		t.Fatal("Server was read as false")
	}
	if c.Server.Port != "1234" {
		t.Fatalf("Port was %s instead of 1234", c.Server.Port)
	}
	if c.Chroma.BasePath != "basepath" {
		t.Fatalf("Read %s from base path", c.Chroma.BasePath)
	}
	if c.Chroma.EmbedderPath != "embedderpath" {
		t.Fatalf("Read %s from embedder path", c.Chroma.EmbedderPath)
	}
	if c.Chroma.EmbedderModel != "embeddermodel" {
		t.Fatalf("Read %s from embedder model", c.Chroma.EmbedderModel)
	}
	if c.Chroma.MaxResult != 10 {
		t.Fatalf("Max results %d", c.Chroma.MaxResult)
	}
	if c.Model.Model != "model" {
		t.Fatalf("Read %s from model", c.Model.Model)
	}
	if c.Model.ModelPath != "modelurl" {
		t.Fatalf("Read %s from model url", c.Model.ModelPath)
	}
	if c.Store.StorePath != "storepath" {
		t.Fatalf("Read %s from store path", c.Store.StorePath)
	}
}
