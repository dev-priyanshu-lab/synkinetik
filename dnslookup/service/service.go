package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/redis/go-redis/v9"

	"dnslookup/common"
	"dnslookup/entity"
	"dnslookup/intel"
	"dnslookup/process"
)

type Service struct {
	Config entity.Configuration
	Store  *entity.Storage
}

func New(configPath string) (*Service, error) {
	initLogger()

	config, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	store, err := initStorage()
	if err != nil {
		return nil, err
	}

	return &Service{
		Config: config,
		Store:  store,
	}, nil
}

func (s *Service) Start() {
	intel.UpdateIntel(s.Store, s.Config)
	process.StartSchedulerAsync(s.Store, s.Config)
	go process.StartUDPServer(s.Store, s.Config)
}

func initLogger() {
	if err := os.MkdirAll(filepath.Dir(common.APP_LOG_FILE), 0755); err != nil {
		slog.Error("failed to create log directory", "path", filepath.Dir(common.APP_LOG_FILE), "error", err)
		panic("failed to create log directory")
	}

	file, err := os.OpenFile(common.APP_LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		slog.Error("failed to open log file", "error", err)
		panic("failed to open log file")
	}

	logger := slog.New(slog.NewJSONHandler(file, nil))
	slog.SetDefault(logger)
}

func loadConfig(configPath string) (entity.Configuration, error) {
	if configPath == "" {
		configPath = common.DNS_LOOKUP_CONFIG_FILE
	}

	var config entity.Configuration
	data, err := os.ReadFile(configPath)
	if err != nil {
		slog.Error("failed to read lookup config", "path", configPath, "error", err)
		return config, fmt.Errorf("failed to read lookup config %q: %w", configPath, err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		slog.Error("failed to parse lookup config", "path", configPath, "error", err)
		return config, fmt.Errorf("failed to parse lookup config %q: %w", configPath, err)
	}

	return config, nil
}

func initStorage() (*entity.Storage, error) {
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

	if err := store.Client.Ping(common.Context).Err(); err != nil {
		return nil, fmt.Errorf("redis not connected: %w", err)
	}

	return store, nil
}
