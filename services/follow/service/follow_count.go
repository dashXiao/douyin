package service

import (
	"tiktok/gen/follow"
	"tiktok/services/follow/dal/cache"
	"tiktok/services/follow/dal/db"
)

func (s *FollowService) FollowCount(req *follow.FollowCountRequest) (int64, error) {
	// 先进入redis中查询
	followCount, err := cache.FollowCount(s.ctx, req.UserId)
	if err != nil {
		return -1, err
	}

	if followCount == 0 { // redis中没查到,进入db中查
		followCount, err = db.FollowCount(s.ctx, req.UserId)
	}
	if err != nil {
		return -1, err
	}
	return followCount, nil
}
