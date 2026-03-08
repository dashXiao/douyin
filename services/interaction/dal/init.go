package dal

import (
	"tiktok/services/interaction/dal/cache"
	"tiktok/services/interaction/dal/db"
)

func Init() {
	db.Init()
	cache.Init()
}
