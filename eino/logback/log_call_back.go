package logback

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/schema"
)

/*
	使用
	callbacks.AppendGlobalHandlers(callback)
*/

func LogCallback() callbacks.Handler {
	builder := callbacks.NewHandlerBuilder()
	builder.OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		fmt.Printf("[view OnStartFn start]:[%s:%s:%s]\n", info.Component, info.Type, info.Name)
		b, _ := json.Marshal(input)
		fmt.Printf("%s\n", string(b))
		return ctx
	})
	builder.OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		fmt.Printf("[view OnEndFn end]:[%s:%s:%s]\n", info.Component, info.Type, info.Name)
		return ctx
	})
	builder.OnErrorFn(func(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
		fmt.Printf("[view OnErrorFn end]:[%s:%s:%s]\n", info.Component, info.Type, info.Name)
		return ctx
	})
	builder.OnStartWithStreamInputFn(func(ctx context.Context, info *callbacks.RunInfo, input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
		fmt.Printf("[view OnStartWithStreamInputFn start]:[%s:%s:%s]\n", info.Component, info.Type, info.Name)
		return ctx
	})
	builder.OnEndWithStreamOutputFn(func(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
		fmt.Printf("[view OnEndWithStreamOutputFn end]:[%s:%s:%s]\n", info.Component, info.Type, info.Name)
		return ctx
	})
	return builder.Build()
}
