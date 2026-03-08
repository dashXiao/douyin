package dal

import (
	"tiktok/services/video/dal/cache"
	"tiktok/services/video/dal/db"
)

func Init() {
	db.Init()
	cache.Init()
}
