package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/planexecute"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func formatPlan(plan *planexecute.ExecutionContext) string {
	if plan == nil {
		return "无"
	}
	b, _ := plan.Plan.MarshalJSON()
	return string(b)
}

func formatExecutedSteps(steps []planexecute.ExecutedStep) string {
	if len(steps) == 0 {
		return "无"
	}
	var sb strings.Builder
	for i, s := range steps {
		sb.WriteString(fmt.Sprintf("%d. 步骤: %v → 结果: %v\n", i+1, s.Step, s.Result))
	}
	return sb.String()
}

func NewPer(ctx context.Context, tool []tool.BaseTool) {
	planner, err := planexecute.NewPlanner(ctx, &planexecute.PlannerConfig{
		ToolCallingChatModel: NewChatModel(ctx),
	})
	if err != nil {
		panic(err)
	}

	executor, err := planexecute.NewExecutor(ctx, &planexecute.ExecutorConfig{
		Model: NewChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tool,
			},
		},
		GenInputFn: func(ctx context.Context, in *planexecute.ExecutionContext) ([]adk.Message, error) {
			planStr := formatPlan(in)
			executedStr := formatExecutedSteps(in.ExecutedSteps)
			stepStr := fmt.Sprintf("%v", in.Plan.FirstStep())

			// 构造系统 + 用户消息
			system := schema.SystemMessage("你是一个认真负责的旅行规划执行器。请严格按照计划中的当前步骤执行，使用可用工具完成任务。在回复中先说明你要执行的步骤，然后调用工具，最后总结执行结果。")
			user := schema.UserMessage(fmt.Sprintf(`
当前计划：
%s

已完成步骤：
%s

你现在需要执行的步骤：
%s

请开始执行。`, planStr, executedStr, stepStr))

			return []adk.Message{system, user}, nil
		},
	})
	if err != nil {
		panic(err)
	}

	replanner, err := planexecute.NewReplanner(ctx, &planexecute.ReplannerConfig{
		ChatModel: NewChatModel(ctx),
	})
	if err != nil {
		panic(err)
	}
	planExecuteAgent, err := planexecute.New(ctx, &planexecute.Config{
		Planner:   planner,
		Executor:  executor,
		Replanner: replanner,
	})
	if err != nil {
		panic(err)
	}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           planExecuteAgent,
		EnableStreaming: true,
	})
	iter := runner.Query(ctx, "帮我规划一个为期3天的北京旅行计划,从今天开始,包括景点推荐，以及未来三天的天气情况。")

	// 记录上一个活动的 Agent，用于判断是否启动了新 Agent
	var (
		currentAgent string
	)

	for {
		event, ok := iter.Next() // 获取下一个事件，阻塞直到有事件或结束
		if !ok {
			fmt.Println("event close")
			break // 迭代器关闭，全部事件已消费
		}
		if event.Err != nil {
			// 处理错误
			fmt.Printf("event err:%s\n", event.Err.Error())
			break
		}

		// 1. 通过 AgentName 判断当前活动的 Agent
		if event.AgentName != currentAgent {
			if currentAgent != "" {
				fmt.Printf("⬅️  Agent 【%s】 工作结束\n", currentAgent)
			}
			fmt.Printf("➡️  启动 Agent: 【%s】\n", event.AgentName)
			currentAgent = event.AgentName
		}

		if event.Output != nil && event.Output.MessageOutput != nil {
			// 处理消息输出（可能是流式）
			msg, err := event.Output.MessageOutput.GetMessage()
			if err != nil {
				fmt.Printf("get msg err:%s\n", err.Error())
				break
			}
			if msg.Content == "" && msg.ReasoningContent == "" {
				fmt.Println("Content ReasoningContent is null")
				continue
			}
			if msg.ReasoningContent != "" {
				//思考内容
				fmt.Printf("%s ReasoningContent: %s \n", currentAgent, msg.ReasoningContent)
			}
			if msg.Content != "" {
				fmt.Printf("%s Content: %s \n", currentAgent, msg.Content)
			}
		}
	}

	// 循环结束后，标记最后一个 Agent 结束工作
	if currentAgent != "" {
		fmt.Printf("⬅️  Agent 【%s】 工作结束\n", currentAgent)
	}
}
