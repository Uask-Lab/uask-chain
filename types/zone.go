package types

type ZoneInfo struct {
	// zone name
	Name        string
	Description string
	Owner       *Role
	Managers    []*Role
}

type ZoneProposal struct {
	ID     string    `json:"id"`
	Info   *ZoneInfo `json:"info,omitempty"`
	Status int       `json:"status"`
}

const (
	Pending = iota
	Approved
	Rejected
)

type Role struct {
	Pubkey       string
	NickName     string
	Introduction string
}
