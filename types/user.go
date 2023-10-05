package types

import "github.com/yu-org/yu/common"

type UserInfo struct {
	Addr            common.Address `json:"addr"`
	NickName        string         `json:"nick_name"`
	ContactMe       string         `json:"contact_me"`
	ReputationValue int64          `json:"reputation_value"`
}

type UserRegisterRequest struct {
	Addr      string `json:"addr"`
	NickName  string `json:"nick_name"`
	ContactMe string `json:"contact_me"`
}
