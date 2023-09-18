package types

import "github.com/yu-org/yu/common"

type UserInfo struct {
	Addr            common.Address `json:"addr"`
	NickName        string         `json:"nick_name"`
	ReputationValue uint64         `json:"reputation_value"`
}
