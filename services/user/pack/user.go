package pack

import (
	"tiktok/gen/user"

	"tiktok/services/user/dal/db"
)

func User(data *db.User) *user.User {
	if data == nil {
		return nil
	}

	return &user.User{
		Id:              data.Id,
		Name:            data.Username,
		Avatar:          data.Avatar,
		BackgroundImage: data.BackgroundImage,
		Signature:       data.Signature,
	}

	// TODO: follow_count, follower_count, is_follow, work_count, favorite_count, total_favorited(string)
}
