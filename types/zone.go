package types

import (
	"github.com/yu-org/yu/core/keypair"
)

type ZoneInfo struct {
	// zone name
	Name        string
	Description string
	Owner       *Role
	Managers    []*Role
}

type ZoneProposal struct {
	ID      string    `json:"id"`
	Operate int       `json:"operate"`
	Info    *ZoneInfo `json:"info,omitempty"`
	Status  int       `json:"status"`
}

const (
	AddZone = iota + 1
	UpdateZone
	DeleteZone
)

const (
	Pending = iota
	Approved
	Rejected
)

type Role struct {
	Pubkey       keypair.PubKey
	NickName     string
	Introduction string
}
