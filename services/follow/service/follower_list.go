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

// FollowerList View fan list
func (s *FollowService) FollowerList(req *follow.FollowerListRequest) (*[]*follow.User, error) {
	// 限流
	if err := cache.Limit(s.ctx, constants.FollowerListRate, constants.Interval); err != nil {
		return nil, err
	}

	userList := make([]*follow.User, 0, 10)
	var wg sync.WaitGroup
	var mu sync.Mutex
	isErr := false

	// 先查redis
	followerList, err := cache.FollowerListAction(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	} else if len(*followerList) == 0 { // redis中查不到再查db
		followerList, err = db.FollowerListAction(s.ctx, req.UserId)
		if errors.Is(err, db.RecordNotFound) { // db中也查不到
			return &userList, nil
		} else if err != nil {
			return nil, err
		}
		// db中查到后写入redis
		err := cache.UpdateFollowerList(s.ctx, req.UserId, followerList)
		if err != nil {
			return nil, err
		}
	}

	// 数据处理
	for _, id := range *followerList {
		wg.Add(1)
		go func(id int64, req *follow.FollowerListRequest, userList *[]*follow.User, wg *sync.WaitGroup, mu *sync.Mutex) {
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
				isErr = true // 报错就修改为true
				mu.Unlock()
				return
			}

			follow := pack.User(user) // 结构体转换

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
