package models

type Role struct {
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
}

type RoleGet struct {
	Id int `json:"id"`
	Role
}

type Tabler interface {
	TableName() string
}

func (RoleGet) TableName() string {
	return "roles"
}
