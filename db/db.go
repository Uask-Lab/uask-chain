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

	d := &Database{db}
	err = d.AutoMigrate(&types.QuestionScheme{}, &types.AnswerScheme{}, &types.CommentScheme{})
	return d, err
}

func (db *Database) AddQuestion(q *types.QuestionScheme) error {
	return db.Create(q).Error
}

func (db *Database) UpdateQuestion(q *types.QuestionScheme) error {
	return db.Model(&types.QuestionScheme{ID: q.ID}).Updates(q).Error
}

func (db *Database) GetQuestion(id string) (question *types.QuestionScheme, err error) {
	question = new(types.QuestionScheme)
	err = db.Model(&types.QuestionScheme{ID: id}).Limit(1).Find(question).Error
	return
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
	return db.Model(&types.AnswerScheme{ID: a.ID}).Updates(a).Error
}

func (db *Database) GetAnswer(id string) (answer *types.AnswerScheme, err error) {
	answer = new(types.AnswerScheme)
	err = db.Model(&types.AnswerScheme{ID: id}).Limit(1).Find(answer).Error
	return
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
	return db.Model(&types.CommentScheme{ID: c.ID}).Updates(c).Error
}

func (db *Database) GetComment(id string) (comment *types.CommentScheme, err error) {
	comment = new(types.CommentScheme)
	err = db.Model(&types.CommentScheme{ID: id}).Limit(1).Find(comment).Error
	return
}

func (db *Database) QueryComments(query interface{}) (comments []*types.CommentScheme, err error) {
	err = db.Where(query).Find(&comments).Error
	return
}

func (db *Database) DeleteComment(id string) error {
	return db.Delete(&types.CommentScheme{ID: id}).Error
}
