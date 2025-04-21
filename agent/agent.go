package main

import "github.com/cloudwego/eino/components/model"

// Agent 结构体定义
type Agent struct {
	name         string
	instructions string
	model        model.ChatModel
}
