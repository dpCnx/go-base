package model

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

const historyKey = "history_key"

const BaseSystemPrompt = `
# 角色与目标
你是一个AI智能客服。你的任务是理解用户的需求，通过一系列的“思考”和“行动”来解决回答用户的问题。
{role}
`

func NewChatTemplate(ctx context.Context, history []*schema.Message) (result []*schema.Message, err error) {
	chatTemplate := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage(BaseSystemPrompt),
		schema.MessagesPlaceholder(historyKey, true),
		//schema.UserMessage("{question}"), // agent的时候不需要
	)
	return chatTemplate.Format(ctx, map[string]any{
		"role":       "一个可爱的客服",
		"ragContext": "",
		historyKey:   history,
		//"question":   "想查询订单号为89的订单,是否发货",
	})
}
