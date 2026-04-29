package cache

import (
	"context"
	"encoding/json"
	"mini-lead-crm/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

var ctx = context.Background()

func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr, //localhost:8000
	})

	//ping it to test the connection, if it fails we log (cuz optional layer)
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return &RedisClient{client: nil} //fallback to db
	}

	return &RedisClient{client: rdb}
}

func (r *RedisClient) Get(id string) (*models.Lead, error) {
	if r.client == nil {
		return nil, redis.Nil
	}

	val, err := r.client.Get(ctx, "lead:"+id).Result()
	if err != nil {
		return nil, err
	}

	var lead models.Lead
	// json string back to Lead struct
	err = json.Unmarshal([]byte(val), &lead)
	if err != nil {
		return nil, err
	}
	return &lead, nil
}

// store a lead in cache TTL = 10m
func (r *RedisClient) Set(lead *models.Lead) {
	if r.client == nil {
		return
	}

	val, err := json.Marshal(lead)
	if err == nil {
		r.client.Set(ctx, "lead:"+lead.ID.String(), val, 10*time.Minute)
	}
}

// remove leads if updated or deleted

func (r *RedisClient) Invalidate(id string) {
	if r.client == nil {
		return
	}
	r.client.Del(ctx, "lead:"+id)
}
