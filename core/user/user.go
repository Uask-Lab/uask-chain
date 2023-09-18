package user

import (
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/tripod"
	"gorm.io/gorm"
	"uask-chain/core/user/orm"
	"uask-chain/types"
)

type User struct {
	*tripod.Tripod
	db *orm.Database
}

func NewUser(db *gorm.DB) *User {
	database, err := orm.NewDB(db)
	if err != nil {
		logrus.Fatal("init user db failed: ", err)
	}
	tri := tripod.NewTripod()
	user := &User{tri, database}
	return user
}

const DefaultReputation = 1

func (u *User) CheckReputation(addr common.Address, need uint64) error {
	user, err := u.db.GetUser(addr)
	if err != nil {
		return err
	}
	if user == nil {
		return checkReputation(DefaultReputation, need)
	}
	return checkReputation(user.ReputationValue, need)
}

func (u *User) IncreaseReputation(addr common.Address, value uint64) error {
	return u.db.IncreaseReputation(addr, value)
}

func (u *User) ReduceReputation(addr common.Address, value uint64) error {
	return u.db.ReduceReputation(addr, value)
}

func checkReputation(have int64, need uint64) error {
	if have >= int64(need) {
		return nil
	}
	return types.ErrReputationValueInsufficient
}
