package service

import (
	"golang.org/x/crypto/bcrypt"
	"tiktok/gen/user"
	"tiktok/pkg/errno"
	"tiktok/services/user/dal/db"
)

// CheckUser validates the login credentials and returns the user record.
func (s *UserService) CheckUser(req *user.LoginRequest) (*db.User, error) {
	userModel, err := db.GetUserByUsername(s.ctx, req.Username)

	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(req.Password)) != nil {
		return nil, errno.AuthorizationFailedError
	}

	return userModel, nil
}
