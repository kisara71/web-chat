package user

import (
	"context"
	"fmt"
	"time"
	"web-chat/api/http_model"
	"web-chat/pkg/auth"
	"web-chat/pkg/logger"
)

const blacklistPrefix = "jwt:blacklist:"

func (l *logicImpl) Logout(req *http_model.LogoutReq) error {
	if req == nil || req.Token == "" {
		return fmt.Errorf("token is required")
	}
	lgr := logger.L()
	claim := &auth.UserClaim{}
	_, err := l.svcCtx.Auth.TrackAuthToken(req.Token, claim)
	if err != nil {
		lgr.Errorf("logout token error: %v", err)
		return err
	}
	ttl := authTokenTTL
	if claim.ExpiresAt > 0 {
		expireAt := time.Unix(claim.ExpiresAt, 0)
		if remain := time.Until(expireAt); remain > 0 {
			ttl = remain
		} else {
			ttl = time.Minute
		}
	}
	key := blacklistPrefix + req.Token
	if err := l.svcCtx.Infra.Redis.Set(context.Background(), key, "1", ttl).Err(); err != nil {
		lgr.Errorf("logout redis error: %v", err)
		return err
	}
	return nil
}
