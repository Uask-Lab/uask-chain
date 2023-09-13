package orm

type ZoneScheme struct {
	Name        string
	Description string
}

func (ZoneScheme) TableName() string {
	return "zone"
}

type RoleScheme struct {
	Pubkey       string
	ZoneName     string
	Role         int
	NickName     string
	Introduction string
}

const (
	Owner = iota
	Manager
)

func (RoleScheme) TableName() string {
	return "role"
}
