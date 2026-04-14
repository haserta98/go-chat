package cmd

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func NewDB() *DB {
	return &DB{db: nil}
}

func (d *DB) Connect() error {
	dsn := "host=localhost user=postgres password=postgres dbname=gorm port=5432 sslmode=disable TimeZone=Europe/Istanbul"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	log.Print("Database connection established")

	pgDB, err := db.DB()
	if err != nil {
		return err
	}

	pgDB.SetConnMaxIdleTime(20)
	pgDB.SetMaxOpenConns(20)
	pgDB.SetConnMaxLifetime(time.Hour)

	d.db = db
	return nil
}

func (d *DB) Migrate(entities ...interface{}) {
	d.db.AutoMigrate(entities...)
}

func (d *DB) GetDB() *gorm.DB {
	return d.db
}

func (d *DB) Close() error {
	db, err := d.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
