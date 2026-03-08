package rpc

import (
	"tiktok/gen/chat/messageservice"
	"tiktok/gen/user/userservice"
)

var (
	userClient userservice.Client
	chatClient messageservice.Client
)

func Init() {
	InitUserRPC()
	InitChatRPC()
}
