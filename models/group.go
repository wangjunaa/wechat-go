package models

import "time"

type GroupBasic struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Name      string      `json:"name"`
	OwnerId   string      `json:"ownerId"`
	Icon      []byte      `json:"icon"`
	Members   []UserBasic `json:"members" gorm:"many2many:user_groups"`
	Type      int         `json:"type"`
}

func (g *GroupBasic) TableName() string {
	return "GroupBasic"
}
