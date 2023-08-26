package common

import (
	"TinyTik/model"
	"TinyTik/utils/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
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

var (
	RedisA *redis.Client
	Redis  *RedisClient
	expire = time.Hour * 24
	// expire = time.Second * 1
)

func (r *RedisClient) UpdateUser(key string, updatedUser model.User) error {
	// 获取旧的用户数据
	userBytes, err := r.GetUser(key)
	if err != nil {
		// 处理获取用户数据错误
		return err
	}

	// 解码旧的用户数据
	var existingUser model.User
	if err := json.Unmarshal(userBytes, &existingUser); err != nil {
		// 处理解码错误
		return err
	}

	// 将更新后的用户数据存储回 Redis
	updatedUserBytes, err := json.Marshal(updatedUser)
	if err != nil {
		// 处理编码错误
		return err
	}

	err = r.SetUser(key, updatedUserBytes)
	if err != nil {
		// 处理存储到 Redis 出错的情况
		return err
	}

	// 更新成功
	return nil
}

// 缓存登录用户
func (r *RedisClient) SetUser(key string, user []byte) error {
	err := r.Set(key, user, expire)
	if err != nil {
		return err
	}
	return nil
}
func (r *RedisClient) GetUser(key string) (user []byte, err error) {
	u, err := r.Get(key)
	user = ([]byte)(u.(string)) //类型强转
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetRedisClient() *RedisClient {
	return Redis
}

func (r *RedisClient) UserLoginInfo(token string) (userInfo model.User, exist bool) {
	userBytes, err := r.GetUser(token)
	if err != nil {
		return userInfo, false
	}
	// json反序列化
	err = json.Unmarshal(userBytes, &userInfo)
	if err != nil {
		return userInfo, false
	}
	return userInfo, true
}

func (r *RedisClient) Set(key string, value any, expire time.Duration) error {
	err := r.client.Set(context.Background(), key, value, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) Get(key string) (value any, err error) {
	value, err = r.client.Get(context.Background(), key).Result()

	if err == redis.Nil {
		// Token不存在
		logger.Fatal("Token不存在")
		return "", errors.New("Token不存在")

	} else if err != nil {
		//错误处理
		logger.Error("获取Token令牌错误")
		return "", err
	}
	return value, nil
}

func (r *RedisClient) Del(key string) error {
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}
	return nil
}

// 存在返回true，反之返回false
func (r *RedisClient) TokenIsExist(key string) bool {
	_, err := r.Get(key)
	return err == nil
}

func RedisSetup(cfg *RedisConfig) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), // RedisClient 服务器地址
		Password: cfg.Password,                                // RedisClient 密码
		DB:       cfg.DB,                                      // 使用的 RedisClient 数据库
	})
	expireString := viper.GetString("redis.expire")
	expireDuration, err := parseDuration(expireString)
	if err != nil {
		panic("RedisSetup error")
	}
	expire = expireDuration
	RedisA = client
	Redis = &RedisClient{client: client}
}

func parseDuration(durationString string) (time.Duration, error) {
	// 将字符串转换为数字和单位部分
	duration, err := strconv.ParseInt(durationString[:len(durationString)-1], 10, 64)
	if err != nil {
		return 0, err
	}

	// 获取单位部分
	unit := durationString[len(durationString)-1:]

	// 根据单位返回对应的 time.Duration
	switch unit {
	case "s":
		return time.Duration(duration) * time.Second, nil
	case "m":
		return time.Duration(duration) * time.Minute, nil
	case "h":
		return time.Duration(duration) * time.Hour, nil
	case "d":
		return time.Duration(duration) * time.Hour * 24, nil
	default:
		return 0, fmt.Errorf("invalid duration format")
	}
}
