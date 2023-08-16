package main

import (
	"TinyTik/common"
	"TinyTik/router"
	"TinyTik/utils/logger"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	once sync.Once
)

func main() {
	// 加载配置
	loadConfig()

	r := gin.Default()
	// 加载路由
	router.InitRouter(r)

	Run(r)
}

// 服务启动
func Run(r *gin.Engine) {
	address := viper.GetString("server.address")
	port := viper.GetInt("server.port")
	r.Run(fmt.Sprintf("%s:%d", address, port))
}

func loadConfig() {
	// 配置viper相关信息
	dir, _ := os.Getwd()
	parentDir := filepath.Dir(dir)

	// 配置文件所在目录
	viper.AddConfigPath(parentDir + "/TinyTik/config/")
	viper.SetConfigName("application_dev")
	// 配置文件类型
	viper.SetConfigType("yml")
	// 读取配置信息
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fail to config viper, %s", err))
	}
	// 使用Once
	once.Do(func() {
		//配置mysql
		common.InitDB()
		// 配置logger
		loadLogger()
		// 配置redis
		loadRedis()
	})

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
	//fmt.Println(timeFormat)
	logger.Setup(&logger.Settings{
		Path:       path,
		Name:       name,
		Ext:        ext,
		TimeFormat: timeFormat,
	})
}
