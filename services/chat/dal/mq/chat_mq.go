package mq

import (
	"time"

	"tiktok/services/chat/dal/cache"
)

func convertForMysql(message *cache.Message, tempMessage *MiddleMessage) error {
	message.Id = tempMessage.Id
	message.ToUserId = tempMessage.ToUserId
	message.FromUserId = tempMessage.FromUserId
	message.Content = tempMessage.Content
	createdAt, err := time.ParseInLocation(time.RFC3339, tempMessage.CreatedAt, time.Local)
	if err != nil {
		return err
	}
	message.CreatedAt = createdAt
	return nil
}
