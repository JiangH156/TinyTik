package model

const (
	UNFOLLOW = iota
	FOLLOW
)

type Relation struct {
	Id       int64 `gorm:"primaryKey;autoIncrement"`
	UserId   int64 `gorm:"primaryKey;foreignKey;reference:users.id"`
	ToUserId int64 `gorm:"primaryKey;foreignKey;reference:users.id"`
	Status   byte
}
