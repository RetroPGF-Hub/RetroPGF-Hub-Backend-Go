package datacenterrepository

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter"
	"context"
	"errors"
	"time"
)

func (r *datacenterRepository) InsertCacheToRedis(pctx context.Context, key string, data string) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	// Zero expiration means the key has no expiration time
	if err := r.redis.Set(ctx, key, data, 0).Err(); err != nil {
		return errors.New("error: insert cache to redis failed" + err.Error())
	}
	return nil
}

func (r *datacenterRepository) InsertManyCacheToRedis(pctx context.Context, pipeData []*datacenter.PipeLineCache) error {

	ctx, cancel := context.WithTimeout(pctx, 25*time.Second)
	defer cancel()

	pipeLine := r.redis.Pipeline()
	for _, v := range pipeData {
		pipeLine.Set(pctx, v.CacheId, v.CacheData, 0)
	}

	// Zero expiration means the key has no expiration time
	_, err := pipeLine.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *datacenterRepository) GetCacheFromRedis(pctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	productJson, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return "", errors.New("error: get cache from redis failed" + err.Error())
	}
	return productJson, err
}

func (r *datacenterRepository) DeleteCacheFromRedis(pctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	if err := r.redis.Del(ctx, key).Err(); err != nil {
		return errors.New("error: delete cache from redis failed" + err.Error())
	}

	return nil
}
