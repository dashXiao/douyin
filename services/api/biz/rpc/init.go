package rpc

import (
	"tiktok/gen/chat/messageservice"
	"tiktok/gen/follow/followservice"
	"tiktok/gen/interaction/interactionservice"
	"tiktok/gen/user/userservice"
	"tiktok/gen/video/videoservice"
)

var (
	userClient        userservice.Client
	followClient      followservice.Client
	interactionClient interactionservice.Client
	chatClient        messageservice.Client
	videoClient       videoservice.Client
)

func Init() {
	InitUserRPC()
	InitFollowRPC()
	InitInteractionRPC()
	InitChatRPC()
	InitVideoRPC()
}
