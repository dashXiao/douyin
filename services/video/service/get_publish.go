package service

import (
	"golang.org/x/sync/errgroup"
	"tiktok/gen/interaction"
	"tiktok/gen/user"
	"tiktok/gen/video"
	"tiktok/services/video/dal/db"
	"tiktok/services/video/rpc"
)

func (s *VideoService) GetPublishVideoInfo(req *video.GetPublishListRequest) ([]db.Video, []*user.User, []int64, []int64, []bool, error) {
	// TODO: share video-enrichment logic and replace parallel return slices with a result type.
	videoList, err := db.GetVideoInfoByUid(s.ctx, req.UserId)

	// Collect errors from concurrent operations.
	var eg errgroup.Group
	type result struct {
		userInfo      *user.User
		favoriteCount int64
		commentCount  int64
		isFavorite    bool
	}

	results := make([]result, len(videoList))

	// Enrich videos concurrently with author details, like counts, comment counts, and like status.
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
			// Retrieve the like status for the requested user.
			isFavorite, err := rpc.GetVideoIsFavorite(s.ctx, &interaction.InteractionServiceIsFavoriteArgs{Req: &interaction.IsFavoriteRequest{
				UserId:  videoList[index].UserID,
				VideoId: videoList[index].Id,
				Token:   req.Token,
			}})
			if err != nil {
				return err
			}
			results[index] = result{
				userInfo:      userInfo,
				favoriteCount: favoriteCount,
				commentCount:  commentCount,
				isFavorite:    isFavorite,
			}
			return nil
		})
	}

	// Wait for all enrichment goroutines to finish.
	if err := eg.Wait(); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	var userList []*user.User
	var favoriteCountList []int64
	var commentCountList []int64
	var isFavoriteList []bool

	for _, result := range results {
		userList = append(userList, result.userInfo)
		favoriteCountList = append(favoriteCountList, result.favoriteCount)
		commentCountList = append(commentCountList, result.commentCount)
		isFavoriteList = append(isFavoriteList, result.isFavorite)
	}

	return videoList, userList, favoriteCountList, commentCountList, isFavoriteList, err
}
