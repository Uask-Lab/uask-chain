package orm

import (
	"gorm.io/gorm"
	aorm "uask-chain/core/answer/orm"
	corm "uask-chain/core/comment/orm"
	"uask-chain/types"
)

type Database struct {
	*gorm.DB
}

func NewDB(db *gorm.DB) (*Database, error) {
	d := &Database{db}
	err := d.AutoMigrate(&QuestionScheme{})
	return d, err
}

func (db *Database) AddQuestion(q *QuestionScheme) error {
	return db.Create(q).Error
}

func (db *Database) UpdateQuestion(q *QuestionScheme) error {
	return db.Model(&QuestionScheme{ID: q.ID}).Updates(q).Error
}

func (db *Database) GetQuestion(id string) (*QuestionScheme, error) {
	question := new(QuestionScheme)
	err := db.Model(&QuestionScheme{ID: id}).Limit(1).Find(question).Error
	if err == gorm.ErrRecordNotFound {
		return nil, types.ErrQuestionNotFound
	}
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (db *Database) ListQuestions(limit, offset int) (qs []*QuestionScheme, err error) {
	err = db.Model(&QuestionScheme{}).Limit(limit).Offset(offset).Order("timestamp desc").Find(&qs).Error
	return
}

func (db *Database) QueryQuestions(query interface{}) (qs []*QuestionScheme, err error) {
	err = db.DB.Where(query).Find(&qs).Error
	return
}

func (db *Database) DeleteQuestion(id string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&QuestionScheme{ID: id}).Error
		if err != nil {
			return err
		}
		// delete all comments of this questions
		err = tx.Where(&corm.CommentScheme{QID: id}).Delete(new(corm.CommentScheme)).Error
		if err != nil {
			return err
		}
		// delete all answers of this questions and these quesions' comments.
		var answersIDs []string
		err = tx.Model(&aorm.AnswerScheme{}).Where(&aorm.AnswerScheme{QID: id}).Pluck("id", &answersIDs).Error
		if err != nil {
			return err
		}
		err = tx.Model(&corm.CommentScheme{}).Where("aid IN ?", answersIDs).Delete(new(corm.CommentScheme)).Error
		if err != nil {
			return err
		}
		return tx.Where(&aorm.AnswerScheme{QID: id}).Delete(new(aorm.AnswerScheme)).Error
	})
}
