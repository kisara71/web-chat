package user

import (
	"fmt"
	"time"
	"web-chat/api/http_model"
	"web-chat/pkg/auth"
	"web-chat/pkg/logger"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const authTokenTTL = 24 * time.Hour

func (l *logicImpl) Login(req *http_model.LoginReq) (string, error) {
	if req == nil {
		return "", fmt.Errorf("login request is nil")
	}
	lgr := logger.L()
	var (
		entity *authUser
		err    error
	)
	entity, err = l.fetchLoginUser(req.Account)
	if err != nil {
		lgr.Printf("login fetch user error: %v", err)
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(entity.Password), []byte(req.Password)); err != nil {
		lgr.Printf("login password invalid: %v", err)
		return "", fmt.Errorf("password is invalid")
	}
	claim := auth.UserClaim{
		UserID: entity.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(authTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token, err := l.svcCtx.Auth.GenAuthToken(&claim)
	if err != nil {
		lgr.Printf("login token sign error: %v", err)
		return "", err
	}
	return token, nil
}

type authUser struct {
	ID       int64
	Password string
}

func (l *logicImpl) fetchLoginUser(account string) (*authUser, error) {
	ok, err := l.svcCtx.Utils.Regexp.ValidatePhone(account)
	if err != nil {
		return nil, err
	}
	if ok {
		entity, err := l.svcCtx.Dao.UserDao.GetUserByPhone(account)
		if err != nil {
			return nil, err
		}
		return &authUser{ID: entity.ID, Password: entity.Password}, nil
	}
	ok, err = l.svcCtx.Utils.Regexp.ValidateEmail(account)
	if err != nil {
		return nil, err
	}
	if ok {
		entity, err := l.svcCtx.Dao.UserDao.GetUserByEmail(account)
		if err != nil {
			return nil, err
		}
		return &authUser{ID: entity.ID, Password: entity.Password}, nil
	}
	return nil, fmt.Errorf("account is invalid")
}
