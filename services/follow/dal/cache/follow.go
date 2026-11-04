package cache

import (
	"context"
	"errors"
	"strconv"

	"github.com/redis/go-redis/v9"
	"tiktok/pkg/errno"
)

// IsFollow checks whether the follow relationship exists in Redis.
func IsFollow(ctx context.Context, uid, tid int64) (bool, error) {
	exist, err := RedisClient.SIsMember(ctx, FollowListKey(uid), strconv.FormatInt(tid, 10)).Result()

	if err != nil {
		return exist, err
	}

	return exist, nil
}

func FollowAction(ctx context.Context, uid, tid int64) error {
	exist, err := IsFollow(ctx, uid, tid) // Check the cached follow relationship.
	if err != nil {
		return err
	} else if exist { // The user already follows this account.
		return errno.AlreadyFollowError
	}

	// Add the follow relationship to both Redis sets.
	err = RedisClient.SAdd(ctx, FollowListKey(uid), strconv.FormatInt(tid, 10)).Err() // The current user's following list
	if err != nil {
		return err
	}

	err = RedisClient.SAdd(ctx, FollowerListKey(tid), strconv.FormatInt(uid, 10)).Err() // The target user's follower list
	if err != nil {
		return err
	}
	return nil
}

func UnFollowAction(ctx context.Context, uid, tid int64) error {
	exist, err := IsFollow(ctx, uid, tid) // Check the cached follow relationship.
	if err != nil {
		return err
	} else if !exist { // The user does not follow this account.
		return errno.NotFollowError
	}

	// Remove the current user from the target user's follower list.
	err = RedisClient.SRem(ctx, FollowerListKey(tid), strconv.FormatInt(uid, 10)).Err()
	if err != nil {
		return err
	}

	// Remove the target user from the current user's following list.
	err = RedisClient.SRem(ctx, FollowListKey(uid), strconv.FormatInt(tid, 10)).Err()
	if err != nil {
		return err
	}
	return nil
}

func FollowListAction(ctx context.Context, uid int64) (*[]int64, error) {
	followList := make([]int64, 0, 10)

	key := FollowListKey(uid)

	// Read the cached user IDs from Redis.
	idList, err := RedisClient.SMembers(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
	} else if err != nil {
		return nil, err
	}

	for _, id := range idList {
		followId, _ := strconv.ParseInt(id, 10, 64)
		followList = append(followList, followId)
	}
	return &followList, nil
}

func UpdateFollowList(ctx context.Context, uid int64, followList *[]int64) error {
	key := FollowListKey(uid)
	idList := make([]interface{}, len(*followList))
	for i, v := range *followList {
		idList[i] = v
	}
	if err := RedisClient.SAdd(ctx, key, idList).Err(); err != nil {
		return err
	}
	return nil
}

func FollowerListAction(ctx context.Context, uid int64) (*[]int64, error) {
	followerList := make([]int64, 0, 10)

	key := FollowerListKey(uid)

	// Read the cached user IDs from Redis.
	idList, err := RedisClient.SMembers(ctx, key).Result()
	if err == redis.Nil {
	} else if err != nil {
		return nil, err
	}

	for _, id := range idList {
		followerId, _ := strconv.ParseInt(id, 10, 64)
		followerList = append(followerList, followerId)
	}
	return &followerList, nil
}

func UpdateFollowerList(ctx context.Context, uid int64, followerList *[]int64) error {
	key := FollowerListKey(uid)
	idList := make([]interface{}, len(*followerList))
	for i, v := range *followerList {
		idList[i] = v
	}
	if err := RedisClient.SAdd(ctx, key, idList).Err(); err != nil {
		return err
	}
	return nil
}

func FriendListAction(ctx context.Context, uid int64) (*[]int64, error) {
	friendList := make([]int64, 0, 10)
	// Retrieve the accounts followed by this user.
	tempList, err := FollowListAction(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Check whether each followed account also follows this user.
	for _, followID := range *tempList {
		b, err := RedisClient.SIsMember(ctx, FollowerListKey(uid), followID).Result()
		if err != nil {
			return nil, err
		} else if !b { // Skip accounts that do not follow this user back.
			continue
		}
		friendList = append(friendList, followID)
	}

	return &friendList, nil
}

func UpdateFriendList(ctx context.Context, uid int64, followList, followerList *[]int64) error {
	userList, _ := FollowListAction(ctx, uid)
	if len(*userList) == 0 {
		err := UpdateFollowList(ctx, uid, followList)
		if err != nil {
			return err
		}
	}

	err := UpdateFollowerList(ctx, uid, followerList)
	if err != nil {
		return err
	}

	return nil
}

func FollowCount(ctx context.Context, uid int64) (int64, error) {
	key := FollowListKey(uid)
	size, err := RedisClient.SCard(ctx, key).Result()
	if err != nil {
		return -1, err
	}
	return size, nil
}

func FollowerCount(ctx context.Context, uid int64) (int64, error) {
	key := FollowerListKey(uid)
	size, err := RedisClient.SCard(ctx, key).Result()
	if err != nil {
		return -1, err
	}
	return size, nil
}
