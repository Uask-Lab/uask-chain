package core

import (
	"errors"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
)

type WhiteList struct {
	*tripod.Tripod
	administrators []common.Address
}

var NoAdminError = errors.New("you have no permission to operate white list")

func NewWhiteList(administrators []common.Address) *WhiteList {
	tri := tripod.NewTripod()
	w := &WhiteList{tri, administrators}
	w.SetWritings(w.AddMember, w.DeleteMember)
	w.SetReadings(w.ListMembers)
	return w
}

func (w *WhiteList) AddMember(ctx *context.WriteContext) error {
	operator := ctx.GetCaller()
	if !w.isAdministrator(operator) {
		return NoAdminError
	}

	return nil
}

func (w *WhiteList) DeleteMember(ctx *context.WriteContext) error {
	operator := ctx.GetCaller()
	if !w.isAdministrator(operator) {
		return NoAdminError
	}
	return nil
}

func (w *WhiteList) ListMembers(ctx *context.ReadContext) {

}

func (w *WhiteList) isAdministrator(op common.Address) bool {
	for _, admin := range w.administrators {
		if admin == op {
			return true
		}
	}
	return false
}
