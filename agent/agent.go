package agent

import (
	"context"
	"log"

	"github.com/cloudwego/eino/schema"
	"github.com/menghuiqiang777/myeino/provider"
)

// Agent 结构体定义
type Agent struct {
	Name         string
	Instructions string
	chatModel    schema.ChatModel
}

// NewAgent 创建一个新的 Agent 实例
func NewAgent(ctx context.Context, name, instructions, vendor, modelName string) (*Agent, error) {
	chatModel, err := provider.GetChatModel(vendor, modelName)
	if err != nil {
		return nil, err
	}
	return &Agent{
		Name:         name,
		Instructions: instructions,
		chatModel:    chatModel,
	}, nil
}

// GenerateResponse 方法用于生成响应
func (a *Agent) GenerateResponse(ctx context.Context, inputText string) schema.StreamReader {
	messages := []*schema.Message{
		schema.SystemMessage(a.Instructions),
		schema.UserMessage(inputText),
	}
	return stream(ctx, a.chatModel, messages)
}

// 原有的 stream 函数
func stream(ctx context.Context, chatModel schema.ChatModel, messages []*schema.Message) schema.StreamReader {
	streamResult, err := chatModel.Stream(ctx, messages)
	if err != nil {
		log.Fatalf("Failed to stream: %v", err)
	}
	return streamResult
}

// Runner 结构体定义
type Runner struct{}

// Run 方法用于运行 Agent 并返回结果
func (r *Runner) Run(ctx context.Context, agent *Agent, inputText string) schema.StreamReader {
	return agent.GenerateResponse(ctx, inputText)
}
