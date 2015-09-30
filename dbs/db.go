package dbs

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wzshiming/server/cfg"
)

type DB struct {
	gorm.DB
}

type Model gorm.Model

func Conn(us cfg.DbConfig) (DB, error) {
	db, err := gorm.Open(us.Dialect, us.Source)
	return DB{db}, err
}
