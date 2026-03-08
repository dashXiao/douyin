package service

import (
	"tiktok/gen/video"
	"tiktok/services/video/dal/db"
)

func (s *VideoService) GetWorkCount(req *video.GetWorkCountRequest) (workCount int64, err error) {
	return db.GetWorkCountByUid(s.ctx, req.UserId)
}
