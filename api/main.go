package main

import (
	"context"
	"fmt"
	apiinterface "ragAPI/pkg/apiInterface"
	knowledgebase "ragAPI/pkg/knowledge-base"
)

func TestMessage(message string) {
	service := apiinterface.NewOpenAIChatModel(
		context.Background(),
		"http://localhost:1234/v1/",
		"deepseek-r1-distill-llama-8b@q6_k")
	response, err := service.Send(message)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Response: %s\n", response.Content)
	}
	models, _ := service.Models()
	fmt.Printf("Models: %v\n", models)
}

func TestChroma() {
	chromaOptions := knowledgebase.ChromaKBOptions{
		BasePath:      "http://localhost:8000",
		EmbedderPath:  "http://localhost:1234/v1/embeddings",
		EmbedderModel: "text-embedding-granite-embedding-278m-multilingual",
		MaxResults:    10,
	}
	chromadb, err := knowledgebase.NewChromaKB(context.Background(), chromaOptions)
	if err != nil {
		panic(err)
	}
	if err = chromadb.CreateColletion("TestingCollection"); err != nil {
		panic(fmt.Sprintf("Creating collection: %s", err))
	}
	if err = chromadb.AddDataToCollection("TestingCollection", []string{"South park es muy divertido", "Un cuento de sueños locos", "Un perro caminando por la calle no es bueno"}); err != nil {
		panic(fmt.Sprintf("Adding data to collection: %s", err))
	}
	result := chromadb.Retrieve("TestingCollection", "Perros")
	fmt.Println(result)
}

func main() {
	//TestMessage("Hello there!")
	TestChroma()
}
