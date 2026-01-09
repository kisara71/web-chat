package user

import (
	"context"
	"fmt"
	"time"
	"web-chat/pkg/auth"

	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
)

const (
	authTokenTTL           = 24 * time.Hour
	refreshTokenTTL        = 7 * 24 * time.Hour
	accessBlacklistPrefix  = "jwt:blacklist:"
	refreshBlacklistPrefix = "jwt:refresh:blacklist:"
)

func (l *logicImpl) issueTokens(userID string) (*authTokens, error) {
	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}
	now := time.Now()
	accessClaim := auth.UserClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(authTokenTTL).Unix(),
			IssuedAt:  now.Unix(),
		},
	}
	refreshClaim := auth.UserClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(refreshTokenTTL).Unix(),
			IssuedAt:  now.Unix(),
		},
	}
	accessToken, err := l.svcCtx.Auth.GenAuthToken(&accessClaim)
	if err != nil {
		return nil, err
	}
	refreshToken, err := l.svcCtx.Auth.GenRefreshToken(&refreshClaim)
	if err != nil {
		return nil, err
	}
	return &authTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type authTokens struct {
	AccessToken  string
	RefreshToken string
}

func (l *logicImpl) blacklistToken(ctx context.Context, prefix, token string, expiresAt int64) error {
	ttl := time.Minute
	if expiresAt > 0 {
		expireAt := time.Unix(expiresAt, 0)
		if remain := time.Until(expireAt); remain > 0 {
			ttl = remain
		}
	}
	key := prefix + token
	return l.svcCtx.Infra.Redis.Set(ctx, key, "1", ttl).Err()
}

func (l *logicImpl) isTokenRevoked(ctx context.Context, prefix, token string) (bool, error) {
	val, err := l.svcCtx.Infra.Redis.Get(ctx, prefix+token).Result()
	if err == nil && val != "" {
		return true, nil
	}
	if err != nil && err != redis.Nil {
		return false, err
	}
	return false, nil
}
