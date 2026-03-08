package service

import (
	"sort"

	"tiktok/gen/chat"
	"tiktok/services/chat/dal/db"
	"tiktok/services/chat/dal/mq"
)

// Get Messages history list
func (c *ChatService) GetMessages(req *chat.MessageListRequest, user_id int64) ([]*db.Message, error) {
	mq.Mu.Lock()
	defer mq.Mu.Unlock()
	// MySQL search
	messages, err := db.GetMessageList(c.ctx, req.ToUserId, user_id)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return messages, nil
	}

	sort.Sort(db.MessageArray(messages))
	return messages, nil
}
