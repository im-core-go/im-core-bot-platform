package infra

import (
	"fmt"
	"web-chat/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newMysql(cfg configs.Config) *gorm.DB {
	mysqlConf := cfg.MysqlConf
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		mysqlConf.User,
		mysqlConf.Password,
		mysqlConf.Host,
		mysqlConf.Port,
		mysqlConf.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
<<<<<<< Updated upstream
=======
	db.AutoMigrate(
		&model.User{},
		&model.Message{},
		&model.Conversation{},
	)
>>>>>>> Stashed changes
	return db
}
