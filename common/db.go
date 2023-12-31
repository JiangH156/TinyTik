package common

import (
	"TinyTik/model"
	"fmt"
	"net/url"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	host := viper.GetString("datasource.host")
	port := viper.GetInt("datasource.port")
	database := viper.GetString("datasource.database")
	charset := viper.GetString("datasource.charset")
	parseTime := viper.GetString("datasource.parseTime")
	loc := viper.GetString("datasource.loc")
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		username, password, host, port, database, charset, parseTime, url.QueryEscape(loc))

	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		panic(fmt.Sprintf("fail to init database, %s\n", err))
	}
	db.AutoMigrate(model.UserAuth{}, model.User{}, model.Message{}, model.Comment{}, model.Video{}, model.Like{}, model.Relation{})
	DB = db
}
func GetDB() *gorm.DB {
	return DB
}
