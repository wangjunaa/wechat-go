package models

type FriendShip struct {
	UserId1 string    `gorm:"primaryKey"`
	UserId2 string    `gorm:"primaryKey"`
	User1   UserBasic `gorm:"foreignKey:UserId1"`
	User2   UserBasic `gorm:"foreignKey:UserId2"`
}

func (f *FriendShip) TableName() string {
	return "FriendShip"
}
