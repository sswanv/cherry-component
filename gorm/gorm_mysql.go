package gorm

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type mysqlConfig struct {
	generalDB
}

func (c *mysqlConfig) Dsn() string {
	return fmt.Sprintf("%v:%v@tcp(%v)/%v?%v", c.UserName, c.Password, c.Host, c.DbName, c.Config)
}

func NewMysql(c mysqlConfig) (*gorm.DB, error) {
	if c.DbName == " " {
		return nil, errors.New("dbname not set")
	}

	config := mysql.Config{
		DSN: c.Dsn(),
	}
	db, err := gorm.Open(mysql.New(config), &gorm.Config{Logger: getLogger()})
	if err != nil {
		return nil, err
	}

	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(c.MaxIdleConnect)
	sqlDb.SetMaxOpenConns(c.MaxOpenConnect)
	sqlDb.SetConnMaxLifetime(time.Minute)
	if err = sqlDb.Ping(); err != nil {
		return nil, err
	}
	if c.LogMode {
		return db.Debug(), nil
	}

	return db, nil
}
