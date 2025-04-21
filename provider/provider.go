package provider

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
)

// Provider 是一个模型提供者的接口
type Provider interface {
	InitModel(modelName string) (interface{}, error)
}

// baseProvider 包含了公共的初始化逻辑
type baseProvider struct {
	apiKeyEnv    string
	baseURLEnv   string
	defaultModel string
	initFunc     func(ctx context.Context, config interface{}) (interface{}, error)
	configFunc   func(apiKey, baseURL, modelName string) interface{}
}

func (bp *baseProvider) InitModel(modelName string) (interface{}, error) {
	var timeout = 30 * time.Second
	// 处理默认值
	if modelName == "" {
		modelName = bp.defaultModel
	}

	// 创建带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 获取环境变量
	apiKey, ok := os.LookupEnv(bp.apiKeyEnv)
	if !ok {
		err := fmt.Errorf("%s environment variable is not set.", bp.apiKeyEnv)
		log.Printf("Error: %s\n", err)
		return nil, err
	}
	var baseURL string
	if bp.baseURLEnv != "" {
		baseURL, ok = os.LookupEnv(bp.baseURLEnv)
		if !ok {
			err := fmt.Errorf("%s environment variable is not set.", bp.baseURLEnv)
			log.Printf("Error: %s\n", err)
			return nil, err
		}
	}

	// 初始化模型
	config := bp.configFunc(apiKey, baseURL, modelName)
	model, err := bp.initFunc(ctx, config)
	if err != nil {
		log.Printf("Error initializing model: %s\n", err)
		return nil, err
	}

	return model, nil
}

// QwenProvider 实现了 Provider 接口，用于调用通义千问模型
type QwenProvider struct {
	baseProvider
}

func NewQwenProvider() *QwenProvider {
	return &QwenProvider{
		baseProvider: baseProvider{
			apiKeyEnv:    "TONGYI_API_KEY",
			baseURLEnv:   "TONGYI_BASE_URL",
			defaultModel: "qwen-plus",
			initFunc: func(ctx context.Context, config interface{}) (interface{}, error) {
				return qwen.NewChatModel(ctx, config.(*qwen.ChatModelConfig))
			},
			configFunc: func(apiKey, baseURL, modelName string) interface{} {
				return &qwen.ChatModelConfig{
					BaseURL: baseURL,
					APIKey:  apiKey,
					Model:   modelName,
				}
			},
		},
	}
}

// ArkProvider 实现了 Provider 接口，用于调用豆包模型
type ArkProvider struct {
	baseProvider
}

func NewArkProvider() *ArkProvider {
	return &ArkProvider{
		baseProvider: baseProvider{
			apiKeyEnv:    "ARK_API_KEY",
			baseURLEnv:   "",
			defaultModel: "doubao-1-5-pro-32k-250115",
			initFunc: func(ctx context.Context, config interface{}) (interface{}, error) {
				return ark.NewChatModel(ctx, config.(*ark.ChatModelConfig))
			},
			configFunc: func(apiKey, _, modelName string) interface{} {
				return &ark.ChatModelConfig{
					APIKey: apiKey,
					Model:  modelName,
				}
			},
		},
	}
}

// GetProvider 根据厂商名称返回对应的模型提供者
func GetProvider(vendor string) (Provider, error) {
	var provider Provider
	var err error
	switch vendor {
	case "qwen":
		provider = NewQwenProvider()
	case "ark":
		provider = NewArkProvider()
	default:
		err = fmt.Errorf("unsupported vendor: %s", vendor)
	}
	if err != nil {
		log.Printf("Error getting provider: %s\n", err)
		return nil, err
	}
	return provider, nil
}

// GetChatModel 根据厂商和模型名称返回 ChatModel 类型的模型
func GetChatModel(vendor, modelName string) (model.ChatModel, error) {
	provider, err := GetProvider(vendor)
	if err != nil {
		log.Printf("Error getting provider: %s\n", err)
		return nil, err
	}
	modelObj, err := provider.InitModel(modelName)
	if err != nil {
		log.Printf("Error initializing model: %s\n", err)
		return nil, err
	}
	chatModel, ok := modelObj.(model.ChatModel)
	if !ok {
		err = fmt.Errorf("the initialized model does not implement the ChatModel interface")
		log.Printf("Error casting model: %s\n", err)
		return nil, err
	}
	return chatModel, nil
}
