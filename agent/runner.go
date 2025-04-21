package agent

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/menghuiqiang777/myeino/llmModel"
)

// Runner 结构体封装了 Agent 的运行态，包含一个 Agent 实例和用户输入
type Runner struct {
	agent *Agent
	input string
}

// NewRunner 创建一个新的 Runner 实例
func NewRunner(agent *Agent, input string) *Runner {
	return &Runner{
		agent: agent,
		input: input,
	}
}

// Process 处理用户输入并以 generate 模式获取模型回复
// 如果 isGenerate 为 true，则调用 llmModel.Generate 方法获取回复
// 返回回复消息和可能的错误
func (r *Runner) Process(ctx context.Context, isGenerate bool) (*schema.Message, error) {
	messages := []*schema.Message{
		schema.SystemMessage(r.agent.Instructions),
		schema.UserMessage(r.input),
	}
	if isGenerate {
		result, err := llmModel.Generate(ctx, r.agent.Model, messages)
		if err != nil {
			return nil, fmt.Errorf("failed to generate response: %w", err)
		}
		return result, nil
	}
	return nil, nil
}

// ProcessStream 处理用户输入并以 stream 模式获取模型回复
// 调用 llmModel.Stream 方法获取流式回复，并逐块输出
// 返回可能的错误
func (r *Runner) ProcessStream(ctx context.Context) error {
	messages := []*schema.Message{
		schema.SystemMessage(r.agent.Instructions),
		schema.UserMessage(r.input),
	}
	streamResult, err := llmModel.Stream(ctx, r.agent.Model, messages)
	if err != nil {
		return fmt.Errorf("failed to start streaming: %w", err)
	}
	defer streamResult.Close()

	for {
		message, err := streamResult.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println()
				break
			}
			return fmt.Errorf("failed to receive stream message: %w", err)
		}
		fmt.Print(message.Content)
	}
	return nil
}
