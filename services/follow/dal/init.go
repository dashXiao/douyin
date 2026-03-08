package dal

import (
	"tiktok/services/follow/dal/cache"
	"tiktok/services/follow/dal/db"
)

func Init() {
	db.Init()
	cache.Init()
}
