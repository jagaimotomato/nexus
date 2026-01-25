package service

import (
	"nexus/internal/utils"
	"time"

	"nexus/internal/data"

	"errors"

	"github.com/mojocn/base64Captcha"
	"gorm.io/gorm"
)

var driver = base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)

var store = &utils.RedisStore{Expiration: 5 * time.Minute}

var captcha = base64Captcha.NewCaptcha(driver, store)

func GenerateCaptcha() (id string, b64s string, err error) {
	id, b64s, _, err = captcha.Generate()
	return id, b64s, err
}

type LoginParams struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CaptchaID string `json:"captchaId" binding:"required"`
	Captcha   string `json:"captcha" binding:"required"`
}

type LoginResponse struct {
	Token    string     `json:"token"`
	UserInfo *data.User `json:"userInfo"`
}

func Login(p LoginParams) (*LoginResponse, error) {
	if !store.Verify(p.CaptchaID, p.Captcha, true) {
		return nil, errors.New("验证码错误或已失效")
	}

	var user data.User
	err := data.DB.Where("username = ?", p.Username).Preload("Roles").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 校验密码
	if !utils.CheckPassword(p.Password, user.Password) {
		return nil, errors.New("账号或密码错误")
	}

	// 校验状态
	if user.Status != 1 {
		return nil, errors.New("账号已被停用")
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		return nil, errors.New("生成token失败")
	}

	return &LoginResponse{
		Token:    token,
		UserInfo: &user,
	}, nil
}
