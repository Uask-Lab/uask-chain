package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"uask-chain/types"
)

type Database struct {
	*gorm.DB
}

func NewDB(dsn string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{CreateBatchSize: 50000})
	if err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

func (db *Database) AddQuestion(q *types.QuestionScheme) error {
	return db.Create(q).Error
}

func (db *Database) UpdateQuestion(q *types.QuestionScheme) error {
	return db.Save(q).Error
}

func (db *Database) QueryQuestions(query interface{}) (qs []*types.QuestionScheme, err error) {
	err = db.DB.Where(query).Find(&qs).Error
	return
}

func (db *Database) DeleteQuestion(id string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&types.QuestionScheme{ID: id}).Error
		if err != nil {
			return err
		}
		// delete all comments of this questions
		err = tx.Where(&types.CommentScheme{QID: id}).Delete(new(types.CommentScheme)).Error
		if err != nil {
			return err
		}
		// delete all answers of this questions and these quesions' comments.
		var answersIDs []string
		err = tx.Select("aid").Where(&types.AnswerScheme{QID: id}).Scan(&answersIDs).Error
		if err != nil {
			return err
		}
		err = tx.Where("aid IN ?", answersIDs).Delete(new(types.CommentScheme)).Error
		if err != nil {
			return err
		}
		return tx.Where(&types.AnswerScheme{QID: id}).Delete(new(types.AnswerScheme)).Error
	})
}

func (db *Database) AddAnswer(a *types.AnswerScheme) error {
	return db.Create(a).Error
}

func (db *Database) UpdateAnswer(a *types.AnswerScheme) error {
	return db.Save(a).Error
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

func (db *Database) AddComment(c *types.CommentScheme) error {
	return db.Create(c).Error
}

func (db *Database) UpdateComment(c *types.CommentScheme) error {
	return db.Save(c).Error
}

func (db *Database) QueryComments(query interface{}) (comments []*types.CommentScheme, err error) {
	err = db.Where(query).Find(&comments).Error
	return
}

func (db *Database) DeleteComment(id string) error {
	return db.Delete(&types.CommentScheme{ID: id}).Error
}
