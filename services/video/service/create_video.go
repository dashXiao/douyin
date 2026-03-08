package service

import (
	"tiktok/gen/video"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
	"tiktok/services/video/dal/db"
)

func (s *VideoService) CreateVideo(req *video.UploadVideoRequest, playURL string, coverURL string) (*db.Video, error) {
	claim, err := utils.CheckToken(req.Token)
	if err != nil {
		return nil, errno.AuthorizationFailedError
	}
	videoModel := &db.Video{
		UserID:   claim.UserId,
		PlayUrl:  playURL,
		CoverUrl: coverURL,
		Title:    req.Title,
	}
	return db.CreateVideo(s.ctx, videoModel)
}
