package user

import (
	"context"
	"fmt"
	"web-chat/api/http_model"
	"web-chat/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

func (l *logicImpl) Login(req *http_model.LoginReq) (*http_model.AuthTokenResp, error) {
	if req == nil {
		return nil, fmt.Errorf("login request is nil")
	}
	lgr := logger.L()
	var (
		entity *authUser
		err    error
	)
	entity, err = l.fetchLoginUser(req.Account)
	if err != nil {
		lgr.Errorf("login fetch user error: %v", err)
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(entity.Password), []byte(req.Password)); err != nil {
		lgr.Errorf("login password invalid: %v", err)
		return nil, fmt.Errorf("password is invalid")
	}
	tokens, err := l.issueTokens(entity.UUID)
	if err != nil {
		lgr.Errorf("login token sign error: %v", err)
		return nil, err
	}
	return &http_model.AuthTokenResp{
		Token:        tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (l *logicImpl) LoginByCode(req *http_model.LoginCodeReq) (*http_model.AuthTokenResp, error) {
	if req == nil {
		return nil, fmt.Errorf("login request is nil")
	}
	lgr := logger.L()
	ok, err := l.svcCtx.Utils.Regexp.ValidateEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if !ok {
		if okPhone, err := l.svcCtx.Utils.Regexp.ValidatePhone(req.Email); err == nil && okPhone {
			return nil, fmt.Errorf("sms is not supported")
		}
		return nil, fmt.Errorf("email is invalid")
	}

	ok, err = l.svcCtx.Utils.Code.VerifyCode(context.Background(), req.Code, req.Email)
	if err != nil {
		lgr.Errorf("login code verify error: %v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("code is invalid")
	}

	entity, err := l.svcCtx.Dao.UserDao.GetUserByEmail(req.Email)
	if err != nil {
		lgr.Errorf("login fetch user error: %v", err)
		return nil, err
	}
	tokens, err := l.issueTokens(entity.UUID)
	if err != nil {
		lgr.Errorf("login token sign error: %v", err)
		return nil, err
	}
	return &http_model.AuthTokenResp{
		Token:        tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

type authUser struct {
	UUID     string
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
		return &authUser{UUID: entity.UUID, Password: entity.Password}, nil
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
		return &authUser{UUID: entity.UUID, Password: entity.Password}, nil
	}
	return nil, fmt.Errorf("account is invalid")
}
