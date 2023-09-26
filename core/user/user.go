package user

import (
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"gorm.io/gorm"
	"uask-chain/core/user/orm"
	"uask-chain/types"
)

type User struct {
	*tripod.Tripod
	db *orm.Database
}

func NewUser(db *gorm.DB, whiteList map[string]uint64) *User {
	database, err := orm.NewDB(db)
	if err != nil {
		logrus.Fatal("init user db failed: ", err)
	}
	tri := tripod.NewTripod()
	user := &User{tri, database}

	user.SetWritings(user.RegisterUser)
	user.SetReadings(user.GetUser)

	for addrStr, reputation := range whiteList {
		addr := common.HexToAddress(addrStr)
		err = user.IncreaseReputation(addr, reputation)
		if err != nil {
			logrus.Fatal("load white list error: ", err)
		}
	}
	return user
}

func (u *User) RegisterUser(ctx *context.WriteContext) error {
	var req types.UserRegisterRequest
	err := ctx.BindJson(&req)
	if err != nil {
		return err
	}
	err = u.db.SetUser(common.HexToAddress(req.Addr), req.NickName, req.ContactMe)
	if err != nil {
		return err
	}
	return ctx.EmitJsonEvent(map[string]string{"register_user": req.Addr})
}

func (u *User) GetUser(ctx *context.ReadContext) {
	user := ctx.GetString("user")
	userSch, err := u.db.GetUser(common.HexToAddress(user))
	if err != nil {
		ctx.JsonOk(types.Error(err))
		return
	}
	ctx.JsonOk(types.Ok(schemeToUser(userSch)))
}

func (u *User) CheckReputation(addr common.Address, need uint64) error {
	user, err := u.db.GetUser(addr)
	if err != nil {
		return err
	}
	if user == nil {
		return checkReputation(types.DefaultReputation, need)
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

func schemeToUser(sch *orm.UserScheme) *types.UserInfo {
	return &types.UserInfo{
		Addr:            common.HexToAddress(sch.Addr),
		NickName:        sch.NickName,
		ContactMe:       sch.ContactMe,
		ReputationValue: sch.ReputationValue,
	}
}
