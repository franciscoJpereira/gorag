package main

import (
	"context"
	"fmt"
	apiinterface "ragAPI/pkg/apiInterface"
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

func main() {
	TestMessage("Hello there!")
}
