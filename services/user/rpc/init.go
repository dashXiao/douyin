package rpc

import (
	"tiktok/gen/follow/followservice"
	"tiktok/gen/interaction/interactionservice"
	"tiktok/gen/video/videoservice"
)

var (
	followClient      followservice.Client
	interactionClient interactionservice.Client
	videoClient       videoservice.Client
)

func Init() {
	InitFollowRPC()
	InitInteractionRPC()
	InitVideoRPC()
}
