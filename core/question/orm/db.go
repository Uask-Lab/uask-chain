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
	err := d.AutoMigrate(&types.QuestionScheme{})
	return d, err
}

func (db *Database) AddQuestion(q *types.QuestionScheme) error {
	return db.Create(q).Error
}

func (db *Database) UpdateQuestion(q *types.QuestionScheme) error {
	return db.Model(&types.QuestionScheme{ID: q.ID}).Updates(q).Error
}

func (db *Database) GetQuestion(id string) (*types.QuestionScheme, error) {
	question := new(types.QuestionScheme)
	err := db.Model(&types.QuestionScheme{ID: id}).Limit(1).Find(question).Error
	if err == gorm.ErrRecordNotFound {
		return nil, types.ErrQuestionNotFound
	}
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (db *Database) ListQuestions(limit, offset int) (qs []*types.QuestionScheme, err error) {
	err = db.Model(&types.QuestionScheme{}).Limit(limit).Offset(offset).Order("timestamp desc").Find(&qs).Error
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
		err = tx.Model(&types.AnswerScheme{}).Where(&types.AnswerScheme{QID: id}).Pluck("id", &answersIDs).Error
		if err != nil {
			return err
		}
		err = tx.Model(&types.CommentScheme{}).Where("aid IN ?", answersIDs).Delete(new(types.CommentScheme)).Error
		if err != nil {
			return err
		}
		return tx.Where(&types.AnswerScheme{QID: id}).Delete(new(types.AnswerScheme)).Error
	})
}
