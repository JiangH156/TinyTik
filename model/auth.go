package model

type UserAuth struct {
	ID       int64  `gorm:"id;primaryKey;autoIncrement;comment:用户id"`
	UserName string `gorm:"user_name;type:varchar(32);not null;comment:用户名称" validate:"required,gte=6,lte=32"`
	Password string `gorm:"password;varchar(32);not null;comment:用户密码" validate:"required,gte=6,lte=32"`
}
