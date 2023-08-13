package orm

type ZoneScheme struct {
}

func (ZoneScheme) TableName() string {
	return "zone"
}
