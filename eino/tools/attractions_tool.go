package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type AttractionResponse struct {
	Attractions []string `json:"attractions"`
}

type AttractionsTool struct {
}

func NewAttractionsTool() *AttractionsTool {
	return &AttractionsTool{}
}

func (w *AttractionsTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	var params map[string]any
	if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
		return "", err
	}
	city, ok := params["city"].(string)
	if !ok {
		return "", fmt.Errorf("city is required")
	}
	category, ok := params["category"].(string)
	if !ok {
		return "", fmt.Errorf("category is required")
	}

	fmt.Println("AttractionsTool:", city, category)

	res := &AttractionResponse{
		Attractions: []string{"故宫", "长城", "天坛"},
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(jsonRes), nil
}

func (w *AttractionsTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:        "search_attractions",
		Desc:        "搜索指定城市的旅游景点.",
		ParamsOneOf: schema.NewParamsOneOfByParams(w.Params()),
	}, nil
}

func (w *AttractionsTool) Params() map[string]*schema.ParameterInfo {
	return map[string]*schema.ParameterInfo{
		"city": {
			Desc:     "城市名称",
			Type:     schema.String,
			Required: true,
		},
		"category": {
			Desc:     "景点类别",
			Type:     schema.String,
			Required: true,
		},
	}
}
