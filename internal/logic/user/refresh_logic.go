package user

import (
	"context"
	"fmt"
	"time"
	"web-chat/api/http_model"
	"web-chat/pkg/auth"
	"web-chat/pkg/logger"
)

func (l *logicImpl) RefreshToken(req *http_model.RefreshTokenReq) (*http_model.AuthTokenResp, error) {
	if req == nil || req.RefreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}
	lgr := logger.L()
	claim := &auth.UserClaim{}
	_, err := l.svcCtx.Auth.TrackRefreshToken(req.RefreshToken, claim)
	if err != nil {
		lgr.Errorf("refresh token verify error: %v", err)
		return nil, err
	}
	revoked, err := l.isTokenRevoked(context.Background(), refreshBlacklistPrefix, req.RefreshToken)
	if err != nil {
		lgr.Errorf("refresh token redis error: %v", err)
		return nil, err
	}
	if revoked {
		return nil, fmt.Errorf("refresh token revoked")
	}
	if claim.UserID == "" {
		return nil, fmt.Errorf("invalid refresh token")
	}
	if claim.ExpiresAt > 0 && time.Now().After(time.Unix(claim.ExpiresAt, 0)) {
		return nil, fmt.Errorf("refresh token expired")
	}
	if err := l.blacklistToken(context.Background(), refreshBlacklistPrefix, req.RefreshToken, claim.ExpiresAt); err != nil {
		lgr.Errorf("refresh token blacklist error: %v", err)
		return nil, err
	}
	tokens, err := l.issueTokens(claim.UserID)
	if err != nil {
		lgr.Errorf("refresh token issue error: %v", err)
		return nil, err
	}
	return &http_model.AuthTokenResp{
		Token:        tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
