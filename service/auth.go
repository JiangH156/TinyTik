package service

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/repositoy"
	"TinyTik/utils/logger"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type AuthService struct {
	DB *gorm.DB
}

// 用户注册
func (u *AuthService) Register(userAuth model.UserAuth) (id int64, lErr common.LError) {
	// 用户是否存在
	authRepository := repositoy.NewAuthRepository()
	_, err := authRepository.GetIDByUsername(userAuth.UserName)
	// 1.正确查询到了用户id
	if err == nil {
		logger.Info("user already exist")
		return 0, common.LError{
			HttpCode: http.StatusBadRequest,
			Msg:      "user already exist",
			Err:      errors.New("user already exist"),
		}
	}
	//2.发生非ErrRecordNotFound错误
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("authRepository.GetIDByUsername error:", err)
		return 0, common.LError{
			HttpCode: http.StatusBadRequest,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}
	// UserAuth 存储信息
	// 1.密码加密
	password, err := bcrypt.GenerateFromPassword([]byte(userAuth.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("bcrypt.GenerateFromPassword error:", err)
		return 0, common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}
	creUserAuth := model.UserAuth{UserName: userAuth.UserName, Password: string(password)}
	tx := u.DB.Begin() // 开启事务
	err = authRepository.CreateAuth(tx, &creUserAuth)
	if err != nil {
		logger.Error("authRepository.CreateAuth error:", err)
		tx.Rollback() // 回滚事务
		return 0, common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}
	// User 存储信息
	creUser := model.User{
		Id:            creUserAuth.ID, // 操作的userauth，存放id
		Name:          userAuth.UserName,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	userRepository := repositoy.NewUserRepository()
	err = userRepository.CreateUser(tx, creUser)
	if err != nil {
		logger.Error("userRepository.CreateUser error:", err)
		tx.Rollback() // 回滚事务
		return 0, common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}
	// 注册成功，事务提交
	tx.Commit()
	return creUser.Id, common.LError{
		HttpCode: 0,
		Msg:      "register success!",
		Err:      nil,
	}
}

func (u *AuthService) Login(userAuth model.UserAuth) (id int64, lErr common.LError) {

}

func NewAuthService() *AuthService {
	return &AuthService{
		DB: common.GetDB(),
	}
}
