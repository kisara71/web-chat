package user

import (
	"web-chat/api/http_model"
	"web-chat/internal/svc"
)

type Logic interface {
	Register(req *http_model.UserRegisterReq) error
	Login(req *http_model.LoginReq) (*http_model.AuthTokenResp, error)
	LoginByCode(req *http_model.LoginCodeReq) (*http_model.AuthTokenResp, error)
	RefreshToken(req *http_model.RefreshTokenReq) (*http_model.AuthTokenResp, error)
	SendEmailCode(req *http_model.SendEmailCodeReq) error
	Update(userID string, req *http_model.UserUpdateReq) error
	UpdateProfile(userID string, req *http_model.UserProfileUpdateReq) error
	Logout(req *http_model.LogoutReq) error
	GetUserInfo(userID string) (*http_model.UserInfoResp, error)
	GetProfile(userID string) (*http_model.UserProfileResp, error)
}

type logicImpl struct {
	svcCtx *svc.Context
}

func NewLogic(svcCtx *svc.Context) Logic {
	return &logicImpl{svcCtx: svcCtx}
}
