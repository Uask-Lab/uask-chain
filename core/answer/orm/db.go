package orm

import (
	"gorm.io/gorm"
	"uask-chain/types"
)

type Database struct {
	*gorm.DB
}

func NewDB(db *gorm.DB) (*Database, error) {
	d := &Database{db}
	err := d.AutoMigrate(&types.AnswerScheme{})
	return d, err
}

func (db *Database) AddAnswer(a *types.AnswerScheme) error {
	return db.Create(a).Error
}

func (db *Database) UpdateAnswer(a *types.AnswerScheme) error {
	return db.Model(&types.AnswerScheme{ID: a.ID}).Updates(a).Error
}

func (db *Database) GetAnswer(id string) (*types.AnswerScheme, error) {
	answer := new(types.AnswerScheme)
	err := db.Model(&types.AnswerScheme{ID: id}).First(answer).Error
	if err == gorm.ErrRecordNotFound {
		return nil, types.ErrAnswerNotFound
	}
	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (db *Database) QueryAnswers(query interface{}) (answers []*types.AnswerScheme, err error) {
	err = db.Where(query).Find(&answers).Error
	return
}

func (db *Database) DeleteAnswer(id string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&types.AnswerScheme{ID: id}).Error
		if err != nil {
			return err
		}
		return tx.Where(&types.CommentScheme{AID: id}).Delete(new(types.CommentScheme)).Error
	})
}
