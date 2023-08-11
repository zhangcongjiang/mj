package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"mj/constant"
	"time"
)

var db *gorm.DB

func Setup() {
	var dialector gorm.Dialector
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", constant.SqlHost, constant.SqlUser, constant.SqlPassword, constant.SqlDbName, constant.SqlPort)
	//dsn := "host=10.67.37.131 user=postgres password=postgres dbname=venus port=30432 sslmode=disable TimeZone=Asia/Shanghai"
	dialector = postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	})
	conn, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Print(err.Error())
	}
	sqlDB, err := conn.DB()
	if err != nil {
		log.Print("connect db server failed.")
	}
	sqlDB.SetMaxIdleConns(10)                   // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxOpenConns(100)                  // SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetConnMaxLifetime(time.Second * 600) // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db = conn
}

func GetDB() *gorm.DB {
	if db == nil {
		Setup()
	}
	sqlDB, err := db.DB()
	if err != nil {
		Setup()
	}
	if err := sqlDB.Ping(); err != nil {
		err := sqlDB.Close()
		if err != nil {
			return nil
		}
		Setup()
	}

	return db
}
