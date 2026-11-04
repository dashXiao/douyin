package pack

import (
	"tiktok/gen/user"

	"tiktok/services/user/dal/db"
)

// User copies account fields; the user service adds social and video statistics.
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
}
