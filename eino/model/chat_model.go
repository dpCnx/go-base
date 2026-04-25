package model

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

func NewChatModel(ctx context.Context) model.ToolCallingChatModel {

	var temperature float32 = 0.7
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:       "spark-x",
		BaseURL:     "https://spark-api-open.xf-yun.com/v2",
		APIKey:      "rYADaYejVoqTOATJMCig:uxBfAyNWTAJDPglZHxzE",
		Temperature: &temperature,
	})
	if err != nil {
		panic(err)
	}
	return chatModel

}
