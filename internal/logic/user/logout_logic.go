package user

import (
	"context"
	"fmt"
	"web-chat/api/http_model"
	"web-chat/pkg/auth"
	"web-chat/pkg/logger"
)

func (l *logicImpl) Logout(req *http_model.LogoutReq) error {
	if req == nil {
		return fmt.Errorf("logout request is nil")
	}
	lgr := logger.L()
	if req.Token == "" && req.RefreshToken == "" {
		return fmt.Errorf("token is required")
	}
	if req.Token != "" {
		claim := &auth.UserClaim{}
		_, err := l.svcCtx.Auth.TrackAuthToken(req.Token, claim)
		if err != nil {
			lgr.Errorf("logout token error: %v", err)
			return err
		}
		if err := l.blacklistToken(context.Background(), accessBlacklistPrefix, req.Token, claim.ExpiresAt); err != nil {
			lgr.Errorf("logout redis error: %v", err)
			return err
		}
	}
	if req.RefreshToken != "" {
		claim := &auth.UserClaim{}
		_, err := l.svcCtx.Auth.TrackRefreshToken(req.RefreshToken, claim)
		if err != nil {
			lgr.Errorf("logout refresh token error: %v", err)
			return err
		}
		if err := l.blacklistToken(context.Background(), refreshBlacklistPrefix, req.RefreshToken, claim.ExpiresAt); err != nil {
			lgr.Errorf("logout refresh redis error: %v", err)
			return err
		}
	}
	return nil
}
