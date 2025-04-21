package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/schema"
)

// Runner 结构体用于封装 Agent 的运行态
type Runner struct {
	agent *Agent
	input string
}

// NewRunner 函数用于创建新的 Runner 实例
func NewRunner(agent *Agent, input string) *Runner {
	return &Runner{
		agent: agent,
		input: input,
	}
}

// Process 方法用于处理用户输入并以 generate 模式获取模型回复
func (r *Runner) Process(ctx context.Context, isGenerate bool) *schema.Message {
	messages := []*schema.Message{
		schema.SystemMessage(r.agent.instructions),
		schema.UserMessage(r.input),
	}
	if isGenerate {
		result := generate(ctx, r.agent.model, messages)
		return result
	}
	return nil
}

// ProcessStream 方法用于处理用户输入并以 stream 模式获取模型回复
func (r *Runner) ProcessStream(ctx context.Context) {
	messages := []*schema.Message{
		schema.SystemMessage(r.agent.instructions),
		schema.UserMessage(r.input),
	}
	streamResult := stream(ctx, r.agent.model, messages)
	defer streamResult.Close()

	for {
		message, err := streamResult.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println()
				break
			}
			log.Fatalf("Failed to receive stream message: %v", err)
		}
		fmt.Print(message.Content)
	}
}
