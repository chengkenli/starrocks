package conn

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type Connect struct {
	Host string
	Port int
	User string
	Pass string
	Base string
}

func (c *Connect) StarRocks() (*gorm.DB, error) {
	newLogger := logger.New(nil,
		logger.Config{
			SlowThreshold: time.Second * 1000, // 控制慢SQL阈值
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.Base,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: newLogger,
	})
	return db, err
}
