package cmd

import (
	"fmt"
	"log"
	"os"
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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	log.Print("Database connection established")

	pgDB, err := db.DB()
	if err != nil {
		return err
	}

	pgDB.SetConnMaxIdleTime(20 * time.Second)
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
