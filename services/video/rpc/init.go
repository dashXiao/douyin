package rpc

import (
	"tiktok/gen/interaction/interactionservice"
	"tiktok/gen/user/userservice"
)

var (
	userClient        userservice.Client
	interactionClient interactionservice.Client
)

func Init() {
	InitUserRPC()
	InitInteractionRPC()
}
