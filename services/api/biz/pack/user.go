package pack

import (
	"tiktok/gen/user"
	"tiktok/services/api/biz/model/api"
)

func User(data *user.User) *api.User {
	return &api.User{
		ID:              data.Id,
		Name:            data.Name,
		FollowCount:     &data.FollowCount,
		FollowerCount:   &data.FollowerCount,
		IsFollow:        data.IsFollow,
		Avatar:          &data.Avatar,
		BackgroundImage: &data.BackgroundImage,
		Signature:       &data.Signature,
		TotalFavorited:  &data.TotalFavorited,
		WorkCount:       &data.WorkCount,
	}
}
