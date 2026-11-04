package service

import (
	"bytes"
	"errors"
	"os/exec"

	"tiktok/config"
	"tiktok/gen/video"
)

func (s *VideoService) UploadCover(req *video.UploadVideoRequest, coverName string) (err error) {
	// Store a user-supplied cover without generating a thumbnail.
	if len(req.CoverFile) > 0 {
		return s.bucket.PutObject(config.OSS.MainDirectory+"/"+coverName, bytes.NewReader(req.CoverFile))
	}

	var imageBuffer bytes.Buffer // Capture the generated image.

	// Extract the first frame through stdin/stdout to avoid blocking on named pipes.
	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-vframes", "1", "-f", "image2pipe", "-vcodec", "png", "pipe:1")
	cmd.Stdin = bytes.NewReader(req.VideoFile)
	cmd.Stdout = &imageBuffer

	if err = cmd.Run(); err != nil {
		return err
	}

	if imageBuffer.Len() == 0 {
		return errors.New("cover generation failed: empty output")
	}

	return s.bucket.PutObject(config.OSS.MainDirectory+"/"+coverName, bytes.NewReader(imageBuffer.Bytes()))
}
