package gorm

import (
	"errors"
	"fmt"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"time"
)

type clickhouseConfig struct {
	generalDB
}

func (c *clickhouseConfig) Dsn() string {
	return fmt.Sprintf("clickhouse://%s:%s@%s/%s?%s", c.UserName, c.Password, c.Host, c.DbName, c.Config)
}

func NewClickhouse(c clickhouseConfig) (*gorm.DB, error) {
	if c.DbName == "" {
		return nil, errors.New("dbname not set")
	}

	if db, err := gorm.Open(clickhouse.Open(c.Dsn()), &gorm.Config{}); err != nil {
		panic(err)
	} else {
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
}
