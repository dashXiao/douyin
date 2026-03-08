package main

import (
	"context"
	"time"

	chat "tiktok/gen/chat"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
	pack "tiktok/services/chat/pack"
	service "tiktok/services/chat/service"
)

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct{}

// MessagePost implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessagePost(ctx context.Context, req *chat.MessagePostRequest) (resp *chat.MessagePostResponse, err error) {
	resp = new(chat.MessagePostResponse)
	claim, err := utils.CheckToken(req.Token)
	if err != nil {
		resp.Base = pack.BuildBaseResp(errno.AuthorizationFailedError)
		return resp, err
	}

	err = service.NewChatService(ctx).SendMessage(req, claim.UserId, time.Now().Format(time.RFC3339))
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, err
	}
	resp.Base = pack.BuildBaseResp(nil)
	return
}

// MessageList implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageList(ctx context.Context, req *chat.MessageListRequest) (resp *chat.MessageListResponse, err error) {
	resp = new(chat.MessageListResponse)
	claim, err := utils.CheckToken(req.Token)
	if err != nil || claim == nil {
		resp.Base = pack.BuildBaseResp(errno.AuthorizationFailedError)
		return resp, err
	}
	// 获取消息列表

	// redis->mysql
	// redis中存在则返回，不存在查询mysql,
	messageList, err := service.NewChatService(ctx).GetMessages(req, claim.UserId)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		resp.MessageList = pack.BuildMessage(nil)
		resp.Total = 0
		return resp, err
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.MessageList = pack.BuildMessage(messageList)
	resp.Total = int64(len(messageList))
	return
}
