package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	"tiktok/config"
	"tiktok/gen/chat"
	"tiktok/gen/chat/messageservice"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"tiktok/pkg/middleware"
)

func InitChatRPC() {
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})

	if err != nil {
		panic(err)
	}

	c, err := messageservice.NewClient(
		constants.ChatServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithMuxConnection(constants.MuxConnection),
		client.WithRPCTimeout(constants.RPCTimeout),
		client.WithConnectTimeout(constants.ConnectTimeout),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithResolver(r),
		client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()),
	)

	if err != nil {
		panic(err)
	}

	chatClient = c
}

func MessageAction(ctx context.Context, req *chat.MessagePostRequest) error {
	resp, err := chatClient.MessagePost(ctx, req)

	if err != nil {
		return err
	}

	if resp.Base.Code != errno.SuccessCode {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}

	return nil
}

func MessageList(ctx context.Context, req *chat.MessageListRequest) ([]*chat.Message, int64, error) {
	resp, err := chatClient.MessageList(ctx, req)

	if err != nil {
		return nil, -1, err
	}

	if resp.Base.Code != errno.SuccessCode {
		return nil, -1, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}

	return resp.MessageList, resp.Total, nil
}
