package dal

import (
	"tiktok/services/user/dal/cache"
	"tiktok/services/user/dal/db"
)

func Init() {
	db.Init()
	cache.Init()
}
