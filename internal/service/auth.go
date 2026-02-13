package service

import (
	"errors"
	"nexus/internal/conf"
	"nexus/internal/data"
	"nexus/internal/utils"
	"time"

	"github.com/mojocn/base64Captcha"
)

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

type AuthService struct {
	userRepo *data.UserRepo
	captcha  *base64Captcha.Captcha
	store    *utils.RedisStore
	jwtCfg   utils.JWTConfig
}

func NewAuthService(d *data.Data, cfg *conf.Config) *AuthService {
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	store := &utils.RedisStore{
		Expiration: 5 * time.Minute,
		Client:     d.RDB,
	}
	return &AuthService{
		userRepo: data.NewUserRepo(d.DB),
		captcha:  base64Captcha.NewCaptcha(driver, store),
		store:    store,
		jwtCfg: utils.JWTConfig{
			Secret: cfg.Jwt.Secret,
			Expire: cfg.Jwt.Expire,
			Issuer: cfg.Jwt.Issuer,
		},
	}
}

func (s *AuthService) GetCaptcha() (id string, b64s string, err error) {
	id, b64s, _, err = s.captcha.Generate()
	return id, b64s, err
}

func (s *AuthService) Login(p LoginParams) (*LoginResponse, error) {
	if !s.store.Verify(p.CaptchaID, p.Captcha, true) {
		return nil, errors.New("验证码错误或已失效")
	}

	user, err := s.userRepo.GetUserByUsername(p.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 校验密码
	if !utils.CheckPassword(p.Password, user.Password) {
		return nil, errors.New("账号或密码错误")
	}

	// 校验状态
	if user.Status != 1 {
		return nil, errors.New("账号已被停用")
	}

	token, err := utils.GenerateToken(user, s.jwtCfg)
	if err != nil {
		return nil, errors.New("生成token失败")
	}

	return &LoginResponse{
		Token:    token,
		UserInfo: user,
	}, nil
}

func (s *AuthService) Logout() {}
