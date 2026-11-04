package service

import (
	"tiktok/gen/follow"
	"tiktok/services/follow/dal/cache"
	"tiktok/services/follow/dal/db"
)

func (s *FollowService) FollowerCount(req *follow.FollowerCountRequest) (int64, error) {
	// Read the count from Redis first.
	followerCount, err := cache.FollowerCount(s.ctx, req.UserId)
	if err != nil {
		return -1, err
	}

	if followerCount == 0 { // Fall back to MySQL when the cached count is zero.
		followerCount, err = db.FollowerCount(s.ctx, req.UserId)
	}
	if err != nil {
		return -1, err
	}
	return followerCount, nil
}
