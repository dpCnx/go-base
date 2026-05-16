package tools

import (
	"context"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type TimeTool struct {
}

func NewTimeTool() *TimeTool {
	return &TimeTool{}
}

func (w *TimeTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	return time.Now().Format(time.DateTime), nil
}

func (w *TimeTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get_time",
		Desc: "获取当前时间",
	}, nil
}
