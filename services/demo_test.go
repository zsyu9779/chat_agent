package services

import (
	"context"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"testing"
	"time"
)

func TestDoubao(t *testing.T) {
	apiKey := "662bf06a-e486-45fa-a458-96d49305137e"
	client := arkruntime.NewClientWithApiKey(
		apiKey,
		arkruntime.WithTimeout(2*time.Minute),
		arkruntime.WithRetryTimes(2),
	)
	ctx := context.Background()

	fmt.Println("----- standard request -----")
	req := model.ChatCompletionRequest{
		Model: "ep-20240901192351-wv5cv",
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("你是豆包，是由字节跳动开发的 AI 人工智能助手"),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("常见的十字花科植物有哪些？"),
				},
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("standard chat error: %v\n", err)
		return
	}
	fmt.Println(*resp.Choices[0].Message.Content.StringValue)
}
