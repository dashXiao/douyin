package service

import (
	"time"

	"tiktok/gen/follow"
	"tiktok/gen/user"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
	"tiktok/services/follow/dal/cache"
	"tiktok/services/follow/dal/db"
	"tiktok/services/follow/rpc"
)

// Action follows or unfollows the target user.
func (s *FollowService) Action(req *follow.ActionRequest) error {
	// Apply the request rate limit.
	if err := cache.Limit(s.ctx, constants.ActionRate, constants.Interval); err != nil {
		return err
	}

	claim, err := utils.CheckToken(req.Token)

	if err != nil {
		return errno.AuthorizationFailedError
	}

	// Reject attempts to follow the current user's own account.
	if claim.UserId == req.ToUserId {
		return errno.FollowYourselfError
	}

	// Verify that the target user exists.
	_, err = rpc.GetUser(s.ctx, &user.InfoRequest{
		UserId: req.ToUserId,
		Token:  req.Token,
	})

	if err != nil {
		return errno.UserNotFoundError
	}

	action := &db.Follow{
		UserID:   claim.UserId,
		ToUserID: req.ToUserId,
	}

	switch req.ActionType {
	case constants.FollowAction:
		// Update the follow relationship in Redis.
		if err := cache.FollowAction(s.ctx, action.UserID, action.ToUserID); err != nil {
			return err
		}
		// Create or reactivate the follow record in MySQL.
		if err = db.FollowAction(s.ctx, action); err != nil {
			return err
		}
	case constants.UnFollowAction:
		// Mark the follow record as inactive in MySQL.
		if err = db.UnFollowAction(s.ctx, action); err != nil {
			return err
		}
		time.Sleep(10 * time.Millisecond) // Delay the cache update after the database write.
		// Remove the follow relationship from Redis.
		if err = cache.UnFollowAction(s.ctx, action.UserID, action.ToUserID); err != nil {
			return err
		}
	default:
		return errno.UnexpectedTypeError
	}

	if err != nil {
		return err
	}

	return nil
}
