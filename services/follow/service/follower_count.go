package service

import (
	"tiktok/gen/follow"
	"tiktok/services/follow/dal/cache"
	"tiktok/services/follow/dal/db"
)

func (s *FollowService) FollowerCount(req *follow.FollowerCountRequest) (int64, error) {
	// 先进入redis中查询
	followerCount, err := cache.FollowerCount(s.ctx, req.UserId)
	if err != nil {
		return -1, err
	}

	if followerCount == 0 { // redis中没查到,进入db中查
		followerCount, err = db.FollowerCount(s.ctx, req.UserId)
	}
	if err != nil {
		return -1, err
	}
	return followerCount, nil
}
