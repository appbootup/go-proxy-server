package db

import (
	_ "github.com/lib/pq"
	
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var (
	RecordNotFound = gorm.ErrRecordNotFound
)

func Get() *gorm.DB {
	return db
}

// To make connection with db. Keep maxIdleConnection at min 5 and
// maxOpenConnections at min 10.
func Connect(url string, maxIdleConnections, maxOpenConnections int) error {
	var err error
	db, err = gorm.Open("postgres", url)
	if err != nil {
		return err
	}
	// Get database connection handle [*sql.DB](http://golang.org/pkg/database/sql/#DB)
	db.DB()
	// Then you could invoke `*sql.DB`'s functions with it
	err = db.DB().Ping()
	if err != nil {
		return err
	}
	db.LogMode(false)
	db.DB().SetMaxIdleConns(maxIdleConnections)
	db.DB().SetMaxOpenConns(maxOpenConnections)
	db.SingularTable(false)
	return nil
}

func Close() {
	db.Close()
}
