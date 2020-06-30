package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shark/src/config"
	"github.com/shark/src/model"
	"log"
)

// 初始化数据库连接
func InitConn(user, password, host, db_name, port string) *gorm.DB {
	connArgs := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, db_name)
	log.Println("conStr is : ", connArgs)
	db, err := gorm.Open("mysql", connArgs)
	if err != nil {
		log.Fatal("conn db error ", err)
	}
	db.SingularTable(true)
	db.LogMode(true)
	return db
}

func QueryTaskList() []model.CustomUserProfileTag {
	conf := config.Config().Dbs["online"]
	conn := InitConn(conf.Username, conf.Password, conf.Host, conf.Database, conf.Port)
	var records = make([]model.CustomUserProfileTag, 0)
	conn.Find(&records)
	return records
}
