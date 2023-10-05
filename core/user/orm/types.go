package orm

type UserScheme struct {
	Addr            string `json:"addr" gorm:"primaryKey;column:addr"`
	NickName        string `json:"nick_name" gorm:"column:nick_name"`
	ContactMe       string `json:"contact_me" gorm:"contact_me"`
	ReputationValue int64  `json:"reputation_value" gorm:"column:reputation_value;default:1"`
}

func (UserScheme) TableName() string {
	return "users"
}
