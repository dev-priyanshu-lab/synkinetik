package main

import (
	"dnslookup/common"
	"dnslookup/entity"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

func initLogger() {
	file, err := os.OpenFile(common.APP_LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		slog.Error("failed to open log file", "error", err)
		panic("failed to open log file")
	}

	logger := slog.New(slog.NewJSONHandler(file, nil))
	slog.SetDefault(logger)
}

func loadConfig() entity.Configuration {
	var config entity.Configuration
	data, err := os.ReadFile(common.DNS_LOOKUP_CONFIG_FILE)
	if err != nil {
		slog.Error("failed to lookup.conf", "error", err)
		panic("failed to open config file")
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		slog.Error("failed to parse lookup.conf", "error", err)
		panic("corrupted config file")
	}

	return config
}

func initStorage() *entity.Storage {
	rdb := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     common.REDIS_ADDR,
		Password: "",
		DB:       0,
		PoolSize: 20,
	})

	store := &entity.Storage{
		Client: rdb,
	}

	err := store.Client.Ping(common.Context).Err()
	if err != nil {
		panic("Redis not connected")
	}

	return store
}

func initialize() (entity.Configuration, *entity.Storage) {
	initLogger()
	config := loadConfig()
	store := initStorage()
	return config, store
}
