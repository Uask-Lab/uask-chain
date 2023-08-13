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
	err := d.AutoMigrate(&CommentScheme{})
	return d, err
}

func (db *Database) AddComment(c *CommentScheme) error {
	return db.Create(c).Error
}

func (db *Database) UpdateComment(c *CommentScheme) error {
	return db.Model(&CommentScheme{ID: c.ID}).Updates(c).Error
}

func (db *Database) GetComment(id string) (*CommentScheme, error) {
	comment := new(CommentScheme)
	err := db.Model(&CommentScheme{ID: id}).First(comment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, types.ErrCommentNotFound
	}
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (db *Database) QueryComments(query interface{}) (comments []*CommentScheme, err error) {
	err = db.Where(query).Find(&comments).Error
	return
}

func (db *Database) DeleteComment(id string) error {
	return db.Delete(&CommentScheme{ID: id}).Error
}
