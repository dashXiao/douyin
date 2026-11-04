package service

import (
	"tiktok/gen/follow"
	"tiktok/gen/interaction"
	"tiktok/gen/user"
	"tiktok/gen/video"
	"tiktok/pkg/utils"
	"tiktok/services/user/dal/db"
	"tiktok/services/user/pack"
	"tiktok/services/user/rpc"
)

// GetUser retrieves account information and aggregates social and video statistics.
func (s *UserService) GetUser(req *user.InfoRequest) (*user.User, error) {
	var userResp *user.User

	// Retrieve the user's account information.
	userModel, err := db.GetUserByID(s.ctx, req.UserId)
	userResp = pack.User(userModel)

	if err != nil {
		return nil, err
	}

	// Number of followed accounts
	// Call path: GetUser -> rpc.GetFollowCount -> followClient.FollowCount -> follow service handler.
	userResp.FollowCount, err = rpc.GetFollowCount(s.ctx, &follow.FollowCountRequest{UserId: userModel.Id, Token: req.Token})

	if err != nil {
		return nil, err
	}

	// Number of followers
	userResp.FollowerCount, err = rpc.GetFollowerCount(s.ctx, &follow.FollowerCountRequest{UserId: userModel.Id, Token: req.Token})

	if err != nil {
		return nil, err
	}

	// Whether the current user follows this account
	claims, err := utils.CheckToken(req.Token) // Extract the current user ID from the token.

	if err != nil {
		return nil, err
	}

	userResp.IsFollow, err = rpc.IsFollow(s.ctx, &follow.IsFollowRequest{UserId: claims.UserId, Token: req.Token, ToUserId: userModel.Id})

	if err != nil {
		return nil, err
	}

	// Number of published videos
	userResp.WorkCount, err = rpc.GetWorkCount(s.ctx, &video.GetWorkCountRequest{UserId: userModel.Id, Token: req.Token})

	if err != nil {
		return nil, err
	}

	// Total likes received
	userResp.TotalFavorited, err = rpc.GetTotalFavorited(s.ctx, &interaction.UserTotalFavoritedRequest{UserId: userModel.Id, Token: req.Token})

	if err != nil {
		return nil, err
	}

	return userResp, nil
}
