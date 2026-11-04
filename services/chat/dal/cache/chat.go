package cache

import (
	"context"
	"time"
)

func MessageInsert(ctx context.Context, key string, revkey string, stamp_key int64, field string) error {
	// For an existing conversation, extend its expiration when fewer than ten days remain.
	if ok := MessageExist(ctx, key); ok != 0 {
		err := RedisDB.HSet(ctx, key, stamp_key, field).Err()
		if err != nil {
			return err
		}
		lastTime, err := RedisDB.TTL(ctx, key).Result()
		if err != nil {
			return err
		}
		if lastTime.Seconds() < 60*60*24*10 {
			err = RedisDB.Expire(ctx, key, time.Hour*24*30).Err()
			if err != nil {
				return err
			}
			err = RedisDB.Expire(ctx, revkey, time.Hour*24*30).Err()
			if err != nil {
				return err
			}
			return nil
		}
	}
	// Insert the message and apply a 30-day expiration to both conversation keys.
	err := RedisDB.HSet(ctx, key, stamp_key, field).Err()
	if err != nil {
		return err
	}
	err = RedisDB.Expire(ctx, key, time.Hour*24*30).Err()
	if err != nil {
		return err
	}
	err = RedisDB.Expire(ctx, revkey, time.Hour*24*30).Err()
	if err != nil {
		return err
	}
	return nil
}

func MessageGet(ctx context.Context, key string) ([]string, error) {
	mem, err := RedisDB.HVals(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return mem, nil
}

func MessageExist(ctx context.Context, key string) int64 {
	ok, err := RedisDB.Exists(ctx, key).Result()
	if err != nil {
		return 0
	}
	return ok
}
