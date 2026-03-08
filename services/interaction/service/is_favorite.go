package service

import (
	"errors"

	"gorm.io/gorm"
	"tiktok/gen/interaction"
	"tiktok/services/interaction/dal/cache"
	"tiktok/services/interaction/dal/db"
)

func (s *InteractionService) IsFavorite(req *interaction.IsFavoriteRequest) (bool, error) {
	// read from redis
	exist, err := cache.IsVideoLikeExist(s.ctx, req.VideoId, req.UserId)
	if err != nil {
		return exist, err
	}
	if exist {
		return exist, nil
	}
	// read from mysql
	err = db.IsFavorited(s.ctx, req.UserId, req.VideoId, 1)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
