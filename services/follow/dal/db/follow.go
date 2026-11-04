package db

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
)

type Follow struct {
	Id        int64
	UserID    int64
	ToUserID  int64
	Status    int64 // 1: follow; 2: unfollow
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

var RecordNotFound = gorm.ErrRecordNotFound

// FollowAction creates or reactivates a follow relationship.
func FollowAction(ctx context.Context, follow *Follow) error {
	followResp := new(Follow)

	// Check whether this pair of users has an existing follow record.
	err := DB.WithContext(ctx).Model(&Follow{}).
		Where("user_id= ? AND to_user_id = ?", follow.UserID, follow.ToUserID).
		First(&followResp).Error

	// Create a follow record if none exists.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		follow.Id = SF.NextVal()
		follow.Status = 1
		return DB.WithContext(ctx).Create(follow).Error
	} else if err != nil {
		return err
	}

	// Reactivate the existing follow record.
	err = DB.WithContext(ctx).Model(&Follow{}).
		Where("user_id= ? AND to_user_id = ?", follow.UserID, follow.ToUserID).
		Update("status", constants.FollowAction).Error
	if err != nil {
		return err
	}

	return nil
}

// UnFollowAction marks the follow relationship as inactive.
func UnFollowAction(ctx context.Context, follow *Follow) error {
	_, err := IsFollow(ctx, follow.UserID, follow.ToUserID)
	if errors.Is(err, RecordNotFound) {
		return errno.NotFollowError
	} else if err != nil {
		return err
	}
	// Update the follow status in MySQL.
	err = DB.WithContext(ctx).Model(&Follow{}).
		Where("user_id= ? AND to_user_id = ?", follow.UserID, follow.ToUserID).
		Update("status", constants.UnFollowAction).Error
	if err != nil {
		return err
	}

	return nil
}

// FollowListAction returns the IDs of accounts followed by the user.
func FollowListAction(ctx context.Context, uid int64) (*[]int64, error) {
	var followList []int64

	// Retrieve the user IDs from MySQL after a cache miss.
	err := DB.WithContext(ctx).Table(constants.FollowTableName).Select("to_user_id").
		Where("user_id = ? AND status = ?", uid, constants.FollowAction).
		Find(&followList).Error
	if err != nil {
		return nil, err
	}

	if len(followList) == 0 { // No matching records were found in MySQL.
		return nil, RecordNotFound
	}

	return &followList, nil
}

// FollowerListAction returns the IDs of the user's followers.
func FollowerListAction(ctx context.Context, uid int64) (*[]int64, error) {
	var followerList []int64

	// Retrieve the user IDs from MySQL after a cache miss.
	err := DB.WithContext(ctx).Table(constants.FollowTableName).Select("user_id").
		Where("to_user_id = ? AND status = ?", uid, constants.FollowAction).
		Find(&followerList).Error
	if err != nil {
		return nil, err
	}

	if len(followerList) == 0 { // No matching records were found in MySQL.
		return nil, RecordNotFound
	}

	return &followerList, nil
}

// FriendListAction returns the IDs of accounts that follow each other.
func FriendListAction(ctx context.Context, uid int64) (*[]int64, error) {
	friendList := make([]int64, 0, 10)
	// Retrieve the accounts followed by this user.
	tempList, err := FollowListAction(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Check in MySQL whether each followed account also follows this user.
	for _, followID := range *tempList {
		var count int64
		err = DB.WithContext(ctx).Model(&Follow{}).
			Where("user_id = ? AND to_user_id = ? AND status = ?", followID, uid, constants.FollowAction).
			Count(&count).Error
		if err != nil {
			return nil, err
		} else if count == 0 { // Skip accounts that do not follow this user back.
			continue
		}

		// Include the account when the follow relationship is mutual.
		friendList = append(friendList, followID)
	}

	if len(friendList) == 0 { // No matching records were found in MySQL.
		return nil, RecordNotFound
	}

	return &friendList, nil
}

func FollowCount(ctx context.Context, uid int64) (int64, error) {
	var count int64
	err := DB.WithContext(ctx).Model(&Follow{}).
		Where("user_id = ? AND status = ?", uid, constants.FollowAction).
		Count(&count).Error
	if err != nil {
		return -1, err
	}
	return count, nil
}

func FollowerCount(ctx context.Context, uid int64) (int64, error) {
	var count int64
	err := DB.WithContext(ctx).Model(&Follow{}).
		Where("to_user_id = ? AND status = ?", uid, constants.FollowAction).
		Count(&count).Error
	if err != nil {
		return -1, err
	}
	return count, nil
}

func IsFollow(ctx context.Context, uid, tid int64) (bool, error) {
	var followModel *Follow

	err := DB.WithContext(ctx).
		Model(&Follow{}).
		Select("status").
		Where("user_id = ? AND to_user_id = ?", uid, tid).
		First(&followModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // No matching records were found in MySQL.
			return false, RecordNotFound
		}
		return false, err
	}

	return followModel.Status == int64(constants.FollowAction), nil
}
