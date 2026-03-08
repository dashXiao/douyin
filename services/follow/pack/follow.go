package pack

import (
	"tiktok/gen/follow"
	"tiktok/gen/user"
)

func User(user *user.User) *follow.User {
	return &follow.User{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     &user.FollowCount,
		FollowerCount:   &user.FollowerCount,
		IsFollow:        user.IsFollow,
		Avatar:          &user.Avatar,
		BackgroundImage: &user.BackgroundImage,
		Signature:       &user.Signature,
		TotalFavorited:  &user.TotalFavorited,
		WorkCount:       &user.WorkCount,
	}
}
