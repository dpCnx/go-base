package main

import (
	"context"
	"fmt"
	"go-base/eino/model"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func main() {

	ctx := context.Background()

	chatModel := model.NewChatModel(ctx)
	//toolInfo, err := model.NewMysqlTool().Info(ctx)
	//if err != nil {
	//	panic(err)
	//}
	//toolCallingChatModel, err := chatModel.WithTools([]*schema.ToolInfo{toolInfo})
	//if err != nil {
	//	panic(err)
	//}
	//
	result, err := model.NewChatTemplate(ctx, []*schema.Message{})
	if err != nil {
		panic(err)
		return
	}
	//stream, err := toolCallingChatModel.Stream(ctx, result)
	//if err != nil {
	//	panic(err)
	//}
	//defer stream.Close()
	//
	//for {
	//	chunk, err := stream.Recv()
	//	if err != nil {
	//		break
	//	}
	//	fmt.Printf("chunk.ReasoningContent:%s \n", chunk.ReasoningContent)
	//	fmt.Printf("chunk.Role:%s \n", chunk.Role)
	//	fmt.Printf("chunk.ToolCalls:%+v \n", chunk.ToolCalls)
	//	fmt.Printf("chunk.ToolCallID:%+v \n", chunk.ToolCallID)
	//	fmt.Printf("chunk.ToolName:%+v \n", chunk.ToolName)
	//	fmt.Printf("chunk.Name:%+v \n", chunk.Name)
	//	fmt.Printf("chunk.Content:%+v \n", chunk.Content)
	//	fmt.Println("-----------------------------------------------------------")
	//}

	// agent
	toolInfo := model.NewMysqlTool()
	agent, err := model.NewAdkChatModel(ctx, chatModel, []tool.BaseTool{toolInfo}, result)
	if err != nil {
		panic(err)
	}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})
	events := runner.Query(ctx, "想查询订单号为89的订单,是否发货")
	for {
		event, ok := events.Next() // 获取下一个事件，阻塞直到有事件或结束
		if !ok {
			fmt.Println("event close")
			break // 迭代器关闭，全部事件已消费
		}
		if event.Err != nil {
			// 处理错误
			fmt.Printf("event err:%s\n", event.Err.Error())
			return
		}
		if event.Output != nil && event.Output.MessageOutput != nil {
			// 处理消息输出（可能是流式）
			msg, err := event.Output.MessageOutput.GetMessage()
			if err != nil {
				return
			}
			if msg.Content == "" && msg.ReasoningContent == "" {
				fmt.Println("Content ReasoningContent is null")
				continue
			}
			if msg.ReasoningContent != "" {
				//思考内容
				fmt.Println("ReasoningContent:", msg.ReasoningContent)
			}
			fmt.Printf("Agent名称[%s], 工具名称:[%s], 模型返回内容: %s \n", msg.Name, msg.ToolName, msg.Content)
			if msg.Content != "" {
				fmt.Println("Content:", msg.Content)
			}
		}
	}
}
