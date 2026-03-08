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

// GetUser check token and get user's info
func (s *UserService) GetUser(req *user.InfoRequest) (*user.User, error) {
	var userResp *user.User

	// 获取用户基本信息
	userModel, err := db.GetUserByID(s.ctx, req.UserId)
	userResp = pack.User(userModel)

	if err != nil {
		return nil, err
	}

	// 关注数量
	// 调用链：here -> services/user/rpc/follow.go  followClient.FollowCount -> services/follow/handler.go FollowCount -> ...
	userResp.FollowCount, err = rpc.GetFollowCount(s.ctx, &follow.FollowCountRequest{UserId: userModel.Id, Token: req.Token})

	if err != nil {
		return nil, err
	}

	// 粉丝数量
	userResp.FollowerCount, err = rpc.GetFollowerCount(s.ctx, &follow.FollowerCountRequest{UserId: userModel.Id, Token: req.Token})

	if err != nil {
		return nil, err
	}

	// 是否关注
	claims, err := utils.CheckToken(req.Token) // 通过解析Token获取当前登录的UserId

	if err != nil {
		return nil, err
	}

	userResp.IsFollow, err = rpc.IsFollow(s.ctx, &follow.IsFollowRequest{UserId: claims.UserId, Token: req.Token, ToUserId: userModel.Id})

	if err != nil {
		return nil, err
	}

	// 作品数量
	userResp.WorkCount, err = rpc.GetWorkCount(s.ctx, &video.GetWorkCountRequest{UserId: userModel.Id, Token: req.Token})

	if err != nil {
		return nil, err
	}

	// 获赞数量
	userResp.TotalFavorited, err = rpc.GetTotalFavorited(s.ctx, &interaction.UserTotalFavoritedRequest{UserId: userModel.Id, Token: req.Token})

	if err != nil {
		return nil, err
	}

	return userResp, nil
}
