package service

import (
	"tiktok/gen/interaction"
	"tiktok/services/interaction/dal/cache"
	"tiktok/services/interaction/dal/db"
)

func (s *InteractionService) GetVideoFavoritedCount(req *interaction.VideoFavoritedCountRequest) (int64, error) {
	// read from redis
	_, likeCount, err := cache.GetVideoLikeCount(s.ctx, req.VideoId)
	if err != nil {
		return 0, err
	}
	if likeCount == 0 {
		// read from mysql
		likeCount, err = db.GetVideoLikeCount(s.ctx, req.VideoId)
		if err != nil {
			return 0, err
		}
		// update redis data
		err = cache.SetVideoLikeCount(s.ctx, req.VideoId, likeCount)
	}
	if likeCount < 0 {
		likeCount = 0
	}
	return likeCount, err
}
