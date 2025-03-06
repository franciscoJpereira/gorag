package apiinterface

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIChatModel struct {
	ctx     context.Context
	chat    *openai.Client
	options []option.RequestOption
	model   string
}

func NewOpenAIChatModel(
	ctx context.Context,
	hostname string,
	model string,
	options ...option.RequestOption) *OpenAIChatModel {

	service := openai.NewClient(
		option.WithBaseURL(hostname),
	)

	return &OpenAIChatModel{
		ctx,
		service,
		options,
		model,
	}
}

func (s *OpenAIChatModel) Models() (modelsNames []string, err error) {
	models, err := s.chat.Models.List(
		s.ctx,
	)
	for _, model := range models.Data {
		modelsNames = append(modelsNames, model.ID)
	}
	return
}

func processMessages(chat ...ChatMessage) (completion []openai.ChatCompletionMessageParamUnion) {
	completion = make([]openai.ChatCompletionMessageParamUnion, len(chat))
	for index, message := range chat {
		var value openai.ChatCompletionMessageParamUnion
		if message.Role == "user" {
			value = openai.UserMessage(message.Content)
		} else {
			value = openai.SystemMessage(message.Content)
		}

		completion[index] = value
	}
	return
}

// / Sends a new message to the model
func (s *OpenAIChatModel) Send(
	message string,
	chat ...ChatMessage,
) (response ChatMessage, err error) {
	completion, err := s.chat.Chat.Completions.New(
		s.ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F(
				append(processMessages(chat...),
					openai.UserMessage(message),
				),
			),
			Model: openai.F(s.model),
		},
		s.options...,
	)
	response = ChatMessage{
		Content: completion.Choices[0].Message.Content,
		Role:    string(completion.Choices[0].Message.Role),
	}

	return
}
