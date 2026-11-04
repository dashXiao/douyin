package service

import (
	"golang.org/x/sync/errgroup"
	"tiktok/gen/interaction"
	"tiktok/gen/user"
	"tiktok/gen/video"
	"tiktok/services/video/dal/db"
	"tiktok/services/video/rpc"
)

func (s *VideoService) GetFavoriteVideoInfo(req *video.GetFavoriteVideoInfoRequest) ([]db.Video, []*user.User, []int64, []int64, error) {
	// TODO: share the video-enrichment logic with FeedVideo.
	videoList, err := db.GetVideoInfoByID(s.ctx, req.VideoId)

	// Collect errors from concurrent operations.
	var eg errgroup.Group
	type result struct {
		userInfo      *user.User
		favoriteCount int64
		commentCount  int64
	}
	results := make([]result, len(videoList))
	// Enrich videos concurrently with author details, like counts, and comment counts.
	for i := 0; i < len(videoList); i++ {
		index := i // Capture a separate index for each goroutine.
		eg.Go(func() error {
			// Retrieve the video author.
			userInfo, err := rpc.GetUser(s.ctx, &user.InfoRequest{
				UserId: videoList[index].UserID,
				Token:  req.Token,
			})
			if err != nil {
				return err
			}

			// Retrieve the video's like count.
			favoriteCount, err := rpc.GetVideoFavoriteCount(s.ctx, &interaction.VideoFavoritedCountRequest{
				VideoId: videoList[index].Id,
				Token:   req.Token,
			})
			if err != nil {
				return err
			}

			// Retrieve the video's comment count.
			commentCount, err := rpc.GetCommentCount(s.ctx, &interaction.CommentCountRequest{
				VideoId: videoList[index].Id,
				Token:   &req.Token,
			})
			if err != nil {
				return err
			}
			results[index] = result{
				userInfo:      userInfo,
				favoriteCount: favoriteCount,
				commentCount:  commentCount,
			}
			return nil
		})
	}
	// Wait for all enrichment goroutines to finish.
	if err := eg.Wait(); err != nil {
		return nil, nil, nil, nil, err
	}

	var userList []*user.User
	var favoriteCountList []int64
	var commentCountList []int64

	for _, result := range results {
		userList = append(userList, result.userInfo)
		favoriteCountList = append(favoriteCountList, result.favoriteCount)
		commentCountList = append(commentCountList, result.commentCount)
	}

	return videoList, userList, favoriteCountList, commentCountList, err
}
