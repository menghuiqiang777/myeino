package llmModel

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// Generate 调用模型的 Generate 方法生成回复消息
func Generate(ctx context.Context, llm model.ChatModel, in []*schema.Message, paramsOneOf *schema.ParamsOneOf) (*schema.Message, error) {
	// 这里需要根据实际情况处理 paramsOneOf
	// 假设需要调用 ToOpenAPIV3 方法
	openAPIV3Schema, err := paramsOneOf.ToOpenAPIV3()
	if err != nil {
		return nil, fmt.Errorf("failed to convert params to OpenAPIV3 schema: %w", err)
	}
	// 这里根据实际情况将 openAPIV3Schema 传递给模型
	result, err := llm.Generate(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("llm generate failed: %w", err)
	}
	return result, nil
}

// Stream 调用模型的 Stream 方法获取流式回复
func Stream(ctx context.Context, llm model.ChatModel, in []*schema.Message, paramsOneOf *schema.ParamsOneOf) (*schema.StreamReader[*schema.Message], error) {
	// 这里需要根据实际情况处理 paramsOneOf
	// 假设需要调用 ToOpenAPIV3 方法
	openAPIV3Schema, err := paramsOneOf.ToOpenAPIV3()
	if err != nil {
		return nil, fmt.Errorf("failed to convert params to OpenAPIV3 schema: %w", err)
	}
	// 这里根据实际情况将 openAPIV3Schema 传递给模型
	result, err := llm.Stream(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("llm stream failed: %w", err)
	}
	return result, nil
}

// GetChatModel 从原 provider 包移到这里
func GetChatModel(vendor, modelName string) (model.ChatModel, error) {
	// 这里应该实现获取 ChatModel 的逻辑
	// 假设原 provider 包中有相应的逻辑
	// 以下是示例，你需要根据实际情况替换
	if vendor != "qwen" || modelName != "qwen-plus" {
		return nil, fmt.Errorf("unsupported vendor or model name")
	}
	// 这里应该返回实际的 ChatModel 实例
	// 示例返回 nil，你需要根据实际情况修改
	var chatModel model.ChatModel = nil
	return chatModel, nil
}
