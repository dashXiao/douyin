package service

import (
	"errors"
	"sync"

	"tiktok/gen/chat"
	"tiktok/gen/follow"
	"tiktok/gen/user"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"tiktok/services/follow/dal/cache"
	"tiktok/services/follow/dal/db"
	"tiktok/services/follow/pack"
	"tiktok/services/follow/rpc"
)

// FriendList Viewing friends list
func (s *FollowService) FriendList(req *follow.FriendListRequest) (*[]*follow.FriendUser, error) {
	// Apply the request rate limit.
	if err := cache.Limit(s.ctx, constants.FriendListRate, constants.Interval); err != nil {
		return nil, err
	}

	friendList := make([]*follow.FriendUser, 0, 10)
	var wg sync.WaitGroup
	var mu sync.Mutex
	isErr := false

	// Read from Redis first.
	userList, err := cache.FriendListAction(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	} else if len(*userList) == 0 { // Fall back to MySQL when the cached list is empty.
		userList, err = db.FriendListAction(s.ctx, req.UserId)
		if errors.Is(err, db.RecordNotFound) { // No matching records were found in MySQL.
			return &friendList, nil
		} else if err != nil {
			return nil, err
		}
		// Populate Redis with the IDs retrieved from MySQL.
		followList, _ := db.FollowListAction(s.ctx, req.UserId)
		followerList, _ := db.FollowerListAction(s.ctx, req.UserId)
		err := cache.UpdateFriendList(s.ctx, req.UserId, followList, followerList)
		if err != nil {
			return nil, err
		}
	}

	// Build the user response objects.
	for _, userID := range *userList {
		wg.Add(1)
		go func(id int64, req *follow.FriendListRequest, userList *[]*follow.FriendUser, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer func() {
				if e := recover(); e != nil {
					mu.Lock()
					isErr = true
					mu.Unlock()
				}
				wg.Done()
			}()

			user, err := rpc.GetUser(s.ctx, &user.InfoRequest{
				UserId: id,
				Token:  req.Token,
			})
			if err != nil {
				mu.Lock()
				isErr = true // Record that a worker failed.
				mu.Unlock()
				return
			}
			friend := pack.User(user) // Convert the RPC user to the response type.

			message, msgType, err := rpc.GetMessage(s.ctx, &chat.MessageListRequest{
				Token:    req.Token,
				ToUserId: req.UserId,
			}, req.UserId, id)
			if err != nil {
				mu.Lock()
				isErr = true // Record that a worker failed.
				mu.Unlock()
				return
			}

			friendUser := &follow.FriendUser{
				User:    friend,
				Message: &message,
				MsgType: msgType,
			}

			mu.Lock()
			*userList = append(*userList, friendUser)
			mu.Unlock()
		}(userID, req, &friendList, &wg, &mu)
	}

	wg.Wait()

	if isErr {
		return nil, errno.ServiceError
	}

	return &friendList, nil
}
