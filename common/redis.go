package common

import (
	"TinyTik/utils/logger"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

type RedisConfig struct {
	Address  string
	Port     int
	Password string
	DB       int
}

var Redis *RedisClient

func RedisSetup(cfg *RedisConfig) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), // RedisClient 服务器地址
		Password: cfg.Password,                                // RedisClient 密码
		DB:       cfg.DB,                                      // 使用的 RedisClient 数据库
	})
	Redis = &RedisClient{client: client}
}

// 缓存Token令牌
func (r *RedisClient) SetToken(tokenString string, expire time.Duration) error {
	err := r.client.Set(context.Background(), tokenString, tokenString, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

// 获取Token令牌
func (r *RedisClient) GetToken(tokenString string) (token string, err error) {
	result, err := r.client.Get(context.Background(), tokenString).Result()
	if err == redis.Nil {
		// Token不存在
		logger.Fatal("Token不存在")
		return "", errors.New("Token不存在")
	} else if err != nil {
		//错误处理
		logger.Error("获取Token令牌错误")
		return "", err
	}
	return result, nil
}

// 删除Token令牌
func (r *RedisClient) DelToken(tokenString string) error {
	err := r.client.Del(context.Background(), tokenString).Err()
	if err != nil {
		return err
	}
	return nil
}
