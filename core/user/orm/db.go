package orm

import (
	"errors"
	"github.com/yu-org/yu/common"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDB(db *gorm.DB) (*Database, error) {
	d := &Database{db}
	err := d.AutoMigrate(&UserScheme{})
	return d, err
}

func (db *Database) GetUser(addr common.Address) (*UserScheme, error) {
	var user UserScheme
	err := db.Model(&UserScheme{Addr: addr.String()}).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, err
}

func (db *Database) SetUser(addr common.Address, nickName, contactMe string) error {
	return db.Create(&UserScheme{Addr: addr.String(), NickName: nickName, ContactMe: contactMe}).Error
}

func (db *Database) SetUserReputation(addr common.Address, value uint64) error {
	return db.Create(&UserScheme{Addr: addr.String(), ReputationValue: int64(value)}).Error
}

func (db *Database) SetUserIfNotExist(addr common.Address) error {
	err := db.First(&UserScheme{}, addr.String()).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return db.Create(&UserScheme{Addr: addr.String()}).Error
		}
		return err
	}
	return nil
}

func (db *Database) IncreaseReputation(addr common.Address, value uint64) error {
	return db.Model(&UserScheme{Addr: addr.String()}).
		UpdateColumn("reputation_value", gorm.Expr("reputation_value + ?", value)).
		Error
}

func (db *Database) ReduceReputation(addr common.Address, value uint64) error {
	return db.Model(&UserScheme{Addr: addr.String()}).
		UpdateColumn("reputation_value", gorm.Expr("reputation_value - ?", value)).
		Error
}
