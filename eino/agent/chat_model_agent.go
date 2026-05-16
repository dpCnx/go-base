package agent

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

const BaseSystemPrompt = `
# 角色与目标
你是一个AI智能助手。你的任务是理解用户的需求，通过一系列的“思考”和“行动”来解决问题。
{role}

# 知识库
-------如果以下知识库有内容，优先匹配知识库数据进行回答----------
{ragContext}
-----------------
`

func NewChatModelAgent(ctx context.Context, tools []tool.BaseTool, history []*schema.Message) *adk.ChatModelAgent {

	var optional bool
	if len(history) == 0 {
		optional = true
	}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "chatModelAgent",
		Description: "chatModelAgent",
		Instruction: BaseSystemPrompt,
		Model:       NewChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
		GenModelInput: func(ctx context.Context, instruction string, input *adk.AgentInput) ([]adk.Message, error) {
			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage(instruction),
				schema.MessagesPlaceholder("history", optional),
			)
			messages, err := template.Format(ctx, map[string]any{
				"role":       "小可爱风格客服",
				"ragContext": "",
				"history":    history,
			})
			if err != nil {
				return nil, err
			}
			messages = append(messages, input.Messages...)
			return messages, nil

		},
	})
	if err != nil {
		panic(err)
	}

	return agent
}

func RunAgent(ctx context.Context, agent *adk.ChatModelAgent) {
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})
	events := runner.Query(ctx, "你好，你是谁")
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
			//msg, err := event.Output.MessageOutput.GetMessage()
			//if err != nil {
			//	return
			//}
			//if msg.Content == "" && msg.ReasoningContent == "" {
			//	fmt.Println("Content ReasoningContent is null")
			//	continue
			//}
			//if msg.ReasoningContent != "" {
			//	//思考内容
			//	fmt.Println("ReasoningContent:", msg.ReasoningContent)
			//}
			//if msg.Content != "" {
			//	fmt.Println("Content:", msg.Content)
			//}

			out := event.Output.MessageOutput
			if out.IsStreaming {
				out.MessageStream.SetAutomaticClose()
				for {
					frame, err := out.MessageStream.Recv()
					if errors.Is(err, io.EOF) {
						fmt.Println("out.MessageStream.Recv() err", err)
						break
					}
					if err != nil {
						fmt.Println("out.MessageStream.Recv() err!= nil", err)
						return
					}
					if frame != nil && frame.ReasoningContent != "" {
						fmt.Println("frame.ReasoningContent", frame.ReasoningContent)
					}
					if frame != nil && frame.Content != "" {
						fmt.Println("frame.Content", frame.Content)
					}
				}
				continue

			}

			if out.Message != nil {
				fmt.Println("out.Message", out.Message)
			}

		}
	}
}
