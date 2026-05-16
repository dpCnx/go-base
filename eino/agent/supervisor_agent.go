package agent

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/supervisor"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

func NewSupervisorAgent(ctx context.Context) *adk.ChatModelAgent {
	sv, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "supervisor",
		Description: "the agent responsible to supervise tasks",
		Instruction: `
			你是一名主管，负责管理两名员工：
			- 一名研究专员。将与研究相关的工作分配给该专员
			- 一名数学专员。将与数学相关的工作分配给该专员
			每次只给一个专员分配工作，不要同时呼叫多个专员。
			不要自己进行任何工作。`,
		Model: NewChatModel(ctx),
	})
	if err != nil {
		panic(err)
	}
	return sv
}

func BuildSearchAgent(ctx context.Context, tools []tool.BaseTool) (adk.Agent, error) {

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "research_agent",
		Description: "the agent responsible to search the internet for info",
		Instruction: `
			你是一名研究人员。
			
			说明：
			- 仅协助完成与研究相关的任务，切勿进行任何数学计算
			- 请勿估算任何数字
			- 完成任务后，请直接回复主管
			- 只需回复您的工作成果，切勿添加任何其他文字。`,
		Model: NewChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
				UnknownToolsHandler: func(ctx context.Context, name, input string) (string, error) {
					return fmt.Sprintf("unknown tool: %s", name), nil
				},
			},
		},
	})
}

func BuildMathAgent(ctx context.Context, tools []tool.BaseTool) (adk.Agent, error) {
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "math_agent",
		Description: "the agent responsible to do math",
		Instruction: `
			你是一名数学专家。
			
			说明：
			- 仅协助处理与数学相关的任务
			- 完成任务后，请直接回复主管
			- 只回复工作结果，切勿添加任何其他文字。`,
		Model: NewChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
				UnknownToolsHandler: func(ctx context.Context, name, input string) (string, error) {
					return fmt.Sprintf("unknown tool: %s", name), nil
				},
			},
		},
	})
}

func Run(ctx context.Context) {
	supervisorAgent := NewSupervisorAgent(ctx)
	searchAgent, err := BuildSearchAgent(ctx, []tool.BaseTool{})
	if err != nil {
		panic(err)
	}
	mathAgent, err := BuildMathAgent(ctx, []tool.BaseTool{})
	if err != nil {
		panic(err)
	}

	agents, err := supervisor.New(ctx, &supervisor.Config{
		Supervisor: supervisorAgent,
		SubAgents:  []adk.Agent{searchAgent, mathAgent},
	})
	if err != nil {
		panic(err)
	}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agents,
		EnableStreaming: true,
	})
	iter := runner.Query(ctx, "还记得我刚才问了什么问题吗")

	for {
		event, ok := iter.Next() // 获取下一个事件，阻塞直到有事件或结束
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
				fmt.Printf("get msg err:%s\n", err.Error())
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
			if msg.Content != "" {
				fmt.Println("Content:", msg.Content)
			}
		}
	}
}
