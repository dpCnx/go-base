package model

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type MysqlTool struct {
}

func NewMysqlTool() *MysqlTool {
	return &MysqlTool{}
}

func (m *MysqlTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	var params map[string]any
	if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
		return "", err
	}
	id, ok := params["id"].(string)
	if !ok {
		return "", fmt.Errorf("id is required")
	}
	return fmt.Sprintf("%s 订单已经发货了", id), nil
}

func (m *MysqlTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:        "GetOrderById",
		Desc:        "数据库查询助手",
		ParamsOneOf: schema.NewParamsOneOfByParams(m.Params()),
	}, nil
}

func (m *MysqlTool) Params() map[string]*schema.ParameterInfo {
	return map[string]*schema.ParameterInfo{
		"id": {
			Desc:     "需要查询的订单id",
			Type:     schema.String,
			Required: true,
		},
		"extensions": {
			Desc: "筛选条件",
			Type: schema.String,
		},
	}
}
