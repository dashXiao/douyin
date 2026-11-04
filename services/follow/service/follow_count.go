package service

import (
	"tiktok/gen/follow"
	"tiktok/services/follow/dal/cache"
	"tiktok/services/follow/dal/db"
)

func (s *FollowService) FollowCount(req *follow.FollowCountRequest) (int64, error) {
	// Read the count from Redis first.
	followCount, err := cache.FollowCount(s.ctx, req.UserId)
	if err != nil {
		return -1, err
	}

	if followCount == 0 { // Fall back to MySQL when the cached count is zero.
		followCount, err = db.FollowCount(s.ctx, req.UserId)
	}
	if err != nil {
		return -1, err
	}
	return followCount, nil
}
