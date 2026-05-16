package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type WeatherResponse struct {
	City        string `json:"city"`
	Temperature int    `json:"temperature"`
	Condition   string `json:"condition"`
	Date        string `json:"date"`
}

type WeatherTool struct {
}

func NewWeatherTool() *WeatherTool {
	return &WeatherTool{}
}

func (w *WeatherTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	var params map[string]any
	if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
		return "", err
	}
	city, ok := params["city"].(string)
	if !ok {
		return "", fmt.Errorf("city is required")
	}
	date, ok := params["date"].(string)
	if !ok {
		return "", fmt.Errorf("date is required")
	}
	res := &WeatherResponse{
		City: city, Temperature: 25, Condition: "晴朗", Date: date,
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(jsonRes), nil
}

func (w *WeatherTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:        "get_weather",
		Desc:        "获取指定城市在指定日期的天气信息.",
		ParamsOneOf: schema.NewParamsOneOfByParams(w.Params()),
	}, nil
}

func (w *WeatherTool) Params() map[string]*schema.ParameterInfo {
	return map[string]*schema.ParameterInfo{
		"city": {
			Desc:     "城市名称",
			Type:     schema.String,
			Required: true,
		},
		"date": {
			Desc:     "日期，格式 YYYY-MM-DD",
			Type:     schema.String,
			Required: true,
		},
	}
}
