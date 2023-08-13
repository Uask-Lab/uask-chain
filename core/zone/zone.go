package zone

import (
	"encoding/json"
	"errors"
	"github.com/yu-org/yu/apps/poa"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"uask-chain/types"
)

var (
	ErrPermissionDenied = errors.New("you have no permission to operate")
	ErrZoneExist        = errors.New("zone exist")
)

type Zone struct {
	*tripod.Tripod
	Poa *poa.Poa `tripod:"poa"`
}

func NewZone() *Zone {
	t := tripod.NewTripod()
	zone := &Zone{Tripod: t}
	zone.SetWritings(
		zone.Apply,
		zone.Revoke,
		zone.TransferOwner,
		zone.AddManagers,
		zone.DeleteManagers,
		zone.ExamineZone,
	)
	return zone
}

func (z *Zone) Apply(ctx *context.WriteContext) error {
	zone := new(types.ZoneInfo)
	err := ctx.BindJson(zone)
	if err != nil {
		return err
	}

	if z.existZone(zone.Name) {
		return ErrZoneExist
	}

	proposal := &types.ZoneProposal{
		ID:     ctx.GetTxnHash().String(),
		Info:   zone,
		Status: types.Pending,
	}

	byt, err := json.Marshal(proposal)
	if err != nil {
		return err
	}
	z.Set([]byte(proposal.ID), byt)

	return nil
}

func (z *Zone) Revoke(ctx *context.WriteContext) error {
	return nil
}

func (z *Zone) ExamineZone(ctx *context.WriteContext) error {
	validator := ctx.GetCaller()
	if !z.Poa.IsValidator(validator) {
		return ErrPermissionDenied
	}
	proposalID := ctx.GetString("id")
	result := ctx.GetInt("result")
	proposalByt, err := z.Get([]byte(proposalID))
	if err != nil {
		return err
	}
	proposal := new(types.ZoneProposal)
	err = json.Unmarshal(proposalByt, proposal)
	if err != nil {
		return err
	}
	proposal.Status = result
	byt, err := json.Marshal(proposal)
	if err != nil {
		return err
	}
	z.Set([]byte(proposalID), byt)
	return nil
}

func (z *Zone) TransferOwner(ctx *context.WriteContext) error {
	zoneName := ctx.GetString("name")
	owner := ctx.GetCaller()
	if !z.isOwner(owner, zoneName) {
		return ErrPermissionDenied
	}
	return nil
}

func (z *Zone) AddManagers(ctx *context.WriteContext) error {
	zoneName := ctx.GetString("name")
	owner := ctx.GetCaller()
	if !z.isOwner(owner, zoneName) {
		return ErrPermissionDenied
	}
	return nil
}

func (z *Zone) DeleteManagers(ctx *context.WriteContext) error {
	zoneName := ctx.GetString("name")
	owner := ctx.GetCaller()
	if !z.isOwner(owner, zoneName) {
		return ErrPermissionDenied
	}
	return nil
}

func (z *Zone) isOwner(addr common.Address, zoneName string) bool {
	return false
}

func (z *Zone) existZone(name string) bool {
	return z.Exist([]byte(name))
}
