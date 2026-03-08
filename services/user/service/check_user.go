package service

import (
	"golang.org/x/crypto/bcrypt"
	"tiktok/gen/user"
	"tiktok/pkg/errno"
	"tiktok/services/user/dal/db"
)

// CheckUser check if user exists and it's password
func (s *UserService) CheckUser(req *user.LoginRequest) (*db.User, error) {
	// 函数名叫Check，却返回User? 应该返回 Bool
	userModel, err := db.GetUserByUsername(s.ctx, req.Username)

	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(req.Password)) != nil {
		return nil, errno.AuthorizationFailedError
	}

	return userModel, nil
}
