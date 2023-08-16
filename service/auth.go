package service

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/repository"
	"TinyTik/utils/logger"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type AuthService struct {
	DB *gorm.DB
}

// 用户注册
func (u *AuthService) Register(regAuth model.UserAuth) (id int64, token string, lErr common.LError) {
	// 用户是否存在
	authRepository := repository.NewAuthRepository()
	_, err := authRepository.GetIDByUsername(regAuth.UserName)
	// 1.正确查询到了用户id
	if err == nil {
		logger.Info("user already exist")
		return 0, "", common.LError{
			HttpCode: http.StatusBadRequest,
			Msg:      "user already exist",
			Err:      errors.New("user already exist"),
		}
	}
	//2.发生非ErrRecordNotFound错误
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("authRepository.GetIDByUsername error:", err)
		return 0, "", common.LError{
			HttpCode: http.StatusBadRequest,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}
	// UserAuth 存储信息
	// 1.密码加密
	password, err := bcrypt.GenerateFromPassword([]byte(regAuth.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("bcrypt.GenerateFromPassword error:", err)
		return 0, "", common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}
	creUserAuth := model.UserAuth{UserName: regAuth.UserName, Password: string(password)}
	tx := u.DB.Begin() // 开启事务
	err = authRepository.CreateAuth(tx, &creUserAuth)
	if err != nil {
		logger.Error("authRepository.CreateAuth error:", err)
		tx.Rollback() // 回滚事务
		return 0, "", common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}
	// User 存储信息
	creUser := model.User{
		Id:            creUserAuth.ID, // 操作的userauth，存放id
		Name:          creUserAuth.UserName,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	userRepository := repository.NewUserRepository()
	err = userRepository.CreateUser(tx, creUser)
	if err != nil {
		logger.Error("userRepository.CreateUser error:", err)
		tx.Rollback() // 回滚事务
		return 0, "", common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}

	// 生成token
	token, err = common.GenToken(creUserAuth.UserName)
	if err != nil {
		logger.Error("common.GenToken error:", err)
		tx.Rollback() // 回滚事务
		return 0, "", common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}
	// 注册的用户，直接添加到redis缓存中
	// json格式化user
	userBytes, err := json.Marshal(creUser)
	if err != nil {
		logger.Error("json.Marshal error:", err)
		tx.Rollback() // 回滚事务
		return 0, "", common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}

	redis := common.GetRedisClient()
	err = redis.SetUser(token, userBytes)
	if err != nil {
		logger.Error("redis.SetUser error:", err)
		tx.Rollback() // 回滚事务
		return 0, "", common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}
	// 注册成功，事务提交
	tx.Commit()
	return creUserAuth.ID, token, common.LError{
		HttpCode: http.StatusOK,
		Msg:      "register success!",
		Err:      nil,
	}
}

func (u *AuthService) Login(loginAuth model.UserAuth) (id int64, token string, lErr common.LError) {
	// 数据库查询数据
	authRepository := repository.NewAuthRepository()
	auth, err := authRepository.GetAuthByUsername(loginAuth.UserName)
	if err != nil {
		//1.发生ErrRecordNotFound错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, "", common.LError{
				HttpCode: http.StatusBadRequest,
				Msg:      "User doesn't exist",
				Err:      errors.New("User doesn't exist"),
			}
		}
		// 2.发生其他错误
		logger.Error("authRepository.GetAuthByUsername error", err)
		return 0, "", common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}
	// 账号密码校验
	if err = bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(loginAuth.Password)); err != nil {
		logger.Error("The account number or password is incorrect")
		return 0, "", common.LError{
			HttpCode: http.StatusBadRequest,
			Msg:      "The account number or password is incorrect",
			Err:      errors.New("The account number or password is incorrect"),
		}
	}

	// 重发发放token，放入redis缓存
	token, err = common.GenToken(loginAuth.UserName)
	if err != nil {
		logger.Error("common.GenToken error:", err)
		return 0, "", common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}
	// 登录的用户，直接添加到redis缓存中
	userRepository := repository.NewUserRepository()
	user, err := userRepository.GetUserById(auth.ID)
	// 用户是存在的，错误不是ErrRecordNotFound错误
	if err != nil {
		logger.Error("userRepository.GetUserById error", err)
		return 0, "", common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}
	// json格式化user
	userBytes, err := json.Marshal(user)
	if err != nil {
		logger.Error("json.Marshal error:", err)
		return 0, "", common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}
	// redis client
	redis := common.GetRedisClient()
	err = redis.SetUser(token, userBytes)
	if err != nil {
		logger.Error("redis.SetUser error:", err)
		return 0, "", common.LError{
			HttpCode: http.StatusInternalServerError,
			Msg:      "Internal server error",
			Err:      errors.New("Internal server error"),
		}
	}

	return auth.ID, token, common.LError{
		HttpCode: http.StatusOK,
		Msg:      "Login success!",
		Err:      nil,
	}

}

func NewAuthService() *AuthService {
	return &AuthService{
		DB: common.GetDB(),
	}
}
