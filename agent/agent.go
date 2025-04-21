package agent

import "github.com/cloudwego/eino/components/model"

// Agent 结构体定义了一个智能代理，包含名称、指令和使用的模型
type Agent struct {
	Name         string
	Instructions string
	Model        model.ChatModel
}
