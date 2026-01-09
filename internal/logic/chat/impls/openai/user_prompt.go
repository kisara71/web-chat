package openai

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

func (l *logicImpl) BuildUserSystemPrompt(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", nil
	}
	profile, err := l.svcCtx.Dao.ProfileDao.GetUserProfileByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	parts := make([]string, 0, 3)
	if strings.TrimSpace(profile.SystemPrompt) != "" {
		parts = append(parts, strings.TrimSpace(profile.SystemPrompt))
	}
	if strings.TrimSpace(profile.ModelPreference) != "" {
		parts = append(parts, "Model preferences: "+strings.TrimSpace(profile.ModelPreference))
	}
	if strings.TrimSpace(profile.Traits) != "" {
		parts = append(parts, "User profile: "+strings.TrimSpace(profile.Traits))
	}
	if len(parts) == 0 {
		return "", nil
	}
	return strings.Join(parts, "\n\n"), nil
}
