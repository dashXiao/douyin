package service

import (
	"errors"

	"tiktok/gen/follow"
	"tiktok/services/follow/dal/cache"
	"tiktok/services/follow/dal/db"
)

func (s *FollowService) IsFollow(req *follow.IsFollowRequest) (bool, error) {
	// 先进入redis中判断是否有关注
	ex1, err := cache.IsFollow(s.ctx, req.UserId, req.ToUserId)

	if err != nil {
		return false, err
	}

	if ex1 {
		return true, nil
	}

	ex2, err := db.IsFollow(s.ctx, req.UserId, req.ToUserId)

	if err != nil {
		if errors.Is(err, db.RecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return ex2, nil
}
