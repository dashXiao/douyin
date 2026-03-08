package dal

import (
	"tiktok/services/chat/dal/cache"
	"tiktok/services/chat/dal/db"
	"tiktok/services/chat/dal/mq"
)

func Init() {
	db.Init()
	cache.Init()
	mq.Init()
}
