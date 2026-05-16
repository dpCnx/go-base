package template

import (
	"context"

	"github.com/cloudwego/eino/schema"
)

func NewChatTemplate(ctx context.Context, history []*schema.Message) (result []*schema.Message, err error) {
	//chatTemplate := prompt.FromMessages(
	//	schema.FString,
	//	schema.SystemMessage(BaseSystemPrompt),
	//	schema.MessagesPlaceholder(historyKey, true),
	//	//schema.UserMessage("{question}"), // agent的时候不需要
	//)
	//return chatTemplate.Format(ctx, map[string]any{
	//	"role":       "一个可爱的客服",
	//	"ragContext": "",
	//	historyKey:   history,
	//	//"question":   "想查询订单号为89的订单,是否发货",
	//})
	return nil, err
}
