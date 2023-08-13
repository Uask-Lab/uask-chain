package orm

import (
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDB(db *gorm.DB) (*Database, error) {
	d := &Database{db}
	err := d.AutoMigrate()
	return d, err
}

func (db *Database) InsertZone() error {
	return nil
}

func (db *Database) UpdateZone() error {
	return nil
}

func (db *Database) DeleteZone() error {
	return nil
}
