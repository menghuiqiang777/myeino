package main

import (
	"fmt"
	"io"

	"github.com/cloudwego/eino/schema"
)

// reportStream 函数用于处理从 StreamReader 接收的消息
// 它会持续接收消息，直接输出，直到遇到 EOF 或发生错误
// 如果发生错误，会返回错误信息给调用者
func reportStream(sr *schema.StreamReader[*schema.Message]) error {
	defer sr.Close()

	for {
		message, err := sr.Recv()
		if err == io.EOF {
			fmt.Println() // 最后添加一个回车
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Print(message.Content)
	}
}
