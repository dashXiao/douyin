package pack

import (
	"tiktok/gen/chat"
	"tiktok/services/api/biz/model/api"
)

func MessageList(list []*chat.Message) []*api.Message {
	resp := make([]*api.Message, 0)

	for _, data := range list {
		resp = append(resp, &api.Message{
			ID:         data.Id,
			ToUserID:   data.ToUserId,
			FromUserID: data.FromUserId,
			Content:    data.Content,
			CreateTime: *data.CreateTime,
		})
	}

	return resp
}
