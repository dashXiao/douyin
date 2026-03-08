package service

import (
	"bytes"

	"tiktok/config"
	"tiktok/gen/video"
)

func (s *VideoService) UploadVideo(req *video.UploadVideoRequest, videoName string) (err error) {
	fileReader := bytes.NewReader(req.VideoFile)
	err = s.bucket.PutObject(config.OSS.MainDirectory+"/"+videoName, fileReader)
	return
}
