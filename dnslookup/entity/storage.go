package entity

import "github.com/redis/go-redis/v9"

type Storage struct {
	Client *redis.Client
}
