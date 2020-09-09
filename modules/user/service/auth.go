package service

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// IssueTokenInput 生成token参数
	IssueTokenInput struct {
		UserID  primitive.ObjectID
		RoleIDs []int
		Conf    config.Config
	}
	// IssueTokenOutput 生成token返回参数
	IssueTokenOutput struct {
		*model.Token
	}
)

// IssueToken 生成token
func IssueToken(in *IssueTokenInput) (*IssueTokenOutput, error) {
	now := time.Now().UTC()

	token := &model.Token{
		AccessTokenExpires:  now.Add(30 * time.Minute),
		RefreshTokenExpires: now.Add(168 * time.Hour),
		UserID:              in.UserID,
		RoleIDs:             in.RoleIDs,
	}

	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessTokenClaims := accessToken.Claims.(jwt.MapClaims)
	accessTokenClaims["iss"] = "logit.co.nz"
	accessTokenClaims["exp"] = token.AccessTokenExpires.Unix()
	accessTokenClaims["sub"] = in.UserID.Hex()
	accessTokenClaims["roles"] = in.RoleIDs
	accessTokenSigningKey, _ := in.Conf.Get("system.jwt.access.key")
	accessTokenSigned, err := accessToken.SignedString([]byte(accessTokenSigningKey))
	if err != nil {
		return nil, err
	}
	token.AccessToken = accessTokenSigned

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["iss"] = "logit.co.nz"
	refreshTokenClaims["exp"] = token.RefreshTokenExpires.Unix()
	refreshTokenClaims["sub"] = in.UserID.Hex()
	refreshTokenSigningKey, _ := in.Conf.Get("system.jwt.refresh.key")
	refreshTokenSigned, err := refreshToken.SignedString([]byte(refreshTokenSigningKey))
	if err != nil {
		return nil, err
	}
	token.RefreshToken = refreshTokenSigned

	return &IssueTokenOutput{token}, nil
}

type (
	// RefreshTokenInput 更新token参数
	RefreshTokenInput struct {
		Token string `json:"token"`
		Conf  config.Config
	}
)

// RefreshToken 更新token
func RefreshToken(in *RefreshTokenInput) (*IssueTokenOutput, error) {

	// 验证参数, 获取uid
	key, _ := in.Conf.Get("system.jwt.refresh.key")
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}
	token, err := jwt.Parse(in.Token, keyFunc)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	uid, err := primitive.ObjectIDFromHex(claims["sub"].(string))
	if err != nil {
		return nil, err
	}

	// 查询用户
	user, err := model.FindUser(model.FindUserOpt{UserID: uid})
	if err != nil {
		return nil, err
	}

	// 生成token
	issueTokenOutput, err := IssueToken(&IssueTokenInput{UserID: user.ID, RoleIDs: user.RoleIDs, Conf: in.Conf})
	if err != nil {
		return nil, err
	}
	return issueTokenOutput, nil
}

// PasswordLoginInput 密码登录参数
type PasswordLoginInput struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	License  string `json:"license"`
	Password string `json:"password"`
	Conf     config.Config
}

// PasswordLogin 密码登录
func PasswordLogin(in *PasswordLoginInput) (*IssueTokenOutput, error) {
	// 如果通过license登录，则先获取Driver信息，再获取User信息
	var (
		driver *model.Driver
		err    error
	)
	if len(in.License) > 0 {
		driver, err = model.FindDriver(model.FindDriverOpt{LicenseNumber: in.License})
		if err != nil {
			return nil, err
		}
	}

	user, err := model.FindUser(model.FindUserOpt{UserID: driver.ID, Phone: in.Phone, Email: in.Email})
	if err != nil {
		return nil, err
	}
	if user.Password != in.Password {
		return nil, errors.New("Invalid credentials")
	}

	// 生成token
	issueTokenOutput, err := IssueToken(&IssueTokenInput{UserID: user.ID, RoleIDs: user.RoleIDs, Conf: in.Conf})
	if err != nil {
		return nil, err
	}
	return issueTokenOutput, nil
}

// ForgetPasswordInput 忘记密码参数
type ForgetPasswordInput struct {
	Phone    string `json:"phone" valid:"numeric,stringlength(8|11),optional"`
	Email    string `json:"email" valid:"email,optional"`
	Token    string `json:"token" valid:"required"`
	Password string `json:"password" valid:"stringlength(6|32)"`
}

// ForgetPassword 忘记密码
func ForgetPassword(in *ForgetPasswordInput) error {
	// 参数验证
	if _, err := valid.ValidateStruct(in); err != nil {
		return err
	}
	// 查找用户
	user, err := model.FindUser(model.FindUserOpt{Phone: in.Phone, Email: in.Email})
	if err != nil {
		return err
	}

	// 验证token
	var redisKey string
	switch {
	case len(in.Phone) > 0:
		redisKey = in.Phone
	case len(in.Email) > 0:
		redisKey = in.Email
	default:
		return errors.New("phone number or email is required")
	}

	if token, err := model.RedisClient.Get(redisKey).Result(); err != nil || in.Token != token {
		return errors.New("token does not match")
	}

	// 更新密码
	user.Password = in.Password
	if err := user.Update(); err != nil {
		return err
	}

	// 验证码过期处理
	model.RedisClient.ExpireAt(redisKey, time.Now())

	return nil
}
