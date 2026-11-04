package service

import (
	"errors"
	"sync"

	"tiktok/gen/follow"
	"tiktok/gen/user"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"tiktok/services/follow/dal/cache"
	"tiktok/services/follow/dal/db"
	"tiktok/services/follow/pack"
	"tiktok/services/follow/rpc"
)

// FollowList View the follow list
func (s *FollowService) FollowList(req *follow.FollowListRequest) (*[]*follow.User, error) {
	// Apply the request rate limit.
	if err := cache.Limit(s.ctx, constants.FollowListRate, constants.Interval); err != nil {
		return nil, err
	}

	userList := make([]*follow.User, 0, 10)
	var wg sync.WaitGroup
	var mu sync.Mutex
	isErr := false

	// Read from Redis first.
	followList, err := cache.FollowListAction(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	} else if len(*followList) == 0 { // Fall back to MySQL when the cached list is empty.
		followList, err = db.FollowListAction(s.ctx, req.UserId)
		if errors.Is(err, db.RecordNotFound) { // No matching records were found in MySQL.
			return &userList, nil
		} else if err != nil {
			return nil, err
		}
		// Populate Redis with the IDs retrieved from MySQL.
		err := cache.UpdateFollowList(s.ctx, req.UserId, followList)
		if err != nil {
			return nil, err
		}
	}

	// Build the user response objects.
	for _, id := range *followList {
		wg.Add(1)
		go func(id int64, req *follow.FollowListRequest, userList *[]*follow.User, wg *sync.WaitGroup, mu *sync.Mutex) {
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

			follow := pack.User(user) // Convert the RPC user to the response type.

			mu.Lock()
			*userList = append(*userList, follow)
			mu.Unlock()
		}(id, req, &userList, &wg, &mu)
	}

	wg.Wait()

	if isErr {
		return nil, errno.ServiceError
	}

	return &userList, nil
}
