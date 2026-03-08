package service

import (
	"golang.org/x/crypto/bcrypt"
	"tiktok/gen/user"
	"tiktok/services/user/dal/db"
)

// CreateUser create user info
func (s *UserService) CreateUser(req *user.RegisterRequest) (*db.User, error) {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	userModel := &db.User{
		Username: req.Username,
		Password: string(hashBytes),
	}

	return db.CreateUser(s.ctx, userModel)
}
