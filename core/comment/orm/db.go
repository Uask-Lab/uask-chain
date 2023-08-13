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
	err := d.AutoMigrate(&types.CommentScheme{})
	return d, err
}

func (db *Database) AddComment(c *types.CommentScheme) error {
	return db.Create(c).Error
}

func (db *Database) UpdateComment(c *types.CommentScheme) error {
	return db.Model(&types.CommentScheme{ID: c.ID}).Updates(c).Error
}

func (db *Database) GetComment(id string) (*types.CommentScheme, error) {
	comment := new(types.CommentScheme)
	err := db.Model(&types.CommentScheme{ID: id}).First(comment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, types.ErrCommentNotFound
	}
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (db *Database) QueryComments(query interface{}) (comments []*types.CommentScheme, err error) {
	err = db.Where(query).Find(&comments).Error
	return
}

func (db *Database) DeleteComment(id string) error {
	return db.Delete(&types.CommentScheme{ID: id}).Error
}
