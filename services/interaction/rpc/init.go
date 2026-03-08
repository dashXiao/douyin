package rpc

import (
	"tiktok/gen/user/userservice"
	"tiktok/gen/video/videoservice"
)

var (
	userClient  userservice.Client
	videoClient videoservice.Client
)

func Init() {
	InitUserRPC()
	InitVideoRPC()
}
