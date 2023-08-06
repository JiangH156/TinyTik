package main

import (
	"TinyTik/common"
	"TinyTik/router"
	"TinyTik/service"
	"TinyTik/utils/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func main() {
	loadConfig()
	go service.RunMessageServer()
	r := gin.Default()

	router.InitRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func loadConfig() {
	// 配置viper相关信息
	dir, _ := os.Getwd()
	parentDir := filepath.Dir(dir)
	// 配置文件所在目录
	viper.AddConfigPath(parentDir + "/config/")
	// 配置文件名
	viper.SetConfigName("application")
	// 配置文件类型
	viper.SetConfigType("yml")
	// 读取配置信息
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fail to config viper, %s", err))
	}
	// 配置logger
	loadLogger()
	// 配置redis
	loadRedis()
}

func loadRedis() {
	address := viper.GetString("redis.address")
	port := viper.GetInt("redis.port")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")
	common.RedisSetup(&common.RedisConfig{
		Address:  address,
		Port:     port,
		Password: password,
		DB:       db,
	})
}

func loadLogger() {
	path := viper.GetString("logger.path")
	name := viper.GetString("logger.name")
	ext := viper.GetString("logger.ext")
	timeFormat := viper.GetString("logger.timeFormat")
	fmt.Println(timeFormat)
	logger.Setup(&logger.Settings{
		Path:       path,
		Name:       name,
		Ext:        ext,
		TimeFormat: timeFormat,
	})
}
