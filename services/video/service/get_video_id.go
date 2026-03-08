package service

import (
	"tiktok/gen/video"
	"tiktok/services/video/dal/db"
)

func (s *VideoService) GetVideoIDByUid(req *video.GetVideoIDByUidRequset) (videoIDList []int64, err error) {
	return db.GetVideoIDByUid(s.ctx, req.UserId)
}
