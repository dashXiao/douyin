package pack

import (
	"fmt"

	"tiktok/gen/chat"
	"tiktok/services/chat/dal/db"
)

type MessageBuildArray []*chat.Message

func BuildMessage(data []*db.Message) []*chat.Message {
	if data == nil {
		return make([]*chat.Message, 0)
	}
	res := make(MessageBuildArray, 0)
	for _, val := range data {
		create_at := fmt.Sprintf("%v", val.CreatedAt.UnixMilli())
		msg := &chat.Message{
			Id:         val.Id,
			ToUserId:   val.ToUserId,
			FromUserId: val.FromUserId,
			Content:    val.Content,
			CreateTime: &create_at,
		}
		res = append(res, msg)
	}
	return res
}
